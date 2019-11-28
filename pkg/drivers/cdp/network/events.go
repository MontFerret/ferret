package network

import "github.com/MontFerret/ferret/pkg/drivers/cdp/events"

var (
	eventLoad      = events.New("load")
	eventFrameLoad = events.New("frame_load")
)
