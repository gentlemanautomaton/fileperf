package fileperf

import (
	"fmt"
	"time"

	"github.com/dustin/go-humanize"
)

// JobStats report scanning tallies during and at the completion of scanning.
type JobStats struct {
	// ElapsedRead is the total time spent reading files.
	ElapsedRead time.Duration

	// TotalBytes is the total number of bytes read.
	TotalBytes int64

	// Read is the number of files read without issue.
	Read int

	// Errors is the number of files that encountered an error.
	Errors int

	// Scanned is the number of files scanned.
	Scanned int

	// Skipped is the number of files not scanned due to filters.
	Skipped int
}

// String returns a string representation of the job statistics.
func (s JobStats) String() string {
	return fmt.Sprintf("%s cumulative read time, %s read, %d files read, %d errors, %d files skipped", s.ElapsedRead, humanize.Bytes(uint64(s.TotalBytes)), s.Read, s.Errors, s.Skipped)
}
