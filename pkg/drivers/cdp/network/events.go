package network

import "github.com/MontFerret/ferret/pkg/drivers/cdp/events"

var (
	eventFrameLoad = events.New("frame_load")
	contentReady   = events.New("content_ready")
)
