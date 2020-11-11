package expressions_test

import (
	"context"
	"testing"

	"github.com/MontFerret/ferret/pkg/runtime/collections"
	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/expressions"
	"github.com/MontFerret/ferret/pkg/runtime/values"
	. "github.com/smartystreets/goconvey/convey"
)

var testIterableCollectionType = core.NewType("TestIterableCollection")

type (
	testIterableCollection struct {
		values collections.IndexedCollection
	}

	testCollectionIterator struct {
		values   collections.IndexedCollection
		position values.Int
	}

	TestDataSourceExpression func(ctx context.Context, scope *core.Scope) (core.Value, error)
)

func (ds TestDataSourceExpression) Exec(ctx context.Context, scope *core.Scope) (core.Value, error) {
	return ds(ctx, scope)
}

func (c *testIterableCollection) MarshalJSON() ([]byte, error) {
	return nil, core.ErrInvalidOperation
}
func (c *testIterableCollection) Type() core.Type {
	return testIterableCollectionType
}
func (c *testIterableCollection) String() string {
	return ""
}
func (c *testIterableCollection) Compare(other core.Value) int64 {
	return 1
}
func (c *testIterableCollection) Unwrap() interface{} {
	return nil
}
func (c *testIterableCollection) Hash() uint64 {
	return 0
}
func (c *testIterableCollection) Copy() core.Value {
	return c
}
func (c *testIterableCollection) Iterate(ctx context.Context) (core.Iterator, error) {
	return &testCollectionIterator{c.values, -1}, nil
}

func (i *testCollectionIterator) Next(ctx context.Context) (core.Value, core.Value, error) {
	i.position++

	if i.position > i.values.Length() {
		return values.None, values.None, core.ErrNoMoreData
	}

	return i.values.Get(i.position), i.position, nil
}

func TestDataSource(t *testing.T) {
	Convey(".Iterate", t, func() {
		Convey("Should return custom iterable collection", func() {
			arr := values.NewArrayWith(
				values.NewInt(1),
				values.NewInt(2),
				values.NewInt(3),
				values.NewInt(4),
				values.NewInt(5),
				values.NewInt(6),
				values.NewInt(7),
				values.NewInt(8),
				values.NewInt(9),
				values.NewInt(10),
			)

			ds, err := expressions.NewForInIterableExpression(
				core.SourceMap{},
				collections.DefaultValueVar,
				collections.DefaultKeyVar,
				TestDataSourceExpression(func(ctx context.Context, scope *core.Scope) (core.Value, error) {
					return &testIterableCollection{arr}, nil
				}),
			)

			So(err, ShouldBeNil)

			rootScope, _ := core.NewRootScope()
			ctx := context.Background()
			scope := rootScope.Fork()
			out, err := ds.Iterate(ctx, scope)

			So(err, ShouldBeNil)

			pos := -1

			nextScope := scope

			for {
				nextScope, err = out.Next(ctx, nextScope.Fork())

				if err != nil {
					if core.IsNoMoreData(err) {
						break
					}

					So(err, ShouldBeNil)
				}

				pos++

				actualV, _ := nextScope.GetVariable(collections.DefaultValueVar)
				actualK, _ := nextScope.GetVariable(collections.DefaultKeyVar)

				expectedV := arr.Get(values.Int(pos))

				So(actualV, ShouldEqual, expectedV)
				So(actualK, ShouldEqual, values.Int(pos))
			}

			So(pos, ShouldEqual, int(arr.Length()))
		})

		Convey("Should stop an execution when context is cancelled", func() {
			arr := values.NewArrayWith(
				values.NewInt(1),
				values.NewInt(2),
				values.NewInt(3),
				values.NewInt(4),
				values.NewInt(5),
				values.NewInt(6),
				values.NewInt(7),
				values.NewInt(8),
				values.NewInt(9),
				values.NewInt(10),
			)

			ds, err := expressions.NewForInIterableExpression(
				core.SourceMap{},
				collections.DefaultValueVar,
				collections.DefaultKeyVar,
				TestDataSourceExpression(func(ctx context.Context, scope *core.Scope) (core.Value, error) {
					return &testIterableCollection{arr}, nil
				}),
			)

			So(err, ShouldBeNil)

			rootScope, _ := core.NewRootScope()
			scope := rootScope.Fork()
			ctx, cancel := context.WithCancel(context.Background())
			cancel()

			_, err = ds.Iterate(ctx, scope)

			So(err, ShouldEqual, core.ErrTerminated)
		})
	})
}
