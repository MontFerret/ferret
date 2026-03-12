package mem

import "github.com/MontFerret/ferret/v2/pkg/runtime"

type RegisterFile struct {
	Values []runtime.Value
}

func NewRegisterFile(size int) *RegisterFile {
	rf := &RegisterFile{}
	rf.Init(size)
	return rf
}

func (rf *RegisterFile) Init(size int) {
	values := make([]runtime.Value, size)
	fillWithNone(values)

	rf.Values = values
}

func (rf *RegisterFile) Size() int {
	return len(rf.Values)
}

func (rf *RegisterFile) Set(idx int, val runtime.Value) {
	rf.Values[idx] = val
}

func (rf *RegisterFile) Get(idx int) runtime.Value {
	return rf.Values[idx]
}

func (rf *RegisterFile) Reset() {
	fillWithNoneAndClose(rf.Values)
}
