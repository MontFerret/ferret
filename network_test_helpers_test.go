package ferret

import (
	"testing"

	ferretnet "github.com/MontFerret/ferret/v2/pkg/net"
)

func mustNewTestNetwork(t testing.TB, setters ...ferretnet.Option) ferretnet.Network {
	t.Helper()

	network, err := ferretnet.New(setters...)
	if err != nil {
		t.Fatalf("create test network: %v", err)
	}

	return network
}
