package network

import "github.com/MontFerret/ferret/pkg/drivers/cdp/events"

var (
	eventFrameLoad   = events.New("frame_load")
	responseReceived = events.New("response_received")
)
