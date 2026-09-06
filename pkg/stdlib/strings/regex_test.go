package strings_test

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"testing"

	"github.com/smarty/assertions"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
	"github.com/MontFerret/ferret/v2/pkg/stdlib/strings"
)

func TestRegexMatchShapes(t *testing.T) {
	for _, tc := range []struct{ text, pattern, want string }{
		{text: "abc", pattern: "b", want: `{"match":"b","groups":[],"named":{}}`},
		{text: "é😀", pattern: "(?P<letter>é)(😀)?", want: `{"match":"é😀","groups":["é","😀"],"named":{"letter":"é"}}`},
		{text: "é", pattern: "(?P<letter>é)(😀)?", want: `{"match":"é","groups":["é",""],"named":{"letter":"é"}}`},
		{text: "ab", pattern: "(?P<x>a)(?P<y>b)", want: `{"match":"ab","groups":["a","b"],"named":{"x":"a","y":"b"}}`},
		{text: "ab", pattern: "(?P<x>a)(?P<X>b)", want: `{"match":"ab","groups":["a","b"],"named":{"x":"a","X":"b"}}`},
		{text: "ab", pattern: "(a)(b)", want: `{"match":"ab","groups":["a","b"],"named":{}}`},
		{text: "b", pattern: "(?P<x>a)?(b)", want: `{"match":"b","groups":["","b"],"named":{"x":""}}`},
		{text: "ABC", pattern: "(?i)(abc)", want: `{"match":"ABC","groups":["ABC"],"named":{}}`},
		{text: "", pattern: "", want: "{\"match\":\"\",\"groups\":[],\"named\":{}}"},
	} {
		for _, fn := range []struct {
			run        func(context.Context, runtime.Value, runtime.Value) (runtime.Value, error)
			name, want string
		}{
			{name: "find", run: strings.RegexFind, want: tc.want},
			{name: "find_all", run: strings.RegexFindAll, want: "[" + tc.want + "]"},
		} {
			t.Run(fn.name+"/"+tc.pattern, func(t *testing.T) {
				got, err := fn.run(context.Background(), runtime.String(tc.text), runtime.String(tc.pattern))
				if err != nil {
					t.Fatal(err)
				}

				if message := assertions.ShouldEqualJSON(got.String(), fn.want); message != "" {
					t.Fatal(message)
				}
			})
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
		{text: "é😀 é", pattern: "(?P<letter>é)(?P<emoji>😀)?", want: `[{"match":"é😀","groups":["é","😀"],"named":{"letter":"é","emoji":"😀"}},{"match":"é","groups":["é",""],"named":{"letter":"é","emoji":""}}]`},
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

func TestRegexDuplicateNamedCaptures(t *testing.T) {
	for _, fn := range []struct {
		run  func(context.Context, runtime.Value, runtime.Value) (runtime.Value, error)
		name string
	}{
		{name: "find", run: strings.RegexFind},
		{name: "find_all", run: strings.RegexFindAll},
	} {
		for _, tc := range []struct{ text, pattern, duplicate string }{
			{text: "ab", pattern: `(?P<value>a)(?P<value>b)`, duplicate: "value"},
			{text: "b", pattern: `(?P<value>a)?(?P<value>b)`, duplicate: "value"},
			{text: "b", pattern: `(?P<value>a)|(?P<value>b)`, duplicate: "value"},
			{text: "z", pattern: `(?P<value>a)(?P<value>b)`, duplicate: "value"},
			{text: "", pattern: `(?P<value>a)(?P<value>b)`, duplicate: "value"},
			{text: "ab", pattern: `(?<value>a)(?<value>b)`, duplicate: "value"},
			{text: "abcd", pattern: `(?P<z>a)(?P<a>b)(?P<z>c)(?P<a>d)`, duplicate: "z"},
		} {
			t.Run(fn.name+"/"+tc.pattern+"/"+tc.text, func(t *testing.T) {
				got, err := fn.run(context.Background(), runtime.String(tc.text), runtime.String(tc.pattern))
				if !errors.Is(err, runtime.ErrInvalidArgument) {
					t.Fatalf("value = %v, error = %v, want invalid argument", got, err)
				}

				position, ok, cause := runtime.InvalidArgumentDetails(err)
				if !ok || position != 1 || cause == nil {
					t.Fatalf("argument details = (%d, %t, %v), want pattern argument", position, ok, cause)
				}

				if want := fmt.Sprintf("duplicate named capture group %q", tc.duplicate); cause.Error() != want {
					t.Fatalf("cause = %q, want %q", cause, want)
				}

				if got != runtime.None {
					t.Fatalf("value = %v, want None on error", got)
				}
			})
		}
	}
}

func TestRegexDuplicateNamesOutsideStructuredCaptures(t *testing.T) {
	ctx := context.Background()
	pattern := runtime.String(`(?P<value>a)(?P<value>b)`)
	got, err := strings.RegexTest(ctx, runtime.String("ab"), pattern)
	if err != nil || got != runtime.True {
		t.Fatalf("regex_test = %v, %v, want true", got, err)
	}

	got, err = strings.RegexReplace(ctx, runtime.String("ab"), pattern, runtime.String("$1:$2"))
	if err != nil || got != runtime.String("a:b") {
		t.Fatalf("regex_replace = %v, %v, want a:b", got, err)
	}

	got, err = strings.RegexSplit(ctx, runtime.String("xabx"), pattern)
	if err != nil {
		t.Fatal(err)
	}

	if message := assertions.ShouldEqualJSON(got.String(), `["x","x"]`); message != "" {
		t.Fatal(message)
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
