package artifact

import (
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/bytecode"
	"github.com/MontFerret/ferret/v2/pkg/compiler"
	"github.com/MontFerret/ferret/v2/pkg/file"
)

func BenchmarkMarshalMessagePack(b *testing.B) {
	program := mustBenchmarkArtifactProgram(b)

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		if _, err := Marshal(program, Options{Format: FormatMsgPack}); err != nil {
			b.Fatalf("Marshal() error = %v", err)
		}
	}
}

func BenchmarkUnmarshalMessagePack(b *testing.B) {
	program := mustBenchmarkArtifactProgram(b)
	data, err := Marshal(program, Options{Format: FormatMsgPack})
	if err != nil {
		b.Fatalf("Marshal() error = %v", err)
	}

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		if _, err := Unmarshal(data); err != nil {
			b.Fatalf("Unmarshal() error = %v", err)
		}
	}
}

func BenchmarkMarshalJSON(b *testing.B) {
	program := mustBenchmarkArtifactProgram(b)

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		if _, err := Marshal(program, Options{Format: FormatJSON}); err != nil {
			b.Fatalf("Marshal() error = %v", err)
		}
	}
}

func BenchmarkUnmarshalJSON(b *testing.B) {
	program := mustBenchmarkArtifactProgram(b)
	data, err := Marshal(program, Options{Format: FormatJSON})
	if err != nil {
		b.Fatalf("Marshal() error = %v", err)
	}

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		if _, err := Unmarshal(data); err != nil {
			b.Fatalf("Unmarshal() error = %v", err)
		}
	}
}

func mustBenchmarkArtifactProgram(b *testing.B) *bytecode.Program {
	b.Helper()

	src := file.NewSource("bench.fql", `
LET users = [
  { gender: "m", age: 31 },
  { gender: "f", age: 25 },
  { gender: "m", age: 45 }
]

FOR u IN users
  COLLECT gender = u.gender
  AGGREGATE minAge = MIN(u.age)
  RETURN { gender, minAge }
`)

	program, err := compiler.New().Compile(src)
	if err != nil {
		b.Fatalf("compile failed: %v", err)
	}

	return program
}
