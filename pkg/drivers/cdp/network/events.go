package network

import "github.com/MontFerret/ferret/pkg/drivers/cdp/events"

var (
	eventLoad      = events.ToType("load")
	eventFrameLoad = events.ToType("frame_load")
)
