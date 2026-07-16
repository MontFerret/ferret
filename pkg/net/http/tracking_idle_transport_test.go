package http

import (
	stdhttp "net/http"
	"sync/atomic"
)

type trackingIdleTransport struct {
	closed atomic.Bool
}

func (t *trackingIdleTransport) RoundTrip(*stdhttp.Request) (*stdhttp.Response, error) {
	return responseWithBody(stdhttp.StatusOK, "ok", nil), nil
}

func (t *trackingIdleTransport) CloseIdleConnections() {
	t.closed.Store(true)
}
