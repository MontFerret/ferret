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
	"github.com/MontFerret/ferret/pkg/drivers/cdp/eval"
	"github.com/MontFerret/ferret/pkg/drivers/cdp/events"
	"github.com/MontFerret/ferret/pkg/drivers/cdp/templates"
	"github.com/MontFerret/ferret/pkg/runtime/core"
	evts "github.com/MontFerret/ferret/pkg/runtime/events"
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
		stopLoop    context.CancelFunc
		interceptor *Interceptor
		cancel      context.CancelFunc
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
	m.cancel = cancel
	m.response = new(sync.Map)

	var err error

	defer func() {
		if err != nil {
			m.cancel()

			if m.stopLoop != nil {
				m.stopLoop()
			}

			if m.interceptor != nil && m.interceptor.IsRunning() {
				if err := m.interceptor.Stop(ctx); err != nil {
					m.logger.Err(err).Msg("failed to stop interceptor")
				}
			}
		}
	}()

	m.loop = events.NewLoop(
		createFrameLoadStreamFactory(client),
		//createResponseReceivedStreamFactory(client),
	)

	m.loop.AddListener(responseReceivedEvent, m.handleResponse)

	if options.Filter != nil && len(options.Filter.Patterns) > 0 {
		m.interceptor = NewInterceptor(logger, client)

		if err := m.interceptor.AddFilter("resources", options.Filter); err != nil {
			return nil, err
		}

		if err = m.interceptor.Start(ctx); err != nil {
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

	cancel, err = m.loop.Run(ctx)

	if err != nil {
		return nil, err
	}

	m.stopLoop = cancel

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

	if m.stopLoop != nil {
		m.stopLoop()
	}

	if m.interceptor != nil && m.interceptor.IsRunning() {
		ctx, cancel := context.WithTimeout(context.Background(), drivers.DefaultWaitTimeout)
		defer cancel()

		if err := m.interceptor.Stop(ctx); err != nil {
			m.logger.Err(err).Msg("failed to stop interceptor")
		}
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
		Msg("getting frame responseReceivedEvent")

	if !found {
		return drivers.HTTPResponse{}, core.ErrNotFound
	}

	resp := value.(*drivers.HTTPResponse)

	return *resp, nil
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

	return m.WaitForNavigation(ctx, EventOptions{})
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

	err = m.WaitForNavigation(ctx, EventOptions{})

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

	err = m.WaitForNavigation(ctx, EventOptions{})

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

func (m *Manager) WaitForNavigation(ctx context.Context, opts EventOptions) error {
	onEvent, err := m.OnFrameNavigation(ctx, opts)

	if err != nil {
		return err
	}

	select {
	case <-onEvent:
		return nil
	case <-ctx.Done():
		var urlPatternStr string

		if opts.URL != nil {
			urlPatternStr = opts.URL.String()
		}

		m.logger.Trace().
			Err(core.ErrTimeout).
			Str("fame_id", opts.FrameID).
			Str("url_pattern", urlPatternStr).
			Msg("navigation has been interrupted")

		return core.ErrTimeout
	}
}

func (m *Manager) OnFrameNavigation(ctx context.Context, opts EventOptions) (<-chan evts.Event, error) {
	ch := make(chan evts.Event)

	var urlPatternStr string

	if opts.URL != nil {
		urlPatternStr = opts.URL.String()
	}

	m.logger.Trace().
		Str("fame_id", opts.FrameID).
		Str("url_pattern", urlPatternStr).
		Msg("starting to wait for frame navigation event")

	m.loop.AddListener(frameNavigatedEvent, func(_ context.Context, message interface{}) bool {
		repl := message.(*page.FrameNavigatedReply)

		log := m.logger.With().
			Str("fame_id", opts.FrameID).
			Str("event_fame_id", string(repl.Frame.ID)).
			Str("event_fame_url", repl.Frame.URL).
			Str("url_pattern", urlPatternStr).
			Logger()

		log.Trace().Msg("received framed navigation event")

		if !isFrameMatched(string(repl.Frame.ID), opts.FrameID) || !isURLMatched(repl.Frame.URL, opts.URL) {
			log.Trace().Msg("frame does not match")

			return true
		}

		defer close(ch)

		log.Trace().Msg("frame does match")

		if ctx.Err() != nil {
			// terminated
			return false
		}

		log.Trace().Msg("creating frame execution context")

		ec, err := eval.Create(ctx, m.logger, m.client, repl.Frame.ID)

		if err != nil {
			log.Trace().Err(err).Msg("failed to create frame execution context")

			ch <- evts.Event{
				Err: err,
			}

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

			ch <- evts.Event{
				Err: err,
			}

			return false
		}

		log.Trace().Msg("DOM is ready")

		ch <- evts.Event{
			Data: values.NewString(repl.Frame.URL),
			Err:  nil,
		}

		m.logger.Trace().
			Str("fame_id", opts.FrameID).
			Str("url_pattern", urlPatternStr).
			Msg("navigation has completed")

		return false
	})

	return ch, nil
}

func (m *Manager) OnRequest(ctx context.Context, opts EventOptions) (<-chan evts.Event, error) {
	var urlPatternStr string

	if opts.URL != nil {
		urlPatternStr = opts.URL.String()
	}

	m.logger.Trace().
		Str("fame_id", opts.FrameID).
		Str("url_pattern", urlPatternStr).
		Msg("starting to wait for request event")

	l := events.NewLoop(createBeforeRequestStreamFactory(m.client))

	stop, err := l.Run(ctx)

	if err != nil {
		return nil, err
	}

	ch := make(chan evts.Event)

	l.AddListener(beforeRequestEvent, func(ctx context.Context, message interface{}) bool {
		msg, ok := message.(*network.RequestWillBeSentReply)

		if !ok {
			return true
		}

		var eventFrameID page.FrameID

		if msg.FrameID != nil {
			eventFrameID = *msg.FrameID
		}

		log := m.logger.With().
			Str("url", msg.Request.URL).
			Str("document_url", msg.DocumentURL).
			Str("resource_type", string(msg.Type)).
			Str("request_id", string(msg.RequestID)).
			Str("frame_id", string(eventFrameID)).
			Logger()

		log.Trace().
			Interface("data", msg.Request).
			Msg("trying to match an intercepted request...")

		if !isFrameMatched(string(eventFrameID), opts.FrameID) || !isURLMatched(msg.Request.URL, opts.URL) {
			log.Trace().Msg("request does not match")

			return true
		}

		defer func() {
			stop()
			close(ch)
		}()

		log.Trace().Msg("request does match")

		ch <- evts.Event{
			Data: toDriverRequest(msg.Request),
		}

		// exit
		return false
	})

	return ch, nil
}

func (m *Manager) OnResponse(ctx context.Context, opts EventOptions) (<-chan evts.Event, error) {
	var urlPatternStr string

	if opts.URL != nil {
		urlPatternStr = opts.URL.String()
	}

	m.logger.Trace().
		Str("fame_id", opts.FrameID).
		Str("url_pattern", urlPatternStr).
		Msg("starting to wait for response event")

	l := events.NewLoop(createResponseReceivedStreamFactory(m.client))

	stop, err := l.Run(ctx)

	if err != nil {
		return nil, err
	}

	ch := make(chan evts.Event)

	l.AddListener(responseReceivedEvent, func(ctx context.Context, message interface{}) bool {
		msg, ok := message.(*network.ResponseReceivedReply)

		if !ok {
			return true
		}

		var eventFrameID page.FrameID

		if msg.FrameID != nil {
			eventFrameID = *msg.FrameID
		}

		log := m.logger.With().
			Str("url", msg.Response.URL).
			Str("resource_type", string(msg.Type)).
			Str("frame_id", string(eventFrameID)).
			Str("request_id", string(msg.RequestID)).
			Logger()

		log.Trace().
			Interface("data", msg.Response).
			Msg("trying to match an intercepted response...")

		if !isFrameMatched(string(eventFrameID), opts.FrameID) || !isURLMatched(msg.Response.URL, opts.URL) {
			log.Trace().Msg("response does not match")

			return true
		}

		defer func() {
			stop()
			close(ch)
		}()

		log.Trace().Msg("response does match")

		ch <- evts.Event{
			Data: toDriverResponse(msg.Response),
		}

		// exit
		return false
	})

	return ch, nil
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

	m.response.Store(*msg.FrameID, toDriverResponse(msg.Response))

	log.Trace().Msg("updated frame responseReceivedEvent information")

	return
}
