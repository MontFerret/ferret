package vm_test

import (
	"context"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/file"
	"github.com/MontFerret/ferret/v2/pkg/stdlib"

	"github.com/MontFerret/ferret/v2/test/integration/base"

	"github.com/MontFerret/ferret/v2/pkg/compiler"
	"github.com/MontFerret/ferret/v2/pkg/parser"
	"github.com/MontFerret/ferret/v2/pkg/runtime"
	"github.com/MontFerret/ferret/v2/pkg/vm"

	. "github.com/smartystreets/goconvey/convey"
)

func TestMember(t *testing.T) {
	RunUseCases(t, []UseCase{
		CaseNil("LET arr = [1,2,3,4] RETURN arr[10]"),
		Case("LET arr = [1,2,3,4] RETURN arr[1]", 2),
		Case("LET arr = [1,2,3,4] LET idx = 1 RETURN arr[idx]", 2),
		Case(`LET obj = { foo: "bar", qaz: "wsx"} RETURN obj["qaz"]`, "wsx"),
		Case(fmt.Sprintf(`
								LET obj = { "foo": "bar", %s: "wsx"}
		
								RETURN obj["qaz"]
							`, "`qaz`"), "wsx"),
		Case(fmt.Sprintf(`
								LET obj = { "foo": "bar", %s: "wsx"}
		
								RETURN obj["let"]
							`, "`let`"),
			"wsx"),
		Case(`LET obj = { foo: "bar", qaz: "wsx"} LET key = "qaz" RETURN obj[key]`, "wsx"),
		Case(`RETURN { foo: "bar" }.foo`, "bar"),
		Case(`LET inexp = 1 IN {'foo': [1]}.foo
			LET ternaryexp = FALSE ? TRUE : {foo: TRUE}.foo
			RETURN inexp && ternaryexp`,
			true),
		Case(`RETURN ["bar", "foo"][0]`, "bar"),
		Case(`LET inexp = 1 IN [[1]][0]
								LET ternaryexp = FALSE ? TRUE : [TRUE][0]
								RETURN inexp && ternaryexp`,
			true),
		Case(`LET obj = {
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
		Case(`LET o1 = {
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
		Case(`LET o1 = {
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
		Case(`LET obj = {
							attributes: {
								'data-index': 1
							}
						}
		
						RETURN obj.attributes['data-index']`,
			1),
		Case(`LET obj = { '[1]': 42 } RETURN obj[[1]]`, 42),
		Case(`LET obj = { '{"a":1}': 7 } RETURN obj[{a:1}]`, 7),
		CaseRuntimeError(`LET obj = NONE RETURN obj.foo`),
		CaseNil(`LET obj = NONE RETURN obj?.foo`),
		CaseObject(`RETURN {first: {second: "third"}}.first`,
			map[string]any{
				"second": "third",
			}),
		SkipCaseObject(`RETURN KEEP_KEYS({first: {second: "third"}}.first, "second")`,
			map[string]any{
				"second": "third",
			}),
		CaseArray(`
					FOR v, k IN {f: {foo: "bar"}}.f
						RETURN [k, v]
				`,
			[]any{
				[]any{"foo", "bar"},
			}),
		Case(`RETURN FIRST([[1, 2]][0])`,
			1),
		CaseArray(`RETURN [[1, 2]][0]`,
			[]any{1, 2}),
		CaseArray(`
					FOR i IN [[1, 2]][0]
						RETURN i
				`,
			[]any{1, 2}),
		Case(`
					LET arr = [{ name: "Bob" }]
		
					RETURN FIRST(arr).name
				`,
			"Bob"),
		Case(`
					LET arr = [{ name: { first: "Bob" } }]
	
					RETURN FIRST(arr)['name'].first
				`,
			"Bob"),
		CaseNil(`
					LET obj = { foo: None }

					RETURN obj.foo?.bar
				`),
	})
}

func TestMemberReservedWords(t *testing.T) {
	Convey("Reserved words as property name", t, func() {
		p := parser.New("RETURN TRUE")

		r := regexp.MustCompile(`\w+`)

		for idx, l := range p.GetLiteralNames() {
			if r.MatchString(l) {
				query := strings.Builder{}
				query.WriteString("LET o = {\n")
				query.WriteString(l[1 : len(l)-1])
				query.WriteString(":")
				query.WriteString(strconv.Itoa(idx))
				query.WriteString(",\n")
				query.WriteString("}\n")
				query.WriteString("RETURN o")

				expected := strings.Builder{}
				expected.WriteString("{")
				expected.WriteString(strings.ReplaceAll(l, "'", "\""))
				expected.WriteString(":")
				expected.WriteString(strconv.Itoa(idx))
				expected.WriteString("}")

				c := compiler.New(compiler.WithOptimizationLevel(compiler.O0))
				prog, err := c.Compile(file.NewAnonymousSource(query.String()))

				So(err, ShouldBeNil)

				out, err := base.Exec(prog, true, vm.WithFunctions(stdlib.New()))

				So(err, ShouldBeNil)
				So(out, ShouldEqual, expected.String())
			}
		}
	})
}

func TestOptionalChaining(t *testing.T) {
	RunUseCases(t, []UseCase{
		Case(
			`
						LET obj = { foo: { bar: "bar" } }
			
						RETURN obj.foo?.bar
					`,
			"bar",
		),
		Case(
			`
						LET obj = { prop1: { prop2: 1 } }

						RETURN obj?.prop1?.prop2
					`,
			1,
		),
		CaseNil(`LET obj = NONE RETURN obj?.prop1?.prop2`),
		Case(
			`
						LET obj = { "prop1": { "prop2": 1 } }

						RETURN obj?.["prop1"]?.["prop2"]
					`,
			1,
		),
		CaseNil(`LET obj = NONE RETURN obj?.["prop1"]?.["prop2"]`),
		CaseNil(`
						LET obj = { foo: None }
			
						RETURN obj.foo?.bar?.[0]
					`),
		Case(
			`
					LET obj = { foo: { bar: ["bar"] } }
		
					RETURN obj.foo?.bar?.[0]
				`,
			"bar"),
		CaseNil(`RETURN FIRST([])?.foo`),
		Case(
			`
							RETURN FIRST([{ foo: "bar" }])?.foo
						`,
			"bar",
		),
		CaseNil(`LET obj = NONE RETURN obj?.foo?.[ERROR()]`),
		CaseRuntimeError(`LET obj = { foo: { bar: 1 } } RETURN obj?.foo?.[ERROR()]`),
		Case(`LET obj = { '[1]': 42 } RETURN obj?.[[1]]`, 42),
		Case(`LET obj = { '{"a":1}': 7 } RETURN obj?.[{a:1}]`, 7),
		CaseNil("RETURN ERROR()?.foo"),
		SkipCaseNil(`LET res = (FOR i IN ERROR() RETURN i)? RETURN res`),

		CaseArray(`LET res = (FOR i IN [1, 2, 3, 4] LET y = ERROR() RETURN y+i)? RETURN res`, []any{}, "Error in array comprehension"),
		CaseArray(`FOR i IN [1, 2, 3, 4] ERROR()? RETURN i`, []any{1, 2, 3, 4}, "Error in FOR loop"),
	}, vm.WithFunction("ERROR", func(ctx context.Context, args ...runtime.Value) (runtime.Value, error) {
		return nil, runtime.ErrNotImplemented
	}))
}

func TestTaggedTypes(t *testing.T) {
	type SomeValue struct {
		StrProp        string     `ferret:"strProp"`
		IntProp        int        `ferret:"intProp"`
		SliceProp      []int      `ferret:"sliceProp"`
		PointerProp    *SomeValue `ferret:"pointerProp"`
		NilPointerProp *SomeValue `ferret:"nilPointerProp"`

		UntaggedProp string

		privateStrProp string
	}

	RunUseCases(t, []UseCase{
		Case("RETURN GET_VALUE().strProp", "test"),
		Case("RETURN GET_VALUE().intProp", 99),
		CaseArray("RETURN GET_VALUE().sliceProp", []any{1, 2, 3}),
		Case("RETURN GET_VALUE().pointerProp.strProp", "nested"),
		CaseNil("RETURN GET_VALUE().nilPointerProp"),
		CaseNil("RETURN GET_VALUE().untagged"),
		CaseNil("RETURN GET_VALUE().privateStrProp"),
	}, vm.WithFunction("GET_VALUE", func(ctx context.Context, args ...runtime.Value) (runtime.Value, error) {
		return runtime.Encode(SomeValue{
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
