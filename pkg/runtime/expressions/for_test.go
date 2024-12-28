package expressions

import (
	"context"
	"testing"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/MontFerret/ferret/pkg/runtime/collections"
	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/expressions/clauses"
	"github.com/MontFerret/ferret/pkg/runtime/values"
)

type (
	testIterator struct {
		values           []*core.Scope
		pos              int
		causeErrorInNext bool
	}
	testIterable struct {
		values              []*core.Scope
		causeErrorInIterate bool
		causeErrorInNext    bool
	}
	testExpression struct {
		causeErrorInExec bool
	}
	testError struct{}
)

func (iterator *testIterator) Next(ctx context.Context, scope *core.Scope) (*core.Scope, error) {
	if iterator.causeErrorInNext {
		return nil, testError{}
	}

	if len(iterator.values) > iterator.pos {
		val := iterator.values[iterator.pos]
		iterator.pos++

		return val, nil
	}

	return nil, core.ErrNoMoreData
}

func (iterable *testIterable) Iterate(ctx context.Context, scope *core.Scope) (collections.Iterator, error) {
	if iterable.causeErrorInIterate {
		return nil, testError{}
	}

	return &testIterator{iterable.values, 0, iterable.causeErrorInNext}, nil
}

func (expression *testExpression) Exec(ctx context.Context, scope *core.Scope) (core.Value, error) {
	return nil, nil
}

func (error testError) Error() string {
	return "error"
}

func testInitTestIterable(values []*core.Scope, causeErrorInIterate, causeErrorInNext bool) *testIterable {
	return &testIterable{values, causeErrorInIterate, causeErrorInNext}
}

func TestNewForExpression(t *testing.T) {
	Convey(".NewForExpression", t, func() {
		Convey("Should return new ForExpresssion.", func() {
			forExp, err := NewForExpression(core.SourceMap{}, testInitTestIterable([]*core.Scope{}, false, false), &testExpression{}, false, false, false)
			So(forExp, ShouldNotBeNil)
			So(err, ShouldBeNil)
		})

		Convey("Should return error when a dataSource is nil", func() {
			forExp, err := NewForExpression(core.SourceMap{}, nil, &testExpression{}, false, false, false)
			So(forExp, ShouldBeNil)
			So(err, ShouldNotBeNil)
			So(err, ShouldEqual, err)
		})

		Convey("Should return error when a predicate is nil", func() {
			forExp, err := NewForExpression(core.SourceMap{}, testInitTestIterable([]*core.Scope{}, false, false), nil, false, false, false)
			So(forExp, ShouldBeNil)
			So(err, ShouldNotBeNil)
			So(err, ShouldEqual, err)
		})
	})
}

func TestAddLimit(t *testing.T) {
	Convey(".AddLimit", t, func() {
		Convey("Should success. (An Error Should be nil.)", func() {
			forExpression, _ := NewForExpression(
				core.SourceMap{},
				testInitTestIterable([]*core.Scope{}, false, false),
				&testExpression{},
				false,
				false,
				false,
			)

			err := forExpression.AddLimit(core.SourceMap{}, &testExpression{}, &testExpression{})

			So(err, ShouldBeNil)
		})

		Convey("Should return an error.", func() {
			forExpression, _ := NewForExpression(
				core.SourceMap{},
				testInitTestIterable([]*core.Scope{}, false, false),
				&testExpression{},
				false,
				false,
				false,
			)
			forExpression.dataSource = nil

			err := forExpression.AddLimit(core.SourceMap{}, &testExpression{}, &testExpression{})

			So(err, ShouldNotBeNil)
		})
	})
}

func TestAddFilter(t *testing.T) {
	Convey(".AddFilter", t, func() {
		Convey("Should success. (An Error Should be nil.)", func() {
			forExpression, _ := NewForExpression(
				core.SourceMap{},
				testInitTestIterable([]*core.Scope{}, false, false),
				&testExpression{},
				false,
				false,
				false,
			)

			err := forExpression.AddFilter(core.SourceMap{}, &testExpression{})

			So(err, ShouldBeNil)
		})

		Convey("Should return an error.", func() {
			forExpression, _ := NewForExpression(
				core.SourceMap{},
				testInitTestIterable([]*core.Scope{}, false, false),
				&testExpression{},
				false,
				false,
				false,
			)
			forExpression.dataSource = nil

			err := forExpression.AddFilter(core.SourceMap{}, &testExpression{})

			So(err, ShouldNotBeNil)
		})
	})

}

func TestAddSort(t *testing.T) {
	Convey(".AddSort", t, func() {
		Convey("Should success.(An Error Should be nil.)", func() {
			forExpression, _ := NewForExpression(
				core.SourceMap{},
				testInitTestIterable([]*core.Scope{}, false, false),
				&testExpression{},
				false,
				false,
				false,
			)

			err := forExpression.AddSort(core.SourceMap{}, &clauses.SorterExpression{})

			So(err, ShouldBeNil)
		})

		Convey("Should return an error.", func() {
			forExpression, _ := NewForExpression(
				core.SourceMap{},
				testInitTestIterable([]*core.Scope{}, false, false),
				&testExpression{},
				false,
				false,
				false,
			)
			forExpression.dataSource = nil

			err := forExpression.AddSort(core.SourceMap{}, &clauses.SorterExpression{})

			So(err, ShouldNotBeNil)
		})
	})
}

func TestAddCollect(t *testing.T) {
	Convey(".AddCollect", t, func() {
		Convey("Should success. (Error Should be nil.)", func() {
			forExpression, _ := NewForExpression(
				core.SourceMap{},
				testInitTestIterable([]*core.Scope{}, false, false),
				&testExpression{},
				false,
				false,
				false,
			)

			err := forExpression.AddCollect(core.SourceMap{}, &clauses.Collect{})

			So(err, ShouldBeNil)
		})

		Convey("Should return an error.", func() {
			forExpression, _ := NewForExpression(
				core.SourceMap{},
				testInitTestIterable([]*core.Scope{}, false, false),
				&testExpression{},
				false,
				false,
				false,
			)
			forExpression.dataSource = nil

			err := forExpression.AddCollect(core.SourceMap{}, &clauses.Collect{})

			So(err, ShouldNotBeNil)
		})
	})
}

func TestAddStatement(t *testing.T) {
	Convey(".AddStatement (Error Should be nil.)", t, func() {
		Convey("Should success.", func() {
			forExpression, _ := NewForExpression(
				core.SourceMap{},
				testInitTestIterable([]*core.Scope{}, false, false),
				&testExpression{},
				false,
				false,
				false,
			)

			err := forExpression.AddStatement(&testExpression{})

			So(err, ShouldBeNil)
		})

		Convey("Should return an error.", func() {
			forExpression, _ := NewForExpression(
				core.SourceMap{},
				testInitTestIterable([]*core.Scope{}, false, false),
				&testExpression{},
				false,
				false,
				false,
			)
			forExpression.dataSource = nil

			err := forExpression.AddStatement(&testExpression{})

			So(err, ShouldNotBeNil)
		})
	})
}

func TestExec(t *testing.T) {
	Convey(".Exec", t, func() {
		Convey("Should success.", func() {
			forExpression, _ := NewForExpression(
				core.SourceMap{},
				testInitTestIterable([]*core.Scope{{}, {}, {}}, false, false),
				&testExpression{},
				false,
				false,
				false,
			)

			_, err := forExpression.Exec(context.Background(), &core.Scope{})

			So(err, ShouldBeNil)
		})

		Convey("Should stop an execution when context is cancelled.", func() {
			forExpression, _ := NewForExpression(
				core.SourceMap{},
				testInitTestIterable([]*core.Scope{{}, {}, {}}, false, false),
				&testExpression{},
				false,
				false,
				false,
			)

			ctx0 := context.Background()
			ctx1, cancelFn := context.WithCancel(ctx0)
			cancelFn()

			_, err := forExpression.Exec(ctx1, &core.Scope{})

			So(err, ShouldNotBeNil)
		})

		Convey("Should return an error when a dataSource expression is invalid.", func() {
			forExpression, _ := NewForExpression(
				core.SourceMap{},
				testInitTestIterable([]*core.Scope{}, true, false),
				&testExpression{},
				false,
				false,
				false,
			)

			result, err := forExpression.Exec(context.Background(), &core.Scope{})
			So(result, ShouldEqual, values.None)
			So(err, ShouldNotBeNil)
		})

		Convey("Should return an error when element expressions of dataSource is invalidated.", func() {
			forExpression, _ := NewForExpression(
				core.SourceMap{},
				testInitTestIterable([]*core.Scope{}, false, true),
				&testExpression{},
				false,
				false,
				false,
			)

			result, err := forExpression.Exec(context.Background(), &core.Scope{})
			So(result, ShouldEqual, values.None)
			So(err, ShouldNotBeNil)
		})

		Convey("Should return an error when an predicate expression is invalidated.", func() {
			forExpression, _ := NewForExpression(
				core.SourceMap{},
				testInitTestIterable([]*core.Scope{}, false, true),
				&testExpression{true},
				false,
				false,
				false,
			)

			result, err := forExpression.Exec(context.Background(), &core.Scope{})
			So(result, ShouldEqual, values.None)
			So(err, ShouldNotBeNil)
		})
	})
}
