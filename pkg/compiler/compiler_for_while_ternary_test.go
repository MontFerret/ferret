package compiler_test

import (
	"context"
	"github.com/MontFerret/ferret/pkg/compiler"
	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestForTernaryWhileExpression(t *testing.T) {
	Convey("RETURN foo ? TRUE : (FOR i WHILE false RETURN i*2)", t, func() {
		c := compiler.New()

		out1, err := c.MustCompile(`
			LET foo = FALSE
			RETURN foo ? TRUE : (FOR i WHILE false RETURN i*2)
		`).Run(context.Background())

		So(err, ShouldBeNil)
		So(string(out1), ShouldEqual, `[]`)
	})

	Convey("RETURN foo ? TRUE : (FOR i WHILE T::FAIL() RETURN i*2)?", t, func() {
		c := compiler.New()

		out1, err := c.MustCompile(`
			LET foo = FALSE
			RETURN foo ? TRUE : (FOR i WHILE T::FAIL() RETURN i*2)?
		`).Run(context.Background())

		So(err, ShouldBeNil)
		So(string(out1), ShouldEqual, `null`)
	})

	Convey("RETURN foo ? TRUE : (FOR i WHILE F() < 10 RETURN i*2)", t, func() {
		c := compiler.New()

		counter := -1
		c.MustRegisterFunction("F", func(ctx context.Context, args ...core.Value) (core.Value, error) {
			counter++
			return values.NewInt(counter), nil
		})

		out1, err := c.MustCompile(`
			LET foo = FALSE
			RETURN foo ? TRUE : (FOR i WHILE F() < 10 RETURN i*2)
		`).Run(context.Background())

		So(err, ShouldBeNil)
		So(string(out1), ShouldEqual, `[0,2,4,6,8,10,12,14,16,18]`)
	})
}
