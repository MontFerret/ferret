package testing

import "testing"

func TestAppendObjectPath(t *testing.T) {
	t.Parallel()

	tests := []struct {
		key      string
		expected string
	}{
		{key: "name", expected: "$.name"},
		{key: "user_name2", expected: "$.user_name2"},
		{key: "2name", expected: `$["2name"]`},
		{key: "content-type", expected: `$["content-type"]`},
		{key: "hello world", expected: `$["hello world"]`},
		{key: `quote"and\slash`, expected: `$["quote\"and\\slash"]`},
		{key: "", expected: `$[""]`},
		{key: "naïve", expected: `$["naïve"]`},
	}

	for _, test := range tests {
		t.Run(test.key, func(t *testing.T) {
			t.Parallel()

			if actual := appendObjectPath("$", test.key); actual != test.expected {
				t.Fatalf("appendObjectPath() = %q, want %q", actual, test.expected)
			}
		})
	}
}
