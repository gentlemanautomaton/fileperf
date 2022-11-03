package fileperf

import (
	"context"
	"time"
)

// Scanner scans a set of files for issues identified by its issue handlers.
type Scanner struct {
	// Handlers examine each file passing through the scanner's filters and
	// determine whether they have issues. They determine what constitutes an "issue".
	//Handlers []IssueHandler

	// Include is a filter that limits the number of files scanned. If
	// provided, only files with names matching at least one pattern will
	// be scanned.
	Include []Pattern

	// Exclude is a filter that limits the number of files scanned. If
	// provided, only files with names that don't match any of its patterns
	// will be scanned.
	Exclude []Pattern
}

// ScanDir causes the scanner to scan the given file system directory.
func (s Scanner) ScanDir(root Dir) *FileIter {
	// Prepare a cancellation function that the iterator can use to stop the
	// job.
	ctx, cancel := context.WithCancel(context.Background())

	// Prepare a communications channel
	ch := make(chan fileIterUpdate)

	// Note the start time
	now := time.Now()

	// Prepare a file iterator that will be returned
	iter := FileIter{
		start:  now,
		ch:     ch,
		cancel: cancel,
		end:    now,
	}

	// Prepare a job
	job := scanJob{
		root:    root,
		ch:      ch,
		cancel:  cancel,
		include: s.Include,
		exclude: s.Exclude,
	}

	// Execute the job
	go executeJob(ctx, job)

	// Return the file iterator
	return &iter
}

// ScanDir scans the given file system directory for issues.
func ScanDir(ctx context.Context, root Dir) *FileIter {
	return Scanner{}.ScanDir(root)
}
