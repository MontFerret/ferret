package sdk_test

import (
	stdjson "encoding/json"
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/sdk"
)

var benchmarkHostValueJSON []byte

func TestHostValueMarshalJSON(t *testing.T) {
	t.Run("ordinary target does not escape HTML", func(t *testing.T) {
		value := sdk.NewHostValue("<item>&value")

		encoded, err := value.MarshalJSON()
		if err != nil {
			t.Fatalf("marshal host value: %v", err)
		}

		if got, want := string(encoded), `"<item>&value"`; got != want {
			t.Fatalf("MarshalJSON() = %s, want %s", got, want)
		}
	})

	t.Run("custom marshaler is delegated", func(t *testing.T) {
		value := sdk.NewHostValue(stdjson.RawMessage(`{"custom":"<item>&value"}`))

		encoded, err := value.MarshalJSON()
		if err != nil {
			t.Fatalf("marshal host value: %v", err)
		}

		if got, want := string(encoded), `{"custom":"<item>&value"}`; got != want {
			t.Fatalf("MarshalJSON() = %s, want %s", got, want)
		}
	})

	t.Run("nil receiver", func(t *testing.T) {
		var value *sdk.HostValue[string]

		encoded, err := value.MarshalJSON()
		if err != nil {
			t.Fatalf("marshal host value: %v", err)
		}

		if got, want := string(encoded), "null"; got != want {
			t.Fatalf("MarshalJSON() = %s, want %s", got, want)
		}
	})

	t.Run("typed nil target", func(t *testing.T) {
		var target *string
		value := sdk.NewHostValue(target)

		encoded, err := value.MarshalJSON()
		if err != nil {
			t.Fatalf("marshal host value: %v", err)
		}

		if got, want := string(encoded), "null"; got != want {
			t.Fatalf("MarshalJSON() = %s, want %s", got, want)
		}
	})
}

func BenchmarkHostValueMarshalJSON(b *testing.B) {
	value := sdk.NewHostValue(struct {
		Name string   `json:"name"`
		Tags []string `json:"tags"`
	}{
		Name: "<item>&value",
		Tags: []string{"one", "two", "three"},
	})

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

		benchmarkHostValueJSON = encoded
	}
}
