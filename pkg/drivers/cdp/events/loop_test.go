package events_test

import (
	"context"
	"github.com/MontFerret/ferret/pkg/drivers/cdp/events"
	"github.com/mafredri/cdp/protocol/dom"
	"github.com/mafredri/cdp/protocol/page"
	"github.com/mafredri/cdp/rpcc"
	. "github.com/smartystreets/goconvey/convey"
	"sync"
	"testing"
	"time"
)

type (
	TestEventStream struct {
		ready   chan struct{}
		message chan interface{}
	}

	TestLoadEventFiredClient struct {
		*TestEventStream
	}

	TestDocumentUpdatedClient struct {
		*TestEventStream
	}

	TestAttributeModifiedClient struct {
		*TestEventStream
	}

	TestAttributeRemovedClient struct {
		*TestEventStream
	}

	TestChildNodeCountUpdatedClient struct {
		*TestEventStream
	}

	TestChildNodeInsertedClient struct {
		*TestEventStream
	}

	TestChildNodeRemovedClient struct {
		*TestEventStream
	}
)

var TestEvent = events.New("test_event")

func NewTestEventStream() *TestEventStream {
	es := new(TestEventStream)
	es.ready = make(chan struct{})
	es.message = make(chan interface{})
	return es
}

func (es *TestEventStream) Ready() <-chan struct{} {
	return es.ready
}

func (es *TestEventStream) RecvMsg(i interface{}) error {
	// NOT IMPLEMENTED
	return nil
}

func (es *TestEventStream) Close() error {
	close(es.message)
	close(es.ready)
	return nil
}

func (es *TestEventStream) Emit(msg interface{}) {
	es.ready <- struct{}{}
	es.message <- msg
}

func (es *TestLoadEventFiredClient) Recv() (*page.LoadEventFiredReply, error) {
	r := <-es.message
	reply := r.(*page.LoadEventFiredReply)

	return reply, nil
}

func (es *TestLoadEventFiredClient) EmitDefault() {
	es.TestEventStream.Emit(&page.LoadEventFiredReply{})
}

func (es *TestDocumentUpdatedClient) Recv() (*dom.DocumentUpdatedReply, error) {
	r := <-es.message
	reply := r.(*dom.DocumentUpdatedReply)

	return reply, nil
}

func (es *TestAttributeModifiedClient) Recv() (*dom.AttributeModifiedReply, error) {
	r := <-es.message
	reply := r.(*dom.AttributeModifiedReply)

	return reply, nil
}

func (es *TestAttributeRemovedClient) Recv() (*dom.AttributeRemovedReply, error) {
	r := <-es.message
	reply := r.(*dom.AttributeRemovedReply)

	return reply, nil
}

func (es *TestChildNodeCountUpdatedClient) Recv() (*dom.ChildNodeCountUpdatedReply, error) {
	r := <-es.message
	reply := r.(*dom.ChildNodeCountUpdatedReply)

	return reply, nil
}

func (es *TestChildNodeInsertedClient) Recv() (*dom.ChildNodeInsertedReply, error) {
	r := <-es.message
	reply := r.(*dom.ChildNodeInsertedReply)

	return reply, nil
}

func (es *TestChildNodeRemovedClient) Recv() (*dom.ChildNodeRemovedReply, error) {
	r := <-es.message
	reply := r.(*dom.ChildNodeRemovedReply)

	return reply, nil
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
			loop := events.NewLoop()
			counter := NewCounter()

			onLoad := &TestLoadEventFiredClient{NewTestEventStream()}
			src := events.NewSource(TestEvent, onLoad, func(_ rpcc.Stream) (i interface{}, e error) {
				return onLoad.Recv()
			})

			loop.AddSource(src)

			ctx, cancel := context.WithCancel(context.Background())

			loop.Run(ctx)
			defer cancel()

			onLoad.EmitDefault()

			wait()

			So(counter.Value(), ShouldEqual, 0)

			loop.AddListener(TestEvent, events.Always(func(ctx context.Context, message interface{}) {
				counter.Increase()
			}))

			wait()

			onLoad.EmitDefault()

			wait()

			So(counter.Value(), ShouldEqual, 1)
		})
	})

	Convey(".RemoveListener", t, func() {
		Convey("Should remove a listener", func() {
			Convey("Should add a new listener", func() {
				loop := events.NewLoop()
				counter := NewCounter()

				onLoad := &TestLoadEventFiredClient{NewTestEventStream()}
				src := events.NewSource(TestEvent, onLoad, func(_ rpcc.Stream) (i interface{}, e error) {
					return onLoad.Recv()
				})

				loop.AddSource(src)
				id := loop.AddListener(TestEvent, events.Always(func(ctx context.Context, message interface{}) {
					counter.Increase()
				}))

				ctx, cancel := context.WithCancel(context.Background())

				loop.Run(ctx)
				defer cancel()

				onLoad.EmitDefault()

				wait()

				So(counter.Value(), ShouldEqual, 1)

				wait()

				loop.RemoveListener(TestEvent, id)

				wait()

				onLoad.EmitDefault()

				wait()

				So(counter.Value(), ShouldEqual, 1)
			})
		})
	})

	Convey(".AddSource", t, func() {
		Convey("Should add a new event source when not started", func() {
			loop := events.NewLoop()
			counter := NewCounter()

			loop.AddListener(TestEvent, events.Always(func(ctx context.Context, message interface{}) {
				counter.Increase()
			}))

			ctx, cancel := context.WithCancel(context.Background())

			loop.Run(ctx)
			defer cancel()

			onLoad := &TestLoadEventFiredClient{NewTestEventStream()}

			go func() {
				onLoad.EmitDefault()
			}()

			wait()

			So(counter.Value(), ShouldEqual, 0)

			src := events.NewSource(TestEvent, onLoad, func(_ rpcc.Stream) (i interface{}, e error) {
				return onLoad.Recv()
			})

			loop.AddSource(src)

			wait()

			So(counter.Value(), ShouldEqual, 1)
		})
	})

	Convey(".RemoveSource", t, func() {
		Convey("Should remove a source", func() {
			loop := events.NewLoop()
			counter := NewCounter()

			ctx, cancel := context.WithCancel(context.Background())

			loop.Run(ctx)
			defer cancel()

			loop.AddListener(TestEvent, events.Always(func(ctx context.Context, message interface{}) {
				counter.Increase()
			}))

			onLoad := &TestLoadEventFiredClient{NewTestEventStream()}
			src := events.NewSource(TestEvent, onLoad, func(_ rpcc.Stream) (i interface{}, e error) {
				return onLoad.Recv()
			})

			loop.AddSource(src)

			wait()

			onLoad.EmitDefault()

			wait()

			So(counter.Value(), ShouldEqual, 1)

			loop.RemoveSource(src)

			wait()

			go func() {
				onLoad.EmitDefault()
			}()

			wait()

			So(counter.Value(), ShouldEqual, 1)
		})
	})

	Convey("Should not call listener once it was removed", t, func() {
		loop := events.NewLoop()
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

		onLoad := &TestLoadEventFiredClient{NewTestEventStream()}

		loop.AddSource(events.NewSource(TestEvent, onLoad, func(_ rpcc.Stream) (i interface{}, e error) {
			return onLoad.Recv()
		}))

		ctx, cancel := context.WithCancel(context.Background())

		loop.Run(ctx)
		defer cancel()

		time.Sleep(time.Duration(100) * time.Millisecond)

		onLoad.Emit(&page.LoadEventFiredReply{})

		time.Sleep(time.Duration(10) * time.Millisecond)

		So(counter.Value(), ShouldEqual, 1)
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
	loop := events.NewLoop()

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

	onLoad := &TestLoadEventFiredClient{NewTestEventStream()}

	loop.AddSource(events.NewSource(TestEvent, onLoad, func(_ rpcc.Stream) (i interface{}, e error) {
		return onLoad.Recv()
	}))

	ctx, cancel := context.WithCancel(context.Background())

	loop.Run(ctx)
	defer cancel()

	for n := 0; n < b.N; n++ {
		onLoad.Emit(&page.LoadEventFiredReply{})
	}
}
