package compiler_test

import (
	"context"
	"github.com/MontFerret/ferret/pkg/compiler"
	"github.com/MontFerret/ferret/pkg/runtime"
	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/events"
	"github.com/MontFerret/ferret/pkg/runtime/values"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
	"time"
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

func (m *MockedObservable) Subscribe(ctx context.Context, subscription events.Subscription) (<-chan events.Event, error) {
	calls, found := m.Args[subscription.EventName]

	if !found {
		calls = make([]*values.Object, 0, 10)
		m.Args[subscription.EventName] = calls
	}

	m.Args[subscription.EventName] = append(calls, subscription.Options)

	return m.subscribers[subscription.EventName], nil
}

func newCompilerWithObservable() *compiler.Compiler {
	c := compiler.New()

	err := c.Namespace("X").
		RegisterFunctions(core.NewFunctionsFromMap(
			map[string]core.Function{
				"CREATE": func(ctx context.Context, args ...core.Value) (core.Value, error) {
					return NewMockedObservable(), nil
				},
				"EMIT": func(ctx context.Context, args ...core.Value) (core.Value, error) {
					if err := core.ValidateArgs(args, 2, 3); err != nil {
						return values.None, err
					}

					observable := args[0].(*MockedObservable)
					eventName := values.ToString(args[1])

					timeout := values.NewInt(100)

					if len(args) > 2 {
						timeout = values.ToInt(args[2])
					}

					observable.Emit(eventName.String(), values.None, nil, int64(timeout))

					return values.None, nil
				},
				"EMIT_WITH": func(ctx context.Context, args ...core.Value) (core.Value, error) {
					if err := core.ValidateArgs(args, 3, 4); err != nil {
						return values.None, err
					}

					observable := args[0].(*MockedObservable)
					eventName := values.ToString(args[1])

					timeout := values.NewInt(100)

					if len(args) > 3 {
						timeout = values.ToInt(args[3])
					}

					observable.Emit(eventName.String(), args[2], nil, int64(timeout))

					return values.None, nil
				},
				"EVENT": func(ctx context.Context, args ...core.Value) (core.Value, error) {
					return values.NewString("test"), nil
				},
			},
		))
	So(err, ShouldBeNil)

	return c
}

func TestWaitforEventExpression(t *testing.T) {
	Convey("WAITFOR EVENT parser", t, func() {
		Convey("Should parse", func() {
			c := newCompilerWithObservable()

			_, err := c.Compile(`
LET obj = X::CREATE()

X::EMIT(obj, "test", 100)
WAITFOR EVENT "test" IN obj

RETURN NONE
`)

			So(err, ShouldBeNil)
		})

		Convey("Should parse 2", func() {
			c := newCompilerWithObservable()

			_, err := c.Compile(`
LET obj = X::CREATE()

X::EMIT(obj, "test", 100)
WAITFOR EVENT "test" IN obj 1000

RETURN NONE
`)

			So(err, ShouldBeNil)
		})

		Convey("Should parse 3", func() {
			c := newCompilerWithObservable()

			_, err := c.Compile(`
LET obj = X::CREATE()

X::EMIT(obj, "test", 100)
LET timeout = 1000
WAITFOR EVENT "test" IN obj timeout

RETURN NONE
`)

			So(err, ShouldBeNil)
		})

		Convey("Should parse 4", func() {
			c := newCompilerWithObservable()

			_, err := c.Compile(`
LET obj = X::CREATE()

X::EMIT(obj, "test", 100)
LET timeout = 1000
WAITFOR EVENT "test" IN obj timeout

X::EMIT(obj, "test", 100)

RETURN NONE
`)

			So(err, ShouldBeNil)
		})
	})
	Convey("WAITFOR EVENT X IN Y runtime", t, func() {
		Convey("Should wait for a given event", func() {
			c := newCompilerWithObservable()

			prog := c.MustCompile(`
LET obj = X::CREATE()

X::EMIT(obj, "test", 100)
WAITFOR EVENT "test" IN obj

RETURN NONE
`)

			_, err := prog.Run(context.Background())

			So(err, ShouldBeNil)
		})

		Convey("Should wait for a given event using variable", func() {
			c := newCompilerWithObservable()

			prog := c.MustCompile(`
LET obj = X::CREATE()
LET eventName = "test"

X::EMIT(obj, eventName, 100)
WAITFOR EVENT eventName IN obj

RETURN NONE
`)

			_, err := prog.Run(context.Background())

			So(err, ShouldBeNil)
		})

		Convey("Should wait for a given event using object", func() {
			c := newCompilerWithObservable()

			prog := c.MustCompile(`
LET obj = X::CREATE()
LET evt = {
   name: "test"
}

X::EMIT(obj, evt.name, 100)
WAITFOR EVENT evt.name IN obj

RETURN NONE
`)

			_, err := prog.Run(context.Background())

			So(err, ShouldBeNil)
		})

		Convey("Should wait for a given event using param", func() {
			c := newCompilerWithObservable()

			prog := c.MustCompile(`
LET obj = X::CREATE()

X::EMIT(obj, @evt, 100)
WAITFOR EVENT @evt IN obj

RETURN NONE
`)

			_, err := prog.Run(context.Background(), runtime.WithParam("evt", "test"))

			So(err, ShouldBeNil)
		})

		Convey("Should use options", func() {
			observable := NewMockedObservable()
			c := newCompilerWithObservable()
			c.Namespace("X").RegisterFunction("SINGLETONE", func(ctx context.Context, args ...core.Value) (core.Value, error) {
				return observable, nil
			})

			prog := c.MustCompile(`
LET obj = X::SINGLETONE()

X::EMIT(obj, "test", 1000)
WAITFOR EVENT "test" IN obj OPTIONS { value: "foo" }

RETURN NONE
`)

			_, err := prog.Run(context.Background())

			So(err, ShouldBeNil)

			options := observable.Args["test"][0]
			So(options, ShouldNotBeNil)
			So(options.MustGet("value").String(), ShouldEqual, "foo")
		})

		Convey("Should timeout", func() {
			c := newCompilerWithObservable()

			prog := c.MustCompile(`
LET obj = X::CREATE()

X::EMIT(obj, @evt, 6000)
WAITFOR EVENT @evt IN obj

RETURN NONE
`)

			_, err := prog.Run(context.Background(), runtime.WithParam("evt", "test"))

			So(err, ShouldNotBeNil)
		})

	})

}
