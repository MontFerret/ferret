package network

import (
	"context"
	"github.com/MontFerret/ferret/pkg/drivers/cdp/dom"
	"github.com/mafredri/cdp/protocol/page"
	"github.com/wI2L/jettison"

	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
)

var NavigationEventType = core.NewType("ferret.drivers.cdp.network.NavigationEvent")

type NavigationEvent struct {
	URL      string
	FrameID  page.FrameID
	MimeType string
}

func (evt *NavigationEvent) MarshalJSON() ([]byte, error) {
	if evt == nil {
		return values.None.MarshalJSON()
	}

	return jettison.MarshalOpts(map[string]string{
		"url":      evt.URL,
		"frame_id": string(evt.FrameID),
	}, jettison.NoHTMLEscaping())
}

func (evt *NavigationEvent) Type() core.Type {
	return NavigationEventType
}

func (evt *NavigationEvent) String() string {
	return evt.URL
}

func (evt *NavigationEvent) Compare(other core.Value) int64 {
	if other.Type() != NavigationEventType {
		return -1
	}

	otherEvt := other.(*NavigationEvent)
	comp := values.NewString(evt.URL).Compare(values.NewString(otherEvt.URL))

	if comp != 0 {
		return comp
	}

	return values.String(evt.FrameID).Compare(values.String(otherEvt.FrameID))
}

func (evt *NavigationEvent) Unwrap() interface{} {
	return evt
}

func (evt *NavigationEvent) Hash() uint64 {
	return values.Parse(evt).Hash()
}

func (evt *NavigationEvent) Copy() core.Value {
	return *(&evt)
}

func (evt *NavigationEvent) GetIn(_ context.Context, path []core.Value) (core.Value, core.PathError) {
	if len(path) == 0 {
		return evt, nil
	}

	switch path[0].String() {
	case "url", "URL":
		return values.NewString(evt.URL), nil
	case "frame":
		return dom.NewFrameID(evt.FrameID), nil
	default:
		return values.None, nil
	}
}
