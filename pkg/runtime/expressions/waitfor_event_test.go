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

func (m *MockedObservable) Subscribe(_ context.Context, sub events.Subscription) (<-chan events.Event, error) {
	calls, found := m.Args[sub.EventName]

	if !found {
		calls = make([]*values.Object, 0, 10)
		m.Args[sub.EventName] = calls
	}

	m.Args[sub.EventName] = append(calls, sub.Options)

	return m.subscribers[sub.EventName], nil
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
		)

		So(err, ShouldBeNil)

		So(expression.SetOptions(literals.NewObjectLiteralWith(prop)), ShouldBeNil)

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
		)

		So(err, ShouldBeNil)

		scope, _ := core.NewRootScope()
		So(scope.SetVariable("observable", mock), ShouldBeNil)

		_, err = expression.Exec(context.Background(), scope)
		So(err, ShouldNotBeNil)
	})

	Convey("Should filter", t, func() {
		mock := NewMockedObservable()
		eventName := "foobar"

		eventNameExp, err := expressions.NewVariableExpression(
			core.NewSourceMap("test", 1, 10),
			"observable",
		)

		So(err, ShouldBeNil)

		sourceMap := core.NewSourceMap("test", 2, 10)

		expression, err := expressions.NewWaitForEventExpression(
			sourceMap,
			literals.NewStringLiteral(eventName),
			eventNameExp,
		)

		So(err, ShouldBeNil)

		evtVar := "CURRENT"

		err = expression.SetFilter(core.SourceMap{}, evtVar, core.AsExpression(func(ctx context.Context, scope *core.Scope) (core.Value, error) {
			out, err := scope.GetVariable(evtVar)

			if err != nil {
				return nil, err
			}

			num := values.ToInt(out)

			return values.NewBoolean(int64(num) > 2), nil
		}))

		So(err, ShouldBeNil)

		mock.Emit(eventName, values.NewInt(0), nil, 0)
		mock.Emit(eventName, values.NewInt(1), nil, 0)
		mock.Emit(eventName, values.NewInt(2), nil, 0)
		mock.Emit(eventName, values.NewInt(3), nil, 0)

		scope, _ := core.NewRootScope()
		So(scope.SetVariable("observable", mock), ShouldBeNil)
		out, err := expression.Exec(context.Background(), scope)

		So(err, ShouldBeNil)
		So(int64(out.(values.Int)), ShouldEqual, 3)
	})
}
