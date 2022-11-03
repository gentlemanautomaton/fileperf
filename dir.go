package fileperf

import (
	"io/fs"
	"os"
	"path"
)

// Dir is a file directory path accessible via operating system API acalls.
type Dir string

// Open opens the named file.
func (dir Dir) Open(name string) (fs.File, error) {
	return os.DirFS(string(dir)).Open(name)
}

// Stat returns a FileInfo describing the file.
func (dir Dir) Stat(name string) (fs.FileInfo, error) {
	return os.DirFS(string(dir)).(fs.StatFS).Stat(name)
}

// FilePath returns the full path of the given file name by joining it
// with dir.
func (dir Dir) FilePath(name string) string {
	return path.Join(string(dir), name)
}
