package network

import (
	"context"
	"sync"

	"github.com/mafredri/cdp"
	"github.com/mafredri/cdp/protocol/network"
	"github.com/mafredri/cdp/protocol/page"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/wI2L/jettison"

	"github.com/MontFerret/ferret/pkg/drivers"
	"github.com/MontFerret/ferret/pkg/drivers/cdp/events"
	"github.com/MontFerret/ferret/pkg/runtime/core"
	rtEvents "github.com/MontFerret/ferret/pkg/runtime/events"
	"github.com/MontFerret/ferret/pkg/runtime/logging"
	"github.com/MontFerret/ferret/pkg/runtime/values"
)

const BlankPageURL = "about:blank"

type (
	FrameLoadedListener = func(ctx context.Context, frame page.Frame)

	Manager struct {
		mu          sync.RWMutex
		logger      zerolog.Logger
		client      *cdp.Client
		headers     *drivers.HTTPHeaders
		loop        *events.Loop
		interceptor *Interceptor
		stop        context.CancelFunc
		response    *sync.Map
	}
)

func New(
	logger zerolog.Logger,
	client *cdp.Client,
	options Options,
) (*Manager, error) {
	ctx, cancel := context.WithCancel(context.Background())

	m := new(Manager)
	m.logger = logging.WithName(logger.With(), "network_manager").Logger()
	m.client = client
	m.headers = drivers.NewHTTPHeaders()
	m.stop = cancel
	m.response = new(sync.Map)

	var err error

	defer func() {
		if err != nil {
			m.stop()
		}
	}()

	m.loop = events.NewLoop(
		createResponseReceivedStreamFactory(client),
	)

	m.loop.AddListener(responseReceivedEvent, m.handleResponse)

	if options.Filter != nil && len(options.Filter.Patterns) > 0 {
		m.interceptor = NewInterceptor(logger, client)

		if err := m.interceptor.AddFilter("resources", options.Filter); err != nil {
			return nil, err
		}

		if err = m.interceptor.Run(ctx); err != nil {
			return nil, err
		}
	}

	if options.Cookies != nil && len(options.Cookies) > 0 {
		for url, cookies := range options.Cookies {
			err = m.setCookiesInternal(ctx, url, cookies)

			if err != nil {
				return nil, err
			}
		}
	}

	if options.Headers != nil && options.Headers.Length() > 0 {
		err = m.setHeadersInternal(ctx, options.Headers)

		if err != nil {
			return nil, err
		}
	}

	if err = m.loop.Run(ctx); err != nil {
		return nil, err
	}

	return m, nil
}

func (m *Manager) Close() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.logger.Trace().Msg("closing")

	if m.stop != nil {
		m.stop()
		m.stop = nil
	}

	return nil
}

func (m *Manager) GetCookies(ctx context.Context) (*drivers.HTTPCookies, error) {
	m.logger.Trace().Msg("starting to get cookies")

	repl, err := m.client.Network.GetAllCookies(ctx)

	if err != nil {
		m.logger.Trace().Err(err).Msg("failed to get cookies")

		return nil, errors.Wrap(err, "failed to get cookies")
	}

	cookies := drivers.NewHTTPCookies()

	if repl.Cookies == nil {
		m.logger.Trace().Msg("no cookies found")

		return cookies, nil
	}

	for _, c := range repl.Cookies {
		cookies.Set(toDriverCookie(c))
	}

	m.logger.Trace().Err(err).Msg("succeeded to get cookies")

	return cookies, nil
}

func (m *Manager) SetCookies(ctx context.Context, url string, cookies *drivers.HTTPCookies) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	return m.setCookiesInternal(ctx, url, cookies)
}

func (m *Manager) setCookiesInternal(ctx context.Context, url string, cookies *drivers.HTTPCookies) error {
	m.logger.Trace().Str("url", url).Msg("starting to set cookies")

	if cookies == nil {
		m.logger.Trace().Msg("nil cookies passed")

		return errors.Wrap(core.ErrMissedArgument, "cookies")
	}

	if cookies.Length() == 0 {
		m.logger.Trace().Msg("no cookies passed")

		return nil
	}

	params := make([]network.CookieParam, 0, cookies.Length())

	cookies.ForEach(func(value drivers.HTTPCookie, _ values.String) bool {
		params = append(params, fromDriverCookie(url, value))

		return true
	})

	err := m.client.Network.SetCookies(ctx, network.NewSetCookiesArgs(params))

	if err != nil {
		m.logger.Trace().Err(err).Msg("failed to set cookies")

		return err
	}

	m.logger.Trace().Msg("succeeded to set cookies")

	return nil
}

func (m *Manager) DeleteCookies(ctx context.Context, url string, cookies *drivers.HTTPCookies) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.logger.Trace().Str("url", url).Msg("starting to delete cookies")

	if cookies == nil {
		m.logger.Trace().Msg("nil cookies passed")

		return errors.Wrap(core.ErrMissedArgument, "cookies")
	}

	if cookies.Length() == 0 {
		m.logger.Trace().Msg("no cookies passed")

		return nil
	}

	var err error

	cookies.ForEach(func(value drivers.HTTPCookie, _ values.String) bool {
		m.logger.Trace().Str("name", value.Name).Msg("deleting a cookie")

		err = m.client.Network.DeleteCookies(ctx, fromDriverCookieDelete(url, value))

		if err != nil {
			m.logger.Trace().Err(err).Str("name", value.Name).Msg("failed to delete a cookie")

			return false
		}

		m.logger.Trace().Str("name", value.Name).Msg("succeeded to delete a cookie")

		return true
	})

	return err
}

func (m *Manager) GetHeaders(_ context.Context) (*drivers.HTTPHeaders, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.headers == nil {
		return drivers.NewHTTPHeaders(), nil
	}

	return m.headers.Clone().(*drivers.HTTPHeaders), nil
}

func (m *Manager) SetHeaders(ctx context.Context, headers *drivers.HTTPHeaders) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	return m.setHeadersInternal(ctx, headers)
}

func (m *Manager) setHeadersInternal(ctx context.Context, headers *drivers.HTTPHeaders) error {
	m.logger.Trace().Msg("starting to set headers")

	if headers.Length() == 0 {
		m.logger.Trace().Msg("no headers passed")

		return nil
	}

	m.headers = headers

	m.logger.Trace().Msg("marshaling headers")

	j, err := jettison.MarshalOpts(headers, jettison.NoHTMLEscaping())

	if err != nil {
		m.logger.Trace().Err(err).Msg("failed to marshal headers")

		return errors.Wrap(err, "failed to marshal headers")
	}

	m.logger.Trace().Msg("sending headers to browser")

	err = m.client.Network.SetExtraHTTPHeaders(
		ctx,
		network.NewSetExtraHTTPHeadersArgs(j),
	)

	if err != nil {
		m.logger.Trace().Err(err).Msg("failed to set headers")

		return errors.Wrap(err, "failed to set headers")
	}

	m.logger.Trace().Msg("succeeded to set headers")

	return nil
}

func (m *Manager) GetResponse(_ context.Context, frameID page.FrameID) (drivers.HTTPResponse, error) {
	value, found := m.response.Load(frameID)

	m.logger.Trace().
		Str("frame_id", string(frameID)).
		Bool("found", found).
		Msg("getting frame response")

	if !found {
		return drivers.HTTPResponse{}, core.ErrNotFound
	}

	return *(value.(*drivers.HTTPResponse)), nil
}

func (m *Manager) Navigate(ctx context.Context, url values.String) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if url == "" {
		url = BlankPageURL
	}

	urlStr := url.String()
	m.logger.Trace().Str("url", urlStr).Msg("starting navigation")

	repl, err := m.client.Page.Navigate(ctx, page.NewNavigateArgs(urlStr))

	if err == nil && repl.ErrorText != nil {
		err = errors.New(*repl.ErrorText)
	}

	if err != nil {
		m.logger.Trace().Err(err).Msg("failed starting navigation")

		return err
	}

	m.logger.Trace().Msg("succeeded starting navigation")

	return m.WaitForNavigation(ctx, WaitEventOptions{})
}

func (m *Manager) NavigateForward(ctx context.Context, skip values.Int) (values.Boolean, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.logger.Trace().
		Int64("skip", int64(skip)).
		Msg("starting forward navigation")

	history, err := m.client.Page.GetNavigationHistory(ctx)

	if err != nil {
		m.logger.Trace().
			Err(err).
			Msg("failed to get navigation history")

		return values.False, err
	}

	length := len(history.Entries)
	lastIndex := length - 1

	// nowhere to go forward
	if history.CurrentIndex == lastIndex {
		m.logger.Trace().
			Int("history_entries", length).
			Int("history_current_index", history.CurrentIndex).
			Int("history_last_index", lastIndex).
			Msg("no forward history. nowhere to navigate. done.")

		return values.False, nil
	}

	if skip < 1 {
		skip = 1
	}

	to := int(skip) + history.CurrentIndex

	if to > lastIndex {
		m.logger.Trace().
			Int("history_entries", length).
			Int("history_current_index", history.CurrentIndex).
			Int("history_last_index", lastIndex).
			Int("history_target_index", to).
			Msg("not enough history items. using the edge index")

		to = lastIndex
	}

	entry := history.Entries[to]
	err = m.client.Page.NavigateToHistoryEntry(ctx, page.NewNavigateToHistoryEntryArgs(entry.ID))

	if err != nil {
		m.logger.Trace().
			Int("history_entries", length).
			Int("history_current_index", history.CurrentIndex).
			Int("history_last_index", lastIndex).
			Int("history_target_index", to).
			Err(err).
			Msg("failed to get navigation history entry")

		return values.False, err
	}

	err = m.WaitForNavigation(ctx, WaitEventOptions{})

	if err != nil {
		m.logger.Trace().
			Int("history_entries", length).
			Int("history_current_index", history.CurrentIndex).
			Int("history_last_index", lastIndex).
			Int("history_target_index", to).
			Err(err).
			Msg("failed to wait for navigation completion")

		return values.False, err
	}

	m.logger.Trace().
		Int("history_entries", length).
		Int("history_current_index", history.CurrentIndex).
		Int("history_last_index", lastIndex).
		Int("history_target_index", to).
		Msg("succeeded to wait for navigation completion")

	return values.True, nil
}

func (m *Manager) NavigateBack(ctx context.Context, skip values.Int) (values.Boolean, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.logger.Trace().
		Int64("skip", int64(skip)).
		Msg("starting backward navigation")

	history, err := m.client.Page.GetNavigationHistory(ctx)

	if err != nil {
		m.logger.Trace().Err(err).Msg("failed to get navigation history")

		return values.False, err
	}

	length := len(history.Entries)

	// we are in the beginning
	if history.CurrentIndex == 0 {
		m.logger.Trace().
			Int("history_entries", length).
			Int("history_current_index", history.CurrentIndex).
			Msg("no backward history. nowhere to navigate. done.")

		return values.False, nil
	}

	if skip < 1 {
		skip = 1
	}

	to := history.CurrentIndex - int(skip)

	if to < 0 {
		m.logger.Trace().
			Int("history_entries", length).
			Int("history_current_index", history.CurrentIndex).
			Int("history_target_index", to).
			Msg("not enough history items. using 0 index")

		to = 0
	}

	entry := history.Entries[to]
	err = m.client.Page.NavigateToHistoryEntry(ctx, page.NewNavigateToHistoryEntryArgs(entry.ID))

	if err != nil {
		m.logger.Trace().
			Int("history_entries", length).
			Int("history_current_index", history.CurrentIndex).
			Int("history_target_index", to).
			Err(err).
			Msg("failed to get navigation history entry")

		return values.False, err
	}

	err = m.WaitForNavigation(ctx, WaitEventOptions{})

	if err != nil {
		m.logger.Trace().
			Int("history_entries", length).
			Int("history_current_index", history.CurrentIndex).
			Int("history_target_index", to).
			Err(err).
			Msg("failed to wait for navigation completion")

		return values.False, err
	}

	m.logger.Trace().
		Int("history_entries", length).
		Int("history_current_index", history.CurrentIndex).
		Int("history_target_index", to).
		Msg("succeeded to wait for navigation completion")

	return values.True, nil
}

func (m *Manager) WaitForNavigation(ctx context.Context, opts WaitEventOptions) error {
	stream, err := m.OnNavigation(ctx)
	if err != nil {
		return err
	}

	defer stream.Close(ctx)

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	for evt := range stream.Read(ctx) {
		if err := ctx.Err(); err != nil {
			return err
		}

		if err := evt.Err(); err != nil {
			return nil
		}

		nav := evt.Value().(*NavigationEvent)

		if !isFrameMatched(nav.FrameID, opts.FrameID) || !isURLMatched(nav.URL, opts.URL) {
			continue
		}

		return nil
	}

	return nil
}

func (m *Manager) OnNavigation(ctx context.Context) (rtEvents.Stream, error) {
	m.logger.Trace().Msg("starting to wait for frame navigation event")

	stream1, err := m.client.Page.FrameNavigated(ctx)

	if err != nil {
		m.logger.Trace().Err(err).Msg("failed to open frame navigation event stream")

		return nil, err
	}

	stream2, err := m.client.Page.NavigatedWithinDocument(ctx)

	if err != nil {
		_ = stream1.Close()
		m.logger.Trace().Err(err).Msg("failed to open within document navigation event streams")

		return nil, err
	}

	return newNavigationEventStream(m.logger, m.client, stream1, stream2), nil
}

func (m *Manager) OnRequest(ctx context.Context) (rtEvents.Stream, error) {
	m.logger.Trace().Msg("starting to receive request event")

	stream, err := m.client.Network.RequestWillBeSent(ctx)

	if err != nil {
		m.logger.Trace().Err(err).Msg("failed to open request event stream")

		return nil, err
	}

	m.logger.Trace().Msg("succeeded to receive request event")

	return newRequestWillBeSentStream(m.logger, stream), nil
}

func (m *Manager) OnResponse(ctx context.Context) (rtEvents.Stream, error) {
	m.logger.Trace().Msg("starting to receive response events")

	stream, err := m.client.Network.ResponseReceived(ctx)

	if err != nil {
		m.logger.Trace().Err(err).Msg("failed to open response event stream")

		return nil, err
	}

	m.logger.Trace().Msg("succeeded to receive response events")

	return newResponseReceivedReader(m.logger, m.client, stream), nil
}

func (m *Manager) handleResponse(_ context.Context, message interface{}) (out bool) {
	out = true
	msg, ok := message.(*network.ResponseReceivedReply)

	if !ok {
		return
	}

	// we are interested in documents only
	if msg.Type != network.ResourceTypeDocument {
		return
	}

	if msg.FrameID == nil {
		return
	}

	log := m.logger.With().
		Str("frame_id", string(*msg.FrameID)).
		Str("request_id", string(msg.RequestID)).
		Str("loader_id", string(msg.LoaderID)).
		Float64("timestamp", float64(msg.Timestamp)).
		Str("url", msg.Response.URL).
		Int("status_code", msg.Response.Status).
		Str("status_text", msg.Response.StatusText).
		Logger()

	log.Trace().Msg("received browser response")

	m.response.Store(*msg.FrameID, toDriverResponse(msg.Response, nil))

	log.Trace().Msg("updated frame response information")

	return
}
