package expressions_test

import (
	"context"
	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/events"
	"github.com/MontFerret/ferret/pkg/runtime/expressions"
	"github.com/MontFerret/ferret/pkg/runtime/expressions/literals"
	"sync"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/MontFerret/ferret/pkg/runtime/values"
)

type MockedObservable struct {
	*values.Object

	subscribers map[string]*MockedEventStream

	Args map[string][]*values.Object
}

type MockedEventStream struct {
	mu     sync.Mutex
	ch     chan events.Message
	closed bool
}

func NewMockedEventStream(ch chan events.Message) *MockedEventStream {
	es := new(MockedEventStream)
	es.ch = ch

	return es
}

func (m *MockedEventStream) Close(ctx context.Context) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	close(m.ch)
	m.closed = true

	return nil
}

func (m *MockedEventStream) Read(_ context.Context) <-chan events.Message {
	m.mu.Lock()
	defer m.mu.Unlock()

	return m.ch
}

func (m *MockedEventStream) Write(ctx context.Context, evt events.Message) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if ctx.Err() != nil {
		return
	}

	m.ch <- evt
}

func (m *MockedEventStream) IsClosed() bool {
	m.mu.Lock()
	defer m.mu.Unlock()

	return m.closed
}

func NewMockedObservable() *MockedObservable {
	return &MockedObservable{
		Object:      values.NewObject(),
		subscribers: make(map[string]*MockedEventStream),
		Args:        make(map[string][]*values.Object),
	}
}

func (m *MockedObservable) Emit(ctx context.Context, eventName string, args core.Value, err error, timeout int64) {
	es, ok := m.subscribers[eventName]

	if !ok {
		return
	}

	go func() {
		<-time.After(time.Millisecond * time.Duration(timeout))

		if ctx.Err() != nil {
			return
		}

		if es.IsClosed() {
			return
		}

		if err == nil {
			es.Write(ctx, events.WithValue(args))
		} else {
			es.Write(ctx, events.WithErr(err))
		}
	}()
}

func (m *MockedObservable) Subscribe(_ context.Context, sub events.Subscription) (events.Stream, error) {
	calls, found := m.Args[sub.EventName]

	if !found {
		calls = make([]*values.Object, 0, 10)
		m.Args[sub.EventName] = calls
	}

	es, found := m.subscribers[sub.EventName]

	if !found {
		es = NewMockedEventStream(make(chan events.Message))
		m.subscribers[sub.EventName] = es
	}

	m.Args[sub.EventName] = append(calls, sub.Options)

	return es, nil
}

func TestWaitForEventExpression(t *testing.T) {
	SkipConvey("Should create a return expression", t, func() {
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

	SkipConvey("Should wait for an event", t, func() {
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

		mock.Emit(context.Background(), eventName, values.None, nil, 100)
		_, err = expression.Exec(context.Background(), scope)
		So(err, ShouldBeNil)
	})

	SkipConvey("Should receive opts", t, func() {
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

		mock.Emit(context.Background(), eventName, values.None, nil, 100)
		_, err = expression.Exec(context.Background(), scope)
		So(err, ShouldBeNil)

		opts := mock.Args[eventName][0]
		So(opts, ShouldNotBeNil)
	})

	SkipConvey("Should return event arg", t, func() {
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
		mock.Emit(context.Background(), eventName, arg, nil, 100)
		out, err := expression.Exec(context.Background(), scope)
		So(err, ShouldBeNil)
		So(out.String(), ShouldEqual, arg.String())
	})

	SkipConvey("Should timeout", t, func() {
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

	SkipConvey("Should filter", t, func() {
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
		ctx := context.Background()
		mock.Emit(ctx, eventName, values.NewInt(0), nil, 0)
		mock.Emit(ctx, eventName, values.NewInt(1), nil, 0)
		mock.Emit(ctx, eventName, values.NewInt(2), nil, 0)
		mock.Emit(ctx, eventName, values.NewInt(3), nil, 0)

		scope, _ := core.NewRootScope()
		So(scope.SetVariable("observable", mock), ShouldBeNil)
		out, err := expression.Exec(context.Background(), scope)

		So(err, ShouldBeNil)
		So(int64(out.(values.Int)), ShouldEqual, 3)
	})
}
