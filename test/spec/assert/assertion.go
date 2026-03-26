package assert

import (
	"fmt"
	"testing"
)

type Assertion func(t *testing.T, actual any, other ...any)

func NewBinaryAssertion(check Binary) Assertion {
	return func(t *testing.T, actual any, other ...any) {
		t.Helper()

		if len(other) == 0 {
			t.Fatal("expected value is required for binary assertion")
		}

		expected := other[0]
		desc := toDesc(other[1:])

		Expect(t, check, actual, expected, desc...)
	}
}

func NewUnaryAssertion(check Unary) Assertion {
	return func(t *testing.T, actual any, other ...any) {
		t.Helper()

		desc := toDesc(other)

		ExpectThat(t, check, actual, desc...)
	}
}

func Expect(t *testing.T, fn Binary, actual, expected any, desc ...string) {
	t.Helper()

	if err := fn(actual, expected); err != nil {
		fatal(t, err, desc)
	}
}

func ExpectThat(t *testing.T, fn Unary, actual any, desc ...string) {
	t.Helper()

	if err := fn(actual); err != nil {
		fatal(t, err, desc)
	}
}

func toDesc(input []any) []string {
	desc := make([]string, len(input))

	for i := range input {
		desc[i] = fmt.Sprintf("%v", input[i])
	}

	return desc
}

func fatal(t *testing.T, err error, desc []string) {
	t.Helper()

	if len(desc) > 0 && desc[0] != "" {
		t.Fatalf("%s: %v", desc[0], err)
	}

	t.Fatal(err)
}
