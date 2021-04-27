package html

import (
	"context"
	"net/http"

	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
	"github.com/MontFerret/ferret/pkg/runtime/values/types"
)

// DOCUMENT_EXISTS returns a boolean value indicating whether a web page exists by a given url.
// @param {String} url - Target url.
// @param {Object} [options] - Request options.
// @param {Object} [options.headers] - Request headers.
// @return {Boolean} - A boolean value indicating whether a web page exists by a given url.
func DocumentExists(ctx context.Context, args ...core.Value) (core.Value, error) {
	if err := core.ValidateArgs(args, 1, 2); err != nil {
		return nil, err
	}

	if err := core.ValidateType(args[0], types.String); err != nil {
		return nil, err
	}

	url := args[0].String()

	client := http.Client{}
	req, err := http.NewRequest("GET", url, nil)

	if err != nil {
		return values.None, err
	}

	if len(args) > 1 {
		if err := core.ValidateType(args[1], types.Object); err != nil {
			return nil, err
		}

		options := args[1].(*values.Object)

		if options.Has("headers") {
			headersOpt := options.MustGet("headers")

			if err := core.ValidateType(headersOpt, types.Object); err != nil {
				return nil, err
			}

			headers := headersOpt.(*values.Object)

			req.Header = http.Header{}

			headers.ForEach(func(value core.Value, key string) bool {
				req.Header.Set(key, value.String())

				return true
			})
		}
	}

	resp, err := client.Do(req.WithContext(ctx))

	if err != nil {
		return values.False, nil
	}

	var exists values.Boolean

	if resp.StatusCode >= 200 && resp.StatusCode <= 299 {
		exists = values.True
	}

	return exists, nil
}
