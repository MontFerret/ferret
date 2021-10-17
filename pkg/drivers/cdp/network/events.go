package network

import (
	"context"
	"github.com/MontFerret/ferret/pkg/drivers/cdp/events"
	"github.com/mafredri/cdp"
	"github.com/mafredri/cdp/protocol/fetch"
	"github.com/mafredri/cdp/protocol/network"
	"github.com/mafredri/cdp/protocol/page"
	"github.com/mafredri/cdp/rpcc"
)

var (
	frameNavigatedEvent   = events.New("frame_navigated")
	responseReceivedEvent = events.New("response_received")
	requestPausedEvent    = events.New("request_paused")
	beforeRequestEvent    = events.New("before_request")
)

func createFrameLoadStreamFactory(client *cdp.Client) events.SourceFactory {
	return events.NewStreamSourceFactory(frameNavigatedEvent, func(ctx context.Context) (rpcc.Stream, error) {
		return client.Page.FrameNavigated(ctx)
	}, func(stream rpcc.Stream) (interface{}, error) {
		return stream.(page.FrameNavigatedClient).Recv()
	})
}

func createResponseReceivedStreamFactory(client *cdp.Client) events.SourceFactory {
	return events.NewStreamSourceFactory(responseReceivedEvent, func(ctx context.Context) (rpcc.Stream, error) {
		return client.Network.ResponseReceived(ctx)
	}, func(stream rpcc.Stream) (interface{}, error) {
		return stream.(network.ResponseReceivedClient).Recv()
	})
}

func createRequestPausedStreamFactory(client *cdp.Client) events.SourceFactory {
	return events.NewStreamSourceFactory(requestPausedEvent, func(ctx context.Context) (rpcc.Stream, error) {
		return client.Fetch.RequestPaused(ctx)
	}, func(stream rpcc.Stream) (interface{}, error) {
		return stream.(fetch.RequestPausedClient).Recv()
	})
}

func createBeforeRequestStreamFactory(client *cdp.Client) events.SourceFactory {
	return events.NewStreamSourceFactory(beforeRequestEvent, func(ctx context.Context) (rpcc.Stream, error) {
		return client.Network.RequestWillBeSent(ctx)
	}, func(stream rpcc.Stream) (interface{}, error) {
		return stream.(network.RequestWillBeSentClient).Recv()
	})
}
