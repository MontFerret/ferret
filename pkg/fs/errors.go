package fs

import "errors"

var (
	ErrReadOnly          = errors.New("filesystem is read-only")
	ErrNotFound          = errors.New("filesystem not found in context")
	ErrRootNotConfigured = errors.New("filesystem root is not configured")
)
