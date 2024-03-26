package compiler_test

import (
	"testing"
)

func TestLogicalOperators(t *testing.T) {
	RunUseCases(t, []UseCase{
		{"RETURN 1 AND 0", 0, nil},
		{"RETURN 1 AND 1", 1, nil},
		{"RETURN 2 > 1 AND 1 > 0", true, nil},
		{"RETURN NONE && true", nil, nil},
		{"RETURN '' && true", "", nil},
		{"RETURN true && 23", 23, nil},
		{"RETURN 2 > 1 OR 1 < 0", true, nil},
		{"RETURN 1 || 7", 1, nil},
		{"RETURN 0 || 7", 7, nil},
		{"RETURN NONE || 'foo'", "foo", nil},
	})

	//
	//Convey("ERROR()? || 'boo'  should return 'boo'", t, func() {
	//	c := compiler.New()
	//	c.RegisterFunction("ERROR", func(ctx context.Context, args ...core.Value) (core.Value, error) {
	//		return nil, errors.New("test")
	//	})
	//
	//	p, err := c.Compile(`
	//		RETURN ERROR()? || 'boo'
	//	`)
	//
	//	So(err, ShouldBeNil)
	//
	//	out, err := p.Run(context.Background())
	//
	//	So(err, ShouldBeNil)
	//	So(string(out), ShouldEqual, `"boo"`)
	//})
	//
	//Convey("!ERROR()? && TRUE should return false", t, func() {
	//	c := compiler.New()
	//	c.RegisterFunction("ERROR", func(ctx context.Context, args ...core.Value) (core.Value, error) {
	//		return nil, errors.New("test")
	//	})
	//
	//	p, err := c.Compile(`
	//		RETURN !ERROR()? && TRUE
	//	`)
	//
	//	So(err, ShouldBeNil)
	//
	//	out, err := p.Run(context.Background())
	//
	//	So(err, ShouldBeNil)
	//	So(string(out), ShouldEqual, `true`)
	//})
	//
	//

	//
	//Convey("NOT TRUE should return false", t, func() {
	//	c := compiler.New()
	//
	//	p, err := c.Compile(`
	//		RETURN NOT TRUE
	//	`)
	//
	//	So(err, ShouldBeNil)
	//
	//	out, err := p.Run(context.Background())
	//
	//	So(err, ShouldBeNil)
	//	So(string(out), ShouldEqual, `false`)
	//})
	//
	//Convey("NOT u.valid should return true", t, func() {
	//	c := compiler.New()
	//
	//	p, err := c.Compile(`
	//		LET u = { valid: false }
	//
	//		RETURN NOT u.valid
	//	`)
	//
	//	So(err, ShouldBeNil)
	//
	//	out, err := p.Run(context.Background())
	//
	//	So(err, ShouldBeNil)
	//	So(string(out), ShouldEqual, `true`)
	//})
}
