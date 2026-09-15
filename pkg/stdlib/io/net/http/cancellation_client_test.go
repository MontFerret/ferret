package http_test

import (
	"context"

	ferrethttp "github.com/MontFerret/ferret/v2/pkg/net/http"
)

type blockingHTTPClient struct {
	started chan context.Context
}

func (client *blockingHTTPClient) Do(ctx context.Context, _ *ferrethttp.Request) (*ferrethttp.Response, error) {
	client.started <- ctx
	<-ctx.Done()

	return nil, ctx.Err()
}
