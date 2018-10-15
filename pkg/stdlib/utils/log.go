package utils

import (
	"context"

	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/logging"
	"github.com/MontFerret/ferret/pkg/runtime/values"
)

// Log writes messages into the system log.
func Log(ctx context.Context, args ...core.Value) (core.Value, error) {
	err := core.ValidateArgs(args, 1, core.MaxArgs)

	if err != nil {
		return values.None, err
	}

	messages := make([]interface{}, 0, len(args))

	for _, input := range args {
		messages = append(messages, input)
	}

	logger := logging.FromContext(ctx)

	logger.Print(messages...)

	return values.None, nil
}
