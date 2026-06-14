package vm

import (
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

func BenchmarkDebugValueAccessTypeNameBuiltin(b *testing.B) {
	value := runtime.NewString("value")
	access := NewDebugValueAccess()

	b.ReportAllocs()
	for b.Loop() {
		_ = access.TypeName(value)
	}
}

func BenchmarkDebugValueAccessDebugInfoCustom(b *testing.B) {
	value := debugMetadataValue{info: runtime.DebugInfo{
		TypeName: "SQL::Connection",
		Display:  "open connection",
	}}
	access := NewDebugValueAccess()

	b.ReportAllocs()
	for b.Loop() {
		_, _ = access.DebugInfo(value)
	}
}
