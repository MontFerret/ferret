package ferret

import (
	"context"
	"sync"

	nethttp "github.com/MontFerret/ferret/v2/pkg/net/http"
)

type recordingHTTPClient struct {
	lastURL string
	body    []byte
	calls   int
	mu      sync.Mutex
}

func (c *recordingHTTPClient) Do(_ context.Context, req *nethttp.Request) (*nethttp.Response, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.calls++
	if req != nil {
		c.lastURL = req.URL
	}

	return &nethttp.Response{
		StatusCode: 200,
		Status:     "200 OK",
		Body:       append([]byte(nil), c.body...),
	}, nil
}
