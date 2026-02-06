package optimization

import "errors"

var (
	ErrCFGBuildFailed = errors.New("failed to build control flow graph")
	ErrPassFailed     = errors.New("optimization pass failed")
)
