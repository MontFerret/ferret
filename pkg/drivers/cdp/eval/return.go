package eval

type ReturnType int

const (
	ReturnNothing ReturnType = iota
	ReturnValue
	ReturnRef
)

func (rt ReturnType) String() string {
	switch rt {
	case ReturnValue:
		return "value"
	case ReturnRef:
		return "reference"
	default:
		return "nothing"
	}
}
