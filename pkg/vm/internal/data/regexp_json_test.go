package data_test

import (
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
	"github.com/MontFerret/ferret/v2/pkg/vm/internal/data"
)

var benchmarkRegexpJSON []byte

func TestRegexpMarshalJSONDoesNotEscapeHTML(t *testing.T) {
	value, err := data.NewRegexp(runtime.NewString(`^<item>&[a-z]+>$`))
	if err != nil {
		t.Fatalf("compile regexp: %v", err)
	}

	encoded, err := value.MarshalJSON()
	if err != nil {
		t.Fatalf("marshal regexp: %v", err)
	}

	if got, want := string(encoded), `"^<item>&[a-z]+>$"`; got != want {
		t.Fatalf("MarshalJSON() = %s, want %s", got, want)
	}
}

func BenchmarkRegexpMarshalJSON(b *testing.B) {
	value, err := data.NewRegexp(runtime.NewString(`^<item>&[a-z]+>$`))
	if err != nil {
		b.Fatal(err)
	}

	encoded, err := value.MarshalJSON()
	if err != nil {
		b.Fatal(err)
	}

	b.ReportAllocs()
	b.SetBytes(int64(len(encoded)))
	b.ResetTimer()

	for b.Loop() {
		encoded, err := value.MarshalJSON()
		if err != nil {
			b.Fatal(err)
		}

		benchmarkRegexpJSON = encoded
	}
}
