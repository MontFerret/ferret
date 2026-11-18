package debugger

import (
	"github.com/MontFerret/ferret/v2/pkg/runtime"
	"github.com/MontFerret/ferret/v2/pkg/vm"
)

type fakeValueAccess struct {
	inner         vm.DebugValueAccess
	describe      func(runtime.Value) vm.DebugValueDescription
	describeCalls int
	lookupCalls   int
	inspectCalls  int
}

func (f *fakeValueAccess) Describe(value runtime.Value) vm.DebugValueDescription {
	f.describeCalls++
	if f.describe != nil {
		return f.describe(value)
	}

	return f.inner.Describe(value)
}

func (f *fakeValueAccess) Lookup(value, key runtime.Value) (runtime.Value, error) {
	f.lookupCalls++

	return f.inner.Lookup(value, key)
}

func (f *fakeValueAccess) Inspect(value runtime.Value, maxItems int) (vm.DebugValueInspection, bool) {
	f.inspectCalls++

	return f.inner.Inspect(value, maxItems)
}
