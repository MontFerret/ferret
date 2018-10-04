package utils

import (
	"context"
	"math/rand"
	"time"

	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/logging"
	"github.com/MontFerret/ferret/pkg/runtime/values"
)

var (
	userAgents = []string{
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/60.0.3112.113 Safari/537.36",
		"Mozilla/5.0 (Windows NT 5.1; rv:7.0.1) Gecko/20100101 Firefox/7.0.1",
		"Mozilla/5.0 (Windows NT 6.1; WOW64; rv:54.0) Gecko/20100101 Firefox/54.0",
		"Mozilla/5.0 (Windows NT 6.1; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/60.0.3112.90 Safari/537.36",
		"Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/44.0.2403.157 Safari/537.36",
		"Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/44.0.2403.157 Safari/537.36",
	}
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

func GetRandomUserAgent() string {
	rand.Seed(time.Now().UnixNano())
	return userAgents[rand.Intn(len(userAgents)-1)]
}
