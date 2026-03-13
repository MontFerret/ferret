package mem

import "github.com/MontFerret/ferret/v2/pkg/runtime"

type RegisterFile []runtime.Value

func NewRegisterFile(size int) RegisterFile {
	values := make([]runtime.Value, size)
	fillWithNone(values)
	return values
}

// Reset scrubs register slots to runtime.None. OwnedResources is responsible
// for closing any frame-owned values before the storage is reused.
func (rf RegisterFile) Reset() {
	fillWithNone(rf)
}
