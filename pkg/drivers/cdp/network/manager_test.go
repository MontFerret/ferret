package network_test

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/mafredri/cdp"
	"github.com/mafredri/cdp/protocol/fetch"
	network2 "github.com/mafredri/cdp/protocol/network"
	"github.com/mafredri/cdp/protocol/page"
	"github.com/rs/zerolog"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/stretchr/testify/mock"

	"github.com/MontFerret/ferret/pkg/drivers"
	"github.com/MontFerret/ferret/pkg/drivers/cdp/network"
)

type (
	PageAPI struct {
		mock.Mock
		cdp.Page
		frameNavigated func(ctx context.Context) (page.FrameNavigatedClient, error)
	}

	NetworkAPI struct {
		mock.Mock
		cdp.Network
		responseReceived    func(ctx context.Context) (network2.ResponseReceivedClient, error)
		setExtraHTTPHeaders func(ctx context.Context, args *network2.SetExtraHTTPHeadersArgs) error
	}

	FetchAPI struct {
		mock.Mock
		cdp.Fetch
		enable        func(context.Context, *fetch.EnableArgs) error
		disable       func(context.Context) error
		requestPaused func(context.Context) (fetch.RequestPausedClient, error)
	}

	TestEventStream struct {
		mock.Mock
		ready   chan struct{}
		message chan interface{}
	}

	FrameNavigatedClient struct {
		*TestEventStream
	}

	ResponseReceivedClient struct {
		*TestEventStream
	}

	RequestPausedClient struct {
		*TestEventStream
	}
)

func (api *PageAPI) FrameNavigated(ctx context.Context) (page.FrameNavigatedClient, error) {
	return api.frameNavigated(ctx)
}

func (api *NetworkAPI) ResponseReceived(ctx context.Context) (network2.ResponseReceivedClient, error) {
	return api.responseReceived(ctx)
}

func (api *NetworkAPI) SetExtraHTTPHeaders(ctx context.Context, args *network2.SetExtraHTTPHeadersArgs) error {
	return api.setExtraHTTPHeaders(ctx, args)
}

func (api *FetchAPI) Enable(ctx context.Context, args *fetch.EnableArgs) error {
	if api.enable == nil {
		return nil
	}

	return api.enable(ctx, args)
}

func (api *FetchAPI) Disable(ctx context.Context) error {
	if api.disable == nil {
		return nil
	}

	return api.disable(ctx)
}

func (api *FetchAPI) RequestPaused(ctx context.Context) (fetch.RequestPausedClient, error) {
	return api.requestPaused(ctx)
}

func NewTestEventStream() *TestEventStream {
	return NewBufferedTestEventStream(0)
}

func NewBufferedTestEventStream(buffer int) *TestEventStream {
	es := new(TestEventStream)
	es.ready = make(chan struct{}, buffer)
	es.message = make(chan interface{}, buffer)
	return es
}

func (stream *TestEventStream) Ready() <-chan struct{} {
	return stream.ready
}

func (stream *TestEventStream) RecvMsg(i interface{}) error {
	return nil
}

func (stream *TestEventStream) Message() interface{} {
	return <-stream.message
}

func (stream *TestEventStream) Close() error {
	stream.Called()
	close(stream.message)
	close(stream.ready)
	return nil
}

func (stream *TestEventStream) Emit(msg interface{}) {
	stream.ready <- struct{}{}
	stream.message <- msg
}

func NewFrameNavigatedClient() *FrameNavigatedClient {
	return &FrameNavigatedClient{
		TestEventStream: NewTestEventStream(),
	}
}

func (stream *FrameNavigatedClient) Recv() (*page.FrameNavigatedReply, error) {
	<-stream.Ready()
	msg := stream.Message()

	repl, ok := msg.(*page.FrameNavigatedReply)

	if !ok {
		panic("Invalid message type")
	}

	return repl, nil
}

func NewResponseReceivedClient() *ResponseReceivedClient {
	return &ResponseReceivedClient{
		TestEventStream: NewTestEventStream(),
	}
}

func (stream *ResponseReceivedClient) Recv() (*network2.ResponseReceivedReply, error) {
	<-stream.Ready()
	msg := stream.Message()

	repl, ok := msg.(*network2.ResponseReceivedReply)

	if !ok {
		panic("Invalid message type")
	}

	return repl, nil
}

func NewRequestPausedClient() *RequestPausedClient {
	return &RequestPausedClient{
		TestEventStream: NewTestEventStream(),
	}
}

func (stream *RequestPausedClient) Recv() (*fetch.RequestPausedReply, error) {
	<-stream.Ready()
	msg := stream.Message()

	repl, ok := msg.(*fetch.RequestPausedReply)

	if !ok {
		panic("Invalid message type")
	}

	return repl, nil
}

func TestManager(t *testing.T) {
	Convey("Network manager", t, func() {

		Convey("New", func() {
			Convey("Should close all resources on Close", func() {
				responseReceivedClient := NewResponseReceivedClient()
				responseReceivedClient.On("Close", mock.Anything).Once().Return(nil)
				networkAPI := new(NetworkAPI)
				networkAPI.responseReceived = func(ctx context.Context) (network2.ResponseReceivedClient, error) {
					return responseReceivedClient, nil
				}
				networkAPI.setExtraHTTPHeaders = func(ctx context.Context, args *network2.SetExtraHTTPHeadersArgs) error {
					return nil
				}

				requestPausedClient := NewRequestPausedClient()
				requestPausedClient.On("Close", mock.Anything).Once().Return(nil)
				fetchAPI := new(FetchAPI)
				fetchAPI.enable = func(ctx context.Context, args *fetch.EnableArgs) error {
					return nil
				}
				fetchAPI.requestPaused = func(ctx context.Context) (fetch.RequestPausedClient, error) {
					return requestPausedClient, nil
				}

				client := &cdp.Client{
					Network: networkAPI,
					Fetch:   fetchAPI,
				}

				mgr, err := network.New(
					zerolog.New(os.Stdout).Level(zerolog.Disabled),
					client,
					network.Options{
						Headers: drivers.NewHTTPHeadersWith(map[string][]string{"x-correlation-id": {"foo"}}),
						Filter: &network.Filter{
							Patterns: []drivers.ResourceFilter{
								{
									URL:  "http://google.com",
									Type: "img",
								},
							},
						},
					},
				)

				So(err, ShouldBeNil)
				So(mgr.Close(), ShouldBeNil)

				time.Sleep(time.Duration(100) * time.Millisecond)

				responseReceivedClient.AssertExpectations(t)
				requestPausedClient.AssertExpectations(t)
			})
		})
	})
}
