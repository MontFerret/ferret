package testing

import (
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

func TestBooleanValueAssertions(t *testing.T) {
	t.Parallel()

	tests := []struct {
		matching   runtime.Value
		mismatch   runtime.Value
		name       string
		positive   string
		negative   string
		descriptor assertion
	}{
		{
			name:       "true",
			descriptor: trueAssertion,
			matching:   runtime.True,
			mismatch:   runtime.False,
			positive:   "assertion error: expected Boolean 'false' to be Boolean 'true'",
			negative:   "assertion error: expected Boolean 'true' not to be Boolean 'true'",
		},
		{
			name:       "false",
			descriptor: falseAssertion,
			matching:   runtime.False,
			mismatch:   runtime.True,
			positive:   "assertion error: expected Boolean 'true' to be Boolean 'false'",
			negative:   "assertion error: expected Boolean 'false' not to be Boolean 'false'",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			requireAssertionSuccess(t, test.descriptor, true, test.matching)
			requireAssertionFailure(t, test.descriptor, false, test.negative, test.matching)
			requireAssertionFailure(t, test.descriptor, true, test.positive, test.mismatch)
			requireAssertionSuccess(t, test.descriptor, false, test.mismatch)
		})
	}

	requireAssertionFailure(t, trueAssertion, true, "assertion error: expected String 'true' to be Boolean 'true'", runtime.NewString("true"))
	requireAssertionFailure(t, falseAssertion, true, "assertion error: expected String 'false' to be Boolean 'false'", runtime.NewString("false"))
}

func TestRuntimeTypeAssertions(t *testing.T) {
	t.Parallel()

	tests := []struct {
		matching   runtime.Value
		mismatch   runtime.Value
		name       string
		positive   string
		negative   string
		descriptor assertion
	}{
		{
			name:       "none",
			descriptor: noneAssertion,
			matching:   runtime.None,
			mismatch:   runtime.NewString("true"),
			positive:   "assertion error: expected String 'true' to be None",
			negative:   "assertion error: expected None 'none' not to be None",
		},
		{
			name:       "string",
			descriptor: stringAssertion,
			matching:   runtime.NewString("hello"),
			mismatch:   runtime.NewInt(1),
			positive:   "assertion error: expected Int '1' to be String",
			negative:   "assertion error: expected String 'hello' not to be String",
		},
		{
			name:       "int",
			descriptor: intAssertion,
			matching:   runtime.NewInt(42),
			mismatch:   runtime.NewFloat(1.5),
			positive:   "assertion error: expected Float '1.5' to be Int",
			negative:   "assertion error: expected Int '42' not to be Int",
		},
		{
			name:       "float",
			descriptor: floatAssertion,
			matching:   runtime.NewFloat(3.14),
			mismatch:   runtime.NewInt(1),
			positive:   "assertion error: expected Int '1' to be Float",
			negative:   "assertion error: expected Float '3.14' not to be Float",
		},
		{
			name:       "datetime",
			descriptor: dateTimeAssertion,
			matching:   runtime.ZeroDateTime,
			mismatch:   runtime.NewString("2023-01-01"),
			positive:   "assertion error: expected String '2023-01-01' to be DateTime",
			negative:   "assertion error: expected DateTime '0001-01-01 00:00:00 +0000 UTC' not to be DateTime",
		},
		{
			name:       "array",
			descriptor: arrayAssertion,
			matching:   runtime.NewArrayWith(runtime.NewInt(1), runtime.NewInt(2)),
			mismatch:   runtime.NewObjectWith(map[string]runtime.Value{}),
			positive:   "assertion error: expected Object '{}' to be Array",
			negative:   "assertion error: expected Array '[1,2]' not to be Array",
		},
		{
			name:       "object",
			descriptor: objectAssertion,
			matching: runtime.NewObjectWith(map[string]runtime.Value{
				"key": runtime.NewString("value"),
			}),
			mismatch: runtime.NewArrayWith(),
			positive: "assertion error: expected Array '[]' to be Object",
			negative: `assertion error: expected Object '{"key":"value"}' not to be Object`,
		},
		{
			name:       "binary",
			descriptor: binaryAssertion,
			matching:   runtime.NewBinary([]byte{1, 2, 3}),
			mismatch:   runtime.NewArrayWith(),
			positive:   "assertion error: expected Array '[]' to be Binary",
			negative:   "assertion error: expected Binary '\x01\x02\x03' not to be Binary",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			requireAssertionSuccess(t, test.descriptor, true, test.matching)
			requireAssertionFailure(t, test.descriptor, false, test.negative, test.matching)
			requireAssertionFailure(t, test.descriptor, true, test.positive, test.mismatch)
			requireAssertionSuccess(t, test.descriptor, false, test.mismatch)
		})
	}
}

func TestNewTypeAssertionRejectsNilType(t *testing.T) {
	t.Parallel()

	defer func() {
		if recover() == nil {
			t.Fatal("newTypeAssertion(nil) did not panic")
		}
	}()

	newTypeAssertion(nil)
}
