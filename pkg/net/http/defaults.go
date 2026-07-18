package http

import "time"

const (
	defaultTimeout               = 30 * time.Second
	defaultMaxResponseHeaderSize = 1 << 20
)
