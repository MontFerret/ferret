package drivers

type ScrollBehavior int

const (
	ScrollBehaviorAuto ScrollBehavior = 0

	ScrollBehaviorSmooth ScrollBehavior = 1

	ScrollBehaviorInstant ScrollBehavior = 2
)

func NewScrollBehavior(value string) ScrollBehavior {
	switch value {
	case "auto":
		return ScrollBehaviorAuto
	case "smooth":
		return ScrollBehaviorSmooth
	case "instant":
		return ScrollBehaviorInstant
	default:
		return ScrollBehaviorAuto
	}
}

func (b ScrollBehavior) String() string {
	switch b {
	case ScrollBehaviorAuto:
		return "auto"
	case ScrollBehaviorSmooth:
		return "smooth"
	case ScrollBehaviorInstant:
		return "instant"
	default:
		return "auto"
	}
}

type ScrollVerticalAlignment int

const (
	ScrollVerticalAlignmentStart   ScrollVerticalAlignment = 0
	ScrollVerticalAlignmentCenter  ScrollVerticalAlignment = 1
	ScrollVerticalAlignmentEnd     ScrollVerticalAlignment = 2
	ScrollVerticalAlignmentNearest ScrollVerticalAlignment = 3
)

func NewScrollVerticalAlignment(value string) ScrollVerticalAlignment {
	switch value {
	case "start":
		return ScrollVerticalAlignmentStart
	case "center":
		return ScrollVerticalAlignmentCenter
	case "end":
		return ScrollVerticalAlignmentEnd
	case "nearest":
		return ScrollVerticalAlignmentNearest
	default:
		return ScrollVerticalAlignmentStart
	}
}

func (a ScrollVerticalAlignment) String() string {
	switch a {
	case ScrollVerticalAlignmentStart:
		return "start"
	case ScrollVerticalAlignmentCenter:
		return "center"
	case ScrollVerticalAlignmentEnd:
		return "end"
	case ScrollVerticalAlignmentNearest:
		return "nearest"
	default:
		return "start"
	}
}

type ScrollHorizontalAlignment int

const (
	ScrollHorizontalAlignmentNearest ScrollHorizontalAlignment = 0
	ScrollHorizontalAlignmentStart   ScrollHorizontalAlignment = 1
	ScrollHorizontalAlignmentCenter  ScrollHorizontalAlignment = 2
	ScrollHorizontalAlignmentEnd     ScrollHorizontalAlignment = 3
)

func NewScrollHorizontalAlignment(value string) ScrollHorizontalAlignment {
	switch value {
	case "nearest":
		return ScrollHorizontalAlignmentNearest
	case "start":
		return ScrollHorizontalAlignmentStart
	case "center":
		return ScrollHorizontalAlignmentCenter
	case "end":
		return ScrollHorizontalAlignmentEnd
	default:
		return ScrollHorizontalAlignmentNearest
	}
}

func (a ScrollHorizontalAlignment) String() string {
	switch a {
	case ScrollHorizontalAlignmentNearest:
		return "nearest"
	case ScrollHorizontalAlignmentStart:
		return "start"
	case ScrollHorizontalAlignmentCenter:
		return "center"
	case ScrollHorizontalAlignmentEnd:
		return "end"
	default:
		return "nearest"
	}
}

type ScrollOptions struct {
	Behavior ScrollBehavior
	Block    ScrollVerticalAlignment
	Inline   ScrollHorizontalAlignment
}
