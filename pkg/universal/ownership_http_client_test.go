package universal

import (
	"context"
	"errors"
	"sync/atomic"

	ferrethttp "github.com/MontFerret/ferret/v2/pkg/net/http"
)

type ownershipHTTPClient struct {
	closes atomic.Int32
}

func (c *ownershipHTTPClient) Do(context.Context, *ferrethttp.Request) (*ferrethttp.Response, error) {
	return nil, errors.New("unexpected HTTP request")
}

func (c *ownershipHTTPClient) CloseIdleConnections() {
	c.closes.Add(1)
}
