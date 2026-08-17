package fs

import (
	"io"
	"io/fs"
	"os"
)

type (
	Reader interface {
		ReadFile(path string) ([]byte, error)
		Open(path string) (ReadableFile, error)
		OpenFile(path string, flag int, perm fs.FileMode) (WritableFile, error)
		Stat(path string) (fs.FileInfo, error)
		Lstat(path string) (fs.FileInfo, error)
		Exists(path string) (bool, error)
	}

	Directories interface {
		Mkdir(path string, perm fs.FileMode) error
		MkdirAll(path string, perm fs.FileMode) error
	}

	Writer interface {
		WriteFile(path string, data []byte, perm fs.FileMode) error
		AppendFile(path string, data []byte, perm fs.FileMode) error
	}

	Remover interface {
		Remove(path string) error
		RemoveAll(path string) error
	}

	// FileSystem combines file reading, directory operations, writing, removal,
	// and lifecycle management. The caller owns every FileSystem returned by New
	// and must close it when it is no longer needed.
	FileSystem interface {
		Reader
		Directories
		Writer
		Remover
		io.Closer
	}
)

// New creates a filesystem from the provided options. The caller owns the
// returned filesystem and must close it when it is no longer needed.
func New(setters ...Option) (FileSystem, error) {
	opts := &options{
		Root:     "",
		ReadOnly: false,
	}

	for _, opt := range setters {
		opt(opts)
	}

	if opts.Root == "" {
		return disabledFileSystem, nil
	}

	r, err := os.OpenRoot(opts.Root)
	if err != nil {
		return nil, err
	}

	return &rootFS{root: r, readOnly: opts.ReadOnly}, nil
}
