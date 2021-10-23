package expressions_test

import (
	"context"
	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/expressions"
	"github.com/MontFerret/ferret/pkg/runtime/values"
	"github.com/MontFerret/ferret/pkg/runtime/values/types"
	"github.com/stretchr/testify/mock"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

type (
	TestObject struct {
		mock.Mock
		*values.Object
		failAt string
	}
)

func NewTestObject() *TestObject {
	o := new(TestObject)
	o.Object = values.NewObject()

	return o
}

func (to *TestObject) GetIn(ctx context.Context, path []core.Value) (core.Value, core.PathError) {
	to.Mock.Called(path)

	var current core.Value = to.Object

	for i, segment := range path {
		if segment.String() == to.failAt {
			return values.None, core.NewPathError(core.ErrInvalidPath, i)
		}

		next, err := values.GetIn(ctx, current, []core.Value{segment})

		if err != nil {
			return values.None, core.NewPathError(err, i)
		}

		current = next
	}

	return current, nil
}

func TestMemberExpression(t *testing.T) {
	Convey(".Exec", t, func() {
		Convey("Should use .Getter interface if a source implements it", func() {
			o := NewTestObject()
			o.Set("foo", values.NewObjectWith(
				values.NewObjectProperty("bar", values.NewObjectWith(
					values.NewObjectProperty("baz", values.NewObject()),
				)),
			))

			args := []core.Value{
				values.NewString("foo"),
				values.NewString("bar"),
				values.NewString("baz"),
			}

			o.On("GetIn", args)

			s1, _ := expressions.NewMemberPathSegment(
				core.AsExpression(func(ctx context.Context, scope *core.Scope) (core.Value, error) {
					return args[0], nil
				}),
				false,
			)

			s2, _ := expressions.NewMemberPathSegment(
				core.AsExpression(func(ctx context.Context, scope *core.Scope) (core.Value, error) {
					return args[1], nil
				}),
				false,
			)

			s3, _ := expressions.NewMemberPathSegment(
				core.AsExpression(func(ctx context.Context, scope *core.Scope) (core.Value, error) {
					return args[2], nil
				}),
				false,
			)

			segments := []*expressions.MemberPathSegment{s1, s2, s3}

			exp, err := expressions.NewMemberExpression(
				core.SourceMap{},
				core.AsExpression(func(ctx context.Context, scope *core.Scope) (core.Value, error) {
					return o, nil
				}),
				segments,
				nil,
			)

			So(err, ShouldBeNil)

			root, cancel := core.NewRootScope()

			defer func() {
				if err := cancel(); err != nil {
					panic(err)
				}
			}()

			out, err := exp.Exec(context.Background(), root.Fork())
			So(err, ShouldBeNil)
			So(out.Type().String(), ShouldNotEqual, types.None.String())

			o.AssertExpectations(t)
		})

		Convey("Should use generic traverse logic if a source does not implement Getter interface", func() {
			o := values.NewString("abcdefg")

			args := []core.Value{
				values.NewInt(0),
			}

			s1, _ := expressions.NewMemberPathSegment(
				core.AsExpression(func(ctx context.Context, scope *core.Scope) (core.Value, error) {
					return args[0], nil
				}),
				false,
			)

			segments := []*expressions.MemberPathSegment{s1}

			exp, err := expressions.NewMemberExpression(
				core.SourceMap{},
				core.AsExpression(func(ctx context.Context, scope *core.Scope) (core.Value, error) {
					return o, nil
				}),
				segments,
				nil,
			)

			So(err, ShouldBeNil)

			root, cancel := core.NewRootScope()

			defer func() {
				if err := cancel(); err != nil {
					panic(err)
				}
			}()

			out, err := exp.Exec(context.Background(), root.Fork())
			So(err, ShouldBeNil)
			So(out.String(), ShouldEqual, "a")
		})

		Convey("When path is not optional", func() {
			Convey("Should return an error if it occurs during path resolution", func() {
				o := NewTestObject()
				o.failAt = "bar"
				o.Set("foo", values.NewObjectWith(
					values.NewObjectProperty("bar", values.NewObjectWith(
						values.NewObjectProperty("baz", values.NewObject()),
					)),
				))

				args := []core.Value{
					values.NewString("foo"),
					values.NewString("bar"),
					values.NewString("baz"),
				}

				o.On("GetIn", mock.Anything)

				s1, _ := expressions.NewMemberPathSegment(
					core.AsExpression(func(ctx context.Context, scope *core.Scope) (core.Value, error) {
						return args[0], nil
					}),
					false,
				)

				s2, _ := expressions.NewMemberPathSegment(
					core.AsExpression(func(ctx context.Context, scope *core.Scope) (core.Value, error) {
						return args[1], nil
					}),
					false,
				)

				s3, _ := expressions.NewMemberPathSegment(
					core.AsExpression(func(ctx context.Context, scope *core.Scope) (core.Value, error) {
						return args[2], nil
					}),
					false,
				)

				segments := []*expressions.MemberPathSegment{s1, s2, s3}

				exp, err := expressions.NewMemberExpression(
					core.SourceMap{},
					core.AsExpression(func(ctx context.Context, scope *core.Scope) (core.Value, error) {
						return o, nil
					}),
					segments,
					nil,
				)

				So(err, ShouldBeNil)

				root, cancel := core.NewRootScope()

				defer func() {
					if err := cancel(); err != nil {
						panic(err)
					}
				}()

				_, err = exp.Exec(context.Background(), root.Fork())
				So(err, ShouldNotBeNil)
				So(err.Error(), ShouldEqual, core.NewPathError(core.ErrInvalidPath, 1).Format(args))

				o.AssertExpectations(t)
			})
		})

		Convey("When path is optional", func() {
			Convey("Should return None if it occurs during path resolution", func() {
				o := NewTestObject()
				o.failAt = "bar"
				o.Set("foo", values.NewObjectWith(
					values.NewObjectProperty("bar", values.NewObjectWith(
						values.NewObjectProperty("baz", values.NewObject()),
					)),
				))

				args := []core.Value{
					values.NewString("foo"),
					values.NewString("bar"),
					values.NewString("baz"),
				}

				o.On("GetIn", mock.Anything)

				s1, _ := expressions.NewMemberPathSegment(
					core.AsExpression(func(ctx context.Context, scope *core.Scope) (core.Value, error) {
						return args[0], nil
					}),
					true,
				)

				s2, _ := expressions.NewMemberPathSegment(
					core.AsExpression(func(ctx context.Context, scope *core.Scope) (core.Value, error) {
						return args[1], nil
					}),
					true,
				)

				s3, _ := expressions.NewMemberPathSegment(
					core.AsExpression(func(ctx context.Context, scope *core.Scope) (core.Value, error) {
						return args[2], nil
					}),
					true,
				)

				segments := []*expressions.MemberPathSegment{s1, s2, s3}

				exp, err := expressions.NewMemberExpression(
					core.SourceMap{},
					core.AsExpression(func(ctx context.Context, scope *core.Scope) (core.Value, error) {
						return o, nil
					}),
					segments,
					nil,
				)

				So(err, ShouldBeNil)

				root, cancel := core.NewRootScope()

				defer func() {
					if err := cancel(); err != nil {
						panic(err)
					}
				}()

				out, err := exp.Exec(context.Background(), root.Fork())
				So(err, ShouldBeNil)
				So(out.Type().String(), ShouldEqual, values.None.Type().String())

				o.AssertExpectations(t)
			})
		})
	})
}
