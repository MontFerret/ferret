package vm

import (
	"context"
	"errors"
	"sync"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

type (
	streamGroupDescriptor struct {
		observable runtime.Observable
		options    runtime.Map
		eventName  runtime.String
	}

	streamGroupSetupResult struct {
		stream runtime.Stream
		err    error
		index  int
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

	setupCtx, cancel := context.WithCancel(ctx)
	results := make(chan streamGroupSetupResult, len(descriptors))

	for idx, descriptor := range descriptors {
		go func() {
			stream, subscribeErr := descriptor.observable.Subscribe(setupCtx, runtime.Subscription{
				EventName: descriptor.eventName,
				Options:   descriptor.options,
			})
			results <- streamGroupSetupResult{stream: stream, err: subscribeErr, index: idx}
		}()
	}

	streams := make([]runtime.Stream, len(descriptors))
	var setupErr error

	for range descriptors {
		result := <-results
		if result.err != nil {
			setupErr = errors.Join(setupErr, result.err)
			cancel()
		}

		if result.stream != nil {
			streams[result.index] = result.stream
		}
	}

	cancel()

	if setupErr == nil {
		return streams, nil
	}

	return nil, errors.Join(setupErr, closeStreamGroupSetup(streams))
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
	var wg sync.WaitGroup
	errs := make(chan error, len(streams))
	for _, stream := range streams {
		if stream == nil {
			continue
		}

		wg.Go(func() {
			if err := stream.Close(); err != nil {
				errs <- err
			}
		})
	}

	wg.Wait()
	close(errs)

	var err error
	for closeErr := range errs {
		err = errors.Join(err, closeErr)
	}

	return err
}
