package network

import (
	"github.com/MontFerret/ferret/pkg/drivers/cdp/streams"
	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
	"github.com/mafredri/cdp/protocol/network"
	"github.com/mafredri/cdp/protocol/page"
	"github.com/mafredri/cdp/rpcc"
	"github.com/rs/zerolog"
)

func newFrameNavigatedReader(logger zerolog.Logger) *streams.Reader {
	return streams.NewReader(func(stream rpcc.Stream) (core.Value, error) {
		repl, err := stream.(page.FrameNavigatedClient).Recv()

		if err != nil {
			logger.Trace().Err(err).Msg("failed to read data from frame navigation event stream")

			return values.None, nil
		}

		logger.Trace().
			Str("url", repl.Frame.URL).
			Str("frame_id", string(repl.Frame.ID)).
			Str("type", string(repl.Type)).
			Msg("received frame navigation event")

		return &NavigationEvent{
			URL:     repl.Frame.URL,
			FrameID: repl.Frame.ID,
		}, nil
	})
}

func newRequestWillBeSentReader(logger zerolog.Logger) *streams.Reader {
	return streams.NewReader(func(stream rpcc.Stream) (core.Value, error) {
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

func newResponseReceivedReader(logger zerolog.Logger) *streams.Reader {
	return streams.NewReader(func(stream rpcc.Stream) (core.Value, error) {
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

		return toDriverResponse(repl.Response), nil
	})
}
