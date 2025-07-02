package vm_test

import (
	"fmt"
	"testing"

	"github.com/MontFerret/ferret/pkg/compiler"
	"github.com/MontFerret/ferret/pkg/runtime"
	"github.com/MontFerret/ferret/pkg/vm"
	. "github.com/MontFerret/ferret/test/integration/base"

	gocontext "context"

	. "github.com/smartystreets/goconvey/convey"
)

func TestVariables(t *testing.T) {
	RunUseCases(t, []UseCase{
		SkipCaseCompilationError(`RETURN foo`, "Should not compile if a variable not defined"),
		SkipCaseCompilationError(`
			LET foo = "bar"
			LET foo = "baz"

			RETURN foo
		`, "Should not compile if a variable is not unique"),
		SkipCaseCompilationError(`			LET _ = (FOR i IN 1..100 RETURN NONE)
	
			RETURN _`, "Should not allow to use ignorable variable name"),
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
		Case(`
			LET _ = (FOR i IN 1..100 RETURN NONE)
			LET _ = (FOR i IN 1..100 RETURN NONE)

			RETURN TRUE
		`, true),
	})

	SkipConvey("Should compile LET i = (FOR i WHILE COUNTER() < 5 RETURN i) RETURN i", t, func() {
		c := compiler.New()

		p, err := c.Compile(`
			LET i = (FOR i WHILE COUNTER() < 5 RETURN i)
			RETURN i
		`)

		So(err, ShouldBeNil)
		So(p, ShouldHaveSameTypeAs, &vm.Program{})

		counter := -1
		out, err := Run(p, vm.WithFunction("COUNTER", func(ctx gocontext.Context, args ...runtime.Value) (runtime.Value, error) {
			counter++

			return runtime.NewInt(counter), nil
		}))

		So(err, ShouldBeNil)
		So(string(out), ShouldEqual, "[0,1,2,3,4]")
	})

	SkipConvey("Should compile LET i = (FOR i WHILE COUNTER() < 5 T::FAIL() RETURN i)? RETURN length(i) == 0", t, func() {
		c := compiler.New()

		p, err := c.Compile(`
			LET i = (FOR i WHILE COUNTER() < 5 T::FAIL() RETURN i)?
			RETURN length(i) == 0
		`)

		So(err, ShouldBeNil)
		So(p, ShouldHaveSameTypeAs, &vm.Program{})

		counter := -1
		out, err := Run(p, vm.WithFunction("COUNTER", func(ctx gocontext.Context, args ...runtime.Value) (runtime.Value, error) {
			counter++

			return runtime.NewInt(counter), nil
		}), vm.WithFunction("T::FAIL", func(ctx gocontext.Context, args ...runtime.Value) (runtime.Value, error) {
			return runtime.None, fmt.Errorf("test")
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
