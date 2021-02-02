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

func (m *MockedObservable) Subscribe(_ context.Context, eventName string, opts *values.Object) <-chan events.Event {
	calls, found := m.Args[eventName]

	if !found {
		calls = make([]*values.Object, 0, 10)
		m.Args[eventName] = calls
	}

	m.Args[eventName] = append(calls, opts)

	return m.subscribers[eventName]
}

func TestWaitforEventStatement(t *testing.T) {
	newCompiler := func() *compiler.Compiler {
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

						opts := values.ToObject(ctx, args[2])

						timeout := values.NewInt(100)

						if len(args) > 3 {
							timeout = values.ToInt(args[2])
						}

						observable.Emit(eventName.String(), opts, nil, int64(timeout))

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

	Convey("WAITFOR EVENT X IN Y", t, func() {

		Convey("Should wait for a given event", func() {
			c := newCompiler()

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
			c := newCompiler()

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
			c := newCompiler()

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
			c := newCompiler()

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
			c := newCompiler()
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
			c := newCompiler()

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
