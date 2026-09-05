package strings_test

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"github.com/smarty/assertions"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
	"github.com/MontFerret/ferret/v2/pkg/stdlib/strings"
)

func TestRegexMatchShapes(t *testing.T) {
	for _, tc := range []struct{ text, pattern, want string }{
		{text: "abc", pattern: "b", want: `{"match":"b","groups":[],"named":{}}`},
		{text: "é😀 é", pattern: "(?P<letter>é)(😀)?", want: `{"match":"é😀","groups":["é","😀"],"named":{"letter":"é"}}`},
		{text: "é", pattern: "(?P<letter>é)(😀)?", want: `{"match":"é","groups":["é",""],"named":{"letter":"é"}}`},
		{text: "ab", pattern: "(?P<x>a)(?P<x>b)", want: `{"match":"ab","groups":["a","b"],"named":{"x":"a"}}`},
		{text: "b", pattern: "(?P<x>a)?(?P<x>b)", want: `{"match":"b","groups":["","b"],"named":{"x":""}}`},
		{text: "ABC", pattern: "(?i)(abc)", want: `{"match":"ABC","groups":["ABC"],"named":{}}`},
		{text: "", pattern: "", want: "{\"match\":\"\",\"groups\":[],\"named\":{}}"},
	} {
		got, err := strings.RegexFind(context.Background(), runtime.String(tc.text), runtime.String(tc.pattern))
		if err != nil {
			t.Fatal(err)
		}

		if message := assertions.ShouldEqualJSON(got.String(), tc.want); message != "" {
			t.Fatal(message)
		}
	}

	got, err := strings.RegexFind(context.Background(), runtime.String("a"), runtime.String("z"))
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(runtime.None, got) {
		t.Fatalf("got %#v, want %#v", got, runtime.None)
	}

	for _, tc := range []struct{ text, pattern, want string }{
		{text: "aba", pattern: "a", want: `[{"match":"a","groups":[],"named":{}},{"match":"a","groups":[],"named":{}}]`},
		{text: "aba", pattern: "z", want: "[]"}, {text: "aaa", pattern: "aa", want: `[{"match":"aa","groups":[],"named":{}}]`},
		{text: "é", pattern: "", want: `[{"match":"","groups":[],"named":{}},{"match":"","groups":[],"named":{}}]`},
		{text: "a", pattern: "a*", want: `[{"match":"a","groups":[],"named":{}}]`},
	} {
		got, err := strings.RegexFindAll(context.Background(), runtime.String(tc.text), runtime.String(tc.pattern))
		if err != nil {
			t.Fatal(err)
		}

		if message := assertions.ShouldEqualJSON(got.String(), tc.want); message != "" {
			t.Fatal(message)
		}
	}

	got, err = strings.RegexTest(context.Background(), runtime.String("ABC"), runtime.String("(?i)abc"))
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(runtime.True, got) {
		t.Fatalf("got %#v, want %#v", got, runtime.True)
	}

	got, err = strings.RegexTest(context.Background(), runtime.String("ABC"), runtime.String("abc"))
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(runtime.False, got) {
		t.Fatalf("got %#v, want %#v", got, runtime.False)
	}

	got, err = strings.RegexReplace(context.Background(), runtime.String("a1 a2"), runtime.String("(?P<letter>a)([0-9])"), runtime.String("${letter}:$2:$$"))
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(runtime.String("a:1:$ a:2:$"), got) {
		t.Fatalf("got %#v, want %#v", got, runtime.String("a:1:$ a:2:$"))
	}

	got, err = strings.RegexReplace(context.Background(), runtime.String("é"), runtime.String(""), runtime.String("x"))
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(runtime.String("xéx"), got) {
		t.Fatalf("got %#v, want %#v", got, runtime.String("xéx"))
	}
}

func TestRegexCompileErrors(t *testing.T) {
	for _, fn := range []func(context.Context, runtime.Value, runtime.Value) (runtime.Value, error){strings.RegexTest, strings.RegexFind, strings.RegexFindAll} {
		_, err := fn(context.Background(), runtime.String("text"), runtime.String("["))
		if !errors.Is(err, runtime.ErrInvalidArgument) {
			t.Fatalf("error = %v, want %v", err, runtime.ErrInvalidArgument)
		}
	}

	_, err := strings.RegexReplace(context.Background(), runtime.String("text"), runtime.String("["), runtime.String(""))
	if !errors.Is(err, runtime.ErrInvalidArgument) {
		t.Fatalf("error = %v, want %v", err, runtime.ErrInvalidArgument)
	}

	_, err = strings.RegexSplit(context.Background(), runtime.String("text"), runtime.String("["), runtime.Int(0))
	if !errors.Is(err, runtime.ErrInvalidArgument) {
		t.Fatalf("error = %v, want %v", err, runtime.ErrInvalidArgument)
	}
}
