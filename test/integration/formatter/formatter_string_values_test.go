package formatter_test

import (
	"fmt"
	"slices"
	"strings"
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/compiler"
	"github.com/MontFerret/ferret/v2/pkg/formatter"
	"github.com/MontFerret/ferret/v2/pkg/runtime"
	"github.com/MontFerret/ferret/v2/test/spec"
)

func TestFormatterStringValues(t *testing.T) {
	tests := []struct {
		name    string
		literal string
		value   string
		double  string
		single  string
	}{
		{"empty single", `''`, "", `""`, `''`},
		{"empty double", `""`, "", `""`, `''`},
		{"plain single", `'abc'`, "abc", `"abc"`, `'abc'`},
		{"plain double", `"abc"`, "abc", `"abc"`, `'abc'`},
		{"unknown escape", `"a\b"`, `a\b`, `"a\b"`, `'a\b'`},
		{"opposite double quote", `'a"b'`, `a"b`, `'a"b'`, `'a"b'`},
		{"opposite single quote", `"a'b"`, `a'b`, `"a'b"`, `"a'b"`},
		{"doubled single quote", `'a''b'`, `a''b`, `"a''b"`, `'a''b'`},
		{"doubled double quote", `"a""b"`, `a""b`, `"a""b"`, `'a""b'`},
		{"escaped double quote", `"a\"b"`, `a\"b`, `"a\"b"`, `'a\"b'`},
		{"escaped single quote", `'a\'b'`, `a\'b`, `"a\'b"`, `'a\'b'`},
		{"both quotes", `'a"b''c'`, `a"b''c`, `'a"b''c'`, `'a"b''c'`},
		{"backslash pair", `"a\\b"`, `a\\b`, `"a\\b"`, `'a\\b'`},
		{"repeated backslashes", `'a\\\\b'`, `a\\\\b`, `"a\\\\b"`, `'a\\\\b'`},
		{"trailing backslash pair", `'a\\'`, `a\\`, `"a\\"`, `'a\\'`},
		{"recognized escapes", `'a\nb\tc'`, "a\nb\tc", `"a\nb\tc"`, `'a\nb\tc'`},
		{"protected escape letters", `'a\\nb\\tc'`, `a\\nb\\tc`, `"a\\nb\\tc"`, `'a\\nb\\tc'`},
		{"odd backslash run", `'a\\\nb'`, "a\\\\\nb", `"a\\\nb"`, `'a\\\nb'`},
		{"unknown escapes", `'\q\r\u0041\U0001F600'`, `\q\r\u0041\U0001F600`, `"\q\r\u0041\U0001F600"`, `'\q\r\u0041\U0001F600'`},
		{"raw newline and tab", "'a\n\tb'", "a\n\tb", "\"a\n\tb\"", "'a\n\tb'"},
		{"raw line terminators", "'a\r\nb\rc\u2028d\u2029e'", "a\r\nb\rc\u2028d\u2029e", "\"a\r\nb\rc\u2028d\u2029e\"", "'a\r\nb\rc\u2028d\u2029e'"},
		{"raw controls", "'a\x00\v\fb'", "a\x00\v\fb", "\"a\x00\v\fb\"", "'a\x00\v\fb'"},
		{"backslash before newline", "'a\\\nb'", "a\\\nb", "\"a\\\nb\"", "'a\\\nb'"},
		{"backslash before tab", "'a\\\tb'", "a\\\tb", "\"a\\\tb\"", "'a\\\tb'"},
		{"unicode adjacency", `'é\中\n😀\t尾'`, "é\\中\n😀\t尾", `"é\中\n😀\t尾"`, `'é\中\n😀\t尾'`},
		{"acute empty", `´´`, "", `""`, `''`},
		{"acute plain", `´abc´`, "abc", `"abc"`, `'abc'`},
		{"acute unicode", `´é\n😀´`, "é\n😀", `"é\n😀"`, `'é\n😀'`},
		{"acute double quote", `´a"b´`, `a"b`, `´a"b´`, `'a"b'`},
		{"acute single quote", `´a'b´`, `a'b`, `"a'b"`, `´a'b´`},
		{"acute escaped delimiter", `´a\´b´`, `a\´b`, `"a\´b"`, `'a\´b'`},
		{"acute trailing backslash", `´a\´`, `a\´`, `´a\´`, `´a\´`},
		{"acute trailing pair", `´a\\´`, `a\\`, `"a\\"`, `'a\\'`},
		{"acute trailing odd run", `´a\\\´`, `a\\\´`, `´a\\\´`, `´a\\\´`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for _, single := range []bool{false, true} {
				t.Run(fmt.Sprintf("single=%t", single), func(t *testing.T) {
					want := tt.double
					if single {
						want = tt.single
					}

					checkStringFormatting(t, "RETURN "+tt.literal, "return "+want, []string{tt.value}, formatter.WithSingleQuote(single))
				})
			}
		})
	}
}

func TestFormatterStringBackslashMatrix(t *testing.T) {
	for _, delimiter := range []string{`"`, `'`, "´"} {
		for count := 0; count <= 6; count++ {
			for _, suffix := range []string{"b", "n", "t", "中", "\n", "\t", ""} {
				// Ordinary quoted tokens cannot end in an unpaired backslash.
				if suffix == "" && count%2 != 0 && delimiter != "´" {
					continue
				}

				literal := delimiter + "é" + strings.Repeat(`\`, count) + suffix + delimiter
				t.Run(fmt.Sprintf("%q", literal), func(t *testing.T) {
					input := "RETURN " + literal
					before := compiledStringValues(t, input)
					if len(before) != 1 {
						t.Fatalf("string constants = %q, want one", before)
					}

					for _, single := range []bool{false, true} {
						checkStringFormatting(t, input, "", before, formatter.WithSingleQuote(single))
					}
				})
			}
		}
	}
}

func TestFormatterStringContexts(t *testing.T) {
	tests := []struct {
		name   string
		input  string
		values []string
	}{
		{"array", `RETURN ['a"b', "a\\b"]`, []string{`a"b`, `a\\b`}},
		{"object key", `RETURN {'a"b': 1, "a'b": 2}`, []string{`a"b`, `a'b`}},
		{"computed key", `RETURN {['a''b']: "a\b"}`, []string{`a''b`, `a\b`}},
		{"interpolation", "RETURN `pre${'a\"b'}post`", []string{`prea"bpost`}},
		{"interpolated backslashes", "RETURN `pre${'a\\\\b'}post`", []string{`prea\\bpost`}},
		{"template chunks", "RETURN `a\\\\b\\n\\t\\${x}\\`é`", []string{"a\\b\n\t${x}`é"}},
		{"multiline array", "RETURN ['a\n\nb', 'c\\\nd']", []string{"a\n\nb", "c\\\nd"}},
		{"multiline object key", "RETURN {'a\nb': 'c\n\td'}", []string{"a\nb", "c\n\td"}},
		{"nested multiline", "RETURN [{v: ['a\n', '\nb', 'c\n    d']}]", []string{"v", "a\n", "\nb", "c\n    d"}},
		{"multiline interpolation", "RETURN [`pre${'a\\\nb'}post`]", []string{"prea\\\nbpost"}},
		{"multiline template", "RETURN [`pre\n${'a\"b'}\npost`]", []string{"pre\na\"b\npost"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for _, single := range []bool{false, true} {
				for _, width := range []uint64{12, 200} {
					t.Run(fmt.Sprintf("single=%t/width=%d", single, width), func(t *testing.T) {
						checkStringFormatting(t, tt.input, "", tt.values, formatter.WithSingleQuote(single), formatter.WithPrintWidth(width))
					})
				}
			}
		})
	}
}

func TestFormatterMultilineStringLayout(t *testing.T) {
	input := "RETURN ['a\nb']"
	want := "return [\n    \"a\nb\"\n]"
	checkStringFormatting(t, input, want, []string{"a\nb"}, formatter.WithPrintWidth(200))
}

func TestFormatterStringWithUnresolvedHostFunction(t *testing.T) {
	const input = `RETURN UNREGISTERED('a"b', 'a\b')`
	const want = `return UNREGISTERED('a"b', "a\b")`

	formatted := formatSource(t, input)
	if formatted != want {
		t.Fatalf("formatted = %q, want %q", formatted, want)
	}

	if second := formatSource(t, formatted); second != formatted {
		t.Fatalf("second format = %q, first = %q", second, formatted)
	}
}

func checkStringFormatting(t *testing.T, input, want string, values []string, options ...formatter.Option) {
	t.Helper()

	before := compiledStringValues(t, input)
	expected := slices.Clone(values)
	slices.Sort(expected)
	if !slices.Equal(before, expected) {
		t.Fatalf("original native strings = %q, want %q", before, expected)
	}

	formatted := formatSource(t, input, options...)
	if want != "" && formatted != want {
		t.Errorf("formatted = %q, want %q", formatted, want)
	}

	after := compiledStringValues(t, formatted)
	if !slices.Equal(before, after) {
		t.Errorf("native strings changed: before %q, after %q; formatted %q", before, after, formatted)
	}

	if second := formatSource(t, formatted, options...); second != formatted {
		t.Errorf("second format = %q, first = %q", second, formatted)
	}
}

func compiledStringValues(t *testing.T, input string) []string {
	t.Helper()

	program, err := spec.Compile(input, compiler.None)
	if err != nil {
		t.Fatalf("compile %q: %v", input, err)
	}

	var values []string
	for _, constant := range program.Constants {
		if value, ok := constant.(runtime.String); ok {
			values = append(values, string(value))
		}
	}

	// Context fixtures can change constant-pool order without changing values.
	slices.Sort(values)

	return values
}
