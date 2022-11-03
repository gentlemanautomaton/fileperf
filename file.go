package fileperf

import (
	"fmt"
	"io/fs"
	"time"

	"github.com/dustin/go-humanize"
)

// File describes a file that has been scanned.
type File struct {
	// Scanned file location
	Root    Dir
	Path    string
	Index   int
	Skipped bool

	// Error handling
	Err error

	// FileInfo values collected during a scan (may be empty)
	Name    string
	Size    int64
	Mode    fs.FileMode
	ModTime time.Time

	// File performance results
	Start, End time.Time
}

// String returns a string representation of f, including its index and path.
func (f File) String() string {
	if f.Err != nil {
		return fmt.Sprintf("[%d]: \"%s\": %v", f.Index, f.Path, f.Err)
	}
	return fmt.Sprintf("[%d]: \"%s\": %s (%s)", f.Index, f.Path, f.End.Sub(f.Start), humanize.Bytes(uint64(f.Size)))
}
