package net

import (
	"context"

	ferrethttp "github.com/MontFerret/ferret/v2/pkg/net/http"
)

type stubHTTPClient struct{}

func (stubHTTPClient) Do(context.Context, *ferrethttp.Request) (*ferrethttp.Response, error) {
	return &ferrethttp.Response{StatusCode: 200}, nil
}
