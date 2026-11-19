package debugger

import (
	"context"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
	"github.com/MontFerret/ferret/v2/pkg/vm"
)

type fakeValueAccess struct {
	inner         vm.DebugValueAccess
	read          func(context.Context, runtime.Value, VariablePage) ([]vm.DebugChild, error)
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

func (f *fakeValueAccess) ReadChildren(ctx context.Context, value runtime.Value, page VariablePage) ([]vm.DebugChild, error) {
	f.inspectCalls++
	if f.read != nil {
		return f.read(ctx, value, page)
	}

	return f.inner.ReadChildren(ctx, value, page)
}
