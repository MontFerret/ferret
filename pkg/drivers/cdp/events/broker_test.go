package events_test

import (
	"context"
	"github.com/MontFerret/ferret/pkg/drivers/cdp/events"
	"github.com/mafredri/cdp/protocol/dom"
	"github.com/mafredri/cdp/protocol/page"
	. "github.com/smartystreets/goconvey/convey"
	"golang.org/x/sync/errgroup"
	"sync/atomic"
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

	TestBroker struct {
		*events.EventBroker
		OnLoad           *TestLoadEventFiredClient
		OnReload         *TestDocumentUpdatedClient
		OnAttrMod        *TestAttributeModifiedClient
		OnAttrRem        *TestAttributeRemovedClient
		OnChildNodeCount *TestChildNodeCountUpdatedClient
		OnChildNodeIns   *TestChildNodeInsertedClient
		OnChildNodeRem   *TestChildNodeRemovedClient
	}
)

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

func NewTestEventBroker() *TestBroker {
	onLoad := &TestLoadEventFiredClient{NewTestEventStream()}
	onReload := &TestDocumentUpdatedClient{NewTestEventStream()}
	onAttrMod := &TestAttributeModifiedClient{NewTestEventStream()}
	onAttrRem := &TestAttributeRemovedClient{NewTestEventStream()}
	onChildCount := &TestChildNodeCountUpdatedClient{NewTestEventStream()}
	onChildIns := &TestChildNodeInsertedClient{NewTestEventStream()}
	onChildRem := &TestChildNodeRemovedClient{NewTestEventStream()}

	b := events.NewEventBroker(
		onLoad,
		onReload,
		onAttrMod,
		onAttrRem,
		onChildCount,
		onChildIns,
		onChildRem,
	)

	return &TestBroker{
		b,
		onLoad,
		onReload,
		onAttrMod,
		onAttrRem,
		onChildCount,
		onChildIns,
		onChildRem,
	}
}

func StressTest(h func() error, count int) error {
	var err error

	for i := 0; i < count; i++ {
		err = h()

		if err != nil {
			return err
		}
	}

	return nil
}

func StressTestAsync(h func() error, count int) error {
	var gr errgroup.Group

	for i := 0; i < count; i++ {
		gr.Go(h)
	}

	return gr.Wait()
}

func TestEventBroker(t *testing.T) {
	Convey(".AddEventListener", t, func() {
		Convey("Should add a new listener when not started", func() {
			b := NewTestEventBroker()

			StressTest(func() error {
				b.AddEventListener(events.EventTypeLoad, func(ctx context.Context, message interface{}) {})

				return nil
			}, 500)
		})

		Convey("Should add a new listener when started", func() {
			b := NewTestEventBroker()
			b.Start()
			defer b.Stop()

			StressTest(func() error {
				b.AddEventListener(events.EventTypeLoad, func(ctx context.Context, message interface{}) {})

				return nil
			}, 500)
		})
	})

	Convey(".RemoveEventListener", t, func() {
		Convey("Should remove a listener when not started", func() {
			b := NewTestEventBroker()

			StressTest(func() error {
				listener := func(ctx context.Context, message interface{}) {}

				b.AddEventListener(events.EventTypeLoad, listener)
				b.RemoveEventListener(events.EventTypeLoad, listener)

				So(b.ListenerCount(events.EventTypeLoad), ShouldEqual, 0)

				return nil
			}, 500)
		})

		Convey("Should add a new listener when started", func() {
			b := NewTestEventBroker()
			b.Start()
			defer b.Stop()

			StressTest(func() error {
				listener := func(ctx context.Context, message interface{}) {}

				b.AddEventListener(events.EventTypeLoad, listener)

				StressTestAsync(func() error {
					b.OnLoad.EmitDefault()

					return nil
				}, 250)

				b.RemoveEventListener(events.EventTypeLoad, listener)

				So(b.ListenerCount(events.EventTypeLoad), ShouldEqual, 0)

				return nil
			}, 250)

		})

		Convey("Should not call listener once it was removed", func() {
			b := NewTestEventBroker()
			b.Start()
			defer b.Stop()

			counter := 0

			var listener events.Handler

			listener = func(ctx context.Context, message interface{}) {
				counter++

				b.RemoveEventListener(events.EventTypeLoad, listener)
			}

			b.AddEventListener(events.EventTypeLoad, listener)
			b.OnLoad.Emit(&page.LoadEventFiredReply{})

			time.Sleep(time.Duration(10) * time.Millisecond)

			StressTestAsync(func() error {
				b.OnLoad.Emit(&page.LoadEventFiredReply{})

				return nil
			}, 250)

			So(b.ListenerCount(events.EventTypeLoad), ShouldEqual, 0)
			So(counter, ShouldEqual, 1)
		})
	})

	Convey(".Stop", t, func() {
		Convey("Should stop emitting sources", func() {
			b := NewTestEventBroker()
			b.Start()

			var counter int64

			b.AddEventListener(events.EventTypeLoad, func(ctx context.Context, message interface{}) {
				atomic.AddInt64(&counter, 1)
				b.Stop()
			})

			b.OnLoad.EmitDefault()

			time.Sleep(time.Duration(5) * time.Millisecond)

			go func() {
				b.OnLoad.EmitDefault()
			}()

			go func() {
				b.OnLoad.EmitDefault()
			}()

			time.Sleep(time.Duration(5) * time.Millisecond)

			So(atomic.LoadInt64(&counter), ShouldEqual, 1)
		})
	})
}

func BenchmarkEventBroker_AddEventListenerSync(b *testing.B) {
	broker := NewTestEventBroker()

	for n := 0; n < b.N; n++ {
		broker.AddEventListener(events.EventTypeLoad, func(ctx context.Context, message interface{}) {})
	}
}

func BenchmarkEventBroker_AddListenerAsync(b *testing.B) {
	broker := NewTestEventBroker()
	broker.Start()
	defer broker.Stop()

	for n := 0; n < b.N; n++ {
		broker.AddEventListener(events.EventTypeLoad, func(ctx context.Context, message interface{}) {})
	}
}

func BenchmarkEventBroker_Start(b *testing.B) {
	broker := NewTestEventBroker()
	broker.Start()
	defer broker.Stop()

	broker.AddEventListener(events.EventTypeLoad, func(ctx context.Context, message interface{}) {
	})

	broker.AddEventListener(events.EventTypeLoad, func(ctx context.Context, message interface{}) {
	})

	broker.AddEventListener(events.EventTypeLoad, func(ctx context.Context, message interface{}) {
	})

	broker.AddEventListener(events.EventTypeLoad, func(ctx context.Context, message interface{}) {
	})

	broker.AddEventListener(events.EventTypeLoad, func(ctx context.Context, message interface{}) {
	})

	broker.AddEventListener(events.EventTypeLoad, func(ctx context.Context, message interface{}) {
	})

	for n := 0; n < b.N; n++ {
		broker.OnLoad.Emit(&page.LoadEventFiredReply{})
	}
}
