package compiler_test

import (
	"github.com/MontFerret/ferret/pkg/compiler"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestVariables(t *testing.T) {
	RunUseCases(t, []UseCase{
		//{
		//	`LET i = NONE RETURN i`,
		//	nil,
		//	nil,
		//},
		//{
		//	`LET a = TRUE RETURN a`,
		//	true,
		//	nil,
		//},
		//{
		//	`LET a = 1 RETURN a`,
		//	1,
		//	nil,
		//},
		//{
		//	`LET a = 1.1 RETURN a`,
		//	1.1,
		//	nil,
		//},
		//{
		//	`LET i = 'foo' RETURN i`,
		//	"foo",
		//	nil,
		//},
		//{
		//	`LET i = [] RETURN i`,
		//	[]any{},
		//	ShouldEqualJSON,
		//},
		//{
		//	`LET i = [1, 2, 3] RETURN i`,
		//	[]any{1, 2, 3},
		//	ShouldEqualJSON,
		//},
		//{
		//	`LET i = [None, FALSE, "foo", 1, 1.1] RETURN i`,
		//	[]any{nil, false, "foo", 1, 1.1},
		//	ShouldEqualJSON,
		//},
		//{
		//	`LET i = {} RETURN i`,
		//	map[string]any{},
		//	ShouldEqualJSON,
		//},
		//{
		//	`LET i = {a: 1, b: 2} RETURN i`,
		//	map[string]any{"a": 1, "b": 2},
		//	ShouldEqualJSON,
		//},
		//{
		//	`LET i = {a: 1, b: [1]} RETURN i`,
		//	map[string]any{"a": 1, "b": []any{1}},
		//	ShouldEqualJSON,
		//},
		//{
		//	`LET i = {a: {c: 1}, b: [1]} RETURN i`,
		//	map[string]any{"a": map[string]any{"c": 1}, "b": []any{1}},
		//	ShouldEqualJSON,
		//},
		//{
		//	`LET i = {a: 'foo', b: 1, c: TRUE, d: [], e: {}} RETURN i`,
		//	map[string]any{"a": "foo", "b": 1, "c": true, "d": []any{}, "e": map[string]any{}},
		//	ShouldEqualJSON,
		//},
		//{
		//	`LET prop = "name" LET i = { [prop]: "foo" } RETURN i`,
		//	map[string]any{"name": "foo"},
		//	ShouldEqualJSON,
		//},
		//{
		//	`LET name="foo" LET i = { name } RETURN i`,
		//	map[string]any{"name": "foo"},
		//	ShouldEqualJSON,
		//},
		//{
		//	`LET i = [{a: {c: 1}, b: [1]}] RETURN i`,
		//	[]any{map[string]any{"a": map[string]any{"c": 1}, "b": []any{1}}},
		//	ShouldEqualJSON,
		//},
		//{
		//	"LET a = 'a' LET b = a LET c = 'c' RETURN b",
		//	"a",
		//	ShouldEqual,
		//},
		//{
		//	"LET i = (FOR i IN [1,2,3] RETURN i) RETURN i",
		//	[]int{1, 2, 3},
		//	ShouldEqualJSON,
		//},
		//{
		//	" LET i = { items: [1,2,3]}  FOR el IN i.items RETURN el",
		//	[]int{1, 2, 3},
		//	ShouldEqualJSON,
		//},
		{
			`LET _ = (FOR i IN 1..100 RETURN NONE)
				RETURN TRUE`,
			true,
			ShouldEqualJSON,
		},
	})

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
	//	c.RegisterFunction("COUNTER", func(ctx context.visitor, args ...core.Value) (core.Value, error) {
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
	//	c.RegisterFunction("COUNTER", func(ctx context.visitor, args ...core.Value) (core.Value, error) {
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

	//Convey("RETURN { enabled: !val}", t, func() {
	//	c := compiler.New()
	//
	//	out1, err := c.MustCompile(`
	//		LET val = ""
	//		RETURN { enabled: !val }
	//	`).Run(context.Background())
	//
	//	So(err, ShouldBeNil)
	//	So(string(out1), ShouldEqual, `{"enabled":true}`)
	//
	//	out2, err := c.MustCompile(`
	//		LET val = ""
	//		RETURN { enabled: !!val }
	//	`).Run(context.Background())
	//
	//	So(err, ShouldBeNil)
	//	So(string(out2), ShouldEqual, `{"enabled":false}`)
	//})
	//
}
