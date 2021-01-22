package expressions

import (
	"context"
	"time"

	"github.com/pkg/errors"

	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/expressions/literals"
	"github.com/MontFerret/ferret/pkg/runtime/values"
)

const DefaultWaitTimeout = 500

type WaitForEventExpression struct {
	src         core.SourceMap
	eventName   core.Expression
	eventSource core.Expression
	timeout     core.Expression
}

func NewWaitForEventExpression(
	src core.SourceMap,
	eventName core.Expression,
	eventSource core.Expression,
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

	observable, ok := eventSource.(core.Observable)

	if !ok {
		return values.None, core.TypeError(eventSource.Type(), core.NewType("Observable"))
	}

	ctx, cancel := context.WithCancel(ctx)

	defer cancel()

	ch, err := observable.Subscribe(ctx, eventName)

	if err != nil {
		return values.None, core.SourceError(e.src, err)
	}

	timeout, err := e.getTimeout(ctx, scope)

	if err !=nil {
		return values.None, core.SourceError(e.src, errors.Wrap(err, "failed to calculate timeout value"))
	}

	timer := time.After(timeout * time.Millisecond)

	select {
	case <-ch:
		return values.None, nil
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

func (e *WaitForEventExpression) getTimeout(ctx context.Context, scope *core.Scope) (time.Duration, error) {
	timeoutValue, err := e.timeout.Exec(ctx, scope)

	if err != nil {
		return 0, err
	}

	timeoutInt := values.ToIntDefault(timeoutValue, DefaultWaitTimeout)

	return time.Duration(timeoutInt), nil
}