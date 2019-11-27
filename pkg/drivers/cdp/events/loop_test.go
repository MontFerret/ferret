package events_test

import (
	"context"
	"github.com/MontFerret/ferret/pkg/drivers/cdp/events"
	"github.com/mafredri/cdp/protocol/dom"
	"github.com/mafredri/cdp/protocol/page"
	"github.com/mafredri/cdp/rpcc"
	. "github.com/smartystreets/goconvey/convey"
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

func TestLoop(t *testing.T) {
	Convey(".AddListener", t, func() {
		Convey("Should add a new listener when not started", func() {
			loop := events.NewLoop()

			loop.AddListener(events.IDLoad, func(ctx context.Context, message interface{}) {})

			So(loop.ListenerCount(events.Any), ShouldEqual, 1)
		})

		Convey("Should add a new listener when started", func() {
			loop := events.NewLoop()
			loop.Start()
			defer loop.Stop()

			loop.AddListener(events.IDLoad, func(ctx context.Context, message interface{}) {})

			time.Sleep(time.Duration(100) * time.Millisecond)

			So(loop.ListenerCount(events.Any), ShouldEqual, 1)
		})
	})

	Convey(".RemoveListener", t, func() {
		Convey("Should remove a listener when not started", func() {
			loop := events.NewLoop()

			listener := func(ctx context.Context, message interface{}) {}

			loop.AddListener(events.IDLoad, listener)
			loop.RemoveListener(events.IDLoad, listener)

			So(loop.ListenerCount(events.Any), ShouldEqual, 0)
		})

		Convey("Should add a new listener when started", func() {
			loop := events.NewLoop()

			listener := func(ctx context.Context, message interface{}) {}

			loop.AddListener(events.IDLoad, listener)

			loop.Start()
			defer loop.Stop()

			loop.RemoveListener(events.IDLoad, listener)

			time.Sleep(time.Duration(100) * time.Millisecond)

			So(loop.ListenerCount(events.Any), ShouldEqual, 0)
		})
	})

	Convey(".AddSource", t, func() {
		Convey("Should add a new event source when not started", func() {
			loop := events.NewLoop()

			onLoad := &TestLoadEventFiredClient{NewTestEventStream()}

			loop.AddSource(events.NewSource(events.IDLoad, onLoad, func(_ rpcc.Stream) (i interface{}, e error) {
				return onLoad.Recv()
			}))

			So(loop.SourceCount(), ShouldEqual, 1)
		})

		Convey("Should add a new listener when started", func() {
			loop := events.NewLoop()
			loop.Start()
			defer loop.Stop()

			onLoad := &TestLoadEventFiredClient{NewTestEventStream()}

			loop.AddSource(events.NewSource(events.IDLoad, onLoad, func(_ rpcc.Stream) (i interface{}, e error) {
				return onLoad.Recv()
			}))

			time.Sleep(time.Duration(100) * time.Millisecond)

			So(loop.SourceCount(), ShouldEqual, 1)
		})
	})

	Convey(".RemoveListener", t, func() {
		Convey("Should remove a listener when not started", func() {
			loop := events.NewLoop()

			onLoad := &TestLoadEventFiredClient{NewTestEventStream()}
			src := events.NewSource(events.IDLoad, onLoad, func(_ rpcc.Stream) (i interface{}, e error) {
				return onLoad.Recv()
			})

			loop.AddSource(src)

			So(loop.SourceCount(), ShouldEqual, 1)

			loop.RemoveSource(src)

			So(loop.SourceCount(), ShouldEqual, 0)
		})

		Convey("Should add a new listener when started", func() {
			loop := events.NewLoop()

			onLoad := &TestLoadEventFiredClient{NewTestEventStream()}
			src := events.NewSource(events.IDLoad, onLoad, func(_ rpcc.Stream) (i interface{}, e error) {
				return onLoad.Recv()
			})

			loop.AddSource(src)
			So(loop.SourceCount(), ShouldEqual, 1)

			loop.Start()
			defer loop.Stop()

			loop.RemoveSource(src)

			time.Sleep(time.Duration(100) * time.Millisecond)

			So(loop.SourceCount(), ShouldEqual, 0)
		})
	})

	Convey("Should not call listener once it was removed", t, func() {
		loop := events.NewLoop()

		counter := 0

		var listener events.Handler

		listener = func(ctx context.Context, message interface{}) {
			counter++

			loop.RemoveListener(events.IDLoad, listener)
		}

		loop.AddListener(events.IDLoad, listener)

		onLoad := &TestLoadEventFiredClient{NewTestEventStream()}

		loop.AddSource(events.NewSource(events.IDLoad, onLoad, func(_ rpcc.Stream) (i interface{}, e error) {
			return onLoad.Recv()
		}))

		loop.Start()
		defer loop.Stop()

		time.Sleep(time.Duration(100) * time.Millisecond)

		onLoad.Emit(&page.LoadEventFiredReply{})

		time.Sleep(time.Duration(10) * time.Millisecond)

		So(loop.ListenerCount(events.IDLoad), ShouldEqual, 0)
		So(counter, ShouldEqual, 1)
	})

	//
	//Convey(".Stop", t, func() {
	//	Convey("Should stop emitting sources", func() {
	//		b := NewTestEventBroker()
	//		b.Start()
	//
	//		var counter int64
	//
	//		b.AddEventListener(sources.IDLoad, func(ctx context.Context, message interface{}) {
	//			atomic.AddInt64(&counter, 1)
	//			b.Stop()
	//		})
	//
	//		b.OnLoad.EmitDefault()
	//
	//		time.Sleep(time.Duration(5) * time.Millisecond)
	//
	//		go func() {
	//			b.OnLoad.EmitDefault()
	//		}()
	//
	//		go func() {
	//			b.OnLoad.EmitDefault()
	//		}()
	//
	//		time.Sleep(time.Duration(5) * time.Millisecond)
	//
	//		So(atomic.LoadInt64(&counter), ShouldEqual, 1)
	//	})
	//})
}

func BenchmarkLoop_AddListenerSync(b *testing.B) {
	loop := events.NewLoop()

	for n := 0; n < b.N; n++ {
		loop.AddListener(events.IDLoad, func(ctx context.Context, message interface{}) {})
	}
}

func BenchmarkLoop_AddListenerAsync(b *testing.B) {
	loop := events.NewLoop()
	loop.Start()
	defer loop.Stop()

	for n := 0; n < b.N; n++ {
		loop.AddListener(events.IDLoad, func(ctx context.Context, message interface{}) {})
	}
}

func BenchmarkLoop_AddListenerAsync2(b *testing.B) {
	loop := events.NewLoop()
	loop.Start()
	defer loop.Stop()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			loop.AddListener(events.IDLoad, func(ctx context.Context, message interface{}) {})
		}
	})
}

func BenchmarkLoop_Start(b *testing.B) {
	loop := events.NewLoop()

	loop.AddListener(events.IDLoad, func(ctx context.Context, message interface{}) {

	})
	loop.AddListener(events.IDLoad, func(ctx context.Context, message interface{}) {

	})

	loop.AddListener(events.IDLoad, func(ctx context.Context, message interface{}) {

	})

	loop.AddListener(events.IDLoad, func(ctx context.Context, message interface{}) {

	})

	loop.AddListener(events.IDLoad, func(ctx context.Context, message interface{}) {

	})

	loop.AddListener(events.IDLoad, func(ctx context.Context, message interface{}) {

	})

	onLoad := &TestLoadEventFiredClient{NewTestEventStream()}

	loop.AddSource(events.NewSource(events.IDLoad, onLoad, func(_ rpcc.Stream) (i interface{}, e error) {
		return onLoad.Recv()
	}))

	loop.Start()
	defer loop.Stop()

	for n := 0; n < b.N; n++ {
		onLoad.Emit(&page.LoadEventFiredReply{})
	}
}
