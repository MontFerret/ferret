package debugger

import (
	"strings"
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
	"github.com/MontFerret/ferret/v2/pkg/vm"
)

func BenchmarkFormatValueString(b *testing.B) {
	const budget = 64
	options := DefaultFormatOptions()
	options.MaxBytes = budget
	access := vm.NewDebugValueAccess()

	for _, size := range []struct {
		name string
		n    int
	}{{"short", 96}, {"large", 4 << 20}} {
		for _, content := range []struct {
			name string
			unit string
			last string
		}{{"plain", "a", "z"}, {"escaped", "\n", "\t"}, {"late_invalid", "a", "\xff"}} {
			b.Run(size.name+"/"+content.name, func(b *testing.B) {
				var value runtime.Value = runtime.String(strings.Repeat(content.unit, size.n-1) + content.last)
				b.ReportAllocs()

				for b.Loop() {
					_ = formatValue(value, access.Describe(value), options.MaxBytes)
				}
			})
		}
	}
}

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
		_ = formatValue(value, access.Describe(value), options.MaxBytes)
	}
}

func BenchmarkFormatValueCustomDisplay(b *testing.B) {
	var value runtime.Value = hostileDebugValue{info: runtime.DebugInfo{Display: "open connection"}}
	access := vm.NewDebugValueAccess()
	options := DefaultFormatOptions()

	b.ReportAllocs()
	for b.Loop() {
		_ = formatValue(value, access.Describe(value), options.MaxBytes)
	}
}
