package events_test

import (
	"github.com/MontFerret/ferret/pkg/html/dynamic/events"
	"github.com/mafredri/cdp/protocol/dom"
	"github.com/mafredri/cdp/protocol/page"
	. "github.com/smartystreets/goconvey/convey"
	"golang.org/x/sync/errgroup"
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
				b.AddEventListener(events.EventLoad, func(message interface{}) {})

				return nil
			}, 500)
		})

		Convey("Should add a new listener when started", func() {
			b := NewTestEventBroker()
			b.Start()
			defer b.Stop()

			StressTest(func() error {
				b.AddEventListener(events.EventLoad, func(message interface{}) {})

				return nil
			}, 500)
		})
	})

	Convey(".RemoveEventListener", t, func() {
		Convey("Should remove a listener when not started", func() {
			b := NewTestEventBroker()

			StressTest(func() error {
				listener := func(message interface{}) {}

				b.AddEventListener(events.EventLoad, listener)
				b.RemoveEventListener(events.EventLoad, listener)

				So(b.ListenerCount(events.EventLoad), ShouldEqual, 0)

				return nil
			}, 500)
		})

		Convey("Should add a new listener when started", func() {
			b := NewTestEventBroker()
			b.Start()
			defer b.Stop()

			StressTest(func() error {
				listener := func(message interface{}) {}

				b.AddEventListener(events.EventLoad, listener)

				StressTestAsync(func() error {
					b.OnLoad.EmitDefault()

					return nil
				}, 250)

				b.RemoveEventListener(events.EventLoad, listener)

				So(b.ListenerCount(events.EventLoad), ShouldEqual, 0)

				return nil
			}, 250)

		})

		Convey("Should not call listener once it was removed", func() {
			b := NewTestEventBroker()
			b.Start()
			defer b.Stop()

			counter := 0

			var listener events.EventListener

			listener = func(message interface{}) {
				counter += 1

				b.RemoveEventListener(events.EventLoad, listener)
			}

			b.AddEventListener(events.EventLoad, listener)
			b.OnLoad.Emit(&page.LoadEventFiredReply{})

			time.Sleep(time.Duration(10) * time.Millisecond)

			StressTestAsync(func() error {
				b.OnLoad.Emit(&page.LoadEventFiredReply{})

				return nil
			}, 250)

			So(b.ListenerCount(events.EventLoad), ShouldEqual, 0)
			So(counter, ShouldEqual, 1)
		})
	})

	Convey(".Stop", t, func() {
		Convey("Should stop emitting events", func() {
			b := NewTestEventBroker()
			b.Start()

			counter := 0
			b.AddEventListener(events.EventLoad, func(message interface{}) {
				counter++
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

			So(counter, ShouldEqual, 1)
		})
	})
}
