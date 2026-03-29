package fs

import (
	"io"
	"io/fs"
	"os"
)

type rootFS struct {
	root    *os.Root
	reaOnly bool
}

func (r *rootFS) Close() error {
	return r.root.Close()
}

func (r *rootFS) ReadFile(path string) ([]byte, error) {
	f, err := r.root.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	return io.ReadAll(f)
}

func (r *rootFS) Open(path string) (fs.File, error) {
	return r.root.Open(path)
}

func (r *rootFS) Stat(path string) (fs.FileInfo, error) {
	return r.root.Stat(path)
}

func (r *rootFS) Exists(path string) (bool, error) {
	_, err := r.root.Stat(path)
	if err == nil {
		return true, nil
	}

	if os.IsNotExist(err) {
		return false, nil
	}

	return false, err
}

func (r *rootFS) WriteFile(path string, data []byte, perm fs.FileMode) error {
	if r.reaOnly {
		return ErrReadOnly
	}

	f, err := r.root.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, perm)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.Write(data)
	return err
}

func (r *rootFS) AppendFile(path string, data []byte, perm fs.FileMode) error {
	if r.reaOnly {
		return ErrReadOnly
	}

	f, err := r.root.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, perm)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.Write(data)
	return err
}

func (r *rootFS) Mkdir(path string, perm fs.FileMode) error {
	if r.reaOnly {
		return ErrReadOnly
	}

	return r.root.Mkdir(path, perm)
}

func (r *rootFS) MkdirAll(path string, perm fs.FileMode) error {
	if r.reaOnly {
		return ErrReadOnly
	}

	return r.root.MkdirAll(path, perm)
}

func (r *rootFS) Remove(path string) error {
	if r.reaOnly {
		return ErrReadOnly
	}

	return r.root.Remove(path)
}

func (r *rootFS) RemoveAll(path string) error {
	if r.reaOnly {
		return ErrReadOnly
	}

	return r.root.RemoveAll(path)
}
