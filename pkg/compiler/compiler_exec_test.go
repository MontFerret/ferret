package compiler_test

import (
	"context"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"testing"

	"github.com/MontFerret/ferret/pkg/parser"

	"github.com/MontFerret/ferret/pkg/compiler"
	"github.com/MontFerret/ferret/pkg/runtime"
	"github.com/MontFerret/ferret/pkg/runtime/core"
	. "github.com/smartystreets/goconvey/convey"
)

func TestString(t *testing.T) {
	RunUseCases(t, []UseCase{
		Case(
			`
			RETURN "
FOO
BAR
"
		`, "\nFOO\nBAR\n", "Should be possible to use multi line string"),

		CaseJSON(
			fmt.Sprintf(`
RETURN %s<!DOCTYPE html>
		<html lang="en">
		<head>
		<meta charset="UTF-8">
		<title>GetTitle</title>
		</head>
		<body>
			Hello world
		</body>
		</html>%s
`, "`", "`"), `<!DOCTYPE html>
		<html lang="en">
		<head>
		<meta charset="UTF-8">
		<title>GetTitle</title>
		</head>
		<body>
			Hello world
		</body>
		</html>`, "Should be possible to use multi line string with nested strings using backtick"),

		CaseJSON(
			fmt.Sprintf(`
RETURN %s<!DOCTYPE html>
		<html lang="en">
		<head>
		<meta charset="UTF-8">
		<title>GetTitle</title>
		</head>
		<body>
			Hello world
		</body>
		</html>%s
`, "´", "´"),
			`<!DOCTYPE html>
		<html lang="en">
		<head>
		<meta charset="UTF-8">
		<title>GetTitle</title>
		</head>
		<body>
			Hello world
		</body>
		</html>`, "Should be possible to use multi line string with nested strings using tick"),
	})
}

func TestVariables(t *testing.T) {
	RunUseCases(t, []UseCase{
		CaseCompilationError(`RETURN foo`, "Should not compile if a variable not defined"),
		CaseCompilationError(`
			LET foo = "bar"
			LET foo = "baz"

			RETURN foo
		`, "Should not compile if a variable is not unique"),
		CaseNil(`LET i = NONE RETURN i`),
		Case(`LET a = TRUE RETURN a`, true),
		Case(`LET a = 1 RETURN a`, 1),
		Case(`LET a = 1.1 RETURN a`, 1.1),
		Case(`LET a = "foo" RETURN a`, "foo"),
		Case(
			`
		LET a = 'foo'
		LET b = a
		RETURN a`,
			"foo",
		),
		CaseArray(`LET i = [] RETURN i`, []any{}),
		CaseArray(`LET i = [1, 2, 3] RETURN i`, []any{1, 2, 3}),
		CaseArray(`LET i = [None, FALSE, "foo", 1, 1.1] RETURN i`, []any{nil, false, "foo", 1, 1.1}),
		CaseArray(`
		LET n = None
		LET b = FALSE
		LET s = "foo"
		LET i = 1
		LET f = 1.1
		LET a = [n, b, s, i, f]
		RETURN a`, []any{nil, false, "foo", 1, 1.1}),
		CaseObject(`LET i = {} RETURN i`, map[string]any{}),
		CaseObject(`LET i = {a: 1, b: 2} RETURN i`, map[string]any{"a": 1, "b": 2}),
		CaseObject(`LET i = {a: 1, b: [1]} RETURN i`, map[string]any{"a": 1, "b": []any{1}}, "Nested array in object"),
		CaseObject(`LET i = {a: {c: 1}, b: [1]} RETURN i`,
			map[string]any{"a": map[string]any{"c": 1}, "b": []any{1}}, "Nested object in object"),
		CaseObject(`LET i = {a: 'foo', b: 1, c: TRUE, d: [], e: {}} RETURN i`,
			map[string]any{"a": "foo", "b": 1, "c": true, "d": []any{}, "e": map[string]any{}}, "Complex object"),
		CaseObject(`LET prop = "name" LET i = { [prop]: "foo" } RETURN i`,
			map[string]any{"name": "foo"}, "Computed property name"),
		CaseObject(`LET name="foo" LET i = { name } RETURN i`,
			map[string]any{"name": "foo"}, "Property name shorthand"),
		CaseArray(`LET i = [{a: {c: 1}, b: [1]}] RETURN i`,
			[]any{map[string]any{"a": map[string]any{"c": 1}, "b": []any{1}}}, "Nested object in array"),
		Case("LET a = 'a' LET b = a LET c = 'c' RETURN b",
			"a", "Variable reference"),
		CaseArray("LET i = (FOR i IN [1,2,3] RETURN i) RETURN i",
			[]any{1, 2, 3}, "arrayList comprehension"),
		CaseArray(" LET i = { items: [1,2,3]}  FOR el IN i.items RETURN el",
			[]any{1, 2, 3}, "hashMap property access for a loop source"),
		Case(`LET _ = (FOR i IN 1..100 RETURN NONE) RETURN TRUE`, true),
		CaseArray(`
			LET src = NONE
			LET i = (FOR i IN src RETURN i)?
			RETURN i
		`,
			[]any{}, "Error handling in array comprehension"),
		CaseCompilationError(`			LET _ = (FOR i IN 1..100 RETURN NONE)
	
			RETURN _`, "Should not allow to use ignorable variable name"),
		Case(`
			LET _ = (FOR i IN 1..100 RETURN NONE)
			LET _ = (FOR i IN 1..100 RETURN NONE)

			RETURN TRUE
		`, true),
	})

	Convey("Should compile LET i = (FOR i WHILE COUNTER() < 5 RETURN i) RETURN i", t, func() {
		c := compiler.New()

		p, err := c.Compile(`
			LET i = (FOR i WHILE COUNTER() < 5 RETURN i)
			RETURN i
		`)

		So(err, ShouldBeNil)
		So(p, ShouldHaveSameTypeAs, &runtime.Program{})

		counter := -1
		out, err := Run(p, runtime.WithFunction("COUNTER", func(ctx context.Context, args ...core.Value) (core.Value, error) {
			counter++

			return core.NewInt(counter), nil
		}))

		So(err, ShouldBeNil)
		So(string(out), ShouldEqual, "[0,1,2,3,4]")
	})

	Convey("Should compile LET i = (FOR i WHILE COUNTER() < 5 T::FAIL() RETURN i)? RETURN length(i) == 0", t, func() {
		c := compiler.New()

		p, err := c.Compile(`
			LET i = (FOR i WHILE COUNTER() < 5 T::FAIL() RETURN i)?
			RETURN length(i) == 0
		`)

		So(err, ShouldBeNil)
		So(p, ShouldHaveSameTypeAs, &runtime.Program{})

		counter := -1
		out, err := Run(p, runtime.WithFunction("COUNTER", func(ctx context.Context, args ...core.Value) (core.Value, error) {
			counter++

			return core.NewInt(counter), nil
		}), runtime.WithFunction("T::FAIL", func(ctx context.Context, args ...core.Value) (core.Value, error) {
			return core.None, fmt.Errorf("test")
		}))

		So(err, ShouldBeNil)
		So(string(out), ShouldEqual, "true")
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
}

func TestParam(t *testing.T) {
	RunUseCases(t,
		[]UseCase{
			CaseRuntimeErrorAs(`RETURN @foo`, core.Error(runtime.ErrMissedParam, "@foo")),
			Case(`RETURN @str`, "bar", "Should return a value of a parameter"),
			Case(`RETURN @int + @int`, 2, "Should return a sum of two parameters"),
			Case(`RETURN @obj.str1 + @obj.str2`, "foobar", "Should return a concatenated string of two parameter properties"),
			CaseArray(`FOR i IN @values1 RETURN i`, []any{1, 2, 3, 4}, "Should iterate over an array parameter"),
			CaseArray(`FOR i IN @values2 SORT i RETURN i`, []any{"a", "b", "c", "d"}, "Should iterate over an object parameter"),
			CaseArray(`FOR i IN @start..@end RETURN i`, []any{1, 2, 3, 4, 5}, "Should iterate over a range parameter"),
			Case(`RETURN @obj.str1`, "foo", "Should be possible to use in member expression"),
		},
		runtime.WithParam("str", "bar"),
		runtime.WithParam("int", 1),
		runtime.WithParam("bool", true),
		runtime.WithParam("obj", map[string]interface{}{"str1": "foo", "str2": "bar"}),
		runtime.WithParam("values1", []int{1, 2, 3, 4}),
		runtime.WithParam("values2", map[string]interface{}{"a": "a", "b": "b", "c": "c", "d": "d"}),
		runtime.WithParam("start", 1),
		runtime.WithParam("end", 5),
	)
}

func TestMathOperators(t *testing.T) {
	RunUseCases(t, []UseCase{
		Case(`RETURN 1 + 1`, 2),
		Case(`RETURN 1 - 1`, 0),
		Case(`RETURN 2 * 2`, 4),
		Case(`RETURN 4 / 2`, 2),
		Case(`RETURN 5 % 2`, 1),
	})
}

func TestUnaryOperators(t *testing.T) {
	RunUseCases(t, []UseCase{
		Case("RETURN !TRUE", false),
		Case("RETURN NOT TRUE", false),
		Case("RETURN !FALSE", true),
		Case("RETURN -1", -1),
		Case("RETURN -1.1", -1.1),
		Case("RETURN +1", 1),
		Case("RETURN +1.1", 1.1),
		Case("LET v = 1 RETURN -v", -1),
		Case("LET v = 1.1 RETURN -v", -1.1),
		Case("LET v = -1 RETURN -v", 1),
		Case("LET v = -1.1 RETURN -v", 1.1),
		Case("LET v = -1 RETURN +v", -1),
		Case("LET v = -1.1 RETURN +v", -1.1),
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
		Case("RETURN 1 == 1", true),
		Case("RETURN 1 == 2", false),
		Case("RETURN 1 != 1", false),
		Case("RETURN 1 != 2", true),
		Case("RETURN 1 > 1", false),
		Case("RETURN 1 >= 1", true),
		Case("RETURN 1 < 1", false),
		Case("RETURN 1 <= 1", true),
	})
}

func TestLogicalOperators(t *testing.T) {
	RunUseCases(t, []UseCase{
		Case("RETURN 1 AND 0", 0),
		Case("RETURN 1 AND 1", 1),
		Case("RETURN 2 > 1 AND 1 > 0", true),
		Case("RETURN NONE && true", nil),
		Case("RETURN '' && true", ""),
		Case("RETURN true && 23", 23),
		Case("RETURN 1 OR 0", 1),
		Case("RETURN 0 OR 1", 1),
		Case("RETURN 2 OR 1", 2),
		Case("RETURN 2 > 1 OR 1 > 0", true),
		Case("RETURN 2 < 1 OR 1 > 0", true),
		Case("RETURN 1 || 7", 1),
		Case("RETURN 0 || 7", 7),
		Case("RETURN NONE || 'foo'", "foo"),
		Case("RETURN '' || 'foo'", "foo"),
		Case(`RETURN ERROR()? || 'boo'`, "boo"),
		Case(`RETURN !ERROR()? && TRUE`, true),
		Case(`LET u = { valid: false } RETURN u.valid || TRUE`, true),
	}, runtime.WithFunction("ERROR", func(ctx context.Context, args ...core.Value) (core.Value, error) {
		return core.None, fmt.Errorf("test")
	}))
}

func TestTernaryOperator(t *testing.T) {
	RunUseCases(t, []UseCase{
		Case("RETURN 1 < 2 ? 3 : 4", 3),
		Case("RETURN 1 > 2 ? 3 : 4", 4),
		Case("RETURN 2 ? : 4", 2),
		Case("LET foo = TRUE RETURN foo ? TRUE : FALSE", true),
		Case("LET foo = FALSE RETURN foo ? TRUE : FALSE", false),
		CaseArray("FOR i IN [1, 2, 3, 4, 5, 6] RETURN i < 3 ? i * 3 : i * 2", []any{3, 6, 6, 8, 10, 12}),
		CaseArray(`FOR i IN [NONE, 2, 3, 4, 5, 6] RETURN i ? : i`, []any{nil, 2, 3, 4, 5, 6}),
		Case(`RETURN 0 && true ? "1" : "some"`, "some"),
		Case(`RETURN length([]) > 0 && true ? "1" : "some"`, "some"),
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
		Case(`RETURN "foo" LIKE "f*"`, true),
		Case(`RETURN "foo" LIKE "b*"`, false),
		Case(`RETURN "foo" NOT LIKE "f*"`, false),
		Case(`RETURN "foo" NOT LIKE "b*"`, true),
		Case(`LET res = "foo" LIKE  "f*" RETURN res`, true),
		Case(`RETURN ("foo" LIKE  "b*") ? "foo" : "bar"`, "bar"),
		Case(`RETURN ("foo" NOT LIKE  "b*") ? "foo" : "bar"`, "foo"),
		Case(`RETURN true ? ("foo" NOT LIKE  "b*") : false`, true),
		Case(`RETURN true ? false : ("foo" NOT LIKE  "b*")`, false),
		Case(`RETURN false ? false : ("foo" NOT LIKE  "b*")`, true),
		CaseArray(`FOR str IN ["foo", "bar", "qaz"] FILTER str LIKE "*a*" RETURN str`, []any{"bar", "qaz"}),
		CaseArray(`FOR str IN ["foo", "bar", "qaz"] FILTER str NOT LIKE "*a*" RETURN str`, []any{"foo"}),
		CaseArray(`FOR str IN ["ar", "bar", "qaz"] FILTER str LIKE "a*" RETURN str`, []any{"ar"}),
		CaseArray(`FOR str IN ["ar", "bar", "qaz", "fa", "da"] FILTER str LIKE "*a" RETURN str`, []any{"fa", "da"}),
	})
}

func TestRegexpOperator(t *testing.T) {
	RunUseCases(t, []UseCase{
		Case(`RETURN "foo" =~ "^f[o].$" `, true),
		Case(`RETURN "foo" !~ "[a-z]+bar$"`, true),
		Case(`RETURN "foo" !~ T::REGEXP()`, true),
	}, runtime.WithFunction("T::REGEXP", func(_ context.Context, _ ...core.Value) (value core.Value, e error) {
		return core.NewString("[a-z]+bar$"), nil
	}))

	// TODO: Fix
	SkipConvey("Should return an error during compilation when a regexp string invalid", t, func() {
		_, err := compiler.New().
			Compile(`
			RETURN "foo" !~ "[ ]\K(?<!\d )(?=(?: ?\d){8})(?!(?: ?\d){9})\d[ \d]+\d" 
		`)

		So(err, ShouldBeError)
	})

	// TODO: Fix
	SkipConvey("Should return an error during compilation when a regexp is not a string", t, func() {
		right := []string{
			"[]",
			"{}",
			"1",
			"1.1",
			"TRUE",
		}

		for _, r := range right {
			_, err := compiler.New().
				Compile(fmt.Sprintf(`
			RETURN "foo" !~ %s 
		`, r))

			So(err, ShouldBeError)
		}
	})
}

func TestInOperator(t *testing.T) {
	RunUseCases(t, []UseCase{
		Case("RETURN 1 IN [1,2,3]", true),
		Case("RETURN 4 IN [1,2,3]", false),
		Case("RETURN 1 NOT IN [1,2,3]", false),
	})
}

func TestArrayAllOperator(t *testing.T) {
	RunUseCases(t, []UseCase{
		Case("RETURN [1,2,3] ALL IN [1,2,3]", true, "All elements are in"),
	})
}

func TestRange(t *testing.T) {
	RunUseCases(t, []UseCase{
		CaseArray("RETURN 1..10", []any{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}),
		CaseArray("RETURN 10..1", []any{10, 9, 8, 7, 6, 5, 4, 3, 2, 1}),
		CaseArray(
			`
		LET start = 1
		LET end = 10
		RETURN start..end
		`,
			[]any{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
		),
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
		Case("RETURN TYPENAME(1)", "int"),
		Case("RETURN TYPENAME(1.1)", "float"),
		Case("WAIT(10) RETURN 1", 1),
		Case("RETURN LENGTH([1,2,3])", 3),
		Case("RETURN CONCAT('a', 'b', 'c')", "abc"),
		Case("RETURN CONCAT(CONCAT('a', 'b'), 'c', CONCAT('d', 'e'))", "abcde", "Nested calls"),
		CaseArray(`
		LET arr = []
		LET a = 1
		LET res = APPEND(arr, a)
		RETURN res
		`,
			[]any{1}, "Append to array"),
		Case("LET duration = 10 WAIT(duration) RETURN 1", 1),
		CaseNil("RETURN (FALSE OR T::FAIL())?"),
		CaseNil("RETURN T::FAIL()?"),
		CaseArray(`FOR i IN [1, 2, 3, 4]
				LET duration = 10
		
				WAIT(duration)
		
				RETURN i * 2`,
			[]any{2, 4, 6, 8}),

		Case(`RETURN FIRST((FOR i IN 1..10 RETURN i * 2))`, 2),
		CaseArray(`RETURN UNION((FOR i IN 0..5 RETURN i), (FOR i IN 6..10 RETURN i))`, []any{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10}),
	})
}

func TestBuiltinFunctions(t *testing.T) {
	RunUseCases(t, []UseCase{
		Case("RETURN LENGTH([1,2,3])", 3),
		Case("RETURN TYPENAME([1,2,3])", "array"),
	})
}

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
		CaseRuntimeError(`LET obj = NONE RETURN obj.foo`),
		CaseNil(`LET obj = NONE RETURN obj?.foo`),
		CaseObject(`RETURN {first: {second: "third"}}.first`,
			map[string]any{
				"second": "third",
			}),
		CaseObject(`RETURN KEEP_KEYS({first: {second: "third"}}.first, "second")`,
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
		Case(
			`
					LET obj = { foo: { bar: "bar" } }
		
					RETURN obj.foo?.bar
				`,
			"bar",
		),
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
		CaseNil("RETURN ERROR()?.foo"),
		CaseArray(`LET res = (FOR i IN ERROR() RETURN i)? RETURN res`, []any{}),

		CaseArray(`LET res = (FOR i IN [1, 2, 3, 4] LET y = ERROR() RETURN y+i)? RETURN res`, []any{}, "Error in arrayList comprehension"),
		CaseArray(`FOR i IN [1, 2, 3, 4] ERROR()? RETURN i`, []any{1, 2, 3, 4}, "Error in FOR loop"),
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
		CaseCompilationError(`
			FOR foo IN foo
				RETURN foo
		`, "Should not compile FOR foo IN foo"),
		CaseArray("FOR i IN 1..5 RETURN i", []any{1, 2, 3, 4, 5}),
		CaseArray(
			`
		FOR i IN 1..5
			LET x = i * 2
			RETURN x
		`,
			[]any{2, 4, 6, 8, 10},
		),
		CaseArray(
			`
		FOR val, counter IN 1..5
			LET x = val
			PRINT(counter)
			LET y = counter
			RETURN [x, y]
				`,
			[]any{[]any{1, 0}, []any{2, 1}, []any{3, 2}, []any{4, 3}, []any{5, 4}},
		),
		CaseArray(
			`FOR i IN [] RETURN i
				`,
			[]any{},
		),
		CaseArray(
			`FOR i IN [1, 2, 3] RETURN i
				`,
			[]any{1, 2, 3},
		),
		CaseArray(
			`FOR i, k IN [1, 2, 3] RETURN k`,
			[]any{0, 1, 2},
		),
		CaseArray(
			`FOR i IN ['foo', 'bar', 'qaz'] RETURN i`,
			[]any{"foo", "bar", "qaz"},
		),
		CaseItems(
			`FOR i IN {a: 'bar', b: 'foo', c: 'qaz'} RETURN i`,
			[]any{"bar", "foo", "qaz"},
		),
		CaseArray(
			`FOR i, k IN {a: 'foo', b: 'bar', c: 'qaz'} RETURN k`,
			[]any{"a", "b", "c"},
		),
		CaseArray(
			`FOR i IN [{name: 'foo'}, {name: 'bar'}, {name: 'qaz'}] RETURN i.name`,
			[]any{"foo", "bar", "qaz"},
		),
		CaseArray(
			`FOR i IN { items: [{name: 'foo'}, {name: 'bar'}, {name: 'qaz'}] }.items RETURN i.name`,
			[]any{"foo", "bar", "qaz"},
		),
		CaseArray(
			`FOR prop IN ["a"]
							FOR val IN [1, 2, 3]
								RETURN {[prop]: val}`,
			[]any{map[string]any{"a": 1}, map[string]any{"a": 2}, map[string]any{"a": 3}},
		),
		CaseArray(
			`FOR val IN 1..3
							FOR prop IN ["a"]
								RETURN {[prop]: val}`,
			[]any{map[string]any{"a": 1}, map[string]any{"a": 2}, map[string]any{"a": 3}},
		),
		CaseArray(
			`FOR prop IN ["a"]
							FOR val IN 1..3
								RETURN {[prop]: val}`,
			[]any{map[string]any{"a": 1}, map[string]any{"a": 2}, map[string]any{"a": 3}},
		),
		CaseArray(
			`FOR prop IN ["a"]
							FOR val IN [1, 2, 3]
								FOR val2 IN [1, 2, 3]
									RETURN { [prop]: [val, val2] }`,
			[]any{map[string]any{"a": []int{1, 1}}, map[string]any{"a": []int{1, 2}}, map[string]any{"a": []int{1, 3}}, map[string]any{"a": []int{2, 1}}, map[string]any{"a": []int{2, 2}}, map[string]any{"a": []int{2, 3}}, map[string]any{"a": []int{3, 1}}, map[string]any{"a": []int{3, 2}}, map[string]any{"a": []int{3, 3}}},
		),
		CaseArray(
			`FOR val IN [1, 2, 3]
							RETURN (
								FOR prop IN ["a", "b", "c"]
									RETURN { [prop]: val }
							)`,
			[]any{[]any{map[string]any{"a": 1}, map[string]any{"b": 1}, map[string]any{"c": 1}}, []any{map[string]any{"a": 2}, map[string]any{"b": 2}, map[string]any{"c": 2}}, []any{map[string]any{"a": 3}, map[string]any{"b": 3}, map[string]any{"c": 3}}},
		),
		CaseArray(
			`FOR val IN [1, 2, 3]
							LET sub = (
								FOR prop IN ["a", "b", "c"]
									RETURN { [prop]: val }
							)
		
							RETURN sub`,
			[]any{[]any{map[string]any{"a": 1}, map[string]any{"b": 1}, map[string]any{"c": 1}}, []any{map[string]any{"a": 2}, map[string]any{"b": 2}, map[string]any{"c": 2}}, []any{map[string]any{"a": 3}, map[string]any{"b": 3}, map[string]any{"c": 3}}},
		),
		CaseArray(
			`FOR i IN [ 1, 2, 3, 4, 1, 3 ]
							RETURN DISTINCT i
		`,
			[]any{1, 2, 3, 4},
		),
	})
}

func TestForTernaryExpression(t *testing.T) {
	RunUseCases(t, []UseCase{
		CaseArray(`
			LET foo = FALSE
			RETURN foo ? TRUE : (FOR i IN 1..5 RETURN i*2)`,
			[]any{2, 4, 6, 8, 10}),
		CaseArray(`
			LET foo = FALSE
			RETURN foo ? TRUE : (FOR i IN 1..5 T::FAIL() RETURN i*2)?`,
			[]any{}),
		CaseArray(`
			LET foo = FALSE
			RETURN foo ? (FOR i IN 1..5 RETURN i) : (FOR i IN 1..5 RETURN i*2)`,
			[]any{2, 4, 6, 8, 10}),
		CaseArray(`
			LET foo = FALSE
			RETURN foo ? (FOR i IN 1..5 RETURN T::FAIL()) : (FOR i IN 1..5 RETURN T::FAIL())?`,
			[]any{}),
		CaseArray(`
			LET foo = TRUE
			RETURN foo ? (FOR i IN 1..5 RETURN T::FAIL())? : (FOR i IN 1..5 RETURN T::FAIL())`,
			[]any{}),
		CaseArray(`
			LET foo = FALSE
			LET res = foo ? TRUE : (FOR i IN 1..5 RETURN i*2) 
			RETURN res`,
			[]any{2, 4, 6, 8, 10}),
		Case(`
			LET foo = TRUE
			LET res = foo ? TRUE : (FOR i IN 1..5 RETURN i*2)
			RETURN res`,
			true),
		CaseArray(`
			LET foo = FALSE
			LET res = foo ? TRUE : (FOR i IN 1..5 RETURN i*2) 
			RETURN res`,
			[]any{2, 4, 6, 8, 10}),
		CaseArray(`
			LET foo = FALSE
			LET res = foo ? (FOR i IN 1..5 RETURN i) : (FOR i IN 1..5 RETURN i*2)
			RETURN res`,
			[]any{2, 4, 6, 8, 10}),
		CaseArray(`
			LET foo = TRUE
			LET res = foo ? (FOR i IN 1..5 RETURN i) : (FOR i IN 1..5 RETURN i*2)
			RETURN res`,
			[]any{1, 2, 3, 4, 5}),
		Case(`
			LET res = LENGTH((FOR i IN 1..5 RETURN T::FAIL())?) ? TRUE : FALSE
			RETURN res`,
			false),
		Case(`
			LET res = (FOR i IN 1..5 RETURN i)? ? TRUE : FALSE
			RETURN res
`,
			true),
	})
}

func TestForWhile(t *testing.T) {
	var untilCounter int
	counter := -1
	RunUseCases(t, []UseCase{
		CaseArray("FOR i WHILE false RETURN i", []any{}),
		CaseArray("FOR i WHILE UNTIL(5) RETURN i", []any{0, 1, 2, 3, 4}),
		CaseArray(`
			FOR i WHILE COUNTER() < 5
				LET y = i + 1
				FOR x IN 1..y
					RETURN i * x
		`, []any{0, 1, 2, 2, 4, 6, 3, 6, 9, 12, 4, 8, 12, 16, 20}),
	}, runtime.WithFunctions(map[string]core.Function{
		"UNTIL": func(ctx context.Context, args ...core.Value) (core.Value, error) {
			if untilCounter < int(core.ToIntSafe(ctx, args[0])) {
				untilCounter++

				return core.True, nil
			}

			return core.False, nil
		},
		"COUNTER": func(ctx context.Context, args ...core.Value) (core.Value, error) {
			counter++
			return core.NewInt(counter), nil
		},
	}))
}

func TestForTernaryWhileExpression(t *testing.T) {
	counter := -1
	RunUseCases(t, []UseCase{
		CaseArray(`
			LET foo = FALSE
			RETURN foo ? TRUE : (FOR i WHILE false RETURN i*2)
		`, []any{}),
		CaseArray(`
			LET foo = FALSE
			RETURN foo ? TRUE : (FOR i WHILE T::FAIL() RETURN i*2)?
		`, []any{}),
		CaseArray(`
			LET foo = FALSE
			RETURN foo ? TRUE : (FOR i WHILE COUNTER() < 10 RETURN i*2)`,
			[]any{0, 2, 4, 6, 8, 10, 12, 14, 16, 18}),
	}, runtime.WithFunctions(map[string]core.Function{
		"COUNTER": func(ctx context.Context, args ...core.Value) (core.Value, error) {
			counter++
			return core.NewInt(counter), nil
		},
	}))
}

func TestForDoWhile(t *testing.T) {
	counter := -1
	counter2 := -1

	RunUseCases(t, []UseCase{
		CaseArray(`
			FOR i DO WHILE false
				RETURN i
		`, []any{0}),
		CaseArray(`
		FOR i DO WHILE COUNTER() < 10
				RETURN i`, []any{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10}),
		CaseArray(`
			FOR i WHILE COUNTER2() < 5
				LET y = i + 1
				FOR x IN 1..y
					RETURN i * x
		`, []any{0, 0, 0, 0, 0, 0, 1, 2, 3, 4, 0, 2, 4, 6, 8, 0, 3, 6, 9, 12, 0, 4, 8, 12, 16}),
	}, runtime.WithFunctions(map[string]core.Function{
		"COUNTER": func(ctx context.Context, args ...core.Value) (core.Value, error) {
			counter++
			return core.NewInt(counter), nil
		},
		"COUNTER2": func(ctx context.Context, args ...core.Value) (core.Value, error) {
			counter2++
			return core.NewInt(counter), nil
		},
	}))
}

func TestForFilter(t *testing.T) {
	counterA := 0
	counterB := 0
	RunUseCases(t, []UseCase{
		CaseArray(
			`
			FOR i IN [ 1, 2, 3, 4, 1, 3 ]
				FILTER i > 2
				RETURN i
		`,
			[]any{3, 4, 3},
		),
		CaseArray(
			`
			FOR i IN [ 1, 2, 3, 4, 1, 3 ]
				FILTER i > 1 AND i < 4
				RETURN i
		`,
			[]any{2, 3, 3},
		),
		CaseArray(
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
		),
		CaseArray(
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
		),
		CaseArray(
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
			"Should compile query with left side expression",
		),
		CaseArray(
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
			"Should compile query with multiple FILTER statements",
		),
		CaseArray(`
			LET users = [
				{
					active: true,
					married: true,
					age: 31,
					gender: "m"
				},
				{
					active: true,
					married: false,
					age: 25,
					gender: "f"
				},
				{
					active: true,
					married: false,
					age: 36,
					gender: "m"
				},
				{
					active: false,
					married: true,
					age: 69,
					gender: "m"
				},
				{
					active: true,
					married: true,
					age: 45,
					gender: "f"
				}
			]
			FOR u IN users
				FILTER u.active AND u.married
				RETURN u
`, []any{map[string]any{"active": true, "age": 31, "gender": "m", "married": true}, map[string]any{"active": true, "age": 45, "gender": "f", "married": true}},
			"Should compile query with multiple left side expression"),
		CaseArray(`
LET users = [
				{
					active: true,
					married: true,
					age: 31,
					gender: "m"
				},
				{
					active: true,
					married: false,
					age: 25,
					gender: "f"
				},
				{
					active: true,
					married: false,
					age: 36,
					gender: "m"
				},
				{
					active: false,
					married: true,
					age: 69,
					gender: "m"
				},
				{
					active: true,
					married: true,
					age: 45,
					gender: "f"
				}
			]
			FOR u IN users
				FILTER !u.active AND u.married
				RETURN u
`, []any{map[string]any{"active": false, "age": 69, "gender": "m", "married": true}},
			"Should compile query with multiple left side expression and with binary operator"),
		CaseArray(`
		LET users = [
				{
					active: true,
					married: true,
					age: 31,
					gender: "m"
				},
				{
					active: true,
					married: false,
					age: 25,
					gender: "f"
				},
				{
					active: true,
					married: false,
					age: 36,
					gender: "m"
				},
				{
					active: false,
					married: true,
					age: 69,
					gender: "m"
				},
				{
					active: true,
					married: true,
					age: 45,
					gender: "f"
				}
			]
			FOR u IN users
				FILTER !u.active AND !u.married
				RETURN u
`, []any{},
			"Should compile query with multiple left side expression and with binary operator 2"),
		CaseArray(`
			FOR i IN [ 1, 2, 3, 4, 1, 3 ]
				LET x = 2
				FILTER i > x
				RETURN i + x
`, []any{5, 6, 5}),
		CaseArray(`
			FOR i IN [ 1, 2, 3, 4, 1, 3 ]
				LET x = 2
				COUNT_A()
				FILTER i > x
				COUNT_B()
				RETURN i + x
`, []any{5, 6, 5}),
	}, runtime.WithFunctions(map[string]core.Function{
		"COUNT_A": func(ctx context.Context, args ...core.Value) (core.Value, error) {
			counterA++

			return core.None, nil
		},
		"COUNT_B": func(ctx context.Context, args ...core.Value) (core.Value, error) {
			counterB++

			return core.None, nil
		},
	}))
}

func TestForLimit(t *testing.T) {
	RunUseCases(t, []UseCase{
		CaseArray(
			`
			FOR i IN [ 1, 2, 3, 4, 1, 3 ]
				LIMIT 2
				RETURN i
		`,
			[]any{1, 2}),
		CaseArray(`
			FOR i IN [ 1,2,3,4,5,6,7,8 ]
				LIMIT 4, 2
				RETURN i
			`, []any{5, 6}),
		CaseArray(`
			FOR i IN [ 1,2,3,4,5,6,7,8 ]
				LET x = i
				LIMIT 2
				RETURN i*x
			`, []any{1, 4},
			"Should be able to reuse values from a source"),
		CaseArray(`
			FOR i IN [ 1,2,3,4,5,6,7,8 ]
				LET x = "foo"
				TYPENAME(x)
				LIMIT 2
				RETURN i
		`, []any{1, 2}, "Should define variables and call functions"),
		CaseArray(`
			FOR i IN [ 1,2,3,4,5,6,7,8 ]
				LIMIT LIMIT_VALUE()
				RETURN i
		`, []any{1, 2}, "Should be able to use function call"),
		CaseArray(`
			LET o = {
				limit: 2
			}
			FOR i IN [ 1,2,3,4,5,6,7,8 ]
				LIMIT o.limit
				RETURN i
		`, []any{1, 2}, "Should be able to use object property"),
		CaseArray(`
			LET o = [1,2]

			FOR i IN [ 1,2,3,4,5,6,7,8 ]
				LIMIT o[1]
				RETURN i
		`, []any{1, 2}, "Should be able to use array element"),
	}, runtime.WithFunction("LIMIT_VALUE", func(ctx context.Context, args ...core.Value) (core.Value, error) {
		return core.NewInt(2), nil
	}))
}

func TestForSort(t *testing.T) {
	RunUseCases(t, []UseCase{
		CaseArray(`
LET strs = ["foo", "bar", "qaz", "abc"]

FOR s IN strs
	SORT s
	RETURN s
`, []any{"abc", "bar", "foo", "qaz"}, "Should sort strings"),
		CaseArray(`
LET users = [
				{
					name: "Ron",
					age: 31,
					gender: "m"
				},
				{
					name: "Angela",
					age: 29,
					gender: "f"
				},
				{
					name: "Bob",
					age: 36,
					gender: "m"
				}
			]
			FOR u IN users
				SORT u.name
				RETURN u
`, []any{
			map[string]any{"name": "Angela", "age": 29, "gender": "f"},
			map[string]any{"name": "Bob", "age": 36, "gender": "m"},
			map[string]any{"name": "Ron", "age": 31, "gender": "m"},
		}, "Should sort objects by name (string)"),
		CaseArray(`
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
				SORT u.age
				RETURN u
`, []any{
			map[string]any{"active": true, "age": 29, "gender": "f"},
			map[string]any{"active": true, "age": 31, "gender": "m"},
			map[string]any{"active": true, "age": 36, "gender": "m"},
		}, "Should sort objects by age (int)"),
		CaseArray(`
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
				SORT u.age DESC
				RETURN u
		`, []any{
			map[string]any{"active": true, "age": 36, "gender": "m"},
			map[string]any{"active": true, "age": 31, "gender": "m"},
			map[string]any{"active": true, "age": 29, "gender": "f"},
		}, "Should execute query with DESC SORT statement"),
		CaseArray(`
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
				SORT u.age ASC
				RETURN u
		`, []any{
			map[string]any{"active": true, "age": 29, "gender": "f"},
			map[string]any{"active": true, "age": 31, "gender": "m"},
			map[string]any{"active": true, "age": 36, "gender": "m"},
		}, "Should compile query with ASC SORT statement"),
		CaseArray(`			LET users = [
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
					age: 31,
					gender: "f"
				},
				{
					active: true,
					age: 36,
					gender: "m"
				}
			]
			FOR u IN users
				SORT u.age, u.gender
				RETURN u`,
			[]any{
				map[string]any{"active": true, "age": 29, "gender": "f"},
				map[string]any{"active": true, "age": 31, "gender": "f"},
				map[string]any{"active": true, "age": 31, "gender": "m"},
				map[string]any{"active": true, "age": 36, "gender": "m"},
			}, "Should compile query with SORT statement with multiple expressions"),
		CaseArray(`
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
					age: 31,
					gender: "f"
				},
				{
					active: true,
					age: 36,
					gender: "m"
				}
			]
			FOR u IN users
				LET x = "foo"
				TEST(x)
				SORT u.age, u.gender
				RETURN u
		`, []any{
			map[string]any{"active": true, "age": 29, "gender": "f"},
			map[string]any{"active": true, "age": 31, "gender": "f"},
			map[string]any{"active": true, "age": 31, "gender": "m"},
			map[string]any{"active": true, "age": 36, "gender": "m"},
		}, "Should define variables and call functions"),
		CaseArray(`
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
				FILTER u.gender == "m"
				SORT u.age
				RETURN u
		`, []any{
			map[string]any{"active": true, "age": 31, "gender": "m"},
			map[string]any{"active": true, "age": 36, "gender": "m"},
		}, "Should compile query with FILTER and SORT statements"),
	}, runtime.WithFunction("TEST", func(ctx context.Context, args ...core.Value) (core.Value, error) {
		return core.None, nil
	}))
}

func TestCollect(t *testing.T) {
	RunUseCases(t, []UseCase{
		CaseCompilationError(`
			LET users = [
				{
					active: true,
					married: true,
					age: 31,
					gender: "m"
				},
				{
					active: true,
					married: false,
					age: 25,
					gender: "f"
				},
				{
					active: true,
					married: false,
					age: 36,
					gender: "m"
				},
				{
					active: false,
					married: true,
					age: 69,
					gender: "m"
				},
				{
					active: true,
					married: true,
					age: 45,
					gender: "f"
				}
			]
			FOR i IN users
				COLLECT gender = i.gender
				RETURN {
					user: i,
					gender: gender
				}
		`, "Should not have access to initial variables"),
		CaseCompilationError(`
			LET users = [
				{
					active: true,
					married: true,
					age: 31,
					gender: "m"
				},
				{
					active: true,
					married: false,
					age: 25,
					gender: "f"
				},
				{
					active: true,
					married: false,
					age: 36,
					gender: "m"
				},
				{
					active: false,
					married: true,
					age: 69,
					gender: "m"
				},
				{
					active: true,
					married: true,
					age: 45,
					gender: "f"
				}
			]
			FOR i IN users
				LET x = "foo"
				COLLECT gender = i.gender
				RETURN {x, gender}
		`, "Should not have access to variables defined before COLLECT"),
		CaseArray(`LET users = [
				{
					active: true,
					married: true,
					age: 31,
					gender: "m"
				},
				{
					active: true,
					married: false,
					age: 25,
					gender: "f"
				},
				{
					active: true,
					married: false,
					age: 36,
					gender: "m"
				},
				{
					active: false,
					married: true,
					age: 69,
					gender: "m"
				},
				{
					active: true,
					married: true,
					age: 45,
					gender: "f"
				}
			]
			g`, []any{"f", "m"}, "Should group result by a single key"),
	})
}
