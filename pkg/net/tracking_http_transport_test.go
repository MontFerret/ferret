package net

import (
	"io"
	stdhttp "net/http"
	"strings"
	"sync/atomic"
)

type trackingHTTPTransport struct {
	body       string
	calls      atomic.Int64
	idleCloses atomic.Int64
}

func (t *trackingHTTPTransport) RoundTrip(req *stdhttp.Request) (*stdhttp.Response, error) {
	t.calls.Add(1)

	return &stdhttp.Response{
		StatusCode: stdhttp.StatusOK,
		Status:     stdhttp.StatusText(stdhttp.StatusOK),
		Header:     make(stdhttp.Header),
		Body:       io.NopCloser(strings.NewReader(t.body)),
		Request:    req,
	}, nil
}

func (t *trackingHTTPTransport) CloseIdleConnections() {
	t.idleCloses.Add(1)
}

func (t *trackingHTTPTransport) callCount() int64 {
	return t.calls.Load()
}

func (t *trackingHTTPTransport) idleCloseCount() int64 {
	return t.idleCloses.Load()
}
