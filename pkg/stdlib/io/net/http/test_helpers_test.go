package http_test

import (
	"context"
	"testing"

	ferretnet "github.com/MontFerret/ferret/v2/pkg/net"
)

func testNetworkContext(t testing.TB) context.Context {
	t.Helper()

	network, err := ferretnet.New(ferretnet.WithHTTPClient(testHTTPClient{}))
	if err != nil {
		t.Fatalf("create test network: %v", err)
	}

	return ferretnet.WithNetwork(
		context.Background(),
		network,
	)
}
