package debugger

import (
	"github.com/MontFerret/ferret/v2/pkg/runtime"
	"github.com/MontFerret/ferret/v2/pkg/vm"
)

type fakeValueAccess struct {
	inner     vm.DebugValueAccess
	typeCalls int
}

func (f *fakeValueAccess) TypeName(value runtime.Value) string {
	f.typeCalls++
	return f.inner.TypeName(value)
}

func (f *fakeValueAccess) Lookup(value, key runtime.Value) (runtime.Value, error) {
	return f.inner.Lookup(value, key)
}

func (f *fakeValueAccess) Inspect(value runtime.Value, maxItems int) (vm.DebugValueInspection, bool) {
	return f.inner.Inspect(value, maxItems)
}
