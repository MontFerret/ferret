package utils

import (
	"context"

	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/logging"
	"github.com/MontFerret/ferret/pkg/runtime/values"
)

// PRINT writes messages into the system log.
// @param {Value, repeated} message - Print message.
func Print(ctx context.Context, args ...core.Value) (core.Value, error) {
	err := core.ValidateArgs(args, 1, core.MaxArgs)

	if err != nil {
		return values.None, err
	}

	messages := make([]interface{}, 0, len(args))

	for idx, input := range args {
		if idx == 0 {
			messages = append(messages, input)
		} else {
			messages = append(messages, " "+input.String())
		}
	}

	logger := logging.FromContext(ctx)

	logger.Print(messages...)

	return values.None, nil
}
