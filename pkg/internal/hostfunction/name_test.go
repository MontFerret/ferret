package hostfunction

import "testing"

func TestCanonicalName(t *testing.T) {
	t.Parallel()

	tests := map[string]string{
		"":                    "",
		"trim":                "trim",
		"TRIM":                "trim",
		"TrIm":                "trim",
		"DB::POSTGRES::QUERY": "db::postgres::query",
		"Db::Postgres::Query": "db::postgres::query",
		"db::POSTGRES::query": "db::postgres::query",
	}

	for input, expected := range tests {
		input, expected := input, expected
		t.Run(input, func(t *testing.T) {
			t.Parallel()

			if actual := CanonicalName(input); actual != expected {
				t.Fatalf("CanonicalName(%q) = %q, want %q", input, actual, expected)
			}
		})
	}
}

func TestHasTerminalName(t *testing.T) {
	t.Parallel()

	tests := map[string]bool{
		"":         false,
		"::":       false,
		"DB::":     false,
		"f":        true,
		"trim":     true,
		"DB::trim": true,
	}

	for input, expected := range tests {
		input, expected := input, expected
		t.Run(input, func(t *testing.T) {
			t.Parallel()

			if actual := HasTerminalName(input); actual != expected {
				t.Fatalf("HasTerminalName(%q) = %t, want %t", input, actual, expected)
			}
		})
	}
}

func BenchmarkCanonicalNameLowercase(b *testing.B) {
	for b.Loop() {
		_ = CanonicalName("db::postgres::query")
	}
}

func BenchmarkCanonicalNameMixedCase(b *testing.B) {
	for b.Loop() {
		_ = CanonicalName("Db::Postgres::QuErY")
	}
}
