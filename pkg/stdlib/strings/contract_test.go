package strings_test

import (
	"context"
	"errors"
	"math"
	"reflect"
	stdstrings "strings"
	"testing"

	"github.com/smarty/assertions"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
	"github.com/MontFerret/ferret/v2/pkg/stdlib/strings"
)

func TestUnicodeSlicingContract(t *testing.T) {
	for name, fn := range map[string]func(context.Context, runtime.Value, runtime.Value) (runtime.Value, error){"left": strings.Left, "right": strings.Right} {
		t.Run(name, func(t *testing.T) {
			for _, count := range []runtime.Int{0, 1, 4, 5, math.MaxInt64} {
				text := runtime.String("😀abc")
				want := runtime.String("😀")
				if name == "right" {
					text, want = "abc😀", "😀"
				}

				if count == 0 {
					want = ""
				} else if count >= 4 {
					want = text
				}

				got, err := fn(context.Background(), text, count)
				if err != nil {
					t.Fatal(err)
				}

				if !reflect.DeepEqual(want, got) {
					t.Fatalf("got %#v, want %#v", got, want)
				}
			}

			for _, count := range []runtime.Int{-1, math.MinInt64} {
				_, err := fn(context.Background(), runtime.String("😀"), count)
				if !errors.Is(err, runtime.ErrInvalidArgument) {
					t.Fatalf("error = %v, want %v", err, runtime.ErrInvalidArgument)
				}
			}

			got, err := fn(context.Background(), runtime.EmptyString, runtime.Int(math.MaxInt64))
			if err != nil {
				t.Fatal(err)
			}

			if !reflect.DeepEqual(runtime.EmptyString, got) {
				t.Fatalf("got %#v, want %#v", got, runtime.EmptyString)
			}
		})
	}

	for _, test := range []struct {
		want   string
		offset runtime.Int
		length runtime.Int
	}{
		{offset: 1, length: 1, want: "é"}, {offset: 1, length: math.MaxInt64, want: "é文"}, {offset: 0, length: 0, want: ""},
		{offset: -1, length: 1, want: ""}, {offset: math.MinInt64, length: 1, want: ""}, {offset: math.MaxInt64, length: 1, want: ""}, {offset: 3, length: 1, want: ""},
	} {
		got, err := strings.Substring(context.Background(), runtime.String("😀é文"), test.offset, test.length)
		if err != nil {
			t.Fatal(err)
		}

		if !reflect.DeepEqual(runtime.String(test.want), got) {
			t.Fatalf("got %#v, want %#v", got, runtime.String(test.want))
		}
	}

	_, err := strings.Substring(context.Background(), runtime.EmptyString, runtime.Int(0), runtime.Int(-1))
	if !errors.Is(err, runtime.ErrInvalidArgument) {
		t.Fatalf("error = %v, want %v", err, runtime.ErrInvalidArgument)
	}
}

func TestWhitespaceAndLiteralCutsets(t *testing.T) {
	for name, fn := range map[string]func(context.Context, ...runtime.Value) (runtime.Value, error){"trim": strings.Trim, "ltrim": strings.LTrim, "rtrim": strings.RTrim} {
		t.Run(name, func(t *testing.T) {
			input := runtime.String("\t\n\u00a0\u2003😀\u2003\u00a0\n\t")
			want := "😀"
			if name == "ltrim" {
				want = "😀\u2003\u00a0\n\t"
			} else if name == "rtrim" {
				want = "\t\n\u00a0\u2003😀"
			}

			got, err := fn(context.Background(), input)
			if err != nil {
				t.Fatal(err)
			}

			if !reflect.DeepEqual(runtime.String(want), got) {
				t.Fatalf("got %#v, want %#v", got, runtime.String(want))
			}

			got, err = fn(context.Background(), input, runtime.EmptyString)
			if err != nil {
				t.Fatal(err)
			}

			if !reflect.DeepEqual(input, got) {
				t.Fatalf("got %#v, want %#v", got, input)
			}

			got, err = fn(context.Background(), runtime.String("é文😀文é"), runtime.String("文é"))
			if err != nil {
				t.Fatal(err)
			}

			cutWant := "😀"
			if name == "ltrim" {
				cutWant = "😀文é"
			} else if name == "rtrim" {
				cutWant = "é文😀"
			}

			if !reflect.DeepEqual(runtime.String(cutWant), got) {
				t.Fatalf("got %#v, want %#v", got, runtime.String(cutWant))
			}
		})
	}
}

func TestStringPredicates(t *testing.T) {
	for name, fn := range map[string]func(context.Context, runtime.Value, runtime.Value) (runtime.Value, error){"contains": strings.Contains, "starts_with": strings.StartsWith, "ends_with": strings.EndsWith} {
		t.Run(name, func(t *testing.T) {
			for _, tc := range []struct {
				text, search string
				want         bool
			}{
				{text: "😀é😀", search: "😀", want: true}, {text: "😀é😀", search: "文", want: false}, {text: "", search: "", want: true}, {text: "x", search: "", want: true}, {text: "", search: "x", want: false},
			} {
				got, err := fn(context.Background(), runtime.String(tc.text), runtime.String(tc.search))
				if err != nil {
					t.Fatal(err)
				}

				if !reflect.DeepEqual(runtime.Boolean(tc.want), got) {
					t.Fatalf("got %#v, want %#v", got, runtime.Boolean(tc.want))
				}
			}
		})
	}
}

func TestSplitAndReplacementLimits(t *testing.T) {
	for name, fn := range map[string]func(context.Context, ...runtime.Value) (runtime.Value, error){"split": strings.Split, "regex_split": strings.RegexSplit} {
		t.Run(name, func(t *testing.T) {
			for _, tc := range []struct {
				want  string
				limit runtime.Int
			}{{limit: 0, want: "[]"}, {limit: 1, want: `["a,b,c"]`}, {limit: 2, want: `["a","b,c"]`}, {limit: math.MaxInt64, want: `["a","b","c"]`}} {
				got, err := fn(context.Background(), runtime.String("a,b,c"), runtime.String(","), tc.limit)
				if err != nil {
					t.Fatal(err)
				}

				if message := assertions.ShouldEqualJSON(got.String(), tc.want); message != "" {
					t.Fatal(message)
				}
			}

			got, err := fn(context.Background(), runtime.String("😀é"), runtime.EmptyString)
			if err != nil {
				t.Fatal(err)
			}

			if message := assertions.ShouldEqualJSON(got.String(), `["😀","é"]`); message != "" {
				t.Fatal(message)
			}

			for _, limit := range []runtime.Int{-1, math.MinInt64} {
				_, err := fn(context.Background(), runtime.EmptyString, runtime.EmptyString, limit)
				if !errors.Is(err, runtime.ErrInvalidArgument) {
					t.Fatalf("error = %v, want %v", err, runtime.ErrInvalidArgument)
				}
			}
		})
	}

	for _, tc := range []struct {
		search      string
		replacement string
		want        string
		limit       runtime.Int
	}{
		{search: "a", replacement: "$", limit: 0, want: "aba"}, {search: "a", replacement: "$", limit: 1, want: "$ba"}, {search: "a", replacement: "$", limit: math.MaxInt64, want: "$b$"},
		{search: "", replacement: "x", limit: math.MaxInt64, want: "xaxbxax"}, {search: "a", replacement: "", limit: math.MaxInt64, want: "b"},
	} {
		got, err := strings.Replace(context.Background(), runtime.String("aba"), runtime.String(tc.search), runtime.String(tc.replacement), tc.limit)
		if err != nil {
			t.Fatal(err)
		}

		if !reflect.DeepEqual(runtime.String(tc.want), got) {
			t.Fatalf("got %#v, want %#v", got, runtime.String(tc.want))
		}
	}

	_, err := strings.Replace(context.Background(), runtime.String("aba"), runtime.String("a"), runtime.String("x"), runtime.Int(-1))
	if !errors.Is(err, runtime.ErrInvalidArgument) {
		t.Fatalf("error = %v, want %v", err, runtime.ErrInvalidArgument)
	}
}

func TestRepeatContract(t *testing.T) {
	for _, tc := range []struct {
		text  string
		want  string
		count runtime.Int
	}{{text: "😀é", count: 3, want: "😀é😀é😀é"}, {text: "x", count: 0, want: ""}, {text: "", count: 100, want: ""}} {
		got, err := strings.Repeat(context.Background(), runtime.String(tc.text), tc.count)
		if err != nil {
			t.Fatal(err)
		}

		if !reflect.DeepEqual(runtime.String(tc.want), got) {
			t.Fatalf("got %#v, want %#v", got, runtime.String(tc.want))
		}
	}

	for _, count := range []runtime.Int{-1, math.MinInt64, math.MaxInt64} {
		_, err := strings.Repeat(context.Background(), runtime.String("😀"), count)
		if !errors.Is(err, runtime.ErrInvalidArgument) {
			t.Fatalf("error = %v, want %v", err, runtime.ErrInvalidArgument)
		}
	}
}

func TestRepeatOutputLimit(t *testing.T) {
	const maximumBytes = 64 << 20

	for _, text := range []runtime.String{"a", "😀"} {
		t.Run("exact_limit/"+string(text), func(t *testing.T) {
			count := runtime.Int(maximumBytes / len(text))
			got, err := strings.Repeat(context.Background(), text, count)
			if err != nil {
				t.Fatal(err)
			}

			result, ok := got.(runtime.String)
			if !ok {
				t.Fatalf("result type = %T, want runtime.String", got)
			}

			if len(result) != maximumBytes || stdstrings.Count(string(result), string(text)) != int(count) {
				t.Fatalf("result does not contain %d copies in %d bytes", count, maximumBytes)
			}
		})
	}

	for _, tc := range []struct {
		name  string
		text  runtime.String
		count runtime.Int
	}{
		{name: "ASCII_above_limit", text: "a", count: maximumBytes + 1},
		{name: "Unicode_above_limit", text: "😀", count: maximumBytes/4 + 1},
		{name: "ASCII_MaxInt64", text: "a", count: math.MaxInt64},
	} {
		t.Run(tc.name, func(t *testing.T) {
			got, err := strings.Repeat(context.Background(), tc.text, tc.count)
			if !errors.Is(err, runtime.ErrInvalidArgument) {
				t.Fatalf("error = %v, want %v", err, runtime.ErrInvalidArgument)
			}

			if got != runtime.None {
				t.Fatal("invalid repetition must return None")
			}

			position, ok, _ := runtime.InvalidArgumentDetails(err)
			if !ok || position != 1 {
				t.Fatalf("argument position = %d, present = %t, want count argument", position, ok)
			}

			// MaxInt64 exceeds the host count range on 32-bit systems first.
			if tc.count <= runtime.Int(int(^uint(0)>>1)) && !stdstrings.Contains(err.Error(), "67108864 bytes") {
				t.Fatalf("error must identify the output byte limit: %v", err)
			}
		})
	}

	t.Run("oversized_input", func(t *testing.T) {
		text := runtime.String(stdstrings.Repeat("a", maximumBytes+1))
		got, err := strings.Repeat(context.Background(), text, runtime.Int(1))
		if !errors.Is(err, runtime.ErrInvalidArgument) || got != runtime.None {
			t.Fatalf("count one must enforce the output limit: %v", err)
		}

		got, err = strings.Repeat(context.Background(), text, runtime.Int(0))
		if err != nil || got != runtime.EmptyString {
			t.Fatalf("count zero must produce an empty string: %v", err)
		}
	})

	t.Run("empty_text", func(t *testing.T) {
		got, err := strings.Repeat(context.Background(), runtime.EmptyString, runtime.Int(int(^uint(0)>>1)))
		if err != nil || got != runtime.EmptyString {
			t.Fatalf("empty text must produce an empty string for a valid count: %v", err)
		}
	})
}

func TestStrictStringArguments(t *testing.T) {
	unary := map[string]func(context.Context, runtime.Value) (runtime.Value, error){"lower": strings.Lower, "upper": strings.Upper}
	binary := map[string]func(context.Context, runtime.Value, runtime.Value) (runtime.Value, error){
		"contains": strings.Contains, "starts_with": strings.StartsWith, "ends_with": strings.EndsWith,
		"regex_test": strings.RegexTest, "regex_find": strings.RegexFind, "regex_find_all": strings.RegexFindAll,
	}
	for _, bad := range []runtime.Value{runtime.None, runtime.True, runtime.Int(1), runtime.Float(1), runtime.Binary("x"), runtime.NewArray(0), runtime.NewObject()} {
		for name, fn := range unary {
			t.Run(name+"/"+runtime.TypeOf(bad).String(), func(t *testing.T) {
				_, err := fn(context.Background(), bad)
				if !errors.Is(err, runtime.ErrInvalidArgument) {
					t.Fatalf("error = %v, want %v", err, runtime.ErrInvalidArgument)
				}
			})
		}

		for name, fn := range binary {
			t.Run(name+"/"+runtime.TypeOf(bad).String(), func(t *testing.T) {
				_, err := fn(context.Background(), bad, runtime.String("x"))
				if !errors.Is(err, runtime.ErrInvalidArgument) {
					t.Fatalf("error = %v, want %v", err, runtime.ErrInvalidArgument)
				}

				_, err = fn(context.Background(), runtime.String("x"), bad)
				if !errors.Is(err, runtime.ErrInvalidArgument) {
					t.Fatalf("error = %v, want %v", err, runtime.ErrInvalidArgument)
				}
			})
		}
	}

	calls := []struct {
		name string
		fn   func(context.Context, ...runtime.Value) (runtime.Value, error)
		args []runtime.Value
	}{
		{"trim", strings.Trim, []runtime.Value{runtime.String("x"), runtime.String("x")}},
		{"ltrim", strings.LTrim, []runtime.Value{runtime.String("x"), runtime.String("x")}},
		{"rtrim", strings.RTrim, []runtime.Value{runtime.String("x"), runtime.String("x")}},
		{"find_first", strings.FindFirst, []runtime.Value{runtime.String("x"), runtime.String("x"), runtime.Int(0), runtime.Int(1)}},
		{"find_last", strings.FindLast, []runtime.Value{runtime.String("x"), runtime.String("x"), runtime.Int(0), runtime.Int(1)}},
		{"substring", strings.Substring, []runtime.Value{runtime.String("x"), runtime.Int(0), runtime.Int(1)}},
		{"split", strings.Split, []runtime.Value{runtime.String("x"), runtime.String("x"), runtime.Int(1)}},
		{"replace", strings.Replace, []runtime.Value{runtime.String("x"), runtime.String("x"), runtime.String("y"), runtime.Int(1)}},
		{"regex_split", strings.RegexSplit, []runtime.Value{runtime.String("x"), runtime.String("x"), runtime.Int(1)}},
		{"like", strings.Like, []runtime.Value{runtime.String("x"), runtime.String("x"), runtime.True}},
	}
	for _, call := range calls {
		t.Run(call.name, func(t *testing.T) {
			for index := range call.args {
				for _, bad := range []runtime.Value{runtime.None, runtime.Float(1), runtime.NewArray(0)} {
					args := append([]runtime.Value(nil), call.args...)
					args[index] = bad
					_, err := call.fn(context.Background(), args...)
					if !errors.Is(err, runtime.ErrInvalidArgument) {
						t.Fatalf("error = %v, want %v", err, runtime.ErrInvalidArgument)
					}
				}
			}
		})
	}

	for _, fn := range []func(context.Context, runtime.Value, runtime.Value) (runtime.Value, error){strings.Left, strings.Right, strings.Repeat} {
		_, err := fn(context.Background(), runtime.Int(1), runtime.Int(1))
		if !errors.Is(err, runtime.ErrInvalidArgument) {
			t.Fatalf("error = %v, want %v", err, runtime.ErrInvalidArgument)
		}

		_, err = fn(context.Background(), runtime.String("x"), runtime.True)
		if !errors.Is(err, runtime.ErrInvalidArgument) {
			t.Fatalf("error = %v, want %v", err, runtime.ErrInvalidArgument)
		}
	}

	_, err := strings.RegexReplace(context.Background(), runtime.String("x"), runtime.String("x"), runtime.Int(1))
	if !errors.Is(err, runtime.ErrInvalidArgument) {
		t.Fatalf("error = %v, want %v", err, runtime.ErrInvalidArgument)
	}
}

func TestLikeCanonicalPatterns(t *testing.T) {
	for _, tc := range []struct {
		text, pattern string
		want          bool
	}{
		{text: "cart", pattern: "ca_t", want: false}, {text: "ca_t", pattern: "ca_t", want: true}, {text: "a%b", pattern: "a%b", want: true}, {text: "acb", pattern: "a%b", want: false},
		{text: "😀", pattern: "?", want: true}, {text: "foo", pattern: "o", want: false}, {text: "foo", pattern: "*o", want: true}, {text: "foo", pattern: "f[oa]o", want: true}, {text: "foo", pattern: "{foo,bar}", want: true},
		{text: "", pattern: "", want: true},
	} {
		got, err := strings.Like(context.Background(), runtime.String(tc.text), runtime.String(tc.pattern))
		if err != nil {
			t.Fatal(err)
		}

		if !reflect.DeepEqual(runtime.Boolean(tc.want), got) {
			t.Fatalf("got %#v, want %#v", got, runtime.Boolean(tc.want))
		}
	}

	_, err := strings.Like(context.Background(), runtime.EmptyString, runtime.String("["))
	if !errors.Is(err, runtime.ErrInvalidArgument) {
		t.Fatalf("error = %v, want %v", err, runtime.ErrInvalidArgument)
	}
}

func TestJoinContract(t *testing.T) {
	for _, tc := range []struct {
		separator string
		want      string
		values    []runtime.Value
	}{
		{values: nil, separator: ",", want: ""}, {values: []runtime.Value{runtime.String("😀")}, separator: ",", want: "😀"},
		{values: []runtime.Value{runtime.String(""), runtime.String("é"), runtime.String("")}, separator: ",", want: ",é,"},
		{values: []runtime.Value{runtime.String("é"), runtime.String("😀")}, separator: "", want: "é😀"},
	} {
		got, err := strings.Join(context.Background(), runtime.NewArrayWith(tc.values...), runtime.String(tc.separator))
		if err != nil {
			t.Fatal(err)
		}

		if !reflect.DeepEqual(runtime.String(tc.want), got) {
			t.Fatalf("got %#v, want %#v", got, runtime.String(tc.want))
		}
	}

	for _, bad := range []runtime.Value{runtime.None, runtime.Int(1), runtime.NewArrayWith(runtime.String("nested"))} {
		_, err := strings.Join(context.Background(), runtime.NewArrayWith(runtime.String("ok"), bad), runtime.String(","))
		if !errors.Is(err, runtime.ErrInvalidArgument) {
			t.Fatalf("error = %v, want %v", err, runtime.ErrInvalidArgument)
		}
	}

	_, err := strings.Join(context.Background(), runtime.String("x"), runtime.String(","))
	if !errors.Is(err, runtime.ErrInvalidArgument) {
		t.Fatalf("error = %v, want %v", err, runtime.ErrInvalidArgument)
	}

	_, err = strings.Join(context.Background(), runtime.NewArray(0), runtime.Int(1))
	if !errors.Is(err, runtime.ErrInvalidArgument) {
		t.Fatalf("error = %v, want %v", err, runtime.ErrInvalidArgument)
	}

	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	_, err = strings.Join(ctx, runtime.NewArray(0), runtime.String(","))
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("error = %v, want %v", err, context.Canceled)
	}

	failure := errors.New("iteration failed")
	values := &joinTestList{List: runtime.NewArrayWith(runtime.String("first")), failure: failure}
	_, err = strings.Join(context.Background(), values, runtime.String(","))
	if !errors.Is(err, failure) {
		t.Fatalf("error = %v, want %v", err, failure)
	}

	ctx, cancel = context.WithCancel(context.Background())
	values = &joinTestList{List: runtime.NewArrayWith(runtime.String("first")), cancel: cancel}
	_, err = strings.Join(ctx, values, runtime.String(","))
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("error = %v, want %v", err, context.Canceled)
	}
}
