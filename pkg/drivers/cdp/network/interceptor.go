package network

import (
	"context"
	"sync"

	"github.com/gobwas/glob"
	"github.com/mafredri/cdp"
	"github.com/mafredri/cdp/protocol/fetch"
	"github.com/mafredri/cdp/protocol/network"
	"github.com/rs/zerolog"

	"github.com/MontFerret/ferret/pkg/drivers"
	"github.com/MontFerret/ferret/pkg/drivers/cdp/events"
	"github.com/MontFerret/ferret/pkg/runtime/logging"
)

type (
	ResourceFilter struct {
		URL          glob.Glob
		ResourceType string
	}

	Interceptor struct {
		mu      sync.RWMutex
		running bool
		logger  zerolog.Logger
		client  *cdp.Client
		filters map[string]*InterceptorFilter
		loop    *events.Loop
	}

	InterceptorFilter struct {
		resources []ResourceFilter
	}

	InterceptorListener func(ctx context.Context, msg *fetch.RequestPausedReply) bool
)

func NewInterceptorFilter(filter *Filter) (*InterceptorFilter, error) {
	interFilter := new(InterceptorFilter)
	interFilter.resources = make([]ResourceFilter, 0, len(filter.Patterns))

	for _, pattern := range filter.Patterns {
		rf := ResourceFilter{
			ResourceType: pattern.Type,
		}

		if pattern.URL != "" {
			p, err := glob.Compile(pattern.URL)

			if err != nil {
				return nil, err
			}

			rf.URL = p
		}

		if rf.ResourceType != "" && rf.URL != nil {
			interFilter.resources = append(interFilter.resources, rf)
		}
	}

	return interFilter, nil
}

func (f *InterceptorFilter) Filter(rt network.ResourceType, req network.Request) bool {
	var result bool

	for _, pattern := range f.resources {
		if pattern.ResourceType != "" && pattern.URL != nil {
			result = string(rt) == pattern.ResourceType && pattern.URL.Match(req.URL)
		} else if pattern.ResourceType != "" {
			result = string(rt) == pattern.ResourceType
		} else if pattern.URL != nil {
			result = pattern.URL.Match(req.URL)
		}

		if result {
			break
		}
	}

	return result
}

func NewInterceptor(logger zerolog.Logger, client *cdp.Client) *Interceptor {
	i := new(Interceptor)
	i.logger = logging.WithName(logger.With(), "network_interceptor").Logger()
	i.client = client
	i.filters = make(map[string]*InterceptorFilter)
	i.loop = events.NewLoop(createRequestPausedStreamFactory(client))
	i.loop.AddListener(requestPausedEvent, events.Always(i.filter))

	return i
}

func (i *Interceptor) IsRunning() bool {
	i.mu.Lock()
	defer i.mu.Unlock()

	return i.running
}

func (i *Interceptor) AddFilter(name string, filter *Filter) error {
	i.mu.Lock()
	defer i.mu.Unlock()

	f, err := NewInterceptorFilter(filter)

	if err != nil {
		return err
	}

	i.filters[name] = f

	return nil
}

func (i *Interceptor) RemoveFilter(name string) {
	i.mu.Lock()
	defer i.mu.Unlock()

	delete(i.filters, name)
}

func (i *Interceptor) AddListener(listener InterceptorListener) events.ListenerID {
	i.mu.Lock()
	defer i.mu.Unlock()

	return i.loop.AddListener(requestPausedEvent, func(ctx context.Context, message interface{}) bool {
		msg, ok := message.(*fetch.RequestPausedReply)

		if !ok {
			return true
		}

		return listener(ctx, msg)
	})
}

func (i *Interceptor) RemoveListener(id events.ListenerID) {
	i.mu.Lock()
	defer i.mu.Unlock()

	i.loop.RemoveListener(requestPausedEvent, id)
}

func (i *Interceptor) Run(ctx context.Context) error {
	i.mu.Lock()
	defer i.mu.Unlock()

	if i.running {
		return nil
	}

	err := i.client.Fetch.Enable(ctx, fetch.NewEnableArgs())
	i.running = err == nil

	if err != nil {
		return err
	}

	if err := i.loop.Run(ctx); err != nil {
		if e := i.client.Fetch.Disable(ctx); e != nil {
			i.logger.Err(err).Msg("failed to disable fetch")
		}

		i.running = false

		return err
	}

	go func() {
		<-ctx.Done()

		nested, cancel := context.WithTimeout(context.Background(), drivers.DefaultWaitTimeout)
		defer cancel()

		i.stop(nested)
	}()

	return nil
}

func (i *Interceptor) stop(ctx context.Context) {
	err := i.client.Fetch.Disable(ctx)
	i.running = false

	if err != nil {
		i.logger.Err(err).Msg("failed to stop interceptor")
	}
}

func (i *Interceptor) filter(ctx context.Context, message interface{}) {
	i.mu.Lock()
	defer i.mu.Unlock()

	msg, ok := message.(*fetch.RequestPausedReply)

	if !ok {
		return
	}

	log := i.logger.With().
		Str("request_id", string(msg.RequestID)).
		Str("frame_id", string(msg.FrameID)).
		Str("resource_type", string(msg.ResourceType)).
		Str("url", msg.Request.URL).
		Logger()

	log.Trace().Msg("trying to block resource loading")

	var reject bool

	for _, filter := range i.filters {
		reject = filter.Filter(msg.ResourceType, msg.Request)

		if reject {
			break
		}
	}

	if !reject {
		err := i.client.Fetch.ContinueRequest(ctx, fetch.NewContinueRequestArgs(msg.RequestID))

		if err != nil {
			i.logger.Err(err).Msg("failed to allow resource loading")
		}

		log.Trace().Msg("succeeded to allow resource loading")

		return
	}

	err := i.client.Fetch.FailRequest(ctx, fetch.NewFailRequestArgs(msg.RequestID, network.ErrorReasonBlockedByClient))

	if err != nil {
		log.Trace().Err(err).Msg("failed to block resource loading")
	}

	log.Trace().Msg("succeeded to block resource loading")
}
