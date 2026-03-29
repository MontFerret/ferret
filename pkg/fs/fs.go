package fs

import (
	"io/fs"
	"os"
)

type (
	Reader interface {
		ReadFile(path string) ([]byte, error)
		Open(path string) (fs.File, error)
		Stat(path string) (fs.FileInfo, error)
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

	// FileSystem is an interface that combines file reading, directory operations, file writing, and file removal capabilities.
	// It provides a unified interface for accessing files and directories in a filesystem-based environment.
	FileSystem interface {
		Reader
		Directories
		Writer
		Remover
	}
)

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
