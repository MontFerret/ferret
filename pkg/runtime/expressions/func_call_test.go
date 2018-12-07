package expressions_test

import (
	"context"
	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/expressions"
	"github.com/MontFerret/ferret/pkg/runtime/expressions/literals"
	"github.com/MontFerret/ferret/pkg/runtime/values"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestFunctionCallExpression(t *testing.T) {
	Convey(".Exec", t, func() {
		Convey("Should execute an underlying function without arguments", func() {
			f, err := expressions.NewFunctionCallExpression(
				core.SourceMap{},
				func(ctx context.Context, args ...core.Value) (value core.Value, e error) {
					So(args, ShouldHaveLength, 0)

					return values.True, nil
				},
			)

			So(err, ShouldBeNil)

			rootScope, _ := core.NewRootScope()

			out, err := f.Exec(context.Background(), rootScope.Fork())

			So(err, ShouldBeNil)
			So(out, ShouldEqual, values.True)
		})

		Convey("Should execute an underlying function with arguments", func() {
			args := []core.Expression{
				literals.NewIntLiteral(1),
				literals.NewStringLiteral("foo"),
			}

			f, err := expressions.NewFunctionCallExpression(
				core.SourceMap{},
				func(ctx context.Context, args ...core.Value) (value core.Value, e error) {
					So(args, ShouldHaveLength, len(args))

					return values.True, nil
				},
				args...,
			)

			So(err, ShouldBeNil)

			rootScope, _ := core.NewRootScope()

			out, err := f.Exec(context.Background(), rootScope.Fork())

			So(err, ShouldBeNil)
			So(out, ShouldEqual, values.True)
		})

		Convey("Should stop an execution when context is cancelled", func() {
			args := []core.Expression{
				literals.NewIntLiteral(1),
				literals.NewStringLiteral("foo"),
			}

			f, err := expressions.NewFunctionCallExpression(
				core.SourceMap{},
				func(ctx context.Context, args ...core.Value) (value core.Value, e error) {
					So(args, ShouldHaveLength, len(args))

					return values.True, nil
				},
				args...,
			)

			So(err, ShouldBeNil)

			rootScope, _ := core.NewRootScope()
			ctx, cancel := context.WithCancel(context.Background())
			cancel()

			_, err = f.Exec(ctx, rootScope.Fork())

			So(err, ShouldEqual, core.ErrTerminated)
		})
	})
}
