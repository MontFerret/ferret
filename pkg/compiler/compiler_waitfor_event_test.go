package compiler_test

import (
	"context"
	"github.com/MontFerret/ferret/pkg/compiler"
	"github.com/MontFerret/ferret/pkg/runtime"
	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/events"
	"github.com/MontFerret/ferret/pkg/runtime/values"
	. "github.com/smartystreets/goconvey/convey"
	"sync/atomic"
	"testing"
	"time"
)

type MockedObservable struct {
	*values.Object

	subscribers map[string]*MockedEventStream

	Args map[string][]*values.Object
}

type MockedEventStream struct {
	ch     chan events.Event
	closed atomic.Value
}

func NewMockedEventStream(ch chan events.Event) *MockedEventStream {
	es := new(MockedEventStream)
	es.ch = ch
	es.closed.Store(false)

	return es
}

func (m *MockedEventStream) Close(ctx context.Context) error {
	close(m.ch)
	m.closed.Store(true)

	return nil
}

func (m *MockedEventStream) Read(_ context.Context) <-chan events.Event {
	return m.ch
}

func (m *MockedEventStream) Write(_ context.Context, evt events.Event) {
	closed := m.closed.Load().(bool)

	if !closed {
		m.ch <- evt
	}
}

func NewMockedObservable() *MockedObservable {
	return &MockedObservable{
		Object:      values.NewObject(),
		subscribers: make(map[string]*MockedEventStream),
		Args:        make(map[string][]*values.Object),
	}
}

func (m *MockedObservable) Emit(ctx context.Context, eventName string, args core.Value, err error, timeout int64) {
	stream, found := m.subscribers[eventName]

	if !found {
		stream = m.addStream(eventName)
	}

	go func() {
		<-time.After(time.Millisecond * time.Duration(timeout))

		if ctx.Err() != nil {
			return
		}

		if err == nil {
			stream.Write(ctx, events.WithValue(args))
		} else {
			stream.Write(ctx, events.WithErr(err))
		}
	}()
}

func (m *MockedObservable) Subscribe(_ context.Context, sub events.Subscription) (events.Stream, error) {
	calls, found := m.Args[sub.EventName]

	if !found {
		calls = make([]*values.Object, 0, 10)
		m.Args[sub.EventName] = calls
	}

	stream, found := m.subscribers[sub.EventName]

	if !found {
		stream = m.addStream(sub.EventName)
		m.subscribers[sub.EventName] = stream
	}

	m.Args[sub.EventName] = append(calls, sub.Options)

	return stream, nil
}

func (m *MockedObservable) addStream(eventName string) *MockedEventStream {
	stream := NewMockedEventStream(make(chan events.Event))
	m.subscribers[eventName] = stream

	return stream
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

					observable.Emit(ctx, eventName.String(), values.None, nil, int64(timeout))

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

					observable.Emit(ctx, eventName.String(), args[2], nil, int64(timeout))

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
	SkipConvey("WAITFOR EVENT parser", t, func() {
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
WAITFOR EVENT "test" IN obj TIMEOUT 1000

RETURN NONE
`)

			So(err, ShouldBeNil)
		})

		Convey("Should parse 3", func() {
			c := newCompilerWithObservable()

			_, err := c.Compile(`
LET obj = X::CREATE()

X::EMIT(obj, "test", 100)
LET tmt = 1000
WAITFOR EVENT "test" IN obj TIMEOUT tmt

RETURN NONE
`)

			So(err, ShouldBeNil)
		})

		Convey("Should parse 4", func() {
			c := newCompilerWithObservable()

			_, err := c.Compile(`
LET obj = X::CREATE()

X::EMIT(obj, "test", 100)
LET tmt = 1000
WAITFOR EVENT "test" IN obj TIMEOUT tmt

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

		Convey("Should use filter", func() {
			c := newCompilerWithObservable()

			prog := c.MustCompile(`
LET obj = X::CREATE()

LET _ = (FOR i IN 0..3
	X::EMIT_WITH(obj, @evt, { counter: i })
	RETURN NONE
)

LET evt = (WAITFOR EVENT @evt IN obj FILTER CURRENT.counter > 2)

T::EQ(evt.counter, 3)

RETURN evt
`)

			out, err := prog.Run(context.Background(), runtime.WithParam("evt", "test"))

			So(err, ShouldBeNil)
			So(string(out), ShouldEqual, `{"counter":3}`)
		})

		Convey("Should use filter and time out", func() {
			c := newCompilerWithObservable()

			prog := c.MustCompile(`
LET obj = X::CREATE()

LET _ = (FOR i IN 0..3
	X::EMIT_WITH(obj, @evt, { counter: i })
	RETURN NONE
)

LET evt = (WAITFOR EVENT @evt IN obj FILTER CURRENT.counter > 4 TIMEOUT 100)

T::EQ(evt.counter, 5)

RETURN evt
`)

			_, err := prog.Run(context.Background(), runtime.WithParam("evt", "test"))

			So(err, ShouldNotBeNil)
		})

		Convey("Should support pseudo-variable in different cases", func() {
			c := newCompilerWithObservable()

			prog := c.MustCompile(`
LET obj = X::CREATE()

LET _ = (FOR i IN 0..3
	X::EMIT_WITH(obj, @evt, { counter: i })
	RETURN NONE
)

LET evt = (WAITFOR EVENT @evt IN obj FILTER CURRENT.counter > 4 && current.counter < 6)

T::EQ(evt.counter, 5)

RETURN evt
`)

			_, err := prog.Run(context.Background(), runtime.WithParam("evt", "test"))

			So(err, ShouldNotBeNil)
		})
	})

}
