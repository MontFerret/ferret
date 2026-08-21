package testing

import (
	"context"
	"errors"
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

func TestEmptyAssertion(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		empty    runtime.Value
		nonEmpty runtime.Value
		emptyErr string
		fullErr  string
	}{
		{
			name:     "string",
			empty:    runtime.NewString(""),
			nonEmpty: runtime.NewString("Foo"),
			emptyErr: "assertion error: expected String '' not to be empty",
			fullErr:  "assertion error: expected String 'Foo' to be empty",
		},
		{
			name:     "array",
			empty:    runtime.NewArrayWith(),
			nonEmpty: runtime.NewArrayWith(runtime.NewInt(1), runtime.NewInt(2), runtime.NewInt(3)),
			emptyErr: "assertion error: expected Array '[]' not to be empty",
			fullErr:  "assertion error: expected Array '[1,2,3]' to be empty",
		},
		{
			name:  "object",
			empty: runtime.NewObjectWith(map[string]runtime.Value{}),
			nonEmpty: runtime.NewObjectWith(map[string]runtime.Value{
				"a": runtime.NewInt(1),
				"b": runtime.NewInt(2),
				"c": runtime.NewInt(3),
			}),
			emptyErr: "assertion error: expected Object '{}' not to be empty",
			fullErr:  `assertion error: expected Object '{"a":1,"b":2,"c":3}' to be empty`,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			requireAssertionSuccess(t, emptyAssertion, true, test.empty)
			requireAssertionFailure(t, emptyAssertion, false, test.emptyErr, test.empty)
			requireAssertionFailure(t, emptyAssertion, true, test.fullErr, test.nonEmpty)
			requireAssertionSuccess(t, emptyAssertion, false, test.nonEmpty)
		})
	}
}

func TestContainsAssertion(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		haystack    runtime.Value
		contained   runtime.Value
		missing     runtime.Value
		positiveErr string
		negativeErr string
	}{
		{
			name:        "string",
			haystack:    runtime.NewString("FooBar"),
			contained:   runtime.NewString("Bar"),
			missing:     runtime.NewString("Baz"),
			positiveErr: "assertion error: expected String 'FooBar' to contain String 'Baz'",
			negativeErr: "assertion error: expected String 'FooBar' not to contain String 'Bar'",
		},
		{
			name:        "array",
			haystack:    runtime.NewArrayWith(runtime.NewInt(1), runtime.NewInt(2), runtime.NewInt(3)),
			contained:   runtime.NewInt(2),
			missing:     runtime.NewInt(4),
			positiveErr: "assertion error: expected Array '[1,2,3]' to contain Int '4'",
			negativeErr: "assertion error: expected Array '[1,2,3]' not to contain Int '2'",
		},
		{
			name: "object",
			haystack: runtime.NewObjectWith(map[string]runtime.Value{
				"a": runtime.NewInt(1),
				"b": runtime.NewInt(2),
				"c": runtime.NewInt(3),
			}),
			contained:   runtime.NewInt(2),
			missing:     runtime.NewInt(4),
			positiveErr: `assertion error: expected Object '{"a":1,"b":2,"c":3}' to contain Int '4'`,
			negativeErr: `assertion error: expected Object '{"a":1,"b":2,"c":3}' not to contain Int '2'`,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			requireAssertionSuccess(t, containsAssertion, true, test.haystack, test.contained)
			requireAssertionFailure(t, containsAssertion, false, test.negativeErr, test.haystack, test.contained)
			requireAssertionFailure(t, containsAssertion, true, test.positiveErr, test.haystack, test.missing)
			requireAssertionSuccess(t, containsAssertion, false, test.haystack, test.missing)
		})
	}

	requireAssertionSuccess(t, containsAssertion, true, runtime.NewRange(1, 3), runtime.NewInt(2))
	requireAssertionSuccess(t, containsAssertion, false, runtime.NewRange(1, 3), runtime.NewInt(4))
}

func TestHasAssertion(t *testing.T) {
	t.Parallel()

	object := runtime.NewObjectWith(map[string]runtime.Value{
		"id":   runtime.NewInt(1),
		"name": runtime.NewString("Ferret"),
		"none": runtime.None,
		"url":  runtime.NewString("https://ferretlang.org"),
	})

	requireAssertionSuccess(t, hasAssertion, true, object, runtime.NewString("id"))
	requireAssertionSuccess(t, hasAssertion, true, object, runtime.NewString("none"))
	requireAssertionSuccess(
		t,
		hasAssertion,
		true,
		object,
		runtime.NewArrayWith(runtime.NewString("id"), runtime.NewString("name"), runtime.NewString("url")),
	)
	requireAssertionSuccess(t, hasAssertion, true, object, runtime.NewArrayWith())
	requireAssertionSuccess(t, hasAssertion, false, object, runtime.NewString("missing"))
	requireAssertionSuccess(
		t,
		hasAssertion,
		false,
		object,
		runtime.NewArrayWith(runtime.NewString("id"), runtime.NewString("missing")),
	)
	requireAssertionFailure(
		t,
		hasAssertion,
		true,
		`assertion error: expected Object '{"id":1,"name":"Ferret","none":null,"url":"https://ferretlang.org"}' to have property String 'missing'`,
		object,
		runtime.NewString("missing"),
	)
	requireAssertionFailure(
		t,
		hasAssertion,
		false,
		`assertion error: expected Object '{"id":1,"name":"Ferret","none":null,"url":"https://ferretlang.org"}' not to have all properties Array '["id","name"]'`,
		object,
		runtime.NewArrayWith(runtime.NewString("id"), runtime.NewString("name")),
	)
	requireAssertionFailure(
		t,
		hasAssertion,
		false,
		`assertion error: expected Object '{"id":1,"name":"Ferret","none":null,"url":"https://ferretlang.org"}' not to have all properties Array '[]'`,
		object,
		runtime.NewArrayWith(),
	)
	requireAssertionFailure(
		t,
		hasAssertion,
		true,
		"assertion error: required property is missing",
		object,
		runtime.NewString("missing"),
		runtime.NewString("required property is missing"),
	)
}

func TestHasAssertionRejectsInvalidUsage(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		args []runtime.Value
	}{
		{name: "invalid target", args: []runtime.Value{runtime.NewInt(1), runtime.NewString("id")}},
		{name: "invalid key", args: []runtime.Value{runtime.NewObject(), runtime.NewInt(1)}},
		{name: "invalid key list element", args: []runtime.Value{runtime.NewObject(), runtime.NewArrayWith(runtime.NewString("id"), runtime.NewInt(1))}},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			out, err := hasAssertion.positive()(context.Background(), test.args...)
			if out != runtime.None {
				t.Fatalf("output = %v, want None", out)
			}
			if err == nil || errors.Is(err, errAssertion) {
				t.Fatalf("error = %v, want propagated invalid usage", err)
			}
		})
	}
}

func TestHasAssertionStopsAfterFirstMissingKey(t *testing.T) {
	t.Parallel()

	target := &countingMap{
		Object: runtime.NewObjectWith(map[string]runtime.Value{
			"present": runtime.NewInt(1),
		}),
	}
	keys := runtime.NewArrayWith(runtime.NewString("missing"), runtime.NewString("present"))

	matched, err := hasAssertion.fn(context.Background(), []runtime.Value{target, keys})
	if err != nil {
		t.Fatal(err)
	}

	if matched {
		t.Fatal("has assertion matched after a required key was missing")
	}

	if len(target.lookups) != 1 || target.lookups[0] != "missing" {
		t.Fatalf("key lookups = %v, want [missing]", target.lookups)
	}
}

func TestLenAssertion(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		value       runtime.Value
		positiveErr string
		negativeErr string
	}{
		{
			name:        "string",
			value:       runtime.NewString("Foo"),
			positiveErr: "assertion error: expected String 'Foo' to has size 1",
			negativeErr: "assertion error: expected String 'Foo' not to has size 3",
		},
		{
			name:        "array",
			value:       runtime.NewArrayWith(runtime.NewInt(1), runtime.NewInt(2), runtime.NewInt(3)),
			positiveErr: "assertion error: expected Array '[1,2,3]' to has size 1",
			negativeErr: "assertion error: expected Array '[1,2,3]' not to has size 3",
		},
		{
			name: "object",
			value: runtime.NewObjectWith(map[string]runtime.Value{
				"a": runtime.NewInt(1),
				"b": runtime.NewInt(2),
				"c": runtime.NewInt(3),
			}),
			positiveErr: `assertion error: expected Object '{"a":1,"b":2,"c":3}' to has size 1`,
			negativeErr: `assertion error: expected Object '{"a":1,"b":2,"c":3}' not to has size 3`,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			requireAssertionFailure(t, lenAssertion, true, test.positiveErr, test.value, runtime.NewInt(1))
			requireAssertionSuccess(t, lenAssertion, false, test.value, runtime.NewInt(1))
			requireAssertionSuccess(t, lenAssertion, true, test.value, runtime.NewInt(3))
			requireAssertionFailure(t, lenAssertion, false, test.negativeErr, test.value, runtime.NewInt(3))
		})
	}
}

func TestMatchAssertion(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		value       runtime.Value
		match       runtime.String
		miss        runtime.String
		positiveErr string
		negativeErr string
	}{
		{
			name:        "string",
			value:       runtime.NewString("hello world"),
			match:       runtime.NewString(`hello\s+world`),
			miss:        runtime.NewString(`^goodbye`),
			positiveErr: "assertion error: expected String 'hello world' to match regular expression",
			negativeErr: "assertion error: expected String 'hello world' not to match regular expression",
		},
		{
			name:        "complex pattern",
			value:       runtime.NewString("abc123def"),
			match:       runtime.NewString(`\d+`),
			miss:        runtime.NewString(`^xyz`),
			positiveErr: "assertion error: expected String 'abc123def' to match regular expression",
			negativeErr: "assertion error: expected String 'abc123def' not to match regular expression",
		},
		{
			name:        "non-string value",
			value:       runtime.NewInt(123),
			match:       runtime.NewString(`\d+`),
			miss:        runtime.NewString(`^[a-z]+$`),
			positiveErr: "assertion error: expected Int '123' to match regular expression",
			negativeErr: "assertion error: expected Int '123' not to match regular expression",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			requireAssertionSuccess(t, matchAssertion, true, test.value, test.match)
			requireAssertionFailure(t, matchAssertion, false, test.negativeErr, test.value, test.match)
			requireAssertionFailure(t, matchAssertion, true, test.positiveErr, test.value, test.miss)
			requireAssertionSuccess(t, matchAssertion, false, test.value, test.miss)
		})
	}
}

func TestCollectionErrorsPropagate(t *testing.T) {
	t.Parallel()

	tests := []struct {
		descriptor assertion
		name       string
		args       []runtime.Value
	}{
		{name: "empty non-measurable", descriptor: emptyAssertion, args: []runtime.Value{runtime.NewInt(1)}},
		{name: "len non-measurable", descriptor: lenAssertion, args: []runtime.Value{runtime.NewInt(1), runtime.NewInt(1)}},
		{name: "contains unsupported", descriptor: containsAssertion, args: []runtime.Value{runtime.NewInt(1), runtime.NewInt(1)}},
		{name: "invalid regular expression", descriptor: matchAssertion, args: []runtime.Value{runtime.NewString("value"), runtime.NewString("[")}},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			out, err := test.descriptor.positive()(context.Background(), test.args...)
			if out != runtime.None {
				t.Fatalf("output = %v, want None", out)
			}
			if err == nil {
				t.Fatal("expected underlying operation error")
			}
			if errors.Is(err, errAssertion) {
				t.Fatalf("operation error %v was replaced with an assertion error", err)
			}
		})
	}
}
