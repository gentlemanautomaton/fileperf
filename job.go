package fileperf

import (
	"context"
	"io"
	"io/fs"
	"path"
	"time"
)

type scanJob struct {
	// Internal job state
	root   Dir
	ch     chan<- fileIterUpdate
	cancel context.CancelFunc

	// External job requirements
	include, exclude []Pattern
	sendSkipped      bool

	// Job statistics and tallies
	stats JobStats
}

func executeJob(ctx context.Context, job scanJob) {
	// Close the update channel when finished
	defer close(job.ch)

	// Make sure the cancellation function always gets triggered as clean up
	defer job.cancel()

	// Walk each file in the directory
	err := fs.WalkDir(job.root, ".", func(p string, d fs.DirEntry, dirErr error) error {
		// Stop walking the directory if the job has been cancelled
		if err := ctx.Err(); err != nil {
			return err
		}

		// Ignore the root directory itself
		if p == "." {
			return nil
		}

		// Prepare the file object with our results
		file := File{
			Root:  job.root,
			Path:  p,
			Index: job.stats.Skipped + job.stats.Scanned,
		}

		// Skip this file if it doesn't pass our file name pattern matching
		// filters
		{
			_, name := path.Split(p)
			skip := false

			// Handle exclusions
			for _, pattern := range job.exclude {
				if pattern.Expression.MatchString(name) {
					skip = true
					break
				}
			}

			// Handle inclusions
			if !skip && len(job.include) > 0 {
				matched := false
				for _, pattern := range job.include {
					if pattern.Expression.MatchString(name) {
						matched = true
						break
					}
				}
				if !matched {
					skip = true
				}
			}

			// Record skipped jobs and carry on
			if skip {
				job.stats.Skipped++
				if job.sendSkipped {
					file.Skipped = true
					select {
					case <-ctx.Done():
						return ctx.Err()
					case job.ch <- fileIterUpdate{file: file, stats: job.stats, updated: time.Now()}:
						return nil
					}
				}
				return nil
			}
		}

		// Increment our scanned file count
		job.stats.Scanned++

		// If an error was reported, such as access denied, record it as a
		// scan error
		if dirErr != nil {
			file.Err = dirErr
		} else {
			// Attempt to collect more information about the file
			info, err := d.Info()
			if err != nil {
				file.Err = err
			} else {
				file.Name = info.Name()
				file.Size = info.Size()
				file.Mode = info.Mode()
				file.ModTime = info.ModTime()
			}

		}

		// If we haven't encountered an error, attempt to read the file
		if file.Err == nil {
			var data []byte
			file.Start = time.Now()
			data, file.Err = fs.ReadFile(job.root, file.Path)
			file.End = time.Now()
			job.stats.ElapsedRead += file.End.Sub(file.Start)
			job.stats.TotalBytes += int64(len(data))
			if file.Err == nil {
				job.stats.Read++
			}
		}

		// Tally error statistics
		if file.Err != nil {
			job.stats.Errors++
		}

		// Send files to the iterator via the job's channel
		select {
		case <-ctx.Done():
			return ctx.Err()
		case job.ch <- fileIterUpdate{file: file, stats: job.stats, updated: time.Now()}:
			return nil
		}
	})

	// If no error was encountered, send io.EOF in the last update so it
	// doesn't get processed as an incoming file update by the file iterator
	if err == nil {
		err = io.EOF
	}

	// Always provide a final update with the completed statistics
	job.ch <- fileIterUpdate{streamErr: err, stats: job.stats, updated: time.Now()}
}
