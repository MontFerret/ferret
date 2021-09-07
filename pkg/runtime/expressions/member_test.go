package expressions_test

import (
	"context"
	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/expressions"
	"github.com/MontFerret/ferret/pkg/runtime/values"
	"github.com/stretchr/testify/mock"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

type TestObject struct {
	*values.Object
	*mock.Mock
}

func NewTestObject() *TestObject {
	o := new(TestObject)

	return o
}

func (to *TestObject) GetIn(ctx context.Context, path []core.Value) (core.Value, error) {
	to.Mock.Called(path)

	return to.Object.GetIn(ctx, path)
}

func TestMemberExpression(t *testing.T) {

	Convey(".Exec", t, func() {
		Convey("Should use .Getter interface if a source implements it", func() {
			o := NewTestObject()

			args := []core.Value{
				values.NewString("foo"),
				values.NewString("bar"),
				values.NewString("baz"),
			}

			s1, _ := expressions.NewMemberPathSegment(
				core.NewExpressionFn(func(ctx context.Context, scope *core.Scope) (core.Value, error) {
					return args[0], nil
				}),
				false,
			)

			s2, _ := expressions.NewMemberPathSegment(
				core.NewExpressionFn(func(ctx context.Context, scope *core.Scope) (core.Value, error) {
					return args[1], nil
				}),
				false,
			)

			s3, _ := expressions.NewMemberPathSegment(
				core.NewExpressionFn(func(ctx context.Context, scope *core.Scope) (core.Value, error) {
					return args[2], nil
				}),
				false,
			)

			segments := []*expressions.MemberPathSegment{s1, s2, s3}

			o.On("GetIn", args)

			exp, err := expressions.NewMemberExpression(
				core.SourceMap{},
				core.NewExpressionFn(func(ctx context.Context, scope *core.Scope) (core.Value, error) {
					return o, nil
				}),
				[]*expressions.MemberPathSegment{s1, s2, s3},
			)

			So(err, ShouldBeNil)

			root, cancel := core.NewRootScope()

			defer cancel()

			exp.Exec(context.Background())

			o.AssertExpectations(t)
		})
	})
}
