package network

import (
	"context"
	"encoding/base64"
	"sync/atomic"

	"github.com/mafredri/cdp"
	"github.com/mafredri/cdp/protocol/network"
	"github.com/mafredri/cdp/protocol/page"
	"github.com/mafredri/cdp/rpcc"
	"github.com/rs/zerolog"

	"github.com/MontFerret/ferret/pkg/drivers/cdp/eval"
	"github.com/MontFerret/ferret/pkg/drivers/cdp/events"
	"github.com/MontFerret/ferret/pkg/drivers/cdp/templates"
	"github.com/MontFerret/ferret/pkg/runtime/core"
	rtEvents "github.com/MontFerret/ferret/pkg/runtime/events"
	"github.com/MontFerret/ferret/pkg/runtime/values"
)

type NavigationEventStream struct {
	logger  zerolog.Logger
	client  *cdp.Client
	tail    atomic.Value
	onFrame page.FrameNavigatedClient
	onDoc   page.NavigatedWithinDocumentClient
}

func newNavigationEventStream(
	logger zerolog.Logger,
	client *cdp.Client,
	onFrame page.FrameNavigatedClient,
	onDoc page.NavigatedWithinDocumentClient,
) rtEvents.Stream {
	es := new(NavigationEventStream)
	es.logger = logger
	es.client = client
	es.onFrame = onFrame
	es.onDoc = onDoc

	return es
}

func (s *NavigationEventStream) Read(ctx context.Context) <-chan rtEvents.Message {
	ch := make(chan rtEvents.Message)

	go func() {
		defer close(ch)

		for {
			select {
			case <-ctx.Done():
				return
			case <-s.onDoc.Ready():
				if ctx.Err() != nil {
					return
				}

				repl, err := s.onDoc.Recv()

				if err != nil {
					ch <- rtEvents.WithErr(err)
					s.logger.Trace().Err(err).Msg("failed to read data from within document navigation event stream")

					return
				}

				evt := NavigationEvent{
					URL:     repl.URL,
					FrameID: repl.FrameID,
				}

				s.logger.Trace().
					Str("url", evt.URL).
					Str("frame_id", string(evt.FrameID)).
					Str("type", evt.MimeType).
					Msg("received withing document navigation event")

				s.tail.Store(evt)

				ch <- rtEvents.WithValue(&evt)
			case <-s.onFrame.Ready():
				if ctx.Err() != nil {
					return
				}

				repl, err := s.onFrame.Recv()

				if err != nil {
					ch <- rtEvents.WithErr(err)
					s.logger.Trace().Err(err).Msg("failed to read data from frame navigation event stream")

					return
				}

				evt := NavigationEvent{
					URL:      repl.Frame.URL,
					FrameID:  repl.Frame.ID,
					MimeType: repl.Frame.MimeType,
				}

				s.logger.Trace().
					Str("url", evt.URL).
					Str("frame_id", string(evt.FrameID)).
					Str("type", evt.MimeType).
					Msg("received frame navigation event")

				s.tail.Store(evt)

				ch <- rtEvents.WithValue(&evt)
			}
		}
	}()

	return ch
}

func (s *NavigationEventStream) Close(ctx context.Context) error {
	val := s.tail.Load()

	evt, ok := val.(NavigationEvent)

	if !ok || evt.FrameID == "" {
		// TODO: err?
		return nil
	}

	_ = s.onFrame.Close()
	_ = s.onDoc.Close()

	s.logger.Trace().
		Str("frame_id", string(evt.FrameID)).
		Str("frame_url", evt.URL).
		Msg("creating frame execution context")

	ec, err := eval.Create(ctx, s.logger, s.client, evt.FrameID)

	if err != nil {
		s.logger.Trace().
			Err(err).
			Str("frame_id", string(evt.FrameID)).
			Str("frame_url", evt.URL).
			Msg("failed to create frame execution context")

		return err
	}

	s.logger.Trace().
		Str("frame_id", string(evt.FrameID)).
		Str("frame_url", evt.URL).
		Msg("starting polling DOM ready event")

	_, err = events.NewEvalWaitTask(
		ec,
		templates.DOMReady(),
		events.DefaultPolling,
	).Run(ctx)

	if err != nil {
		s.logger.Trace().
			Err(err).
			Str("frame_id", string(evt.FrameID)).
			Str("frame_url", evt.URL).
			Msg("failed to poll DOM ready event")

		return err
	}

	s.logger.Trace().
		Str("frame_id", string(evt.FrameID)).
		Str("frame_url", evt.URL).
		Msg("DOM is ready. Navigation has completed")

	return nil
}

func newRequestWillBeSentStream(logger zerolog.Logger, input network.RequestWillBeSentClient) rtEvents.Stream {
	return events.NewEventStream(input, func(_ context.Context, stream rpcc.Stream) (core.Value, error) {
		repl, err := stream.(network.RequestWillBeSentClient).Recv()

		if err != nil {
			logger.Trace().Err(err).Msg("failed to read data from request event stream")

			return values.None, nil
		}

		var frameID string

		if repl.FrameID != nil {
			frameID = string(*repl.FrameID)
		}

		logger.Trace().
			Str("url", repl.Request.URL).
			Str("document_url", repl.DocumentURL).
			Str("frame_id", frameID).
			Interface("data", repl.Request).
			Msg("received request event")

		return toDriverRequest(repl.Request), nil
	})
}

func newResponseReceivedReader(logger zerolog.Logger, client *cdp.Client, input network.ResponseReceivedClient) rtEvents.Stream {
	return events.NewEventStream(input, func(ctx context.Context, stream rpcc.Stream) (core.Value, error) {
		repl, err := stream.(network.ResponseReceivedClient).Recv()

		if err != nil {
			logger.Trace().Err(err).Msg("failed to read data from request event stream")

			return values.None, nil
		}

		var frameID string

		if repl.FrameID != nil {
			frameID = string(*repl.FrameID)
		}

		logger.Trace().
			Str("url", repl.Response.URL).
			Str("frame_id", frameID).
			Str("request_id", string(repl.RequestID)).
			Interface("data", repl.Response).
			Msg("received response event")

		var body []byte

		resp, err := client.Network.GetResponseBody(ctx, network.NewGetResponseBodyArgs(repl.RequestID))

		if err == nil {
			body = make([]byte, 0, 0)

			if resp.Base64Encoded {
				body, err = base64.StdEncoding.DecodeString(resp.Body)

				if err != nil {
					logger.Warn().
						Str("url", repl.Response.URL).
						Str("frame_id", frameID).
						Str("request_id", string(repl.RequestID)).
						Interface("data", repl.Response).
						Msg("failed to decode response body")
				}
			} else {
				body = []byte(resp.Body)
			}
		} else {
			logger.Warn().
				Str("url", repl.Response.URL).
				Str("frame_id", frameID).
				Str("request_id", string(repl.RequestID)).
				Interface("data", repl.Response).
				Msg("failed to get response body")
		}

		return toDriverResponse(repl.Response, body), nil
	})
}
