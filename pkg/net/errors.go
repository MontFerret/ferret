package net

import "errors"

var (
	// ErrNotFound indicates no Network is available in the provided context.
	ErrNotFound = errors.New("network not found in context")
)
