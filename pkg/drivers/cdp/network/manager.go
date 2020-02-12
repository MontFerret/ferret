package network

import (
	"context"
	"encoding/json"
	"io"
	"regexp"
	"sync"

	"github.com/mafredri/cdp"
	"github.com/mafredri/cdp/protocol/network"
	"github.com/mafredri/cdp/protocol/page"
	"github.com/mafredri/cdp/rpcc"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"

	"github.com/MontFerret/ferret/pkg/drivers"
	"github.com/MontFerret/ferret/pkg/drivers/cdp/events"
	"github.com/MontFerret/ferret/pkg/drivers/common"
	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
)

const BlankPageURL = "about:blank"

type (
	FrameLoadedListener = func(ctx context.Context, frame page.Frame)

	Manager struct {
		mu                 sync.Mutex
		logger             *zerolog.Logger
		client             *cdp.Client
		headers            drivers.HTTPHeaders
		eventLoop          *events.Loop
		cancel             context.CancelFunc
		responseListenerID events.ListenerID
		response           *sync.Map
	}
)

func New(
	logger *zerolog.Logger,
	client *cdp.Client,
	eventLoop *events.Loop,
) (*Manager, error) {
	ctx, cancel := context.WithCancel(context.Background())

	m := new(Manager)
	m.logger = logger
	m.client = client
	m.headers = make(drivers.HTTPHeaders)
	m.eventLoop = eventLoop
	m.cancel = cancel
	m.response = new(sync.Map)

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

	return m, nil
}

func (m *Manager) Close() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.cancel != nil {
		m.cancel()
		m.cancel = nil
	}

	return nil
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
	m.mu.Lock()
	defer m.mu.Unlock()

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
	m.mu.Lock()
	defer m.mu.Unlock()

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
	m.mu.Lock()
	defer m.mu.Unlock()

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

func (m *Manager) GetResponse(_ context.Context, frameID page.FrameID) (drivers.HTTPResponse, error) {
	value, found := m.response.Load(frameID)

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

	repl, err := m.client.Page.Navigate(ctx, page.NewNavigateArgs(urlStr))

	if err != nil {
		return err
	}

	if repl.ErrorText != nil {
		return errors.New(*repl.ErrorText)
	}

	return m.WaitForNavigation(ctx, nil)
}

func (m *Manager) NavigateForward(ctx context.Context, skip values.Int) (values.Boolean, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	history, err := m.client.Page.GetNavigationHistory(ctx)

	if err != nil {
		return values.False, err
	}

	length := len(history.Entries)
	lastIndex := length - 1

	// nowhere to go forward
	if history.CurrentIndex == lastIndex {
		return values.False, nil
	}

	if skip < 1 {
		skip = 1
	}

	to := int(skip) + history.CurrentIndex

	if to > lastIndex {
		// TODO: Return error?
		return values.False, nil
	}

	entry := history.Entries[to]
	err = m.client.Page.NavigateToHistoryEntry(ctx, page.NewNavigateToHistoryEntryArgs(entry.ID))

	if err != nil {
		return values.False, err
	}

	err = m.WaitForNavigation(ctx, nil)

	if err != nil {
		return values.False, err
	}

	return values.True, nil
}

func (m *Manager) NavigateBack(ctx context.Context, skip values.Int) (values.Boolean, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	history, err := m.client.Page.GetNavigationHistory(ctx)

	if err != nil {
		return values.False, err
	}

	// we are in the beginning
	if history.CurrentIndex == 0 {
		return values.False, nil
	}

	if skip < 1 {
		skip = 1
	}

	to := history.CurrentIndex - int(skip)

	if to < 0 {
		// TODO: Return error?
		return values.False, nil
	}

	entry := history.Entries[to]
	err = m.client.Page.NavigateToHistoryEntry(ctx, page.NewNavigateToHistoryEntryArgs(entry.ID))

	if err != nil {
		return values.False, err
	}

	err = m.WaitForNavigation(ctx, nil)

	if err != nil {
		return values.False, err
	}

	return values.True, nil
}

func (m *Manager) WaitForNavigation(ctx context.Context, pattern *regexp.Regexp) error {
	return m.WaitForFrameNavigation(ctx, "", pattern)
}

func (m *Manager) WaitForFrameNavigation(ctx context.Context, frameID page.FrameID, urlPattern *regexp.Regexp) error {
	onEvent := make(chan struct{})

	defer func() {
		close(onEvent)
	}()

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

		if matched {
			if ctx.Err() == nil {
				onEvent <- struct{}{}
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

	response := drivers.HTTPResponse{
		StatusCode: msg.Response.Status,
		Status:     msg.Response.StatusText,
		Headers:    make(drivers.HTTPHeaders),
	}

	deserialized := make(map[string]string)

	if len(msg.Response.Headers) > 0 {
		err := json.Unmarshal(msg.Response.Headers, &deserialized)

		if err != nil {
			m.logger.Error().Err(err).Msg("failed to deserialize response headers")
		}
	}

	for key, value := range deserialized {
		response.Headers.Set(key, value)
	}

	m.response.Store(*msg.FrameID, response)

	return
}
