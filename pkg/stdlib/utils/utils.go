package utils

import (
	"context"
	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/logging"
	"github.com/MontFerret/ferret/pkg/runtime/values"
	"time"
)

func Wait(_ context.Context, inputs ...core.Value) (core.Value, error) {
	err := core.ValidateArgs(inputs, 1, 1)

	if err != nil {
		return values.None, nil
	}

	arg := values.ZeroInt

	err = core.ValidateType(inputs[0], core.IntType)

	if err != nil {
		return values.None, err
	}

	arg = inputs[0].(values.Int)

	time.Sleep(time.Millisecond * time.Duration(arg))

	return values.None, nil
}

func Log(ctx context.Context, inputs ...core.Value) (core.Value, error) {
	args := make([]interface{}, 0, len(inputs)+1)

	for _, input := range inputs {
		args = append(args, input)
	}

	logger := logging.FromContext(ctx)

	logger.Print(args...)

	return values.None, nil
}
