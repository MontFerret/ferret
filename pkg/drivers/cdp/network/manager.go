package network

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/mafredri/cdp/protocol/page"
	"regexp"

	"github.com/mafredri/cdp"
	"github.com/mafredri/cdp/protocol/network"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"

	"github.com/MontFerret/ferret/pkg/drivers"
	"github.com/MontFerret/ferret/pkg/drivers/cdp/events"
)

type Manager struct {
	logger    *zerolog.Logger
	client    *cdp.Client
	headers   drivers.HTTPHeaders
	eventLoop *events.Loop
}

func New(
	logger *zerolog.Logger,
	client *cdp.Client,
) (*Manager, error) {
	m := new(Manager)
	m.logger = logger
	m.client = client
	m.headers = make(drivers.HTTPHeaders)
	m.eventLoop = events.NewLoop()

	return m, m.eventLoop.Start()
}

func (m *Manager) GetCookies(ctx context.Context) (drivers.HTTPCookies, error) {
	repl, err := m.client.Network.GetAllCookies(ctx)

	if err != nil {
		return nil, errors.Wrap(err, "failed to get cookies")
	}

	cookies := make(drivers.HTTPCookies)

	if repl.Cookies == nil {
		return cookies, nil
	}

	for _, c := range repl.Cookies {
		cookies[c.Name] = toDriverCookie(c)
	}

	return cookies, nil
}

func (m *Manager) SetCookies(ctx context.Context, url string, cookies drivers.HTTPCookies) error {
	if len(cookies) == 0 {
		return nil
	}

	params := make([]network.CookieParam, 0, len(cookies))

	for _, c := range cookies {
		params = append(params, fromDriverCookie(url, c))
	}

	return m.client.Network.SetCookies(ctx, network.NewSetCookiesArgs(params))
}

func (m *Manager) DeleteCookies(ctx context.Context, url string, cookies drivers.HTTPCookies) error {
	if len(cookies) == 0 {
		return nil
	}

	var err error

	for _, c := range cookies {
		err = m.client.Network.DeleteCookies(ctx, fromDriverCookieDelete(url, c))

		if err != nil {
			break
		}
	}

	return err
}

func (m *Manager) GetHeaders(_ context.Context) (drivers.HTTPHeaders, error) {
	copied := make(drivers.HTTPHeaders)

	for k, v := range m.headers {
		copied[k] = v
	}

	return copied, nil
}

func (m *Manager) SetHeaders(ctx context.Context, headers drivers.HTTPHeaders) error {
	if len(headers) == 0 {
		return nil
	}

	m.headers = headers

	j, err := json.Marshal(headers)

	if err != nil {
		return errors.Wrap(err, "failed to marshal headers")
	}

	err = m.client.Network.SetExtraHTTPHeaders(
		ctx,
		network.NewSetExtraHTTPHeadersArgs(j),
	)

	if err != nil {
		return errors.Wrap(err, "failed to set headers")
	}

	return nil
}

func (m *Manager) WaitForNavigation(ctx context.Context, targetURL *regexp.Regexp) error {
	return m.WaitForFrameNavigation(ctx, "", targetURL)
}

func (m *Manager) WaitForFrameNavigation(ctx context.Context, frameID page.FrameID, targetURL *regexp.Regexp) error {
	stream, err := m.client.Page.FrameNavigated(ctx)

	if err != nil {
		return errors.Wrap(err, "failed to create load event hook")
	}

	onEvent := make(chan struct{})

	listener := func(ctx context.Context, message interface{}) {
		fmt.Println("frame loaded")
		repl := message.(*page.FrameNavigatedReply)
		//fmt.Println(repl.Frame.URL)

		var matched bool

		// if frameID is empty string or equals to the current one
		if len(frameID) == 0 || repl.Frame.ID == frameID {
			// if a target URL is provided
			if targetURL != nil {
				matched = targetURL.Match([]byte(repl.Frame.URL))
			} else {
				// otherwise just notify
				matched = true
			}
		}

		if matched {
			// we need to wait till the DOM becomes ready
			onDOMReady, err := m.client.Page.DOMContentEventFired(ctx)

			if err != nil {
				// oops, well, there is really nothing we can do here
				// just log the error try luck
				m.logger.Error().Err(err).Msg("failed to create load event hook")

				onEvent <- struct{}{}

				return
			}

			_, err = onDOMReady.Recv()

			onEvent <- struct{}{}
			fmt.Println("fired event")
		}
	}

	src := events.NewSource(EventLoad, stream, func() (i interface{}, e error) {
		return stream.Recv()
	})

	m.eventLoop.AddSource(src)
	m.eventLoop.AddListener(EventLoad, listener)

	defer func() {
		m.eventLoop.RemoveListener(EventLoad, listener)
		m.eventLoop.RemoveSource(src)
		close(onEvent)

		if err := stream.Close(); err != nil {
			m.logger.Error().Err(err).Msg("failed to close frame navigated event stream")
		}
	}()

	select {
	case <-onEvent:
		fmt.Println("received event")
		return nil
	case <-ctx.Done():
		return core.ErrTimeout
	}
}
