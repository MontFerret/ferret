package runtime_test

import (
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

func TestNormalizeRegisteredName(t *testing.T) {
	t.Parallel()

	tests := map[string]string{
		"":                    "",
		"trim":                "trim",
		"TRIM":                "trim",
		"TrIm":                "trim",
		"DB::POSTGRES::QUERY": "db::postgres::query",
		"Db::Postgres::Query": "db::postgres::query",
		"db::POSTGRES::query": "db::postgres::query",
		"DB::Ä::QUERY":        "db::Ä::query",
	}

	for input, expected := range tests {
		input, expected := input, expected
		t.Run(input, func(t *testing.T) {
			t.Parallel()

			if actual := runtime.NormalizeRegisteredName(input); actual != expected {
				t.Fatalf("NormalizeRegisteredName(%q) = %q, want %q", input, actual, expected)
			}
		})
	}
}

func TestHasTerminalFunctionName(t *testing.T) {
	t.Parallel()

	tests := map[string]bool{
		"":                   false,
		"::":                 false,
		"DB::":               false,
		"DB::POSTGRES::":     false,
		"f":                  true,
		"trim":               true,
		"DB::trim":           true,
		"DB::POSTGRES::trim": true,
		"::trim":             true,
	}

	for input, expected := range tests {
		input, expected := input, expected
		t.Run(input, func(t *testing.T) {
			t.Parallel()

			if actual := runtime.HasTerminalFunctionName(input); actual != expected {
				t.Fatalf("HasTerminalFunctionName(%q) = %t, want %t", input, actual, expected)
			}
		})
	}
}

func BenchmarkNormalizeRegisteredNameLowercase(b *testing.B) {
	for b.Loop() {
		_ = runtime.NormalizeRegisteredName("db::postgres::query")
	}
}

func BenchmarkNormalizeRegisteredNameMixedCase(b *testing.B) {
	for b.Loop() {
		_ = runtime.NormalizeRegisteredName("Db::Postgres::QuErY")
	}
}
