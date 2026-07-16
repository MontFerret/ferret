package ferret

import ferretnet "github.com/MontFerret/ferret/v2/pkg/net"

type idleConnectionCloser interface {
	CloseIdleConnections()
}

func closeIdleNetworkConnections(network ferretnet.Network) {
	if closer, ok := network.(idleConnectionCloser); ok {
		closer.CloseIdleConnections()
	}
}
