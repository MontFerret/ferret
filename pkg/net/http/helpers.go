package http

import (
	"context"
	"io"
	"math"
	stdhttp "net/http"
	"net/url"
	"strings"
)

func toStdRequest(ctx context.Context, req *Request, p *Policy) (*stdhttp.Request, error) {
	if p == nil {
		p = &Policy{}
	}

	return p.prepareRequest(ctx, req)
}

func parseRequestURL(raw string) (*url.URL, error) {
	rawURL := strings.TrimSpace(raw)
	if rawURL == "" {
		return nil, &URLValidationError{Field: "url", Reason: "is required"}
	}

	u, err := url.Parse(rawURL)
	if err != nil {
		return nil, &URLParseError{Err: err}
	}

	if u.Scheme == "" {
		return nil, &URLValidationError{Field: "scheme", Reason: "is required"}
	}

	if u.Host == "" {
		return nil, &URLValidationError{Field: "host", Reason: "is required"}
	}

	u.Scheme = asciiLower(u.Scheme)
	u.Host = asciiLower(u.Host)

	return u, nil
}

func fromStdResponse(res *stdhttp.Response, p *Policy) (*Response, error) {
	if res == nil {
		return nil, ErrNilResponse
	}

	if res.Body == nil {
		return &Response{
			StatusCode: res.StatusCode,
			Status:     res.Status,
			Headers:    copyHeaders(res.Header),
		}, nil
	}
	defer res.Body.Close()

	body, err := readResponseBody(res.Body, p.maxResponseSize)
	if err != nil {
		return nil, err
	}

	return &Response{
		StatusCode: res.StatusCode,
		Status:     res.Status,
		Headers:    copyHeaders(res.Header),
		Body:       body,
	}, nil
}

func readResponseBody(body io.Reader, limit int64) ([]byte, error) {
	if limit <= 0 {
		return io.ReadAll(body)
	}

	readLimit := saturatedIncrement(limit)
	data, err := io.ReadAll(io.LimitReader(body, readLimit))

	if int64(len(data)) > limit {
		return nil, &ResponseBodyLimitError{
			Size:  saturatedIncrement(limit),
			Limit: limit,
		}
	}

	if err != nil {
		return nil, err
	}

	return data, nil
}

func saturatedIncrement(value int64) int64 {
	if value == math.MaxInt64 {
		return math.MaxInt64
	}

	return value + 1
}

func copyHeaders(src stdhttp.Header) Headers {
	if len(src) == 0 {
		return nil
	}

	dst := make(Headers, len(src))

	for key, values := range src {
		dst[key] = append([]string(nil), values...)
	}

	return dst
}

func isValidMethod(method string) bool {
	if method == "" {
		return false
	}

	for _, r := range method {
		if r <= ' ' || r >= 127 || strings.ContainsRune("()<>@,;:\\\"/[]?={}", r) {
			return false
		}
	}

	return true
}
