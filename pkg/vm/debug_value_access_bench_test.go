package vm

import (
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

func BenchmarkDebugValueAccessDescribeBuiltin(b *testing.B) {
	value := runtime.NewString("value")
	access := NewDebugValueAccess()

	b.ReportAllocs()
	for b.Loop() {
		_ = access.Describe(value)
	}
}

func BenchmarkDebugValueAccessDescribeCustom(b *testing.B) {
	value := debugMetadataValue{info: runtime.DebugInfo{
		TypeName: "SQL::Connection",
		Display:  "open connection",
	}}
	access := NewDebugValueAccess()

	b.ReportAllocs()
	for b.Loop() {
		_ = access.Describe(value)
	}
}
