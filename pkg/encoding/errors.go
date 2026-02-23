package encoding

import "errors"

var (
	ErrNilRegistry      = errors.New("encoding: registry is nil")
	ErrNilCodec         = errors.New("encoding: codec is nil")
	ErrEmptyContentType = errors.New("encoding: content type is empty")
	ErrCodecNotFound    = errors.New("encoding: codec not found")
	ErrRegistryNotFound = errors.New("encoding: registry not found in context")
)
