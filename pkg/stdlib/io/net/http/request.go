package http

import (
	"bytes"
	"context"
	"io"
	h "net/http"

	ferretencoding "github.com/MontFerret/ferret/v2/pkg/encoding"
	encodingjson "github.com/MontFerret/ferret/v2/pkg/encoding/json"
	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

type Params struct {
	Method  runtime.String
	URL     runtime.String
	Headers runtime.Map
	Body    runtime.Binary
}

// REQUEST makes a HTTP request.
// @param {Map} params - Request parameters.
// @param {String} params.method - HTTP method
// @param {String} params.url - Target url
// @param {Binary} params.body - Request data
// @param {Map} [params.headers] - HTTP headers
// @return {Binary} - Response in binary format
func REQUEST(ctx context.Context, arg runtime.Value) (runtime.Value, error) {
	return execMethod(ctx, "", arg)
}

func execMethod(ctx context.Context, method runtime.String, arg runtime.Value) (runtime.Value, error) {
	params, err := runtime.CastArg[runtime.Map](arg, 0)

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

	methodKey := runtime.String("method")
	hasMethod, err := obj.ContainsKey(ctx, methodKey)
	if err != nil {
		return Params{}, err
	}
	if hasMethod {
		method, err := obj.Get(ctx, methodKey)
		if err != nil {
			return Params{}, err
		}
		p.Method = runtime.ToString(method)
	}

	urlKey := runtime.String("url")
	hasURL, err := obj.ContainsKey(ctx, urlKey)
	if err != nil {
		return Params{}, err
	}
	if !hasURL {
		return Params{}, runtime.Error(runtime.ErrMissedArgument, ".url")
	}

	url, err := obj.Get(ctx, urlKey)
	if err != nil {
		return Params{}, err
	}
	p.URL = runtime.String(url.String())

	headersKey := runtime.String("headers")
	hasHeaders, err := obj.ContainsKey(ctx, headersKey)
	if err != nil {
		return Params{}, err
	}
	if hasHeaders {
		headers, err := obj.Get(ctx, headersKey)
		if err != nil {
			return Params{}, err
		}

		if err := runtime.ValidateArgType(headers, 0, runtime.TypeObject, runtime.TypeMap); err != nil {
			return Params{}, runtime.Error(err, ".headers")
		}

		p.Headers = headers.(runtime.Map)
	}

	bodyKey := runtime.String("body")
	hasBody, err := obj.ContainsKey(ctx, bodyKey)
	if err != nil {
		return Params{}, err
	}
	if hasBody {
		body, err := obj.Get(ctx, bodyKey)
		if err != nil {
			return Params{}, err
		}

		bin, ok := body.(runtime.Binary)

		if ok {
			p.Body = bin
		} else {
			encoder := ferretencoding.Encoder(encodingjson.Default)
			if selected, resolverErr := ferretencoding.EncoderFrom(ctx, encodingjson.ContentType); resolverErr == nil {
				encoder = selected
			}

			j, err := encoder.Encode(body)

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
