package compiler_test

import (
	"context"
	"fmt"
	"github.com/MontFerret/ferret/pkg/parser"
	"regexp"
	"strconv"
	"strings"
	"testing"

	"github.com/MontFerret/ferret/pkg/compiler"
	"github.com/MontFerret/ferret/pkg/runtime"
	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"

	. "github.com/smartystreets/goconvey/convey"
)

func TestVariables(t *testing.T) {
	RunUseCases(t, []UseCase{
		{
			`LET i = NONE RETURN i`,
			nil,
			nil,
		},
		{
			`LET a = TRUE RETURN a`,
			true,
			nil,
		},
		{
			`LET a = 1 RETURN a`,
			1,
			nil,
		},
		{
			`LET a = 1.1 RETURN a`,
			1.1,
			nil,
		},
		{
			`LET i = 'foo' RETURN i`,
			"foo",
			nil,
		},
		{
			`
		LET a = 'foo'
		LET b = a
		RETURN a`,
			"foo",
			nil,
		},
		{
			`LET i = [] RETURN i`,
			[]any{},
			ShouldEqualJSON,
		},
		{
			`LET i = [1, 2, 3] RETURN i`,
			[]any{1, 2, 3},
			ShouldEqualJSON,
		},
		{
			`LET i = [None, FALSE, "foo", 1, 1.1] RETURN i`,
			[]any{nil, false, "foo", 1, 1.1},
			ShouldEqualJSON,
		},
		{
			`
		LET n = None
		LET b = FALSE
		LET s = "foo"
		LET i = 1
		LET f = 1.1
		LET a = [n, b, s, i, f]
		RETURN a`,
			[]any{nil, false, "foo", 1, 1.1},
			ShouldEqualJSON,
		},
		{
			`LET i = {} RETURN i`,
			map[string]any{},
			ShouldEqualJSON,
		},
		{
			`LET i = {a: 1, b: 2} RETURN i`,
			map[string]any{"a": 1, "b": 2},
			ShouldEqualJSON,
		},
		{
			`LET i = {a: 1, b: [1]} RETURN i`,
			map[string]any{"a": 1, "b": []any{1}},
			ShouldEqualJSON,
		},
		{
			`LET i = {a: {c: 1}, b: [1]} RETURN i`,
			map[string]any{"a": map[string]any{"c": 1}, "b": []any{1}},
			ShouldEqualJSON,
		},
		{
			`LET i = {a: 'foo', b: 1, c: TRUE, d: [], e: {}} RETURN i`,
			map[string]any{"a": "foo", "b": 1, "c": true, "d": []any{}, "e": map[string]any{}},
			ShouldEqualJSON,
		},
		{
			`LET prop = "name" LET i = { [prop]: "foo" } RETURN i`,
			map[string]any{"name": "foo"},
			ShouldEqualJSON,
		},
		{
			`LET name="foo" LET i = { name } RETURN i`,
			map[string]any{"name": "foo"},
			ShouldEqualJSON,
		},
		{
			`LET i = [{a: {c: 1}, b: [1]}] RETURN i`,
			[]any{map[string]any{"a": map[string]any{"c": 1}, "b": []any{1}}},
			ShouldEqualJSON,
		},
		{
			"LET a = 'a' LET b = a LET c = 'c' RETURN b",
			"a",
			ShouldEqual,
		},
		{
			"LET i = (FOR i IN [1,2,3] RETURN i) RETURN i",
			[]int{1, 2, 3},
			ShouldEqualJSON,
		},
		{
			" LET i = { items: [1,2,3]}  FOR el IN i.items RETURN el",
			[]int{1, 2, 3},
			ShouldEqualJSON,
		},
		{
			`LET _ = (FOR i IN 1..100 RETURN NONE)
				RETURN TRUE`,
			true,
			ShouldEqualJSON,
		},
		{
			`
			LET src = NONE
			LET i = (FOR i IN src RETURN i)?
			RETURN i
		`,
			[]any{},
			nil,
		},
	})

	//Convey("Should compile LET i = (FOR i WHILE COUNTER() < 5 RETURN i) RETURN i", t, func() {
	//	c := compiler.New()
	//	counter := -1
	//	c.RegisterFunction("COUNTER", func(ctx context.visitor, args ...core.Second) (core.Second, error) {
	//		counter++
	//
	//		return values.NewInt(counter), nil
	//	})
	//
	//	p, err := c.Compile(`
	//		LET i = (FOR i WHILE COUNTER() < 5 RETURN i)
	//		RETURN i
	//	`)
	//
	//	So(err, ShouldBeNil)
	//	So(p, ShouldHaveSameTypeAs, &runtime.Program{})
	//
	//	out, err := p.Run(context.Background())
	//
	//	So(err, ShouldBeNil)
	//	So(string(out), ShouldEqual, "[0,1,2,3,4]")
	//})
	//
	//Convey("Should compile LET i = (FOR i WHILE COUNTER() < 5 T::FAIL() RETURN i)? RETURN i == NONE", t, func() {
	//	c := compiler.New()
	//	counter := -1
	//	c.RegisterFunction("COUNTER", func(ctx context.visitor, args ...core.Second) (core.Second, error) {
	//		counter++
	//
	//		return values.NewInt(counter), nil
	//	})
	//
	//	p, err := c.Compile(`
	//		LET i = (FOR i WHILE COUNTER() < 5 T::FAIL() RETURN i)?
	//		RETURN i == NONE
	//	`)
	//
	//	So(err, ShouldBeNil)
	//	So(p, ShouldHaveSameTypeAs, &runtime.Program{})
	//
	//	out, err := p.Run(context.Background())
	//
	//	So(err, ShouldBeNil)
	//	So(string(out), ShouldEqual, "true")
	//})

	Convey("Should not compile FOR foo IN foo", t, func() {
		c := compiler.New()

		_, err := c.Compile(`
			FOR foo IN foo
				RETURN foo
		`)

		So(err, ShouldNotBeNil)
	})

	Convey("Should not compile if a variable not defined", t, func() {
		c := compiler.New()

		_, err := c.Compile(`
			RETURN foo
		`)

		So(err, ShouldNotBeNil)
	})

	Convey("Should not compile if a variable is not unique", t, func() {
		c := compiler.New()

		_, err := c.Compile(`
			LET foo = "bar"
			LET foo = "baz"
	
			RETURN foo
		`)

		So(err, ShouldNotBeNil)
	})

	//SkipConvey("Should use value returned from WAITFOR EVENT", t, func() {
	//	out, err := newCompilerWithObservable().MustCompile(`
	//		LET obj = X::VAL("event", ["data"])
	//
	//		LET res = (WAITFOR EVENT "event" IN obj)
	//
	//		RETURN res
	//	`).Run(context.Background())
	//
	//	So(err, ShouldBeNil)
	//	So(string(out), ShouldEqual, `"data"`)
	//})
	//
	//SkipConvey("Should handle error from WAITFOR EVENT", t, func() {
	//	out, err := newCompilerWithObservable().MustCompile(`
	//		LET obj = X::VAL("foo", ["data"])
	//
	//		LET res = (WAITFOR EVENT "event" IN obj TIMEOUT 100)?
	//
	//		RETURN res == NONE
	//	`).Run(context.Background())
	//
	//	So(err, ShouldBeNil)
	//	So(string(out), ShouldEqual, `true`)
	//})
	//
	//SkipConvey("Should compare result of handled error", t, func() {
	//	out, err := newCompilerWithObservable().MustCompile(`
	//		LET obj = X::VAL("event", ["foo"], 1000)
	//
	//		LET res = (WAITFOR EVENT "event" IN obj TIMEOUT 100)? != NONE
	//
	//		RETURN res
	//	`).Run(context.Background())
	//
	//	So(err, ShouldBeNil)
	//	So(string(out), ShouldEqual, `false`)
	//})
	//

	Convey("Should not allow to use ignorable variable name", t, func() {
		c := compiler.New()

		_, err := c.Compile(`
			LET _ = (FOR i IN 1..100 RETURN NONE)
	
			RETURN _
		`)

		So(err, ShouldNotBeNil)
	})
}

func TestMathOperators(t *testing.T) {
	RunUseCases(t, []UseCase{
		{"RETURN 1 + 1", 2, nil},
		{"RETURN 1 - 1", 0, nil},
		{"RETURN 2 * 2", 4, nil},
		{"RETURN 4 / 2", 2, nil},
		{"RETURN 5 % 2", 1, nil},
	})
}

func TestUnaryOperators(t *testing.T) {
	RunUseCases(t, []UseCase{
		{"RETURN NOT TRUE", false, nil},
		{"RETURN !TRUE", false, nil},
		{"RETURN !FALSE", true, nil},
		{"RETURN -1", -1, nil},
		{"RETURN -1.1", -1.1, nil},
		{"RETURN +1", 1, nil},
		{"RETURN +1.1", 1.1, nil},
		{`LET v = 1 RETURN -v`, -1, nil},
		{`LET v = 1.1 RETURN -v`, -1.1, nil},
		{`LET v = -1 RETURN -v`, 1, nil},
		{`LET v = -1.1 RETURN -v`, 1.1, nil},
		{`LET v = -1 RETURN +v`, -1, nil},
		{`LET v = -1.1 RETURN +v`, -1.1, nil},
	})

	Convey("RETURN { enabled: !val}", t, func() {
		c := compiler.New()

		p1 := c.MustCompile(`
			LET val = ""
			RETURN { enabled: !val }
		`)

		v1, err := runtime.NewVM(p1).Run(context.Background(), nil)

		So(err, ShouldBeNil)

		out1, err := v1.MarshalJSON()
		So(err, ShouldBeNil)

		So(string(out1), ShouldEqual, `{"enabled":true}`)

		p2 := c.MustCompile(`
			LET val = ""
			RETURN { enabled: !!val }
		`)

		v2, err := runtime.NewVM(p2).Run(context.Background(), nil)

		So(err, ShouldBeNil)

		out2, err := v2.MarshalJSON()

		So(err, ShouldBeNil)
		So(string(out2), ShouldEqual, `{"enabled":false}`)
	})
}

func TestEqualityOperators(t *testing.T) {
	RunUseCases(t, []UseCase{
		{"RETURN 1 == 1", true, nil},
		{"RETURN 1 == 2", false, nil},
		{"RETURN 1 != 1", false, nil},
		{"RETURN 1 != 2", true, nil},
		{"RETURN 1 > 1", false, nil},
		{"RETURN 1 >= 1", true, nil},
		{"RETURN 1 < 1", false, nil},
		{"RETURN 1 <= 1", true, nil},
	})
}

func TestLogicalOperators(t *testing.T) {
	RunUseCases(t, []UseCase{
		{"RETURN 1 AND 0", 0, nil},
		{"RETURN 1 AND 1", 1, nil},
		{"RETURN 2 > 1 AND 1 > 0", true, nil},
		{"RETURN NONE && true", nil, nil},
		{"RETURN '' && true", "", nil},
		{"RETURN true && 23", 23, nil},
		{"RETURN 0 OR 1", 1, nil},
		{"RETURN 2 OR 1", 2, nil},
		{"RETURN 2 > 1 OR 1 > 0", true, nil},
		{"RETURN 2 < 1 OR 1 > 0", true, nil},
		{"RETURN 1 || 7", 1, nil},
		{"RETURN 0 || 7", 7, nil},
		{"RETURN NONE || 'foo'", "foo", nil},
		{
			`
			RETURN ERROR()? || 'boo'
		`,
			"boo",
			nil,
		},
		{
			`
			RETURN !ERROR()? && TRUE
		`,
			true,
			nil,
		},
		{
			`			LET u = { valid: false }
	
			RETURN NOT u.valid`,
			true,
			nil,
		},
	}, runtime.WithFunction("ERROR", func(ctx context.Context, args ...core.Value) (core.Value, error) {
		return values.None, fmt.Errorf("test")
	}))
}

func TestTernaryOperator(t *testing.T) {
	RunUseCases(t, []UseCase{
		{"RETURN 1 < 2 ? 3 : 4", 3, nil},
		{"RETURN 1 > 2 ? 3 : 4", 4, nil},
		{"RETURN 2 ? : 4", 2, nil},
		{`
				LET foo = TRUE
				RETURN foo ? TRUE : FALSE
				`, true, nil},
		{`
				LET foo = FALSE
				RETURN foo ? TRUE : FALSE
				`, false, nil},
		{
			`
					FOR i IN [1, 2, 3, 4, 5, 6]
						RETURN i < 3 ? i * 3 : i * 2
				`,
			[]int{3, 6, 6, 8, 10, 12},
			ShouldEqualJSON,
		},
		{
			`
					FOR i IN [1, 2, 3, 4, 5, 6]
						RETURN i < 3 ? : i * 2
				`,
			[]any{true, true, 6, 8, 10, 12},
			ShouldEqualJSON,
		},
		{
			`
					FOR i IN [NONE, 2, 3, 4, 5, 6]
						RETURN i ? : i
		`,
			[]any{nil, 2, 3, 4, 5, 6},
			ShouldEqualJSON,
		},
		{
			`RETURN 0 && true ? "1" : "some"`,
			"some",
			nil,
		},
		{
			`RETURN length([]) > 0 && true ? "1" : "some"`,
			"some",
			nil,
		},
	})

	Convey("Should compile ternary operator with default values", t, func() {
		vals := []string{
			"0",
			"0.0",
			"''",
			"NONE",
			"FALSE",
		}

		c := compiler.New()

		for _, val := range vals {
			p, err := c.Compile(fmt.Sprintf(`
			FOR i IN [%s, 1, 2, 3]
				RETURN i ? i * 2 : 'no value'
		`, val))

			So(err, ShouldBeNil)

			out, err := Run(p)

			So(err, ShouldBeNil)

			So(string(out), ShouldEqual, `["no value",2,4,6]`)
		}
	})
}

func TestLikeOperator(t *testing.T) {
	RunUseCases(t, []UseCase{
		{`RETURN "foo" LIKE "f*"`, true, nil},
		{`RETURN "foo" LIKE "b*"`, false, nil},
		{`RETURN "foo" NOT LIKE "f*"`, false, nil},
		{`RETURN "foo" NOT LIKE "b*"`, true, nil},
		{`LET res = "foo" LIKE  "f*"
			RETURN res`, true, nil},
		{`RETURN ("foo" LIKE  "b*") ? "foo" : "bar"`, `bar`, nil},
		{`RETURN ("foo" NOT LIKE  "b*") ? "foo" : "bar"`, `foo`, nil},
		{`RETURN true ? ("foo" NOT LIKE  "b*") : false`, true, nil},
		{`RETURN true ? false : ("foo" NOT LIKE  "b*")`, false, nil},
		{`RETURN false ? false : ("foo" NOT LIKE  "b*")`, true, nil},
	})

	//Convey("FOR IN LIKE", t, func() {
	//	c := compiler.New()
	//
	//	out1, err := c.MustCompile(`
	//		FOR str IN ["foo", "bar", "qaz"]
	//			FILTER str LIKE "*a*"
	//			RETURN str
	//	`).Run(context.Background())
	//
	//	So(err, ShouldBeNil)
	//	So(string(out1), ShouldEqual, `["bar","qaz"]`)
	//})
	//
	//Convey("FOR IN LIKE 2", t, func() {
	//	c := compiler.New()
	//
	//	out1, err := c.MustCompile(`
	//		FOR str IN ["foo", "bar", "qaz"]
	//			FILTER str LIKE "*a*"
	//			RETURN str
	//	`).Run(context.Background())
	//
	//	So(err, ShouldBeNil)
	//	So(string(out1), ShouldEqual, `["bar","qaz"]`)
	//})

}

func TestRange(t *testing.T) {
	RunUseCases(t, []UseCase{
		{
			"RETURN 1..10",
			[]any{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
			ShouldEqualJSON,
		},
		{
			"RETURN 10..1",
			[]any{10, 9, 8, 7, 6, 5, 4, 3, 2, 1},
			ShouldEqualJSON,
		},
		{
			`
		LET start = 1
		LET end = 10
		RETURN start..end
		`,
			[]any{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
			ShouldEqualJSON,
		},
		//{
		//	`
		//LET start = @start
		//LET end = @end
		//RETURN start..end
		//`,
		//	[]any{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
		//	ShouldEqualJSON,
		//},
	})
}

func TestFunctionCall(t *testing.T) {
	RunUseCases(t, []UseCase{
		{
			"RETURN TYPENAME(1)",
			"int",
			nil,
		},
		{
			"WAIT(10) RETURN 1",
			1,
			nil,
		},
		{
			"RETURN LENGTH([1,2,3])",
			3,
			nil,
		},
		{
			"RETURN CONCAT('a', 'b', 'c')",
			"abc",
			nil,
		},
		{
			"RETURN CONCAT(CONCAT('a', 'b'), 'c', CONCAT('d', 'e'))",
			"abcde",
			nil,
		},
		{
			`
		LET arr = []
		LET a = 1
		LET res = APPEND(arr, a)
		RETURN res
		`,
			[]any{1},
			ShouldEqualJSON,
		},
		{
			"LET duration = 10 WAIT(duration) RETURN 1",
			1,
			nil,
		},
		{
			"RETURN (FALSE OR T::FAIL())?",
			nil,
			nil,
		},
		{
			"RETURN T::FAIL()?",
			nil,
			nil,
		},
		{
			`FOR i IN [1, 2, 3, 4]
				LET duration = 10
		
				WAIT(duration)
		
				RETURN i * 2`,
			[]int{2, 4, 6, 8},
			ShouldEqualJSON,
		},
		{
			`RETURN FIRST((FOR i IN 1..10 RETURN i * 2))`,
			2,
			nil,
		},
		{
			`RETURN UNION((FOR i IN 0..5 RETURN i), (FOR i IN 6..10 RETURN i))`,
			[]int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
			ShouldEqualJSON,
		},
	})
}

func TestMember(t *testing.T) {
	RunUseCases(t, []UseCase{
		{
			"LET arr = [1,2,3,4] RETURN arr[10]",
			nil,
			nil,
		},
		{
			"LET arr = [1,2,3,4] RETURN arr[1]",
			2,
			nil,
		},
		{
			"LET arr = [1,2,3,4] LET idx = 1 RETURN arr[idx]",
			2,
			nil,
		},
		{
			`LET obj = { foo: "bar", qaz: "wsx"} RETURN obj["qaz"]`,
			"wsx",
			nil,
		},
		{
			fmt.Sprintf(`
								LET obj = { "foo": "bar", %s: "wsx"}
		
								RETURN obj["qaz"]
							`, "`qaz`"),
			"wsx",
			nil,
		},
		{
			fmt.Sprintf(`
								LET obj = { "foo": "bar", %s: "wsx"}
		
								RETURN obj["let"]
							`, "`let`"),
			"wsx",
			nil,
		},
		{
			`LET obj = { foo: "bar", qaz: "wsx"} LET key = "qaz" RETURN obj[key]`,
			"wsx",
			nil,
		},
		{
			`RETURN { foo: "bar" }.foo`,
			"bar",
			nil,
		},
		{
			`LET inexp = 1 IN {'foo': [1]}.foo
			LET ternaryexp = FALSE ? TRUE : {foo: TRUE}.foo
			RETURN inexp && ternaryexp`,
			true,
			nil,
		},
		{
			`RETURN ["bar", "foo"][0]`,
			"bar",
			nil,
		},
		{
			`LET inexp = 1 IN [[1]][0]
								LET ternaryexp = FALSE ? TRUE : [TRUE][0]
								RETURN inexp && ternaryexp`,
			true,
			nil,
		},
		{
			`LET obj = {
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
			true,
			nil,
		},
		{
			`LET o1 = {
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

			true,
			nil,
		},
		{
			`LET o1 = {
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

			true,
			nil,
		},
		{
			`LET obj = {
							attributes: {
								'data-index': 1
							}
						}
		
						RETURN obj.attributes['data-index']`,
			1,
			nil,
		},
		{
			`LET obj = NONE RETURN obj.foo`,
			nil,
			ShouldBeError,
		},
		{
			`LET obj = NONE RETURN obj?.foo`,
			nil,
			nil,
		},
		{
			`RETURN {first: {second: "third"}}.first`,
			map[string]any{
				"second": "third",
			},
			nil,
		},
		{
			`RETURN KEEP_KEYS({first: {second: "third"}}.first, "second")`,
			map[string]any{
				"second": "third",
			},
			nil,
		},
		{
			`
					FOR v, k IN {f: {foo: "bar"}}.f
						RETURN [k, v]
				`,
			[]any{
				[]any{"foo", "bar"},
			},
			nil,
		},
		{
			`RETURN FIRST([[1, 2]][0])`,
			1,
			nil,
		},
		{
			`RETURN [[1, 2]][0]`,
			[]any{1, 2},
			ShouldEqualJSON,
		},
		{
			`
					FOR i IN [[1, 2]][0]
						RETURN i
				`,
			[]any{1, 2},
			ShouldEqualJSON,
		},
		{
			`
					LET arr = [{ name: "Bob" }]
		
					RETURN FIRST(arr).name
				`,
			"Bob",
			nil,
		},
		{
			`
					LET arr = [{ name: { first: "Bob" } }]
	
					RETURN FIRST(arr)['name'].first
				`,
			"Bob",
			nil,
		},
		{
			`
					LET obj = { foo: None }
	
					RETURN obj.foo?.bar
				`,
			nil,
			nil,
		},
	})
}

func TestMemberReservedWords(t *testing.T) {
	RunUseCases(t, []UseCase{})

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

				c := compiler.New()
				prog, err := c.Compile(query.String())

				So(err, ShouldBeNil)

				out, err := Exec(prog, true, runtime.WithFunctions(c.Functions().Unwrap()))

				So(err, ShouldBeNil)
				So(out, ShouldEqual, expected.String())
			}
		}
	})
}

func TestOptionalChaining(t *testing.T) {
	RunUseCases(t, []UseCase{
		{
			`
					LET obj = { foo: { bar: "bar" } }
	
					RETURN obj.foo?.bar
				`,
			"bar",
			nil,
		},
		{
			`
					LET obj = { foo: None }
	
					RETURN obj.foo?.bar?.[0]
				`,
			nil,
			nil,
		},
		{
			`
					LET obj = { foo: { bar: ["bar"] } }
	
					RETURN obj.foo?.bar?.[0]
				`,
			"bar",
			nil,
		},
		{
			`
					RETURN FIRST([])?.foo
				`,
			nil,
			nil,
		},
		{
			`
					RETURN FIRST([{ foo: "bar" }])?.foo
				`,
			"bar",
			nil,
		},
	})

	RunUseCases(t, []UseCase{
		{
			`
					RETURN ERROR()?.foo
				`,
			nil,
			nil,
		},
	}, runtime.WithFunction("ERROR", func(ctx context.Context, args ...core.Value) (core.Value, error) {
		return nil, core.ErrNotImplemented
	}))
}

func TestFor(t *testing.T) {
	// Should not allocate memory if NONE is a return statement
	//{
	//	`FOR i IN 0..100
	//		RETURN NONE`,
	//	[]any{},
	//	ShouldEqualJSON,
	//},
	RunUseCases(t, []UseCase{
		{
			"FOR i IN 1..5 RETURN i",
			[]any{1, 2, 3, 4, 5},
			ShouldEqualJSON,
		},
		{
			`
		FOR i IN 1..5
			LET x = i * 2
			RETURN x
		`,
			[]any{2, 4, 6, 8, 10},
			ShouldEqualJSON,
		},
		{
			`
		FOR val, counter IN 1..5
			LET x = val
			PRINT(counter)
			LET y = counter
			RETURN [x, y]
				`,
			[]any{[]any{1, 0}, []any{2, 1}, []any{3, 2}, []any{4, 3}, []any{5, 4}},
			ShouldEqualJSON,
		},
		{
			`FOR i IN [] RETURN i
				`,
			[]any{},
			ShouldEqualJSON,
		},
		{
			`FOR i IN [1, 2, 3] RETURN i
				`,
			[]any{1, 2, 3},
			ShouldEqualJSON,
		},

		{
			`FOR i, k IN [1, 2, 3] RETURN k`,
			[]any{0, 1, 2},
			ShouldEqualJSON,
		},
		{
			`FOR i IN ['foo', 'bar', 'qaz'] RETURN i`,
			[]any{"foo", "bar", "qaz"},
			ShouldEqualJSON,
		},
		{
			`FOR i IN {a: 'bar', b: 'foo', c: 'qaz'} RETURN i`,
			[]any{"foo", "bar", "qaz"},
			ShouldHaveSameItems,
		},
		{
			`FOR i, k IN {a: 'foo', b: 'bar', c: 'qaz'} RETURN k`,
			[]any{"a", "b", "c"},
			ShouldHaveSameItems,
		},
		{
			`FOR i IN [{name: 'foo'}, {name: 'bar'}, {name: 'qaz'}] RETURN i.name`,
			[]any{"foo", "bar", "qaz"},
			ShouldHaveSameItems,
		},
		{
			`FOR i IN { items: [{name: 'foo'}, {name: 'bar'}, {name: 'qaz'}] }.items RETURN i.name`,
			[]any{"foo", "bar", "qaz"},
			ShouldHaveSameItems,
		},
		{
			`FOR prop IN ["a"]
							FOR val IN [1, 2, 3]
								RETURN {[prop]: val}`,
			[]any{map[string]any{"a": 1}, map[string]any{"a": 2}, map[string]any{"a": 3}},
			ShouldEqualJSON,
		},
		{
			`FOR val IN 1..3
							FOR prop IN ["a"]
								RETURN {[prop]: val}`,
			[]any{map[string]any{"a": 1}, map[string]any{"a": 2}, map[string]any{"a": 3}},
			ShouldEqualJSON,
		},
		{
			`FOR prop IN ["a"]
							FOR val IN 1..3
								RETURN {[prop]: val}`,
			[]any{map[string]any{"a": 1}, map[string]any{"a": 2}, map[string]any{"a": 3}},
			ShouldEqualJSON,
		},
		{
			`FOR prop IN ["a"]
							FOR val IN [1, 2, 3]
								FOR val2 IN [1, 2, 3]
									RETURN { [prop]: [val, val2] }`,
			[]any{map[string]any{"a": []int{1, 1}}, map[string]any{"a": []int{1, 2}}, map[string]any{"a": []int{1, 3}}, map[string]any{"a": []int{2, 1}}, map[string]any{"a": []int{2, 2}}, map[string]any{"a": []int{2, 3}}, map[string]any{"a": []int{3, 1}}, map[string]any{"a": []int{3, 2}}, map[string]any{"a": []int{3, 3}}},
			ShouldEqualJSON,
		},
		{
			`FOR val IN [1, 2, 3]
							RETURN (
								FOR prop IN ["a", "b", "c"]
									RETURN { [prop]: val }
							)`,
			[]any{[]any{map[string]any{"a": 1}, map[string]any{"b": 1}, map[string]any{"c": 1}}, []any{map[string]any{"a": 2}, map[string]any{"b": 2}, map[string]any{"c": 2}}, []any{map[string]any{"a": 3}, map[string]any{"b": 3}, map[string]any{"c": 3}}},
			ShouldEqualJSON,
		},
		{
			`FOR val IN [1, 2, 3]
							LET sub = (
								FOR prop IN ["a", "b", "c"]
									RETURN { [prop]: val }
							)
		
							RETURN sub`,
			[]any{[]any{map[string]any{"a": 1}, map[string]any{"b": 1}, map[string]any{"c": 1}}, []any{map[string]any{"a": 2}, map[string]any{"b": 2}, map[string]any{"c": 2}}, []any{map[string]any{"a": 3}, map[string]any{"b": 3}, map[string]any{"c": 3}}},
			ShouldEqualJSON,
		},
		{
			`FOR i IN [ 1, 2, 3, 4, 1, 3 ]
							RETURN DISTINCT i
		`,
			[]any{1, 2, 3, 4},
			ShouldEqualJSON,
		},
	})
}

func TestForWhile(t *testing.T) {
	var counter int64
	RunUseCases(t, []UseCase{
		//{
		//	"FOR i WHILE false RETURN i",
		//	[]any{},
		//	ShouldEqualJSON,
		//},
		{
			"FOR i WHILE UNTIL(5) RETURN i",
			[]any{0, 1, 2, 3, 4},
			ShouldEqualJSON,
		},
	}, runtime.WithFunctions(map[string]core.Function{
		"UNTIL": func(ctx context.Context, args ...core.Value) (core.Value, error) {
			if counter < int64(values.ToInt(args[0])) {
				counter++

				return values.True, nil
			}

			return values.False, nil
		},
	}))
}

func TestForFilter(t *testing.T) {
	RunUseCases(t, []UseCase{
		{
			`
			FOR i IN [ 1, 2, 3, 4, 1, 3 ]
				FILTER i > 2
				RETURN i
		`,
			[]any{3, 4, 3},
			ShouldEqualJSON,
		},
		{
			`
			FOR i IN [ 1, 2, 3, 4, 1, 3 ]
				FILTER i > 1 AND i < 4
				RETURN i
		`,
			[]any{2, 3, 3},
			ShouldEqualJSON,
		},
		{
			`
			LET users = [
				{
					age: 31,
					gender: "m",
					name: "Josh"
				},
				{
					age: 29,
					gender: "f",
					name: "Mary"
				},
				{
					age: 36,
					gender: "m",
					name: "Peter"
				}
			]
			FOR u IN users
				FILTER u.name =~ "r"
				RETURN u
		`,
			[]any{map[string]any{"age": 29, "gender": "f", "name": "Mary"}, map[string]any{"age": 36, "gender": "m", "name": "Peter"}},
			ShouldEqualJSON,
		},
		{
			`
					LET users = [
						{
							active: true,
							age: 31,
							gender: "m"
						},
						{
							active: true,
							age: 29,
							gender: "f"
						},
						{
							active: true,
							age: 36,
							gender: "m"
						}
					]
					FOR u IN users
						FILTER u.active == true
						FILTER u.age < 35
						RETURN u
				`,
			[]any{map[string]any{"active": true, "gender": "m", "age": 31}, map[string]any{"active": true, "gender": "f", "age": 29}},
			ShouldEqualJSON,
		},
		{
			`
			LET users = [
				{
					active: true,
					age: 31,
					gender: "m"
				},
				{
					active: true,
					age: 29,
					gender: "f"
				},
				{
					active: true,
					age: 36,
					gender: "m"
				},
				{
					active: false,
					age: 69,
					gender: "m"
				}
			]
			FOR u IN users
				FILTER u.active
				RETURN u
				`,
			[]any{map[string]any{"active": true, "gender": "m", "age": 31}, map[string]any{"active": true, "gender": "f", "age": 29}, map[string]any{"active": true, "gender": "m", "age": 36}},
			ShouldEqualJSON,
		},
		{
			`
			LET users = [
				{
					active: true,
					age: 31,
					gender: "m"
				},
				{
					active: true,
					age: 29,
					gender: "f"
				},
				{
					active: true,
					age: 36,
					gender: "m"
				},
				{
					active: false,
					age: 69,
					gender: "m"
				}
			]
			FOR u IN users
				FILTER u.active == true
				LIMIT 2
				FILTER u.gender == "m"
				RETURN u
		`,
			[]any{map[string]any{"active": true, "gender": "m", "age": 31}},
			ShouldEqualJSON,
		},
	})
}

func TestForLimit(t *testing.T) {
	RunUseCases(t, []UseCase{
		{
			`
			FOR i IN [ 1, 2, 3, 4, 1, 3 ]
				LIMIT 2
				RETURN i
		`,
			[]any{1, 2},
			ShouldEqualJSON,
		},
	})
}
