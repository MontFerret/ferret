package http

import (
	"bytes"
	"context"
	"io"
	h "net/http"

	"github.com/MontFerret/ferret/pkg/runtime"
)

type Params struct {
	Method  runtime.String
	URL     runtime.String
	Headers *runtime.Object
	Body    runtime.Binary
}

// REQUEST makes a HTTP request.
// @param {hashMap} params - Request parameters.
// @param {String} params.method - HTTP method
// @param {String} params.url - Target url
// @param {Binary} params.body - Request data
// @param {hashMap} [params.headers] - HTTP headers
// @return {Binary} - Response in binary format
func REQUEST(ctx context.Context, args ...runtime.Value) (runtime.Value, error) {
	return execMethod(ctx, "", args)
}

func execMethod(ctx context.Context, method runtime.String, args []runtime.Value) (runtime.Value, error) {
	if err := runtime.ValidateArgs(args, 1, 1); err != nil {
		return runtime.None, err
	}

	params, err := runtime.CastMap(args[0])

	if err != nil {
		return runtime.None, err
	}

	p, err := newParamsFrom(ctx, params)

	if err != nil {
		return runtime.None, err
	}

	if method != "" {
		p.Method = method
	}

	return makeRequest(ctx, p)
}

func makeRequest(ctx context.Context, params Params) (runtime.Value, error) {
	client := h.Client{}
	req, err := h.NewRequest(params.Method.String(), params.URL.String(), bytes.NewBuffer(params.Body))

	if err != nil {
		return runtime.None, err
	}

	req.Header = h.Header{}

	if params.Headers != nil {
		params.Headers.ForEach(ctx, func(c context.Context, value, key runtime.Value) (runtime.Boolean, error) {
			req.Header.Set(key.String(), value.String())

			return true, nil
		})
	}

	resp, err := client.Do(req.WithContext(ctx))

	if err != nil {
		return runtime.None, err
	}

	data, err := io.ReadAll(resp.Body)

	if err != nil {
		return runtime.None, err
	}

	defer resp.Body.Close()

	return runtime.NewBinary(data), nil
}

func newParamsFrom(ctx context.Context, obj runtime.Map) (Params, error) {
	p := Params{}

	method, err := obj.Get(ctx, runtime.String("method"))

	if err == nil {
		p.Method = runtime.ToString(method)
	}

	url, err := obj.Get(ctx, runtime.String("url"))

	if err != nil {
		return Params{}, runtime.Error(runtime.ErrMissedArgument, ".url")
	}

	p.URL = runtime.String(url.String())

	headers, err := obj.Get(ctx, runtime.String("headers"))

	if err == nil {
		if err := runtime.ValidateType(headers, runtime.TypeObject, runtime.TypeMap); err != nil {
			return Params{}, runtime.Error(err, ".headers")
		}

		p.Headers = headers.(*runtime.Object)
	}

	body, err := obj.Get(ctx, runtime.String("body"))

	if err == nil {
		bin, ok := body.(runtime.Binary)

		if ok {
			p.Body = bin
		} else {
			j, err := body.MarshalJSON()

			if err != nil {
				return Params{}, runtime.Error(err, ".body")
			}

			p.Body = runtime.NewBinary(j)

			if p.Headers == nil {
				p.Headers = runtime.NewObject()
			}

			_ = p.Headers.Set(ctx, runtime.String("Content-Type"), runtime.String("application/json"))
		}
	}

	return p, nil
}
