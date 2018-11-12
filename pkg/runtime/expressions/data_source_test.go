package expressions_test

import (
	"context"
	"github.com/MontFerret/ferret/pkg/runtime/expressions"
	"github.com/MontFerret/ferret/pkg/runtime/values"
	"testing"

	"github.com/MontFerret/ferret/pkg/runtime/collections"
	"github.com/MontFerret/ferret/pkg/runtime/core"
	. "github.com/smartystreets/goconvey/convey"
)

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
	return core.Type(11)
}
func (c *testIterableCollection) String() string {
	return ""
}
func (c *testIterableCollection) Compare(other core.Value) int {
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
func (c *testIterableCollection) Iterate(ctx context.Context) (collections.CollectionIterator, error) {
	return &testCollectionIterator{c.values, -1}, nil
}

func (i *testCollectionIterator) Next(ctx context.Context) (core.Value, core.Value, error) {
	i.position++

	if i.position > i.values.Length() {
		return values.None, values.None, nil
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

			ds, err := expressions.NewDataSource(
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
				pos++
				nextScope, err = out.Next(ctx, nextScope.Fork())

				So(err, ShouldBeNil)

				if nextScope == nil {
					break
				}

				actualV, _ := nextScope.GetVariable(collections.DefaultValueVar)
				actualK, _ := nextScope.GetVariable(collections.DefaultKeyVar)

				expectedV := arr.Get(values.Int(pos))

				So(actualV, ShouldEqual, expectedV)
				So(actualK, ShouldEqual, values.Int(pos))
			}

			So(pos, ShouldEqual, int(arr.Length()))
		})
	})
}
