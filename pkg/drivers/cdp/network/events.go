package network

import "github.com/MontFerret/ferret/pkg/drivers/cdp/events"

var (
	EventLoad   = events.ToType("load")
	EventReload = events.ToType("reload")
)
