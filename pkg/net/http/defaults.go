package http

import "time"

const (
	defaultTimeout                         = 30 * time.Second
	defaultMaxIdleConnections              = 100
	defaultMaxIdleConnectionsPerHost       = 10
	defaultMaxConnectionsPerHost           = 20
	defaultMaxRedirects                    = 10
	defaultMaxRequestSize            int64 = 16 << 20
	defaultMaxResponseSize           int64 = 16 << 20
	defaultMaxResponseHeaderSize     int64 = 1 << 20
)
