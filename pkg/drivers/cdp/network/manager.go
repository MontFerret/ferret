package network

import (
	"context"
	"encoding/json"
	"io"
	"regexp"
	"sync"

	"github.com/mafredri/cdp"
	"github.com/mafredri/cdp/protocol/fetch"
	"github.com/mafredri/cdp/protocol/network"
	"github.com/mafredri/cdp/protocol/page"
	"github.com/mafredri/cdp/rpcc"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/wI2L/jettison"

	"github.com/MontFerret/ferret/pkg/drivers"
	"github.com/MontFerret/ferret/pkg/drivers/cdp/eval"
	"github.com/MontFerret/ferret/pkg/drivers/cdp/events"
	"github.com/MontFerret/ferret/pkg/drivers/cdp/templates"
	"github.com/MontFerret/ferret/pkg/drivers/common"
	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
)

const BlankPageURL = "about:blank"

type (
	FrameLoadedListener = func(ctx context.Context, frame page.Frame)

	Manager struct {
		mu                 sync.RWMutex
		logger             *zerolog.Logger
		client             *cdp.Client
		headers            *drivers.HTTPHeaders
		eventLoop          *events.Loop
		cancel             context.CancelFunc
		responseListenerID events.ListenerID
		filterListenerID   events.ListenerID
		response           *sync.Map
	}
)

func New(
	ctx context.Context,
	logger *zerolog.Logger,
	client *cdp.Client,
	options Options,
) (*Manager, error) {
	ctx, cancel := context.WithCancel(ctx)

	m := new(Manager)
	m.logger = logger
	m.client = client
	m.headers = drivers.NewHTTPHeaders()
	m.eventLoop = events.NewLoop()
	m.cancel = cancel
	m.response = new(sync.Map)

	if options.Cookies != nil && len(options.Cookies) > 0 {
		for url, cookies := range options.Cookies {
			if err := m.setCookiesInternal(ctx, url, cookies); err != nil {
				return nil, err
			}
		}
	}

	if options.Headers != nil && options.Headers.Length() > 0 {
		if err := m.setHeadersInternal(ctx, options.Headers); err != nil {
			return nil, err
		}
	}

	var err error

	closers := make([]io.Closer, 0, 10)

	defer func() {
		if err != nil {
			common.CloseAll(logger, closers, "failed to close a DOM event stream")
		}
	}()

	frameNavigatedStream, err := m.client.Page.FrameNavigated(ctx)

	if err != nil {
		return nil, err
	}

	responseReceivedStream, err := m.client.Network.ResponseReceived(ctx)

	if err != nil {
		return nil, err
	}

	m.eventLoop.AddSource(events.NewSource(eventFrameLoad, frameNavigatedStream, func(stream rpcc.Stream) (interface{}, error) {
		return stream.(page.FrameNavigatedClient).Recv()
	}))

	m.eventLoop.AddSource(events.NewSource(responseReceived, responseReceivedStream, func(stream rpcc.Stream) (interface{}, error) {
		return stream.(network.ResponseReceivedClient).Recv()
	}))

	m.responseListenerID = m.eventLoop.AddListener(responseReceived, m.onResponse)

	if options.Filter != nil && len(options.Filter.Patterns) > 0 {
		el2 := events.NewLoop()

		err = m.client.Fetch.Enable(ctx, toFetchArgs(options.Filter.Patterns))

		if err != nil {
			return nil, err
		}

		requestPausedStream, err := m.client.Fetch.RequestPaused(ctx)

		if err != nil {
			return nil, err
		}

		el2.AddSource(events.NewSource(requestPaused, requestPausedStream, func(stream rpcc.Stream) (interface{}, error) {
			return stream.(fetch.RequestPausedClient).Recv()
		}))

		m.filterListenerID = el2.AddListener(requestPaused, m.onRequestPaused)

		// run in a separate loop in order to get higher priority
		// TODO: Consider adding support of event priorities to EventLoop
		el2.Run(ctx)
	}

	m.eventLoop.Run(ctx)

	return m, nil
}

func (m *Manager) Close() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.logger.Trace().Msg("closing")

	if m.cancel != nil {
		m.cancel()
		m.cancel = nil
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
	log := m.logger.With().Str("url", url).Logger()
	log.Trace().Msg("starting to set cookies")

	if cookies == nil {
		log.Trace().Msg("nil cookies passed")

		return errors.Wrap(core.ErrMissedArgument, "cookies")
	}

	if cookies.Length() == 0 {
		log.Trace().Msg("no cookies passed")

		return nil
	}

	params := make([]network.CookieParam, 0, cookies.Length())

	cookies.ForEach(func(value drivers.HTTPCookie, _ values.String) bool {
		params = append(params, fromDriverCookie(url, value))

		return true
	})

	err := m.client.Network.SetCookies(ctx, network.NewSetCookiesArgs(params))

	if err != nil {
		log.Trace().Err(err).Msg("failed to set cookies")

		return err
	}

	log.Trace().Msg("succeeded to set cookies")

	return nil
}

func (m *Manager) DeleteCookies(ctx context.Context, url string, cookies *drivers.HTTPCookies) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	log := m.logger.With().Str("url", url).Logger()
	log.Trace().Msg("starting to delete cookies")

	if cookies == nil {
		log.Trace().Msg("nil cookies passed")

		return errors.Wrap(core.ErrMissedArgument, "cookies")
	}

	if cookies.Length() == 0 {
		log.Trace().Msg("no cookies passed")

		return nil
	}

	var err error

	cookies.ForEach(func(value drivers.HTTPCookie, _ values.String) bool {
		log.Trace().Str("name", value.Name).Msg("deleting a cookie")

		err = m.client.Network.DeleteCookies(ctx, fromDriverCookieDelete(url, value))

		if err != nil {
			log.Trace().Err(err).Str("name", value.Name).Msg("failed to delete a cookie")

			return false
		}

		log.Trace().Str("name", value.Name).Msg("succeeded to delete a cookie")

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
	log := m.logger
	log.Trace().Msg("starting to set headers")

	if headers.Length() == 0 {
		log.Trace().Msg("no headers passed")

		return nil
	}

	m.headers = headers

	log.Trace().Msg("marshaling headers")

	j, err := jettison.MarshalOpts(headers, jettison.NoHTMLEscaping())

	if err != nil {
		log.Trace().Err(err).Msg("failed to marshal headers")

		return errors.Wrap(err, "failed to marshal headers")
	}

	log.Trace().Msg("sending headers to browser")

	err = m.client.Network.SetExtraHTTPHeaders(
		ctx,
		network.NewSetExtraHTTPHeadersArgs(j),
	)

	if err != nil {
		log.Trace().Err(err).Msg("failed to set headers")

		return errors.Wrap(err, "failed to set headers")
	}

	log.Trace().Msg("succeeded to set headers")

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

	return value.(drivers.HTTPResponse), nil
}

func (m *Manager) Navigate(ctx context.Context, url values.String) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if url == "" {
		url = BlankPageURL
	}

	urlStr := url.String()
	log := m.logger.With().Str("url", urlStr).Logger()
	log.Trace().Msg("starting navigation")

	repl, err := m.client.Page.Navigate(ctx, page.NewNavigateArgs(urlStr))

	if err == nil && repl.ErrorText != nil {
		err = errors.New(*repl.ErrorText)
	}

	if err != nil {
		log.Trace().Err(err).Msg("failed starting navigation")

		return err
	}

	log.Trace().Msg("succeeded starting navigation")

	return m.WaitForNavigation(ctx, nil)
}

func (m *Manager) NavigateForward(ctx context.Context, skip values.Int) (values.Boolean, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	log := m.logger.With().
		Int64("skip", int64(skip)).
		Logger()

	history, err := m.client.Page.GetNavigationHistory(ctx)

	if err != nil {
		log.Trace().
			Err(err).
			Msg("failed to get navigation history")

		return values.False, err
	}

	length := len(history.Entries)
	lastIndex := length - 1

	// nowhere to go forward
	if history.CurrentIndex == lastIndex {
		log.Trace().
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
		log.Trace().
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
		log.Trace().
			Int("history_entries", length).
			Int("history_current_index", history.CurrentIndex).
			Int("history_last_index", lastIndex).
			Int("history_target_index", to).
			Err(err).
			Msg("failed to get navigation history entry")

		return values.False, err
	}

	err = m.WaitForNavigation(ctx, nil)

	if err != nil {
		log.Trace().
			Int("history_entries", length).
			Int("history_current_index", history.CurrentIndex).
			Int("history_last_index", lastIndex).
			Int("history_target_index", to).
			Err(err).
			Msg("failed to wait for navigation completion")

		return values.False, err
	}

	log.Trace().
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

	log := m.logger.With().Int64("skip", int64(skip)).Logger()

	log.Trace().Msg("starting backward navigation")

	history, err := m.client.Page.GetNavigationHistory(ctx)

	if err != nil {
		log.Trace().Err(err).Msg("failed to get navigation history")

		return values.False, err
	}

	length := len(history.Entries)

	// we are in the beginning
	if history.CurrentIndex == 0 {
		log.Trace().
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
		log.Trace().
			Int("history_entries", length).
			Int("history_current_index", history.CurrentIndex).
			Int("history_target_index", to).
			Msg("not enough history items. using 0 index")

		to = 0
	}

	entry := history.Entries[to]
	err = m.client.Page.NavigateToHistoryEntry(ctx, page.NewNavigateToHistoryEntryArgs(entry.ID))

	if err != nil {
		log.Trace().
			Int("history_entries", length).
			Int("history_current_index", history.CurrentIndex).
			Int("history_target_index", to).
			Err(err).
			Msg("failed to get navigation history entry")

		return values.False, err
	}

	err = m.WaitForNavigation(ctx, nil)

	if err != nil {
		log.Trace().
			Int("history_entries", length).
			Int("history_current_index", history.CurrentIndex).
			Int("history_target_index", to).
			Err(err).
			Msg("failed to wait for navigation completion")

		return values.False, err
	}

	log.Trace().
		Int("history_entries", length).
		Int("history_current_index", history.CurrentIndex).
		Int("history_target_index", to).
		Msg("succeeded to wait for navigation completion")

	return values.True, nil
}

func (m *Manager) WaitForNavigation(ctx context.Context, pattern *regexp.Regexp) error {
	return m.WaitForFrameNavigation(ctx, "", pattern)
}

func (m *Manager) WaitForFrameNavigation(ctx context.Context, frameID page.FrameID, urlPattern *regexp.Regexp) error {
	onEvent := make(chan struct{})

	var urlPatternStr string

	if urlPattern != nil {
		urlPatternStr = urlPattern.String()
	}

	m.logger.Trace().
		Str("fame_id", string(frameID)).
		Str("url_pattern", urlPatternStr).
		Msg("started waiting for frame navigation event")

	m.eventLoop.AddListener(eventFrameLoad, func(_ context.Context, message interface{}) bool {
		repl := message.(*page.FrameNavigatedReply)

		var matched bool

		// if frameID is empty string or equals to the current one
		if len(frameID) == 0 || repl.Frame.ID == frameID {
			// if a URL pattern is provided
			if urlPattern != nil {
				matched = urlPattern.Match([]byte(repl.Frame.URL))
			} else {
				// otherwise just notify
				matched = true
			}
		}

		log := m.logger.With().
			Str("fame_id", string(frameID)).
			Str("event_fame_id", string(repl.Frame.ID)).
			Str("event_fame_url", repl.Frame.URL).
			Str("url_pattern", urlPatternStr).
			Bool("is_matched", matched).
			Logger()

		log.Trace().Msg("received framed navigation event")

		if matched {
			if ctx.Err() == nil {
				log.Trace().Msg("creating frame execution context")

				ec, err := eval.NewExecutionContextFrom(ctx, m.client, repl.Frame)

				if err != nil {
					log.Trace().Err(err).Msg("failed to create frame execution context")

					close(onEvent)

					return false
				}

				log.Trace().Err(err).Msg("starting polling DOM ready event")

				_, err = events.NewEvalWaitTask(
					ec,
					templates.DOMReady(),
					events.DefaultPolling,
				).Run(ctx)

				if err != nil {
					log.Trace().Err(err).Msg("failed to poll DOM ready event")

					close(onEvent)

					return false
				}

				log.Trace().Msg("DOM is ready. navigation completed")

				onEvent <- struct{}{}
				close(onEvent)
			}
		}

		// if not matched - continue listening
		return !matched
	})

	select {
	case <-onEvent:
		return nil
	case <-ctx.Done():
		return core.ErrTimeout
	}
}

func (m *Manager) AddFrameLoadedListener(listener FrameLoadedListener) events.ListenerID {
	return m.eventLoop.AddListener(eventFrameLoad, func(ctx context.Context, message interface{}) bool {
		repl := message.(*page.FrameNavigatedReply)

		listener(ctx, repl.Frame)

		return true
	})
}

func (m *Manager) RemoveFrameLoadedListener(id events.ListenerID) {
	m.eventLoop.RemoveListener(eventFrameLoad, id)
}

func (m *Manager) onResponse(_ context.Context, message interface{}) (out bool) {
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

	response := drivers.HTTPResponse{
		URL:          msg.Response.URL,
		StatusCode:   msg.Response.Status,
		Status:       msg.Response.StatusText,
		Headers:      drivers.NewHTTPHeaders(),
		ResponseTime: float64(msg.Response.ResponseTime),
	}

	deserialized := make(map[string]string)

	if len(msg.Response.Headers) > 0 {
		err := json.Unmarshal(msg.Response.Headers, &deserialized)

		if err != nil {
			log.Trace().Err(err).Msg("failed to deserialize response headers")
		}
	}

	for key, value := range deserialized {
		response.Headers.Set(key, value)
	}

	m.response.Store(*msg.FrameID, response)

	log.Trace().Msg("updated frame response data")

	return
}

func (m *Manager) onRequestPaused(ctx context.Context, message interface{}) (out bool) {
	out = true
	msg, ok := message.(*fetch.RequestPausedReply)

	if !ok {
		return
	}

	log := m.logger.With().
		Str("request_id", string(msg.RequestID)).
		Str("frame_id", string(msg.FrameID)).
		Str("resource_type", string(msg.ResourceType)).
		Str("url", msg.Request.URL).
		Logger()

	log.Trace().Msg("trying to block resource loading")

	err := m.client.Fetch.FailRequest(ctx, &fetch.FailRequestArgs{
		RequestID:   msg.RequestID,
		ErrorReason: network.ErrorReasonBlockedByClient,
	})

	if err != nil {
		log.Trace().Err(err).Msg("failed to block resource loading")
	}

	log.Trace().Msg("succeeded to block resource loading")

	return
}
