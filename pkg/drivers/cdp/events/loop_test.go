package events_test

import (
	"context"
	"github.com/MontFerret/ferret/pkg/drivers/cdp/events"
	"github.com/mafredri/cdp/rpcc"
	. "github.com/smartystreets/goconvey/convey"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

type TestEventStream struct {
	closed   atomic.Value
	ready    chan struct{}
	messages chan string
}

var TestEvent = events.New("test_event")

func NewTestEventStream() *TestEventStream {
	return NewBufferedTestEventStream(0)
}

func NewBufferedTestEventStream(buffer int) *TestEventStream {
	es := new(TestEventStream)
	es.ready = make(chan struct{}, buffer)
	es.messages = make(chan string, buffer)
	es.closed.Store(false)

	return es
}

func (es *TestEventStream) IsClosed() bool {
	return es.closed.Load().(bool)
}

func (es *TestEventStream) Ready() <-chan struct{} {
	return es.ready
}

func (es *TestEventStream) RecvMsg(i interface{}) error {
	// NOT IMPLEMENTED
	return nil
}

func (es *TestEventStream) Recv() (interface{}, error) {
	msg := <-es.messages

	return msg, nil
}

func (es *TestEventStream) Close() error {
	es.closed.Store(true)
	close(es.messages)
	close(es.ready)
	return nil
}

func (es *TestEventStream) EmitP(msg string, skipCheck bool) {
	if !skipCheck {
		isClosed := es.closed.Load().(bool)

		if isClosed {
			return
		}
	}

	es.ready <- struct{}{}
	es.messages <- msg
}

func (es *TestEventStream) Emit(msg string) {
	es.EmitP(msg, false)
}

func (es *TestEventStream) EmitDefault() {
	es.Emit("")
}

func wait() {
	time.Sleep(time.Duration(50) * time.Millisecond)
}

type Counter struct {
	mu    sync.Mutex
	value int64
}

func NewCounter() *Counter {
	return new(Counter)
}

func (c *Counter) Increase() *Counter {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.value++

	return c
}

func (c *Counter) Decrease() *Counter {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.value--

	return c
}

func (c *Counter) Value() int64 {
	c.mu.Lock()
	defer c.mu.Unlock()

	return c.value
}

func TestLoop(t *testing.T) {
	Convey(".AddListener", t, func() {
		Convey("Should add a new listener", func() {
			counter := NewCounter()

			var tes *TestEventStream

			loop := events.NewLoop(events.NewStreamSourceFactory(TestEvent, func(ctx context.Context) (rpcc.Stream, error) {
				tes = NewTestEventStream()
				return tes, nil
			}, func(stream rpcc.Stream) (interface{}, error) {
				return stream.(*TestEventStream).Recv()
			}))

			ctx, cancel := context.WithCancel(context.Background())

			err := loop.Run(ctx)
			defer cancel()

			So(err, ShouldBeNil)

			tes.EmitDefault()

			wait()

			So(counter.Value(), ShouldEqual, 0)

			loop.AddListener(TestEvent, events.Always(func(ctx context.Context, message interface{}) {
				counter.Increase()
			}))

			wait()

			tes.EmitDefault()

			wait()

			So(counter.Value(), ShouldEqual, 1)
		})
	})

	Convey(".RemoveListener", t, func() {
		Convey("Should remove a listener", func() {
			Convey("Should add a new listener", func() {
				counter := NewCounter()

				var test *TestEventStream

				loop := events.NewLoop(events.NewStreamSourceFactory(TestEvent, func(ctx context.Context) (rpcc.Stream, error) {
					test = NewTestEventStream()
					return test, nil
				}, func(stream rpcc.Stream) (interface{}, error) {
					return stream.(*TestEventStream).Recv()
				}))

				id := loop.AddListener(TestEvent, events.Always(func(ctx context.Context, message interface{}) {
					counter.Increase()
				}))

				ctx, cancel := context.WithCancel(context.Background())

				err := loop.Run(ctx)
				defer cancel()

				So(err, ShouldBeNil)

				test.EmitDefault()

				wait()

				So(counter.Value(), ShouldEqual, 1)

				wait()

				loop.RemoveListener(TestEvent, id)

				wait()

				test.EmitDefault()

				wait()

				So(counter.Value(), ShouldEqual, 1)
			})
		})
	})

	Convey("Should not call listener once it was removed", t, func() {
		var tes *TestEventStream

		loop := events.NewLoop(events.NewStreamSourceFactory(TestEvent, func(ctx context.Context) (rpcc.Stream, error) {
			tes = NewTestEventStream()
			return tes, nil
		}, func(stream rpcc.Stream) (interface{}, error) {
			return stream.(*TestEventStream).Recv()
		}))

		onEvent := make(chan struct{})

		counter := NewCounter()

		id := loop.AddListener(TestEvent, events.Always(func(ctx context.Context, message interface{}) {
			counter.Increase()

			onEvent <- struct{}{}
		}))

		go func() {
			<-onEvent

			loop.RemoveListener(TestEvent, id)
		}()

		ctx, cancel := context.WithCancel(context.Background())

		err := loop.Run(ctx)
		So(err, ShouldBeNil)
		defer cancel()

		time.Sleep(time.Duration(100) * time.Millisecond)

		tes.EmitDefault()

		time.Sleep(time.Duration(10) * time.Millisecond)

		So(counter.Value(), ShouldEqual, 1)
	})

	Convey("Should stop on Context.Done", t, func() {
		eventsToFire := 5
		counter := NewCounter()

		var tes *TestEventStream

		loop := events.NewLoop(events.NewStreamSourceFactory(TestEvent, func(ctx context.Context) (rpcc.Stream, error) {
			tes = NewTestEventStream()
			return tes, nil
		}, func(stream rpcc.Stream) (interface{}, error) {
			return stream.(*TestEventStream).Recv()
		}))

		loop.AddListener(TestEvent, events.Always(func(ctx context.Context, message interface{}) {
			counter.Increase()
		}))

		ctx, cancel := context.WithCancel(context.Background())
		err := loop.Run(ctx)
		So(err, ShouldBeNil)

		for i := 0; i <= eventsToFire; i++ {
			time.Sleep(time.Duration(100) * time.Millisecond)

			tes.EmitDefault()
		}

		// Stop the loop
		cancel()

		time.Sleep(time.Duration(100) * time.Millisecond)

		So(tes.IsClosed(), ShouldBeTrue)
	})
}

func BenchmarkLoop_AddListenerSync(b *testing.B) {
	loop := events.NewLoop()

	for n := 0; n < b.N; n++ {
		loop.AddListener(TestEvent, events.Always(func(ctx context.Context, message interface{}) {}))
	}
}

func BenchmarkLoop_AddListenerAsync(b *testing.B) {
	loop := events.NewLoop()
	ctx, cancel := context.WithCancel(context.Background())

	loop.Run(ctx)
	defer cancel()

	for n := 0; n < b.N; n++ {
		loop.AddListener(TestEvent, events.Always(func(ctx context.Context, message interface{}) {}))
	}
}

func BenchmarkLoop_AddListenerAsync2(b *testing.B) {
	loop := events.NewLoop()
	ctx, cancel := context.WithCancel(context.Background())

	loop.Run(ctx)
	defer cancel()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			loop.AddListener(TestEvent, events.Always(func(ctx context.Context, message interface{}) {}))
		}
	})
}

func BenchmarkLoop_Start(b *testing.B) {
	var tes *TestEventStream

	loop := events.NewLoop(events.NewStreamSourceFactory(TestEvent, func(ctx context.Context) (rpcc.Stream, error) {
		tes = NewTestEventStream()
		return tes, nil
	}, func(stream rpcc.Stream) (interface{}, error) {
		return stream.(*TestEventStream).Recv()
	}))

	loop.AddListener(TestEvent, events.Always(func(ctx context.Context, message interface{}) {

	}))
	loop.AddListener(TestEvent, events.Always(func(ctx context.Context, message interface{}) {

	}))

	loop.AddListener(TestEvent, events.Always(func(ctx context.Context, message interface{}) {

	}))

	loop.AddListener(TestEvent, events.Always(func(ctx context.Context, message interface{}) {

	}))

	loop.AddListener(TestEvent, events.Always(func(ctx context.Context, message interface{}) {

	}))

	loop.AddListener(TestEvent, events.Always(func(ctx context.Context, message interface{}) {

	}))

	ctx, cancel := context.WithCancel(context.Background())

	if err := loop.Run(ctx); err != nil {
		panic(err)
	}

	defer cancel()

	for n := 0; n < b.N; n++ {
		tes.EmitP("", true)
	}
}

func BenchmarkLoop_StartAsync(b *testing.B) {
	var tes *TestEventStream

	loop := events.NewLoop(events.NewStreamSourceFactory(TestEvent, func(ctx context.Context) (rpcc.Stream, error) {
		tes = NewBufferedTestEventStream(b.N)
		return tes, nil
	}, func(stream rpcc.Stream) (interface{}, error) {
		return stream.(*TestEventStream).Recv()
	}))

	loop.AddListener(TestEvent, events.Always(func(ctx context.Context, message interface{}) {

	}))
	loop.AddListener(TestEvent, events.Always(func(ctx context.Context, message interface{}) {

	}))

	loop.AddListener(TestEvent, events.Always(func(ctx context.Context, message interface{}) {

	}))

	loop.AddListener(TestEvent, events.Always(func(ctx context.Context, message interface{}) {

	}))

	loop.AddListener(TestEvent, events.Always(func(ctx context.Context, message interface{}) {

	}))

	loop.AddListener(TestEvent, events.Always(func(ctx context.Context, message interface{}) {

	}))

	ctx, cancel := context.WithCancel(context.Background())

	if err := loop.Run(ctx); err != nil {
		panic(err)
	}

	defer cancel()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			tes.EmitP("", true)
		}
	})
}
