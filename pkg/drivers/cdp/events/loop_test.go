package events_test

import (
	"context"
	"github.com/MontFerret/ferret/pkg/drivers/cdp/events"
	"github.com/mafredri/cdp/protocol/page"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
	"time"
)

func TestLoop(t *testing.T) {
	Convey(".AddListener", t, func() {
		Convey("Should add a new listener when not started", func() {
			loop := events.NewLoop()

			loop.AddListener(events.EventTypeLoad, func(ctx context.Context, message interface{}) {})

			So(loop.ListenerCount(events.EventTypeAny), ShouldEqual, 1)
		})

		Convey("Should add a new listener when started", func() {
			loop := events.NewLoop()
			loop.Start()
			defer loop.Stop()

			loop.AddListener(events.EventTypeLoad, func(ctx context.Context, message interface{}) {})

			time.Sleep(time.Duration(100) * time.Millisecond)

			So(loop.ListenerCount(events.EventTypeAny), ShouldEqual, 1)
		})
	})

	Convey(".RemoveListener", t, func() {
		Convey("Should remove a listener when not started", func() {
			loop := events.NewLoop()

			listener := func(ctx context.Context, message interface{}) {}

			loop.AddListener(events.EventTypeLoad, listener)
			loop.RemoveListener(events.EventTypeLoad, listener)

			So(loop.ListenerCount(events.EventTypeAny), ShouldEqual, 0)
		})

		Convey("Should add a new listener when started", func() {
			loop := events.NewLoop()

			listener := func(ctx context.Context, message interface{}) {}

			loop.AddListener(events.EventTypeLoad, listener)

			loop.Start()
			defer loop.Stop()

			loop.RemoveListener(events.EventTypeLoad, listener)

			time.Sleep(time.Duration(100) * time.Millisecond)

			So(loop.ListenerCount(events.EventTypeAny), ShouldEqual, 0)
		})
	})

	Convey("Should not call listener once it was removed", t, func() {
		loop := events.NewLoop()

		counter := 0

		var listener events.Handler

		listener = func(ctx context.Context, message interface{}) {
			counter++

			loop.RemoveListener(events.EventTypeLoad, listener)
		}

		loop.AddListener(events.EventTypeLoad, listener)

		onLoad := &TestLoadEventFiredClient{NewTestEventStream()}

		loop.AddSource(events.NewSource(events.EventTypeLoad, onLoad, func() (i interface{}, e error) {
			return onLoad.Recv()
		}))

		loop.Start()
		defer loop.Stop()

		time.Sleep(time.Duration(100) * time.Millisecond)

		onLoad.Emit(&page.LoadEventFiredReply{})

		time.Sleep(time.Duration(10) * time.Millisecond)

		So(loop.ListenerCount(events.EventTypeLoad), ShouldEqual, 0)
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
	//		b.AddEventListener(sources.EventTypeLoad, func(ctx context.Context, message interface{}) {
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
		loop.AddListener(events.EventTypeLoad, func(ctx context.Context, message interface{}) {})
	}
}

func BenchmarkLoop_AddListenerAsync(b *testing.B) {
	loop := events.NewLoop()
	loop.Start()
	defer loop.Stop()

	for n := 0; n < b.N; n++ {
		loop.AddListener(events.EventTypeLoad, func(ctx context.Context, message interface{}) {})
	}
}

func BenchmarkLoop_AddListenerAsync2(b *testing.B) {
	loop := events.NewLoop()
	loop.Start()
	defer loop.Stop()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			loop.AddListener(events.EventTypeLoad, func(ctx context.Context, message interface{}) {})
		}
	})
}

func BenchmarkLoop_Start(b *testing.B) {
	loop := events.NewLoop()

	loop.AddListener(events.EventTypeLoad, func(ctx context.Context, message interface{}) {

	})
	loop.AddListener(events.EventTypeLoad, func(ctx context.Context, message interface{}) {

	})

	loop.AddListener(events.EventTypeLoad, func(ctx context.Context, message interface{}) {

	})

	loop.AddListener(events.EventTypeLoad, func(ctx context.Context, message interface{}) {

	})

	loop.AddListener(events.EventTypeLoad, func(ctx context.Context, message interface{}) {

	})

	loop.AddListener(events.EventTypeLoad, func(ctx context.Context, message interface{}) {

	})

	onLoad := &TestLoadEventFiredClient{NewTestEventStream()}

	loop.AddSource(events.NewSource(events.EventTypeLoad, onLoad, func() (i interface{}, e error) {
		return onLoad.Recv()
	}))

	loop.Start()
	defer loop.Stop()

	for n := 0; n < b.N; n++ {
		onLoad.Emit(&page.LoadEventFiredReply{})
	}
}
