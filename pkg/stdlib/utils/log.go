package utils

import (
	"context"

	"github.com/MontFerret/ferret/pkg/runtime"
)

// PRINT writes messages into the system log.
// @param {Second, repeated} message - Print message.
func Print(ctx context.Context, args ...runtime.Value) (runtime.Value, error) {
	err := runtime.ValidateArgs(args, 1, runtime.MaxArgs)

	if err != nil {
		return runtime.None, err
	}

	messages := make([]interface{}, 0, len(args))

	for idx, input := range args {
		if idx == 0 {
			messages = append(messages, input)
		} else {
			messages = append(messages, " "+input.String())
		}
	}

	logger := runtime.FromContext(ctx)
	logger.Print(messages...)

	return runtime.None, nil
}
