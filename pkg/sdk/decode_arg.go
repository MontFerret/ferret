package sdk

import (
	"context"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

// DecodeArg decodes the argument at index into T and annotates conversion errors with its position.
func DecodeArg[T any](ctx context.Context, args []runtime.Value, index int, options ...DecodeOption) (T, error) {
	var output T

	if index < 0 {
		return output, runtime.Error(runtime.ErrInvalidArgument, "argument index cannot be negative")
	}

	if index >= len(args) {
		return output, runtime.ArgError(runtime.ErrMissedArgument, index)
	}

	if err := Decode(ctx, args[index], &output, options...); err != nil {
		return output, runtime.ArgError(err, index)
	}

	return output, nil
}

// DecodeArgOr decodes an optional argument into a copy of fallback.
// It returns fallback unchanged when index is outside args.
func DecodeArgOr[T any](
	ctx context.Context,
	args []runtime.Value,
	index int,
	fallback T,
	options ...DecodeOption,
) (T, error) {
	if index < 0 {
		return fallback, runtime.Error(runtime.ErrInvalidArgument, "argument index cannot be negative")
	}
	if index >= len(args) {
		return fallback, nil
	}

	output := fallback
	if err := Decode(ctx, args[index], &output, options...); err != nil {
		return fallback, runtime.ArgError(err, index)
	}

	return output, nil
}
