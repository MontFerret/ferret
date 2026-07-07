package http_test

import (
	"context"

	ferretnet "github.com/MontFerret/ferret/v2/pkg/net"
)

func testNetworkContext() context.Context {
	return ferretnet.WithNetwork(context.Background(), ferretnet.New())
}
