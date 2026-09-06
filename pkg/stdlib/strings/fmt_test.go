package strings_test

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
	"github.com/MontFerret/ferret/v2/pkg/stdlib/strings"
)

func TestFmtGrammar(t *testing.T) {
	for _, tc := range []struct {
		template, want string
		args           []runtime.Value
	}{
		{template: "", want: "", args: nil}, {template: "literal", want: "literal", args: nil}, {template: "{{}}", want: "{}", args: nil},
		{template: "{} {}", want: "é 42", args: []runtime.Value{runtime.String("é"), runtime.Int(42)}},
		{template: "{1}: {0} {1}", want: "42: 😀 42", args: []runtime.Value{runtime.String("😀"), runtime.Int(42)}},
		{template: "{{{0}}}", want: "{x}", args: []runtime.Value{runtime.String("x")}},
		{template: "{00}{0}", want: "xx", args: []runtime.Value{runtime.String("x")}},
		{template: "{}", want: "{}", args: []runtime.Value{runtime.String("{}")}},
		{template: "{}", want: "[1]", args: []runtime.Value{runtime.NewArrayWith(runtime.Int(1))}},
	} {
		t.Run(tc.template, func(t *testing.T) {
			args := append([]runtime.Value{runtime.String(tc.template)}, tc.args...)
			got, err := strings.Fmt(context.Background(), args...)
			if err != nil {
				t.Fatal(err)
			}

			if !reflect.DeepEqual(runtime.String(tc.want), got) {
				t.Fatalf("got %#v, want %#v", got, runtime.String(tc.want))
			}
		})
	}

	for _, args := range [][]runtime.Value{
		nil, {runtime.Int(1)}, {runtime.String("{}")},
		{runtime.String("literal"), runtime.Int(1)},
		{runtime.String("{1}"), runtime.Int(1), runtime.Int(2)},
		{runtime.String("{0} {}"), runtime.Int(1)},
		{runtime.String("{} {0}"), runtime.Int(1)},
		{runtime.String("{0}"), runtime.Int(1), runtime.Int(2)},
	} {
		_, err := strings.Fmt(context.Background(), args...)
		if err == nil {
			t.Fatal("expected an error")
		}
	}

	for _, template := range []string{"{", "}", "{x}", "{-1}", "{ 0}", "{0:2}", "{{}", "{0}{", "{9999999999999999999999999999999999999999}"} {
		_, err := strings.Fmt(context.Background(), runtime.String(template), runtime.Int(1))
		if !errors.Is(err, runtime.ErrInvalidArgument) {
			t.Fatalf("error = %v, want %v", err, runtime.ErrInvalidArgument)
		}
	}
}
