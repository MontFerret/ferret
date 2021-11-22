package compiler_test

import (
	"context"
	"github.com/MontFerret/ferret/pkg/compiler"
	"github.com/MontFerret/ferret/pkg/runtime"
	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/events"
	"github.com/MontFerret/ferret/pkg/runtime/values"
	"github.com/MontFerret/ferret/pkg/runtime/values/types"
	"github.com/pkg/errors"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
	"time"
)

type (
	TestObservable struct {
		*values.Object

		eventName string
		messages  []events.Message
		delay     time.Duration
		calls     []events.Subscription
	}

	TestStream chan events.Message
)

func NewTestStream(ch chan events.Message) events.Stream {
	return TestStream(ch)
}

func (s TestStream) Close(_ context.Context) error {
	return nil
}

func (s TestStream) Read(ctx context.Context) <-chan events.Message {
	proxy := make(chan events.Message)

	go func() {
		defer close(proxy)

		for {
			select {
			case <-ctx.Done():
				return
			case evt := <-s:
				if ctx.Err() != nil {
					return
				}

				proxy <- evt
			}
		}
	}()

	return s
}

func NewTestObservable(eventName string, messages []events.Message, delay time.Duration) *TestObservable {
	return &TestObservable{
		Object:    values.NewObject(),
		eventName: eventName,
		messages:  messages,
		delay:     delay,
		calls:     make([]events.Subscription, 0, 10),
	}
}

func (m *TestObservable) Subscribe(_ context.Context, sub events.Subscription) (events.Stream, error) {
	m.calls = append(m.calls, sub)

	if sub.EventName != m.eventName {
		ch := make(chan events.Message)

		return NewTestStream(ch), nil
	}

	ch := make(chan events.Message)

	go func() {
		if m.delay > 0 {
			<-time.After(m.delay * time.Millisecond)
		}

		for _, e := range m.messages {
			ch <- e
		}

		close(ch)
	}()

	return NewTestStream(ch), nil
}

func newCompilerWithObservable() *compiler.Compiler {
	c := compiler.New()

	err := c.Namespace("X").
		RegisterFunctions(core.NewFunctionsFromMap(
			map[string]core.Function{
				"VAL": func(ctx context.Context, args ...core.Value) (core.Value, error) {
					if err := core.ValidateArgs(args, 2, 3); err != nil {
						return values.None, nil
					}

					if err := core.ValidateType(args[0], types.String); err != nil {
						return values.None, err
					}

					if err := core.ValidateType(args[1], types.Array); err != nil {
						return values.None, err
					}

					name := values.ToString(args[0])
					arr := values.ToArray(ctx, args[1])
					num := values.Int(0)

					if len(args) == 3 {
						if err := core.ValidateType(args[2], types.Int); err != nil {
							return values.None, err
						}

						num = values.ToInt(args[2])
					}

					evts := make([]events.Message, 0, int(arr.Length()))

					arr.ForEach(func(value core.Value, idx int) bool {
						evts = append(evts, events.WithValue(value))

						return true
					})

					return NewTestObservable(name.String(), evts, time.Duration(num)), nil
				},
				"ERR": func(ctx context.Context, args ...core.Value) (core.Value, error) {
					if err := core.ValidateArgs(args, 3, 3); err != nil {
						return values.None, nil
					}

					if err := core.ValidateType(args[0], types.String); err != nil {
						return values.None, err
					}

					if err := core.ValidateType(args[1], types.String); err != nil {
						return values.None, err
					}

					if err := core.ValidateType(args[2], types.Int); err != nil {
						return values.None, err
					}

					name := values.ToString(args[0])
					str := values.ToString(args[1])
					num := values.ToInt(args[1])

					return NewTestObservable(name.String(), []events.Message{events.WithErr(errors.New(str.String()))}, time.Duration(num)*time.Millisecond), nil

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
LET obj = {}

WAITFOR EVENT "test" IN obj

RETURN NONE
`)

			So(err, ShouldBeNil)
		})

		Convey("Should parse 2", func() {
			c := newCompilerWithObservable()

			_, err := c.Compile(`
LET obj = {}

WAITFOR EVENT "test" IN obj TIMEOUT 1000

RETURN NONE
`)

			So(err, ShouldBeNil)
		})

		Convey("Should parse 3", func() {
			c := newCompilerWithObservable()

			_, err := c.Compile(`
LET obj = {}

LET tmt = 1000
WAITFOR EVENT "test" IN obj TIMEOUT tmt

RETURN NONE
`)

			So(err, ShouldBeNil)
		})

		SkipConvey("Should parse 4", func() {
			c := newCompilerWithObservable()

			_, err := c.Compile(`
LET obj = {}

LET tmt = 1000
WAITFOR EVENT "test" IN obj TIMEOUT tmt

RETURN NONE
`)

			So(err, ShouldBeNil)
		})
	})

	Convey("WAITFOR EVENT X IN Y runtime", t, func() {
		Convey("Should wait for a given event", func() {
			c := newCompilerWithObservable()

			prog := c.MustCompile(`
LET obj = X::VAL("test", ["foo"], 10)

WAITFOR EVENT "test" IN obj

RETURN NONE
`)

			_, err := prog.Run(context.Background())

			So(err, ShouldBeNil)
		})

		Convey("Should wait for a given event using variable", func() {
			c := newCompilerWithObservable()

			prog := c.MustCompile(`
LET eventName = "test"
LET obj = X::VAL(eventName, ["foo"], 10)

WAITFOR EVENT eventName IN obj

RETURN NONE
`)

			_, err := prog.Run(context.Background())

			So(err, ShouldBeNil)
		})

		Convey("Should wait for a given event using object", func() {
			c := newCompilerWithObservable()

			prog := c.MustCompile(`
LET evt = {
   name: "test"
}
LET obj = X::VAL(evt.name, [1], 10)

WAITFOR EVENT evt.name IN obj

RETURN NONE
`)

			_, err := prog.Run(context.Background())

			So(err, ShouldBeNil)
		})

		Convey("Should wait for a given event using param", func() {
			c := newCompilerWithObservable()

			prog := c.MustCompile(`
LET obj = X::VAL(@evt, [1], 10)

WAITFOR EVENT @evt IN obj

RETURN NONE
`)

			_, err := prog.Run(context.Background(), runtime.WithParam("evt", "test"))

			So(err, ShouldBeNil)
		})

		Convey("Should use options", func() {
			observable := NewTestObservable("test", []events.Message{events.WithValue(values.NewInt(1))}, 0)
			c := newCompilerWithObservable()
			c.Namespace("X").RegisterFunction("SINGLETONE", func(ctx context.Context, args ...core.Value) (core.Value, error) {
				return observable, nil
			})

			prog := c.MustCompile(`
			LET obj = X::SINGLETONE()

			WAITFOR EVENT "test" IN obj OPTIONS { value: "foo" }
			
			RETURN NONE
			`)

			_, err := prog.Run(context.Background())

			So(err, ShouldBeNil)

			sub := observable.calls[0]
			So(sub, ShouldNotBeNil)
			So(sub.Options.MustGet("value").String(), ShouldEqual, "foo")
		})

		Convey("Should timeout", func() {
			c := newCompilerWithObservable()

			prog := c.MustCompile(`
LET obj = X::VAL(@evt, [1], 6000)

WAITFOR EVENT @evt IN obj

RETURN NONE
`)

			_, err := prog.Run(context.Background(), runtime.WithParam("evt", "test"))

			So(err, ShouldNotBeNil)
		})

		Convey("Should use filter", func() {
			c := newCompilerWithObservable()

			prog := c.MustCompile(`
LET obj = X::VAL(@evt, [0, 1, 2, 3, 4, 5], 5)

LET evt = (WAITFOR EVENT @evt IN obj FILTER CURRENT > 3)

T::EQ(evt, 4)

RETURN evt
`)

			out, err := prog.Run(context.Background(), runtime.WithParam("evt", "test"))

			So(err, ShouldBeNil)
			So(string(out), ShouldEqual, `4`)
		})

		Convey("Should use filter and time out", func() {
			c := newCompilerWithObservable()

			prog := c.MustCompile(`
LET obj = X::VAL(@evt, [0, 1, 2, 3, 4, 5], 400)

LET evt = (WAITFOR EVENT @evt IN obj FILTER CURRENT > 4 TIMEOUT 100)

T::EQ(evt, 5)

RETURN evt
`)

			_, err := prog.Run(context.Background(), runtime.WithParam("evt", "test"))

			So(err, ShouldNotBeNil)
		})

		Convey("Should support pseudo-variable in different cases", func() {
			c := newCompilerWithObservable()

			prog := c.MustCompile(`
LET obj = X::VAL(@evt, [0,1,2,3,4,5,6], 10)

LET evt = (WAITFOR EVENT @evt IN obj FILTER CURRENT > 4 && current < 6)

T::EQ(evt, 5)

RETURN evt
`)

			_, err := prog.Run(context.Background(), runtime.WithParam("evt", "test"))

			So(err, ShouldBeNil)
		})
	})

}
