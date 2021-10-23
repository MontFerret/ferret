package expressions

import (
	"context"
	"time"

	"github.com/pkg/errors"

	"github.com/MontFerret/ferret/pkg/runtime/collections"
	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/events"
	"github.com/MontFerret/ferret/pkg/runtime/expressions/clauses"
	"github.com/MontFerret/ferret/pkg/runtime/expressions/literals"
	"github.com/MontFerret/ferret/pkg/runtime/values"
	"github.com/MontFerret/ferret/pkg/runtime/values/types"
)

const DefaultWaitTimeout = 5000

type WaitForEventExpression struct {
	src            core.SourceMap
	eventName      core.Expression
	eventSource    core.Expression
	options        core.Expression
	timeout        core.Expression
	filterSrc      core.SourceMap
	filter         core.Expression
	filterVariable string
}

func NewWaitForEventExpression(
	src core.SourceMap,
	eventName core.Expression,
	eventSource core.Expression,
) (*WaitForEventExpression, error) {
	if eventName == nil {
		return nil, core.Error(core.ErrInvalidArgument, "event name")
	}

	if eventSource == nil {
		return nil, core.Error(core.ErrMissedArgument, "event source")
	}

	return &WaitForEventExpression{
		src:         src,
		eventName:   eventName,
		eventSource: eventSource,
		timeout:     literals.NewIntLiteral(DefaultWaitTimeout),
	}, nil
}

func (e *WaitForEventExpression) SetOptions(options core.Expression) error {
	if options == nil {
		return core.ErrInvalidArgument
	}

	e.options = options

	return nil
}

func (e *WaitForEventExpression) SetTimeout(timeout core.Expression) error {
	if timeout == nil {
		return core.ErrInvalidArgument
	}

	e.timeout = timeout

	return nil
}

func (e *WaitForEventExpression) SetFilter(src core.SourceMap, variable string, exp core.Expression) error {
	if variable == "" {
		return core.ErrInvalidArgument
	}

	if exp == nil {
		return core.ErrInvalidArgument
	}

	e.filterSrc = src
	e.filterVariable = variable
	e.filter = exp

	return nil
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

	opts, err := e.getOptions(ctx, scope)

	if err != nil {
		return values.None, core.SourceError(e.src, err)
	}

	timeout, err := e.getTimeout(ctx, scope)

	if err != nil {
		return values.None, core.SourceError(e.src, errors.Wrap(err, "failed to calculate timeout value"))
	}

	subscription := events.Subscription{
		EventName: eventName,
		Options:   opts,
	}

	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	if e.filter == nil {
		return e.consumeFirst(ctx, observable, subscription)
	}

	return e.consumeFiltered(ctx, scope, observable, subscription)
}

func (e *WaitForEventExpression) consumeFirst(ctx context.Context, observable events.Observable, subscription events.Subscription) (core.Value, error) {
	ch, err := observable.Subscribe(ctx, subscription)

	if err != nil {
		return values.None, core.SourceError(e.src, err)
	}

	select {
	case evt := <-ch:
		if evt.Err != nil {
			return values.None, core.SourceError(e.src, evt.Err)
		}

		return evt.Data, nil
	case <-ctx.Done():
		return values.None, core.SourceError(e.src, core.ErrTimeout)
	}
}

func (e *WaitForEventExpression) consumeFiltered(ctx context.Context, scope *core.Scope, observable events.Observable, subscription events.Subscription) (core.Value, error) {
	ch, err := observable.Subscribe(ctx, subscription)

	if err != nil {
		return values.None, core.SourceError(e.src, err)
	}

	iterable, err := clauses.NewFilterClause(
		e.filterSrc,
		collections.AsIterable(func(ctx context.Context, scope *core.Scope) (collections.Iterator, error) {
			return collections.FromCoreIterator(e.filterVariable, "", events.NewIterator(ch))
		}),
		e.filter,
	)

	if err != nil {
		return values.None, core.SourceError(e.src, err)
	}

	iter, err := iterable.Iterate(ctx, scope)

	if err != nil {
		return values.None, core.SourceError(e.src, err)
	}

	out, err := iter.Next(ctx, scope)

	if err != nil {
		return values.None, core.SourceError(e.src, err)
	}

	return out.GetVariable(e.filterVariable)
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

	return time.Duration(timeoutInt) * time.Millisecond, nil
}
