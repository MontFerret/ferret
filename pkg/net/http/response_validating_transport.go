package http

import stdhttp "net/http"

// responseValidatingTransport turns a broken RoundTripper's nil response into
// Ferret's stable structural error before net/http replaces it with a generic
// contract error.
type responseValidatingTransport struct {
	next stdhttp.RoundTripper
}

func newResponseValidatingTransport(next stdhttp.RoundTripper) stdhttp.RoundTripper {
	if next == nil {
		next = stdhttp.DefaultTransport
	}

	// The standard transport enforces its own response contract. Keeping it
	// unwrapped also preserves net/http's optimized cancellation path.
	switch next.(type) {
	case *stdhttp.Transport, *responseValidatingTransport:
		return next
	}

	return &responseValidatingTransport{next: next}
}

func (t *responseValidatingTransport) RoundTrip(req *stdhttp.Request) (*stdhttp.Response, error) {
	res, err := t.next.RoundTrip(req)
	if res == nil && err == nil {
		return nil, ErrNilResponse
	}

	return res, err
}

func (t *responseValidatingTransport) CloseIdleConnections() {
	if closer, ok := t.next.(interface{ CloseIdleConnections() }); ok {
		closer.CloseIdleConnections()
	}
}

// CancelRequest preserves cancellation for legacy custom transports that do
// not consume Request.Context.
func (t *responseValidatingTransport) CancelRequest(req *stdhttp.Request) {
	if canceler, ok := t.next.(interface{ CancelRequest(*stdhttp.Request) }); ok {
		canceler.CancelRequest(req)
	}
}
