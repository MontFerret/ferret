package vm

import (
	"context"
	"errors"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

var errStreamGroupSetupAborted = errors.New("stream group setup aborted")

type (
	streamGroupDescriptor struct {
		observable runtime.Observable
		options    runtime.Map
		outcome    *streamGroupSetupOutcome
		eventName  runtime.String
	}

	streamGroupSetupResult struct {
		stream     runtime.Stream
		err        error
		panicValue any
		panicked   bool
	}

	streamGroupSetupOutcome struct {
		err        error
		panicValue any
		panicked   bool
	}

	streamGroupCloseResult struct {
		err        error
		panicValue any
		panicked   bool
	}
)

func subscribeStreamGroup(
	ctx context.Context,
	sources, names, options runtime.Value,
) ([]runtime.Stream, error) {
	descriptors, err := coerceStreamGroupDescriptors(ctx, sources, names, options)
	if err != nil {
		return nil, err
	}

	setupCtx, cancel := context.WithCancelCause(ctx)
	defer cancel(nil)
	streams := make([]runtime.Stream, len(descriptors))
	completed := make(chan int, len(descriptors))

	for idx, descriptor := range descriptors {
		go func(idx int, descriptor streamGroupDescriptor) {
			result := subscribeStreamGroupDescriptor(setupCtx, descriptor)
			streams[idx] = result.stream

			if result.panicked || result.err != nil {
				descriptors[idx].outcome = &streamGroupSetupOutcome{
					err:        result.err,
					panicValue: result.panicValue,
					panicked:   result.panicked,
				}
			}

			completed <- idx
		}(idx, descriptor)
	}

	setupFailed := false

	for range descriptors {
		idx := <-completed
		outcome := descriptors[idx].outcome

		if !setupFailed && outcome != nil {
			setupFailed = true
			cancel(errStreamGroupSetupAborted)
		}
	}

	if !setupFailed {
		return streams, nil
	}

	var setupErr error
	var setupPanic *streamGroupSetupOutcome
	for idx := range descriptors {
		outcome := descriptors[idx].outcome
		if outcome == nil {
			continue
		}

		if outcome.panicked && setupPanic == nil {
			setupPanic = outcome
		}

		if outcome.err != nil {
			setupErr = errors.Join(setupErr, outcome.err)
		}
	}

	if ctxErr := ctx.Err(); ctxErr != nil && !errors.Is(setupErr, ctxErr) {
		setupErr = errors.Join(setupErr, ctxErr)
	}

	closeResult := closeStreamGroupSetupSafely(streams)
	if setupPanic != nil {
		panic(setupPanic.panicValue)
	}

	if closeResult.panicked {
		panic(closeResult.panicValue)
	}

	return nil, errors.Join(setupErr, closeResult.err)
}

func subscribeStreamGroupDescriptor(
	ctx context.Context,
	descriptor streamGroupDescriptor,
) (result streamGroupSetupResult) {
	defer func() {
		if recovered := recover(); recovered != nil {
			result.panicValue = recovered
			result.panicked = true
		}

		if result.err != nil && errors.Is(context.Cause(ctx), errStreamGroupSetupAborted) {
			result.err = stripStreamGroupInternalCancellation(result.err)
		}
	}()

	result.stream, result.err = descriptor.observable.Subscribe(ctx, runtime.Subscription{
		EventName: descriptor.eventName,
		Options:   descriptor.options,
	})

	return result
}

func coerceStreamGroupDescriptors(
	ctx context.Context,
	sources, names, options runtime.Value,
) ([]streamGroupDescriptor, error) {
	sourceArray, ok := sources.(*runtime.Array)
	if !ok {
		return nil, runtime.TypeErrorOf(sources, runtime.TypeArray)
	}

	nameArray, ok := names.(*runtime.Array)
	if !ok {
		return nil, runtime.TypeErrorOf(names, runtime.TypeArray)
	}

	optionsArray, ok := options.(*runtime.Array)
	if !ok {
		return nil, runtime.TypeErrorOf(options, runtime.TypeArray)
	}

	length, err := sourceArray.Length(ctx)
	if err != nil {
		return nil, err
	}

	nameLength, err := nameArray.Length(ctx)
	if err != nil {
		return nil, err
	}

	optionsLength, err := optionsArray.Length(ctx)
	if err != nil {
		return nil, err
	}

	if length == 0 {
		return nil, runtime.Error(runtime.ErrInvalidArgument, "stream group must contain at least one arm")
	}

	if nameLength != length || optionsLength != length {
		return nil, runtime.Error(runtime.ErrInvalidArgument, "stream group descriptor arrays must have equal lengths")
	}

	descriptors := make([]streamGroupDescriptor, int(length))
	for idx := runtime.ZeroInt; idx < length; idx++ {
		source, err := sourceArray.At(ctx, idx)
		if err != nil {
			return nil, err
		}

		name, err := nameArray.At(ctx, idx)
		if err != nil {
			return nil, err
		}

		option, err := optionsArray.At(ctx, idx)
		if err != nil {
			return nil, err
		}

		observable, eventName, opts, err := coerceSubscribeArgs(source, name, option)
		if err != nil {
			return nil, err
		}

		descriptors[int(idx)] = streamGroupDescriptor{
			observable: observable,
			eventName:  eventName,
			options:    opts,
		}
	}

	return descriptors, nil
}

func closeStreamGroupSetup(streams []runtime.Stream) error {
	result := closeStreamGroupSetupSafely(streams)
	if result.panicked {
		panic(result.panicValue)
	}

	return result.err
}

// closeStreamGroupSetupSafely closes every established stream while retaining the first panic.
func closeStreamGroupSetupSafely(streams []runtime.Stream) streamGroupCloseResult {
	var result streamGroupCloseResult
	for _, stream := range streams {
		if stream == nil {
			continue
		}

		closeResult := closeStreamGroupStream(stream)
		if closeResult.err != nil {
			result.err = errors.Join(result.err, closeResult.err)
		}

		if closeResult.panicked && !result.panicked {
			result.panicValue = closeResult.panicValue
			result.panicked = true
		}
	}

	return result
}

func closeStreamGroupStream(stream runtime.Stream) (result streamGroupCloseResult) {
	defer func() {
		if recovered := recover(); recovered != nil {
			result.panicValue = recovered
			result.panicked = true
		}
	}()

	result.err = stream.Close()

	return result
}

// stripStreamGroupInternalCancellation removes only cancellation leaves caused by setup fail-fast.
func stripStreamGroupInternalCancellation(err error) error {
	if err == nil {
		return nil
	}

	if joined, ok := err.(interface{ Unwrap() []error }); ok {
		var remaining error
		for _, nested := range joined.Unwrap() {
			if cleaned := stripStreamGroupInternalCancellation(nested); cleaned != nil {
				remaining = errors.Join(remaining, cleaned)
			}
		}

		return remaining
	}

	if wrapped, ok := err.(interface{ Unwrap() error }); ok {
		if errors.Is(err, context.Canceled) || errors.Is(err, errStreamGroupSetupAborted) {
			return stripStreamGroupInternalCancellation(wrapped.Unwrap())
		}
	}

	if errors.Is(err, context.Canceled) || errors.Is(err, errStreamGroupSetupAborted) {
		return nil
	}

	return err
}
