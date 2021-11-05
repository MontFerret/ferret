package cdp

import (
	"context"
	"github.com/MontFerret/ferret/pkg/runtime/events"

	"github.com/mafredri/cdp"
	"github.com/mafredri/cdp/protocol/emulation"
	"github.com/mafredri/cdp/protocol/network"
	"github.com/mafredri/cdp/protocol/page"
	"golang.org/x/sync/errgroup"

	"github.com/MontFerret/ferret/pkg/drivers"
	"github.com/MontFerret/ferret/pkg/drivers/common"
)

type (
	batchFunc = func() error

	closer func(ctx context.Context) error

	pageNavigationEventStream struct {
		stream events.Stream
		closer
	}
)

func newPageNavigationEventStream(stream events.Stream, closer closer) events.Stream {
	return &pageNavigationEventStream{stream, closer}
}

func (p *pageNavigationEventStream) Close(ctx context.Context) error {
	if err := p.stream.Close(ctx); err != nil {
		return err
	}

	return p.closer(ctx)
}

func (p *pageNavigationEventStream) Read(ctx context.Context) <-chan events.Message {
	return p.stream.Read(ctx)
}

func runBatch(funcs ...batchFunc) error {
	eg := errgroup.Group{}

	for _, f := range funcs {
		eg.Go(f)
	}

	return eg.Wait()
}

func enableFeatures(ctx context.Context, client *cdp.Client, params drivers.Params) error {
	if err := client.Page.Enable(ctx); err != nil {
		return err
	}

	return runBatch(
		func() error {
			return client.Page.SetLifecycleEventsEnabled(
				ctx,
				page.NewSetLifecycleEventsEnabledArgs(true),
			)
		},

		func() error {
			return client.DOM.Enable(ctx)
		},

		func() error {
			return client.Runtime.Enable(ctx)
		},

		func() error {
			ua := common.GetUserAgent(params.UserAgent)

			// do not use custom user agent
			if ua == "" {
				return nil
			}

			return client.Emulation.SetUserAgentOverride(
				ctx,
				emulation.NewSetUserAgentOverrideArgs(ua),
			)
		},

		func() error {
			return client.Network.Enable(ctx, network.NewEnableArgs())
		},

		func() error {
			return client.Page.SetBypassCSP(ctx, page.NewSetBypassCSPArgs(true))
		},

		func() error {
			if params.Viewport == nil {
				return nil
			}

			orientation := emulation.ScreenOrientation{}

			if !params.Viewport.Landscape {
				orientation.Type = "portraitPrimary"
				orientation.Angle = 0
			} else {
				orientation.Type = "landscapePrimary"
				orientation.Angle = 90
			}

			scaleFactor := params.Viewport.ScaleFactor

			if scaleFactor <= 0 {
				scaleFactor = 1
			}

			deviceArgs := emulation.NewSetDeviceMetricsOverrideArgs(
				params.Viewport.Width,
				params.Viewport.Height,
				scaleFactor,
				params.Viewport.Mobile,
			).SetScreenOrientation(orientation)

			return client.Emulation.SetDeviceMetricsOverride(
				ctx,
				deviceArgs,
			)
		},
	)
}
