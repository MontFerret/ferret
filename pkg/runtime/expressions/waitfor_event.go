package expressions

import (
	"context"
	"github.com/MontFerret/ferret/pkg/runtime/events"
	"github.com/MontFerret/ferret/pkg/runtime/values/types"
	"time"

	"github.com/pkg/errors"

	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/expressions/literals"
	"github.com/MontFerret/ferret/pkg/runtime/values"
)

const DefaultWaitTimeout = 5000

type WaitForEventExpression struct {
	src         core.SourceMap
	eventName   core.Expression
	eventSource core.Expression
	options     core.Expression
	timeout     core.Expression
}

func NewWaitForEventExpression(
	src core.SourceMap,
	eventName core.Expression,
	eventSource core.Expression,
	options core.Expression,
	timeout core.Expression,
) (*WaitForEventExpression, error) {
	if eventName == nil {
		return nil, core.Error(core.ErrInvalidArgument, "event name")
	}

	if eventSource == nil {
		return nil, core.Error(core.ErrMissedArgument, "event source")
	}

	if timeout == nil {
		timeout = literals.NewIntLiteral(DefaultWaitTimeout)
	}

	return &WaitForEventExpression{
		src:         src,
		eventName:   eventName,
		eventSource: eventSource,
		options:     options,
		timeout:     timeout,
	}, nil
}

func (e *WaitForEventExpression) Exec(ctx context.Context, scope *core.Scope) (core.Value, error) {
	eventName, err := e.getEventName(ctx, scope)

	if err != nil {
		return values.None, core.SourceError(e.src, errors.Wrap(err, "unable to calculate event name"))
	}

	eventSource, err := e.eventSource.Exec(ctx, scope)

	if err != nil {
		return values.None, core.SourceError(e.src, err)
	}

	observable, ok := eventSource.(events.Observable)

	if !ok {
		return values.None, core.TypeError(eventSource.Type(), core.NewType("Observable"))
	}

	ctx, cancel := context.WithCancel(ctx)

	defer cancel()

	opts, err := e.getOptions(ctx, scope)

	if err != nil {
		return values.None, core.SourceError(e.src, err)
	}

	ch := observable.Subscribe(ctx, eventName, opts)

	timeout, err := e.getTimeout(ctx, scope)

	if err != nil {
		return values.None, core.SourceError(e.src, errors.Wrap(err, "failed to calculate timeout value"))
	}

	timer := time.After(timeout * time.Millisecond)

	select {
	case evt := <-ch:
		if evt.Err != nil {
			return values.None, core.SourceError(e.src, evt.Err)
		}

		return evt.Data, nil
	case <-timer:
		return values.None, core.SourceError(e.src, core.ErrTimeout)
	}
}

func (e *WaitForEventExpression) getEventName(ctx context.Context, scope *core.Scope) (string, error) {
	eventName, err := e.eventName.Exec(ctx, scope)

	if err != nil {
		return "", err
	}

	return eventName.String(), nil
}

func (e *WaitForEventExpression) getOptions(ctx context.Context, scope *core.Scope) (*values.Object, error) {
	if e.options == nil {
		return nil, nil
	}

	options, err := e.options.Exec(ctx, scope)

	if err != nil {
		return nil, err
	}

	if err := core.ValidateType(options, types.Object); err != nil {
		return nil, err
	}

	return options.(*values.Object), nil
}

func (e *WaitForEventExpression) getTimeout(ctx context.Context, scope *core.Scope) (time.Duration, error) {
	timeoutValue, err := e.timeout.Exec(ctx, scope)

	if err != nil {
		return 0, err
	}

	timeoutInt := values.ToIntDefault(timeoutValue, DefaultWaitTimeout)

	return time.Duration(timeoutInt), nil
}
