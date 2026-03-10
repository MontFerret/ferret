package mem

import "github.com/MontFerret/ferret/v2/pkg/runtime"

type RegisterFile struct {
	Values  []runtime.Value
	isDirty bool
}

func NewRegisterFile(size int) *RegisterFile {
	values := make([]runtime.Value, size)
	for i := range values {
		values[i] = runtime.None
	}

	return &RegisterFile{
		Values: values,
	}
}

func (rf *RegisterFile) IsDirty() bool {
	return rf.isDirty
}

func (rf *RegisterFile) MarkDirty() {
	rf.isDirty = true
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
	for i := range rf.Values {
		rf.Values[i] = runtime.None
	}

	rf.isDirty = false
}
