package mem

import "github.com/MontFerret/ferret/v2/pkg/runtime"

type RegisterFile []runtime.Value

func NewRegisterFile(size int) RegisterFile {
	values := make([]runtime.Value, size)
	fillWithNone(values)
	return values
}

func (rf RegisterFile) Reset() {
	fillWithNoneAndClose(rf)
}
