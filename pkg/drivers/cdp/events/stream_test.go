package events_test

import (
	"context"
	events2 "github.com/MontFerret/ferret/pkg/drivers/cdp/events"
	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/events"
	"github.com/MontFerret/ferret/pkg/runtime/values"
	"github.com/mafredri/cdp/rpcc"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/mock"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
)

type (
	TestStream struct {
		mock.Mock
		ready   chan struct{}
		message chan events.Message
	}
)

func NewTestStream() *TestStream {
	return NewBufferedTestStream(0)
}

func NewBufferedTestStream(buffer int) *TestStream {
	es := new(TestStream)
	es.ready = make(chan struct{}, buffer)
	es.message = make(chan events.Message, buffer)
	return es
}

func (ts *TestStream) Ready() <-chan struct{} {
	return ts.ready
}

func (ts *TestStream) RecvMsg(m interface{}) error {
	return nil
}

func (ts *TestStream) Close() error {
	ts.Called()
	close(ts.message)
	close(ts.ready)
	return nil
}

func (ts *TestStream) Emit(val core.Value) {
	ts.ready <- struct{}{}
	ts.message <- events.WithValue(val)
}

func (ts *TestStream) EmitError(err error) {
	ts.ready <- struct{}{}
	ts.message <- events.WithErr(err)
}

func (ts *TestStream) Recv() (core.Value, error) {
	msg := <-ts.message

	return msg.Value(), msg.Err()
}

func TestStreamReader(t *testing.T) {
	Convey("StreamReader", t, func() {
		Convey("Should read data from Stream", func() {
			ctx, cancel := context.WithCancel(context.Background())

			stream := NewTestStream()
			stream.On("Close", mock.Anything).Maybe().Return(nil)

			go func() {
				stream.Emit(values.NewString("foo"))
				stream.Emit(values.NewString("bar"))
				stream.Emit(values.NewString("baz"))
				cancel()
			}()

			data := make([]string, 0, 3)

			es := events2.NewEventStream(stream, func(_ context.Context, stream rpcc.Stream) (core.Value, error) {
				return stream.(*TestStream).Recv()
			})

			for evt := range es.Read(ctx) {
				So(evt.Err(), ShouldBeNil)
				So(evt.Value(), ShouldNotBeNil)

				data = append(data, evt.Value().String())
			}

			So(data, ShouldResemble, []string{"foo", "bar", "baz"})

			stream.AssertExpectations(t)

			So(es.Close(context.Background()), ShouldBeNil)
		})

		Convey("Should handle error but do not close Stream", func() {
			ctx := context.Background()

			stream := NewTestStream()
			stream.On("Close", mock.Anything).Maybe().Return(nil)

			go func() {
				stream.EmitError(errors.New("foo"))
			}()

			reader := events2.NewEventStream(stream, func(_ context.Context, stream rpcc.Stream) (core.Value, error) {
				return stream.(*TestStream).Recv()
			})

			ch := reader.Read(ctx)
			evt := <-ch
			So(evt.Err(), ShouldNotBeNil)

			time.Sleep(time.Duration(100) * time.Millisecond)

			stream.AssertExpectations(t)
		})

		Convey("Should not close Stream when Context is cancelled", func() {
			stream := NewTestStream()
			stream.On("Close", mock.Anything).Maybe().Return(nil)

			reader := events2.NewEventStream(stream, func(_ context.Context, stream rpcc.Stream) (core.Value, error) {
				return values.EmptyArray(), nil
			})

			ctx, cancel := context.WithCancel(context.Background())

			_ = reader.Read(ctx)

			time.Sleep(time.Duration(100) * time.Millisecond)

			cancel()

			time.Sleep(time.Duration(100) * time.Millisecond)

			stream.AssertExpectations(t)
		})
	})
}
