package vm_test

import (
	"context"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/file"
	"github.com/MontFerret/ferret/v2/pkg/sdk"
	"github.com/MontFerret/ferret/v2/pkg/stdlib"

	"github.com/MontFerret/ferret/v2/pkg/compiler"
	"github.com/MontFerret/ferret/v2/pkg/parser"
	"github.com/MontFerret/ferret/v2/pkg/runtime"
	"github.com/MontFerret/ferret/v2/pkg/vm"

	spec "github.com/MontFerret/ferret/v2/test/spec"
	. "github.com/MontFerret/ferret/v2/test/spec/exec"
)

func TestMember(t *testing.T) {
	RunSpecs(t, []Spec{
		Nil("LET arr = [1,2,3,4] RETURN arr[10]"),
		S("LET arr = [1,2,3,4] RETURN arr[1]", 2),
		S("LET arr = [1,2,3,4] LET idx = 1 RETURN arr[idx]", 2),
		S(`LET obj = { foo: "bar", qaz: "wsx"} RETURN obj["qaz"]`, "wsx"),
		S(fmt.Sprintf(`
								LET obj = { "foo": "bar", %s: "wsx"}

								RETURN obj["qaz"]
							`, "`qaz`"), "wsx"),
		S(fmt.Sprintf(`
								LET obj = { "foo": "bar", %s: "wsx"}

								RETURN obj["let"]
							`, "`let`"),
			"wsx"),
		S(`LET obj = { foo: "bar", qaz: "wsx"} LET key = "qaz" RETURN obj[key]`, "wsx"),
		S(`RETURN { foo: "bar" }.foo`, "bar"),
		S(`LET inexp = 1 IN {'foo': [1]}.foo
			LET ternaryexp = FALSE ? TRUE : {foo: TRUE}.foo
			RETURN inexp && ternaryexp`,
			true),
		S(`RETURN ["bar", "foo"][0]`, "bar"),
		S(`LET inexp = 1 IN [[1]][0]
								LET ternaryexp = FALSE ? TRUE : [TRUE][0]
								RETURN inexp && ternaryexp`,
			true),
		S(`LET obj = {
							first: {
								second: {
									third: {
										fourth: {
											fifth: {
												bottom: true
											}
										}
									}
								}
							}
						}

						RETURN obj.first.second.third.fourth.fifth.bottom`,
			true),
		S(`LET o1 = {
		first: {
		  second: {
		      ["third"]: {
		          fourth: {
		              fifth: {
		                  bottom: true
		              }
		          }
		      }
		  }
		}
		}

		LET o2 = { prop: "third" }

		RETURN o1["first"]["second"][o2.prop]["fourth"]["fifth"].bottom`,

			true),
		S(`LET o1 = {
		first: {
		 second: {
		     third: {
		         fourth: {
		             fifth: {
		                 bottom: true
		             }
		         }
		     }
		 }
		}
		}

		LET o2 = { prop: "third" }

		RETURN o1.first["second"][o2.prop].fourth["fifth"]["bottom"]`,

			true),
		S(`LET obj = {
							attributes: {
								'data-index': 1
							}
						}

						RETURN obj.attributes['data-index']`,
			1),
		S(`LET obj = { '[1]': 42 } RETURN obj[[1]]`, 42),
		S(`LET obj = { '{"a":1}': 7 } RETURN obj[{a:1}]`, 7),
		Error(`LET obj = NONE RETURN obj.foo`),
		Nil(`LET obj = NONE RETURN obj?.foo`),
		Object(`RETURN {first: {second: "third"}}.first`,
			map[string]any{
				"second": "third",
			}),
		Object(`RETURN KEEP_KEYS({first: {second: "third"}}.first, "second")`,
			map[string]any{
				"second": "third",
			}),
		Array(`
					FOR v, k IN {f: {foo: "bar"}}.f
						RETURN [k, v]
				`,
			[]any{
				[]any{"foo", "bar"},
			}),
		S(`RETURN FIRST([[1, 2]][0])`,
			1),
		Array(`RETURN [[1, 2]][0]`,
			[]any{1, 2}),
		Array(`
					FOR i IN [[1, 2]][0]
						RETURN i
				`,
			[]any{1, 2}),
		S(`
					LET arr = [{ name: "Bob" }]

					RETURN FIRST(arr).name
				`,
			"Bob"),
		S(`
					LET arr = [{ name: { first: "Bob" } }]

					RETURN FIRST(arr)['name'].first
				`,
			"Bob"),
		Nil(`
					LET obj = { foo: None }

					RETURN obj.foo?.bar
				`),
	})
}

func TestMemberReservedWords(t *testing.T) {
	p := parser.New("RETURN TRUE")
	r := regexp.MustCompile(`\w+`)
	c := compiler.New(compiler.WithOptimizationLevel(compiler.O0))

	for idx, literal := range p.GetLiteralNames() {
		if !r.MatchString(literal) || literal == "'FUNC'" {
			continue
		}

		idx := idx
		literal := literal

		t.Run(literal, func(t *testing.T) {
			query := strings.Builder{}
			query.WriteString("LET o = {\n")
			query.WriteString(literal[1 : len(literal)-1])
			query.WriteString(":")
			query.WriteString(strconv.Itoa(idx))
			query.WriteString(",\n")
			query.WriteString("}\n")
			query.WriteString("RETURN o")

			expected := strings.Builder{}
			expected.WriteString("{")
			expected.WriteString(strings.ReplaceAll(literal, "'", "\""))
			expected.WriteString(":")
			expected.WriteString(strconv.Itoa(idx))
			expected.WriteString("}")

			prog, err := c.Compile(file.NewAnonymousSource(query.String()))
			if err != nil {
				t.Fatalf("compile failed: %v", err)
			}

			out, err := spec.Exec(prog, true, vm.WithNamespace(stdlib.New()))
			if err != nil {
				t.Fatalf("exec failed: %v", err)
			}

			if out != expected.String() {
				t.Fatalf("unexpected output: got %s, want %s", out, expected.String())
			}
		})
	}
}

func TestOptionalChaining(t *testing.T) {
	RunSpecs(t, []Spec{
		S(
			`
						LET obj = { foo: { bar: "bar" } }

						RETURN obj.foo?.bar
					`,
			"bar",
		),
		S(
			`
						LET obj = { prop1: { prop2: 1 } }

						RETURN obj?.prop1?.prop2
					`,
			1,
		),
		Nil(`LET obj = NONE RETURN obj?.prop1?.prop2`),
		S(
			`
						LET obj = { "prop1": { "prop2": 1 } }

						RETURN obj?.["prop1"]?.["prop2"]
					`,
			1,
		),
		Nil(`LET obj = NONE RETURN obj?.["prop1"]?.["prop2"]`),
		Nil(`
						LET obj = { foo: None }

						RETURN obj.foo?.bar?.[0]
					`),
		S(
			`
					LET obj = { foo: { bar: ["bar"] } }

					RETURN obj.foo?.bar?.[0]
				`,
			"bar"),
		Nil(`RETURN FIRST([])?.foo`),
		S(
			`
							RETURN FIRST([{ foo: "bar" }])?.foo
						`,
			"bar",
		),
		Nil(`LET obj = NONE RETURN obj?.foo?.[ERROR()]`),
		Error(`LET obj = { foo: { bar: 1 } } RETURN obj?.foo?.[ERROR()]`),
		S(`LET obj = { '[1]': 42 } RETURN obj?.[[1]]`, 42),
		S(`LET obj = { '{"a":1}': 7 } RETURN obj?.[{a:1}]`, 7),
		Nil("RETURN ERROR()?.foo"),
		Nil(`LET res = (FOR i IN ERROR() RETURN i)? RETURN res`),

		Array(`LET res = (FOR i IN [1, 2, 3, 4] LET y = ERROR() RETURN y+i)? RETURN res`, []any{}, "Error in array comprehension"),
		Array(`FOR i IN [1, 2, 3, 4] ERROR()? RETURN i`, []any{1, 2, 3, 4}, "Error in FOR loop"),
	}, vm.WithFunction("ERROR", func(ctx context.Context, args ...runtime.Value) (runtime.Value, error) {
		return nil, runtime.ErrNotImplemented
	}))
}

func TestTaggedTypes(t *testing.T) {
	type SomeValue struct {
		PointerProp    *SomeValue `ferret:"pointerProp"`
		NilPointerProp *SomeValue `ferret:"nilPointerProp"`
		StrProp        string     `ferret:"strProp"`
		UntaggedProp   string
		privateStrProp string
		SliceProp      []int `ferret:"sliceProp"`
		IntProp        int   `ferret:"intProp"`
	}

	RunSpecs(t, []Spec{
		S("RETURN GET_VALUE().strProp", "test"),
		S("RETURN GET_VALUE().intProp", 99),
		Array("RETURN GET_VALUE().sliceProp", []any{1, 2, 3}),
		S("RETURN GET_VALUE().pointerProp.strProp", "nested"),
		Nil("RETURN GET_VALUE().nilPointerProp"),
		Nil("RETURN GET_VALUE().untagged"),
		Nil("RETURN GET_VALUE().privateStrProp"),
	}, vm.WithFunction("GET_VALUE", func(ctx context.Context, args ...runtime.Value) (runtime.Value, error) {
		return sdk.Encode(SomeValue{
			StrProp:        "test",
			IntProp:        99,
			SliceProp:      []int{1, 2, 3},
			PointerProp:    &SomeValue{StrProp: "nested"},
			NilPointerProp: nil,
			UntaggedProp:   "untagged",
			privateStrProp: "private",
		}), nil
	}))
}
