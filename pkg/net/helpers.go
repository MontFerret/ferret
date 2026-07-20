package net

import ferrethttp "github.com/MontFerret/ferret/v2/pkg/net/http"

func CloseIdleNetworkConnections(network Network) {
	if closer, ok := network.(ferrethttp.IdleConnectionCloser); ok && closer != nil {
		closer.CloseIdleConnections()
	}
}
