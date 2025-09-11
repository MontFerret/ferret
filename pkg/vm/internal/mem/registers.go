package mem

import "github.com/MontFerret/ferret/pkg/runtime"

type RegisterFile struct {
	Values  []runtime.Value
	isDirty bool
}

func NewRegisterFile(size int) *RegisterFile {
	return &RegisterFile{
		Values: make([]runtime.Value, size),
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
		rf.Values[i] = nil
	}

	rf.isDirty = false
}
