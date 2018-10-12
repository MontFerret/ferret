package events_test

import (
	"github.com/MontFerret/ferret/pkg/html/dynamic/events"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
	"time"
)

type (
	TestMessage struct {
		data string
	}
	TestEventStream struct {
		ready   chan struct{}
		message chan interface{}
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
	m := <-es.message

	input := i.(*TestMessage)
	message := m.(*TestMessage)

	input.data = message.data

	return nil
}

func (es *TestEventStream) Close() error {
	close(es.message)
	close(es.ready)

	return nil
}

func (es *TestEventStream) Emit(msg *TestMessage) {
	es.ready <- struct{}{}
	es.message <- msg
}

func TestEventBroker_AddEventListener(t *testing.T) {
	Convey("When is not started", t, func() {
		Convey("Should add an event listener", func() {
			es := NewTestEventStream()

			broker := events.NewEventBroker(
				events.NewEventStream("test", es, func() interface{} {
					return &TestMessage{}
				}),
			)

			var fired bool

			broker.AddEventListener("test", func(message interface{}) {
				fired = true
			})

			err := broker.Start()

			So(err, ShouldBeNil)

			go es.Emit(&TestMessage{})

			time.Sleep(time.Millisecond * 250)

			So(fired, ShouldBeTrue)
		})
	})

	Convey("When is started", t, func() {
		Convey("Should add an event listener", func() {
			es := NewTestEventStream()

			broker := events.NewEventBroker(
				events.NewEventStream("test", es, func() interface{} {
					return &TestMessage{}
				}),
			)

			err := broker.Start()

			var fired bool

			broker.AddEventListener("test", func(message interface{}) {
				fired = true
			})

			So(err, ShouldBeNil)

			go es.Emit(&TestMessage{})

			time.Sleep(time.Millisecond * 250)

			So(fired, ShouldBeTrue)
		})
	})
}

func TestEventBroker_RemoveEventListener(t *testing.T) {
	Convey("When is not started", t, func() {
		Convey("Should remove an event listener", func() {
			es := NewTestEventStream()

			broker := events.NewEventBroker(
				events.NewEventStream("test", es, func() interface{} {
					return &TestMessage{}
				}),
			)

			var counter int

			listener := func(message interface{}) {
				counter++
			}

			broker.AddEventListener("test", listener)

			err := broker.Start()

			So(err, ShouldBeNil)

			go es.Emit(&TestMessage{})

			time.Sleep(time.Millisecond * 250)

			So(counter, ShouldEqual, 1)

			err = broker.Stop()

			So(err, ShouldBeNil)

			broker.RemoveEventListener("test", listener)

			err = broker.Start()

			So(err, ShouldBeNil)

			go es.Emit(&TestMessage{})

			time.Sleep(time.Millisecond * 250)

			So(counter, ShouldEqual, 1)
		})
	})

	Convey("When is started", t, func() {
		Convey("Should remove an event listener", func() {
			es := NewTestEventStream()

			broker := events.NewEventBroker(
				events.NewEventStream("test", es, func() interface{} {
					return &TestMessage{}
				}),
			)

			err := broker.Start()
			So(err, ShouldBeNil)

			var counter int

			listener := func(message interface{}) {
				counter++
			}

			broker.AddEventListener("test", listener)

			go es.Emit(&TestMessage{})

			time.Sleep(time.Millisecond * 250)

			So(counter, ShouldEqual, 1)

			broker.RemoveEventListener("test", listener)

			go es.Emit(&TestMessage{})

			time.Sleep(time.Millisecond * 250)

			So(counter, ShouldEqual, 1)
		})
	})
}
