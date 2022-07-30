package http

import (
	"bytes"
	"context"
	"io/ioutil"
	h "net/http"

	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
	"github.com/MontFerret/ferret/pkg/runtime/values/types"
)

type Params struct {
	Method  values.String
	URL     values.String
	Headers *values.Object
	Body    values.Binary
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

func execMethod(ctx context.Context, method values.String, args []core.Value) (core.Value, error) {
	if err := core.ValidateArgs(args, 1, 1); err != nil {
		return values.None, err
	}

	arg := args[0]

	if err := core.ValidateType(arg, types.Object); err != nil {
		return values.None, err
	}

	p, err := newParamsFrom(arg.(*values.Object))

	if err != nil {
		return values.None, err
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
		return values.None, err
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
		return values.None, err
	}

	data, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return values.None, err
	}

	defer resp.Body.Close()

	return values.NewBinary(data), nil
}

func newParamsFrom(obj *values.Object) (Params, error) {
	p := Params{}

	method, exists := obj.Get("method")

	if exists {
		p.Method = values.ToString(method)
	}

	url, exists := obj.Get("url")

	if !exists {
		return Params{}, core.Error(core.ErrMissedArgument, ".url")
	}

	p.URL = values.NewString(url.String())

	headers, exists := obj.Get("headers")

	if exists {
		if err := core.ValidateType(headers, types.Object); err != nil {
			return Params{}, core.Error(err, ".headers")
		}

		p.Headers = headers.(*values.Object)
	}

	body, exists := obj.Get("body")

	if exists {
		if core.IsTypeOf(body, types.Binary) {
			p.Body = body.(values.Binary)
		} else {
			j, err := body.MarshalJSON()

			if err != nil {
				return Params{}, core.Error(err, ".body")
			}

			p.Body = values.NewBinary(j)

			if p.Headers == nil {
				p.Headers = values.NewObject()
			}

			p.Headers.Set("Content-Type", values.NewString("application/json"))
		}
	}

	return p, nil
}
