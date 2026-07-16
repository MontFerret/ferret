package http

import (
	stdhttp "net/http"
	"time"
)

func newPolicyTransport(dialer *policyDialer) *stdhttp.Transport {
	return &stdhttp.Transport{
		Proxy:                 nil,
		DialContext:           dialer.DialContext,
		ForceAttemptHTTP2:     true,
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: time.Second,
	}
}
