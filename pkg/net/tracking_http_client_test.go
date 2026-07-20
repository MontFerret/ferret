package net

import (
	"context"
	"sync/atomic"

	ferrethttp "github.com/MontFerret/ferret/v2/pkg/net/http"
)

type trackingHTTPClient struct {
	idleCloses atomic.Int64
}

func (c *trackingHTTPClient) Do(context.Context, *ferrethttp.Request) (*ferrethttp.Response, error) {
	return &ferrethttp.Response{StatusCode: 200}, nil
}

func (c *trackingHTTPClient) CloseIdleConnections() {
	c.idleCloses.Add(1)
}

func (c *trackingHTTPClient) idleCloseCount() int64 {
	return c.idleCloses.Load()
}
