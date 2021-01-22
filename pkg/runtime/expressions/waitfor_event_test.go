package expressions_test

import (
	"context"
	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/expressions"
	"github.com/MontFerret/ferret/pkg/runtime/expressions/literals"
	"testing"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/MontFerret/ferret/pkg/runtime/values"
)

type MockedObservable struct {
	*values.Object

	subscribers map[string]chan struct{}
}

func NewMockedObservable() *MockedObservable {
	return &MockedObservable{Object: values.NewObject(), subscribers: map[string]chan struct{}{}}
}

func (m *MockedObservable) Emit(eventName string) {
	ch := make(chan struct{})
	close(ch)
	m.subscribers[eventName] = ch
}

func (m *MockedObservable) Subscribe(_ context.Context, eventName string) (<-chan struct{}, error) {
	return m.subscribers[eventName], nil
}

func TestWaitForEventExpression(t *testing.T) {
	Convey("Should create a return expression", t, func() {
		variable, err := expressions.NewVariableExpression(core.NewSourceMap("test", 1, 10), "test")

		So(err, ShouldBeNil)

		sourceMap := core.NewSourceMap("test", 2, 10)
		expression, err := expressions.NewWaitForEventExpression(
			sourceMap,
			literals.NewStringLiteral("test"),
			variable,
			nil,
		)
		So(err, ShouldBeNil)
		So(expression, ShouldNotBeNil)
	})

	Convey("Should wait for an event", t, func() {
		mock := NewMockedObservable()
		eventName := "foobar"
		variable, err := expressions.NewVariableExpression(
			core.NewSourceMap("test", 1, 10),
			"observable",
		)

		So(err, ShouldBeNil)

		sourceMap := core.NewSourceMap("test", 2, 10)
		expression, err := expressions.NewWaitForEventExpression(
			sourceMap,
			literals.NewStringLiteral(eventName),
			variable,
			nil,
		)

		So(err, ShouldBeNil)

		scope, _ := core.NewRootScope()
		So(scope.SetVariable("observable", mock), ShouldBeNil)

		mock.Emit(eventName)
		_, err = expression.Exec(context.Background(), scope)
		So(err, ShouldBeNil)
	})

	Convey("Should timeout", t, func() {
		mock := NewMockedObservable()
		eventName := "foobar"
		variable, err := expressions.NewVariableExpression(
			core.NewSourceMap("test", 1, 10),
			"observable",
		)

		So(err, ShouldBeNil)

		sourceMap := core.NewSourceMap("test", 2, 10)
		expression, err := expressions.NewWaitForEventExpression(
			sourceMap,
			literals.NewStringLiteral(eventName),
			variable,
			nil,
		)

		So(err, ShouldBeNil)

		scope, _ := core.NewRootScope()
		So(scope.SetVariable("observable", mock), ShouldBeNil)

		_, err = expression.Exec(context.Background(), scope)
		So(err, ShouldNotBeNil)
	})
}
