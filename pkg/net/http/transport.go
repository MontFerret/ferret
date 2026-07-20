package http

import (
	stdhttp "net/http"
	"time"
)

func newPolicyTransport(dialer *policyDialer, maxResponseHeaderSize int64) *stdhttp.Transport {
	return &stdhttp.Transport{
		Proxy:                  nil,
		DialContext:            dialer.DialContext,
		ForceAttemptHTTP2:      true,
		MaxIdleConns:           defaultMaxIdleConnections,
		MaxIdleConnsPerHost:    defaultMaxIdleConnectionsPerHost,
		MaxConnsPerHost:        defaultMaxConnectionsPerHost,
		MaxResponseHeaderBytes: maxResponseHeaderSize,
		IdleConnTimeout:        90 * time.Second,
		TLSHandshakeTimeout:    10 * time.Second,
		ExpectContinueTimeout:  time.Second,
	}
}
