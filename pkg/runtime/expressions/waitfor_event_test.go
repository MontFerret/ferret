package expressions_test

import (
	"context"
	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/events"
	"github.com/MontFerret/ferret/pkg/runtime/expressions"
	"github.com/MontFerret/ferret/pkg/runtime/expressions/literals"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/MontFerret/ferret/pkg/runtime/values"
)

type MockedObservable struct {
	*values.Object

	subscribers map[string]chan events.Event

	Args map[string][]*values.Object
}

func NewMockedObservable() *MockedObservable {
	return &MockedObservable{
		Object:      values.NewObject(),
		subscribers: map[string]chan events.Event{},
		Args:        map[string][]*values.Object{},
	}
}

func (m *MockedObservable) Emit(eventName string, args core.Value, err error, timeout int64) {
	ch := make(chan events.Event)
	m.subscribers[eventName] = ch

	go func() {
		<-time.After(time.Millisecond * time.Duration(timeout))
		ch <- events.Event{
			Data: args,
			Err:  err,
		}
	}()
}

func (m *MockedObservable) Subscribe(_ context.Context, eventName string, opts *values.Object) <-chan events.Event {
	calls, found := m.Args[eventName]

	if !found {
		calls = make([]*values.Object, 0, 10)
		m.Args[eventName] = calls
	}

	m.Args[eventName] = append(calls, opts)

	return m.subscribers[eventName]
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
			nil,
		)

		So(err, ShouldBeNil)

		scope, _ := core.NewRootScope()
		So(scope.SetVariable("observable", mock), ShouldBeNil)

		mock.Emit(eventName, values.None, nil, 100)
		_, err = expression.Exec(context.Background(), scope)
		So(err, ShouldBeNil)
	})

	Convey("Should receive opts", t, func() {
		mock := NewMockedObservable()
		eventName := "foobar"
		variable, err := expressions.NewVariableExpression(
			core.NewSourceMap("test", 1, 10),
			"observable",
		)

		So(err, ShouldBeNil)

		prop, err := literals.NewObjectPropertyAssignment(
			literals.NewStringLiteral("value"),
			literals.NewStringLiteral("bar"),
		)

		So(err, ShouldBeNil)

		sourceMap := core.NewSourceMap("test", 2, 10)
		expression, err := expressions.NewWaitForEventExpression(
			sourceMap,
			literals.NewStringLiteral(eventName),
			variable,
			literals.NewObjectLiteralWith(prop),
			nil,
		)

		So(err, ShouldBeNil)

		scope, _ := core.NewRootScope()
		So(scope.SetVariable("observable", mock), ShouldBeNil)

		mock.Emit(eventName, values.None, nil, 100)
		_, err = expression.Exec(context.Background(), scope)
		So(err, ShouldBeNil)

		opts := mock.Args[eventName][0]
		So(opts, ShouldNotBeNil)
	})

	Convey("Should return event arg", t, func() {
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
			nil,
		)

		So(err, ShouldBeNil)

		scope, _ := core.NewRootScope()
		So(scope.SetVariable("observable", mock), ShouldBeNil)

		arg := values.NewString("foo")
		mock.Emit(eventName, arg, nil, 100)
		out, err := expression.Exec(context.Background(), scope)
		So(err, ShouldBeNil)
		So(out.String(), ShouldEqual, arg.String())
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
			nil,
		)

		So(err, ShouldBeNil)

		scope, _ := core.NewRootScope()
		So(scope.SetVariable("observable", mock), ShouldBeNil)

		_, err = expression.Exec(context.Background(), scope)
		So(err, ShouldNotBeNil)
	})
}
