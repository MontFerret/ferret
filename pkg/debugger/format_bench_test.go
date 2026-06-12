package debugger

import (
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
	"github.com/MontFerret/ferret/v2/pkg/vm"
)

func BenchmarkFormatValueBuiltin(b *testing.B) {
	value := runtime.NewArrayWith(
		runtime.NewInt(1),
		runtime.NewString("value"),
		runtime.True,
	)
	access := vm.NewDebugValueAccess()
	options := DefaultFormatOptions()

	b.ReportAllocs()
	for b.Loop() {
		_ = formatValue(value, access, options)
	}
}

func BenchmarkFormatValueCustomDisplay(b *testing.B) {
	value := hostileDebugValue{info: runtime.DebugInfo{Display: "open connection"}}
	access := vm.NewDebugValueAccess()
	options := DefaultFormatOptions()

	b.ReportAllocs()
	for b.Loop() {
		_ = formatValue(value, access, options)
	}
}
