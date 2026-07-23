package http

import (
	"bytes"
	"context"
	"math"
	stdhttp "net/http"
	"sort"
	"strings"
)

func toStdRequest(ctx context.Context, req *Request) (*stdhttp.Request, error) {
	if req == nil {
		return nil, ErrNilRequest
	}

	method := normalizeRequestMethod(req.Method)
	if !isValidMethod(method) {
		return nil, &InvalidMethodError{Method: req.Method}
	}

	rawURL := strings.TrimSpace(req.URL)
	stdReq, err := stdhttp.NewRequestWithContext(ctx, method, rawURL, bytes.NewReader(req.Body))
	if err != nil {
		if ctx == nil {
			return nil, &RequestBuildError{Err: err}
		}

		return nil, &URLParseError{Err: err}
	}

	if rawURL == "" {
		return nil, &URLValidationError{Field: "url", Reason: "is required"}
	}

	if stdReq.URL.Scheme == "" {
		return nil, &URLValidationError{Field: "scheme", Reason: "is required"}
	}

	if stdReq.URL.Host == "" {
		return nil, &URLValidationError{Field: "host", Reason: "is required"}
	}

	stdReq.URL.Scheme = asciiLower(stdReq.URL.Scheme)
	stdReq.URL.Host = asciiLower(stdReq.URL.Host)
	stdReq.Host = stdReq.URL.Host
	stdReq.Header = copyRequestHeaders(req.Headers)

	return stdReq, nil
}

func copyRequestHeaders(src Headers) stdhttp.Header {
	dst := make(stdhttp.Header, len(src))
	var keyBuffer [8]string
	keys := keyBuffer[:0]

	if len(src) > len(keyBuffer) {
		keys = make([]string, 0, len(src))
	}

	for key := range src {
		keys = append(keys, key)
	}

	sort.Strings(keys)

	for _, key := range keys {
		canonicalKey := stdhttp.CanonicalHeaderKey(key)
		dst[canonicalKey] = append(dst[canonicalKey], src[key]...)
	}

	return dst
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

	body, err := p.ReadResponseBody(res.Body)
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
