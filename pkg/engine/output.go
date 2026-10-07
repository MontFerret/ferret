package engine

import (
	"context"
	"errors"
	"sync"

	"github.com/MontFerret/ferret/v2/pkg/encoding"
	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

type (
	outputState uint8

	// output owns detached bytes and consumption coordination, never a VM or
	// session. The invocation remains immutable for admission validation even
	// after finalization. Its metadata is independent of mutable Content.
	output struct {
		invocation    context.Context
		operationErr  error
		cleanupErr    error
		content       *encoding.Content
		done          chan struct{}
		cancel        context.CancelCauseFunc
		metadata      encoding.Metadata
		mu            sync.Mutex
		state         outputState
		explicitClose bool
	}

	outputConsumption struct {
		ctx            context.Context
		caller         context.Context
		cancel         context.CancelCauseFunc
		stopInvocation func() bool
		deadlineCancel context.CancelFunc
	}
)

const (
	outputReady outputState = iota
	outputConsuming
	outputStopping
	outputFinalizing
	outputFinalized
)

var _ encoding.Output = (*output)(nil)

func newOutput(invocation context.Context, outcome runOutcome) *output {
	metadata := encoding.Metadata{}
	if outcome.content != nil {
		metadata = outcome.content.Metadata
	}

	return &output{invocation: invocation, content: outcome.content, metadata: metadata,
		operationErr: outcome.operationErr, cleanupErr: outcome.cleanupErr, done: make(chan struct{})}
}

func (o *output) Metadata() encoding.Metadata {
	return o.metadata
}

func (o *output) Consume(ctx context.Context, consumer encoding.Consumer) (err error) {
	if ctx == nil {
		return runtime.Error(runtime.ErrInvalidArgument, "context is required")
	}

	if consumer == nil {
		return runtime.Error(runtime.ErrInvalidArgument, "consumer is required")
	}

	consumption, err := o.admit(ctx)
	if err != nil {
		return err
	}

	// This also settles the handle during panic unwinding, without recovering.
	defer func() { err = o.finish(consumption, err) }()

	o.mu.Lock()
	content := o.content
	deliver := o.state == outputConsuming && consumption.ctx.Err() == nil
	o.mu.Unlock()
	if deliver && content != nil {
		return consumer(consumption.ctx, content.Data)
	}

	return nil
}

func (o *output) Collect(ctx context.Context) (content *encoding.Content, err error) {
	if ctx == nil {
		return nil, runtime.Error(runtime.ErrInvalidArgument, "context is required")
	}

	consumption, err := o.admit(ctx)
	if err != nil {
		return nil, err
	}

	defer func() { err = o.finish(consumption, err) }()

	// The encoder transferred ownership before Run returned. Clearing the
	// handle's reference in finish cannot mutate the transferred content.
	return o.content, nil
}

func (o *output) Close() error {
	o.mu.Lock()
	switch o.state {
	case outputReady:
		o.explicitClose = true
		o.state = outputFinalizing
		o.mu.Unlock()
		o.finalize()
	case outputConsuming:
		o.explicitClose = true
		o.state = outputStopping
		o.cancel(encoding.ErrOutputClosed)
		o.mu.Unlock()
		<-o.done
	default:
		o.mu.Unlock()
		<-o.done
	}

	return o.cleanupErr
}

func (o *output) admit(ctx context.Context) (outputConsumption, error) {
	o.mu.Lock()
	defer o.mu.Unlock()

	if err := outputContextError(ctx); err != nil {
		return outputConsumption{}, err
	}

	if err := outputContextError(o.invocation); err != nil {
		// Its lifetime cannot be revived. Preserve terminal eager failures even
		// though this rejected call must leave the handle unclaimed.
		return outputConsumption{}, joinOutputErrors(err, o.operationErr, o.cleanupErr)
	}

	switch o.state {
	case outputConsuming:
		return outputConsumption{}, encoding.ErrOutputInUse
	case outputReady:
	default:
		return outputConsumption{}, encoding.ErrOutputClosed
	}

	base := ctx
	var deadlineCancel context.CancelFunc
	if deadline, ok := o.invocation.Deadline(); ok {
		if current, present := ctx.Deadline(); !present || deadline.Before(current) {
			base, deadlineCancel = context.WithDeadline(ctx, deadline)
		}
	}

	effective, cancel := context.WithCancelCause(base)
	invocation := o.invocation
	var stop func() bool
	// Shared cancellation is already propagated by the consumption parent.
	// Avoid registering a second callback in the ordinary same-context path.
	if invocation.Done() != nil && invocation.Done() != ctx.Done() {
		stop = context.AfterFunc(invocation, func() { cancel(context.Cause(invocation)) })
	}

	o.cancel = cancel
	o.state = outputConsuming

	return outputConsumption{ctx: effective, caller: ctx, cancel: cancel, stopInvocation: stop, deadlineCancel: deadlineCancel}, nil
}

func (o *output) finish(consumption outputConsumption, deliveryErr error) error {
	o.mu.Lock()
	// This is the outcome commitment point. Close after this point only waits.
	err := joinOutputErrors(o.operationErr, deliveryErr, outputContextError(consumption.ctx), outputContextError(consumption.caller), outputContextError(o.invocation), o.cleanupErr)
	if err == nil && o.content != nil && o.metadata.LengthKnown && int64(len(o.content.Data)) != o.metadata.Length {
		err = runtime.Error(runtime.ErrInvalidOperation, "encoded output length does not match metadata")
	}

	if o.explicitClose && !errors.Is(err, encoding.ErrOutputClosed) {
		err = errors.Join(err, encoding.ErrOutputClosed)
	}

	o.state = outputFinalizing
	o.mu.Unlock()
	if consumption.stopInvocation != nil {
		consumption.stopInvocation()
	}

	consumption.cancel(nil)
	if consumption.deadlineCancel != nil {
		consumption.deadlineCancel()
	}

	o.finalize()

	return err
}

func (o *output) finalize() {
	o.mu.Lock()
	o.content = nil
	o.operationErr = nil
	o.cancel = nil
	o.state = outputFinalized
	close(o.done)
	o.mu.Unlock()
}
