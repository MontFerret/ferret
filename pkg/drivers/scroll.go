package drivers

import (
	"github.com/MontFerret/ferret/pkg/runtime/values"
	"github.com/wI2L/jettison"
	"strings"
)

// ScrollBehavior defines the transition animation.
// In HTML specification, default value is auto, but in Ferret it's instant.
// More details here https://developer.mozilla.org/en-US/docs/Web/API/Element/scrollIntoView
type ScrollBehavior int

const (
	ScrollBehaviorInstant ScrollBehavior = 0
	ScrollBehaviorSmooth  ScrollBehavior = 1
	ScrollBehaviorAuto    ScrollBehavior = 2
)

func NewScrollBehavior(value string) ScrollBehavior {
	switch strings.ToLower(value) {
	case "instant":
		return ScrollBehaviorInstant
	case "smooth":
		return ScrollBehaviorSmooth
	default:
		return ScrollBehaviorAuto
	}
}

func (b ScrollBehavior) MarshalJSON() ([]byte, error) {
	return jettison.MarshalOpts(b.String(), jettison.NoHTMLEscaping())
}

func (b ScrollBehavior) String() string {
	switch b {
	case ScrollBehaviorInstant:
		return "instant"
	case ScrollBehaviorSmooth:
		return "smooth"
	default:
		return "auto"
	}
}

// ScrollVerticalAlignment defines vertical alignment after scrolling.
// In HTML specification, default value is start, but in Ferret it's center.
// More details here https://developer.mozilla.org/en-US/docs/Web/API/Element/scrollIntoView
type ScrollVerticalAlignment int

const (
	ScrollVerticalAlignmentCenter  ScrollVerticalAlignment = 0
	ScrollVerticalAlignmentStart   ScrollVerticalAlignment = 1
	ScrollVerticalAlignmentEnd     ScrollVerticalAlignment = 2
	ScrollVerticalAlignmentNearest ScrollVerticalAlignment = 3
)

func NewScrollVerticalAlignment(value string) ScrollVerticalAlignment {
	switch strings.ToLower(value) {
	case "center":
		return ScrollVerticalAlignmentCenter
	case "start":
		return ScrollVerticalAlignmentStart
	case "end":
		return ScrollVerticalAlignmentEnd
	case "nearest":
		return ScrollVerticalAlignmentNearest
	default:
		return ScrollVerticalAlignmentCenter
	}
}

func (a ScrollVerticalAlignment) MarshalJSON() ([]byte, error) {
	return jettison.MarshalOpts(a.String(), jettison.NoHTMLEscaping())
}

func (a ScrollVerticalAlignment) String() string {
	switch a {
	case ScrollVerticalAlignmentCenter:
		return "center"
	case ScrollVerticalAlignmentStart:
		return "start"
	case ScrollVerticalAlignmentEnd:
		return "end"
	case ScrollVerticalAlignmentNearest:
		return "nearest"
	default:
		return "center"
	}
}

// ScrollHorizontalAlignment defines horizontal alignment after scrolling.
// In HTML specification, default value is nearest, but in Ferret it's center.
// More details here https://developer.mozilla.org/en-US/docs/Web/API/Element/scrollIntoView
type ScrollHorizontalAlignment int

const (
	ScrollHorizontalAlignmentCenter  ScrollHorizontalAlignment = 0
	ScrollHorizontalAlignmentStart   ScrollHorizontalAlignment = 1
	ScrollHorizontalAlignmentEnd     ScrollHorizontalAlignment = 2
	ScrollHorizontalAlignmentNearest ScrollHorizontalAlignment = 3
)

func NewScrollHorizontalAlignment(value string) ScrollHorizontalAlignment {
	switch strings.ToLower(value) {
	case "center":
		return ScrollHorizontalAlignmentCenter
	case "start":
		return ScrollHorizontalAlignmentStart
	case "end":
		return ScrollHorizontalAlignmentEnd
	case "nearest":
		return ScrollHorizontalAlignmentNearest
	default:
		return ScrollHorizontalAlignmentCenter
	}
}

func (a ScrollHorizontalAlignment) MarshalJSON() ([]byte, error) {
	return jettison.MarshalOpts(a.String(), jettison.NoHTMLEscaping())
}

func (a ScrollHorizontalAlignment) String() string {
	switch a {
	case ScrollHorizontalAlignmentCenter:
		return "center"
	case ScrollHorizontalAlignmentNearest:
		return "nearest"
	case ScrollHorizontalAlignmentStart:
		return "start"
	case ScrollHorizontalAlignmentEnd:
		return "end"
	default:
		return "center"
	}
}

// ScrollOptions defines how scroll animation should be performed.
type ScrollOptions struct {
	Top      values.Float              `json:"top"`
	Left     values.Float              `json:"left"`
	Behavior ScrollBehavior            `json:"behavior"`
	Block    ScrollVerticalAlignment   `json:"block"`
	Inline   ScrollHorizontalAlignment `json:"inline"`
}
