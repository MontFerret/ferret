package vm

import (
	"strconv"
	"testing"

	"github.com/MontFerret/api/debugger"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

func BenchmarkDebugBufferedPage(b *testing.B) {
	for _, kind := range []string{"array", "dataset"} {
		for _, size := range []int{100, 10_000, 1_000_000} {
			b.Run(kind+"/"+strconv.Itoa(size), func(b *testing.B) {
				value := bufferedDebugFixture(b, kind, size)
				access := NewDebugValueAccess()
				for _, offset := range []struct {
					name  string
					start int64
				}{{"fixed", 80}, {"tail", int64(size - 20)}} {
					b.Run(offset.name, func(b *testing.B) {
						page := debugger.VariablePage{Filter: debugger.VariableFilterIndexed, Start: offset.start, Count: 20}
						b.ReportAllocs()

						for b.Loop() {
							got, err := access.ReadChildren(b.Context(), value, page)
							if err != nil || len(got) != 20 {
								b.Fatalf("page: %v, %v", got, err)
							}
						}
					})
				}
			})
		}
	}
}

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
