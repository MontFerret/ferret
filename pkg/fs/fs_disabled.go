package fs

import (
	"io/fs"
)

type disabledFS struct{}

var disabledFileSystem = disabledFS{}

func (n disabledFS) ReadFile(_ string) ([]byte, error) {
	return nil, ErrRootNotConfigured
}

func (n disabledFS) Open(_ string) (fs.File, error) {
	return nil, ErrRootNotConfigured
}

func (n disabledFS) Stat(_ string) (fs.FileInfo, error) {
	return nil, ErrRootNotConfigured
}

func (n disabledFS) Exists(_ string) (bool, error) {
	return false, ErrRootNotConfigured
}

func (n disabledFS) Mkdir(_ string, _ fs.FileMode) error {
	return ErrRootNotConfigured
}

func (n disabledFS) MkdirAll(_ string, _ fs.FileMode) error {
	return ErrRootNotConfigured
}

func (n disabledFS) WriteFile(_ string, _ []byte, _ fs.FileMode) error {
	return ErrRootNotConfigured
}

func (n disabledFS) AppendFile(_ string, _ []byte, _ fs.FileMode) error {
	return ErrRootNotConfigured
}

func (n disabledFS) Remove(_ string) error {
	return ErrRootNotConfigured
}

func (n disabledFS) RemoveAll(_ string) error {
	return ErrRootNotConfigured
}
