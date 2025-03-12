package http

import (
	"bytes"
	"context"
	"github.com/MontFerret/ferret/pkg/runtime/internal"
	"io"
	h "net/http"

	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values/types"
)

type Params struct {
	Method  core.String
	URL     core.String
	Headers *internal.Object
	Body    core.Binary
}

// REQUEST makes a HTTP request.
// @param {Object} params - Request parameters.
// @param {String} params.method - HTTP method
// @param {String} params.url - Target url
// @param {Binary} params.body - Request data
// @param {Object} [params.headers] - HTTP headers
// @return {Binary} - Response in binary format
func REQUEST(ctx context.Context, args ...core.Value) (core.Value, error) {
	return execMethod(ctx, "", args)
}

func execMethod(ctx context.Context, method core.String, args []core.Value) (core.Value, error) {
	if err := core.ValidateArgs(args, 1, 1); err != nil {
		return core.None, err
	}

	params, err := core.CastMap(args[0])

	if err != nil {
		return core.None, err
	}

	p, err := newParamsFrom(params)

	if err != nil {
		return core.None, err
	}

	if method != "" {
		p.Method = method
	}

	return makeRequest(ctx, p)
}

func makeRequest(ctx context.Context, params Params) (core.Value, error) {
	client := h.Client{}
	req, err := h.NewRequest(params.Method.String(), params.URL.String(), bytes.NewBuffer(params.Body))

	if err != nil {
		return core.None, err
	}

	req.Header = h.Header{}

	if params.Headers != nil {
		params.Headers.ForEach(func(value core.Value, key string) bool {
			req.Header.Set(key, value.String())

			return true
		})
	}

	resp, err := client.Do(req.WithContext(ctx))

	if err != nil {
		return core.None, err
	}

	data, err := io.ReadAll(resp.Body)

	if err != nil {
		return core.None, err
	}

	defer resp.Body.Close()

	return core.NewBinary(data), nil
}

func newParamsFrom(obj *internal.Object) (Params, error) {
	p := Params{}

	method, exists := obj.Get("method")

	if exists {
		p.Method = internal.ToString(method)
	}

	url, exists := obj.Get("url")

	if !exists {
		return Params{}, core.Error(core.ErrMissedArgument, ".url")
	}

	p.URL = core.NewString(url.String())

	headers, exists := obj.Get("headers")

	if exists {
		if err := core.ValidateType(headers, types.Object); err != nil {
			return Params{}, core.Error(err, ".headers")
		}

		p.Headers = headers.(*internal.Object)
	}

	body, exists := obj.Get("body")

	if exists {
		bin, ok := body.(core.Binary)

		if ok {
			p.Body = bin
		} else {
			j, err := body.MarshalJSON()

			if err != nil {
				return Params{}, core.Error(err, ".body")
			}

			p.Body = core.NewBinary(j)

			if p.Headers == nil {
				p.Headers = internal.NewObject()
			}

			p.Headers.Set("Content-Type", core.NewString("application/json"))
		}
	}

	return p, nil
}
