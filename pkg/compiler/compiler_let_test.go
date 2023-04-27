package compiler_test

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/MontFerret/ferret/pkg/compiler"
)

func TestLet(t *testing.T) {
	RunUseCases(t, compiler.New(), []UseCase{
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
	})

	//
	//Convey("Should compile LET i = (FOR i IN [1,2,3] RETURN i) RETURN i", t, func() {
	//	c := compiler.New()
	//
	//	p, err := c.Compile(`
	//		LET i = (FOR i IN [1,2,3] RETURN i)
	//		RETURN i
	//	`)
	//
	//	So(err, ShouldBeNil)
	//	So(p, ShouldHaveSameTypeAs, &runtime.Program{})
	//
	//	out, err := p.Run(context.Background())
	//
	//	So(err, ShouldBeNil)
	//	So(string(out), ShouldEqual, "[1,2,3]")
	//})
	//
	//Convey("Should compile LET src = NONE LET i = (FOR i IN NONE RETURN i)? RETURN i == NONE", t, func() {
	//	c := compiler.New()
	//
	//	p, err := c.Compile(`
	//		LET src = NONE
	//		LET i = (FOR i IN src RETURN i)?
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
	//
	//Convey("Should compile LET i = (FOR i WHILE COUNTER() < 5 RETURN i) RETURN i", t, func() {
	//	c := compiler.New()
	//	counter := -1
	//	c.RegisterFunction("COUNTER", func(ctx context.Context, args ...core.Value) (core.Value, error) {
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
	//	c.RegisterFunction("COUNTER", func(ctx context.Context, args ...core.Value) (core.Value, error) {
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
	//
	//Convey("Should compile LET i = { items: [1,2,3]}  FOR el IN i.items RETURN i", t, func() {
	//	c := compiler.New()
	//
	//	p, err := c.Compile(`
	//		LET obj = { items: [1,2,3] }
	//
	//		FOR i IN obj.items
	//			RETURN i
	//	`)
	//
	//	So(err, ShouldBeNil)
	//	So(p, ShouldHaveSameTypeAs, &runtime.Program{})
	//
	//	out, err := p.Run(context.Background())
	//
	//	So(err, ShouldBeNil)
	//	So(string(out), ShouldEqual, "[1,2,3]")
	//})
	//
	//Convey("Should not compile FOR foo IN foo", t, func() {
	//	c := compiler.New()
	//
	//	_, err := c.Compile(`
	//		FOR foo IN foo
	//			RETURN foo
	//	`)
	//
	//	So(err, ShouldNotBeNil)
	//})
	//

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
	//Convey("Should use ignorable variable name", t, func() {
	//	out, err := newCompilerWithObservable().MustCompile(`
	//		LET _ = (FOR i IN 1..100 RETURN NONE)
	//
	//		RETURN TRUE
	//	`).Run(context.Background())
	//
	//	So(err, ShouldBeNil)
	//	So(string(out), ShouldEqual, `true`)
	//})
	//
	//Convey("Should allow to declare a variable name using _", t, func() {
	//	c := compiler.New()
	//
	//	out, err := c.MustCompile(`
	//		LET _ = (FOR i IN 1..100 RETURN NONE)
	//		LET _ = (FOR i IN 1..100 RETURN NONE)
	//
	//		RETURN TRUE
	//	`).Run(context.Background())
	//
	//	So(err, ShouldBeNil)
	//	So(string(out), ShouldEqual, `true`)
	//})

	Convey("Should not allow to use ignorable variable name", t, func() {
		c := compiler.New()

		_, err := c.Compile(`
			LET _ = (FOR i IN 1..100 RETURN NONE)
	
			RETURN _
		`)

		So(err, ShouldNotBeNil)
	})
}
