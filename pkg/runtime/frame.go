package runtime

type Frame struct {
	Operands  *Stack
	Variables *Stack
	State     *Stack
	ReturnPC  int
}

func NewFrame(size int) *Frame {
	return &Frame{
		Operands:  NewStack(size),
		Variables: NewStack(size),
		State:     NewStack(size),
	}
}
