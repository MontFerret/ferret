package http

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	stdhttp "net/http"
	"net/url"
	"strings"
)

func toStdRequest(ctx context.Context, req *Request, p *Policies) (*stdhttp.Request, error) {
	if err := p.Eval(req); err != nil {
		return nil, err
	}

	method := strings.TrimSpace(req.Method)
	if method == "" {
		method = stdhttp.MethodGet
	}

	u, err := parseRequestURL(req.URL)
	if err != nil {
		return nil, err
	}

	stdReq, err := stdhttp.NewRequestWithContext(
		ctx,
		method,
		u.String(),
		bytes.NewReader(req.Body),
	)

	if err != nil {
		return nil, fmt.Errorf("http: build request: %w", err)
	}

	for key, values := range req.Headers {
		key = strings.TrimSpace(key)
		if key == "" {
			continue
		}

		canonicalKey := stdhttp.CanonicalHeaderKey(key)
		if p.isBlockedHeader(canonicalKey) {
			continue
		}

		for _, value := range values {
			stdReq.Header.Add(canonicalKey, value)
		}
	}

	for key, value := range p.defaultHeaders {
		if stdReq.Header.Get(key) == "" && !p.isBlockedHeader(key) {
			stdReq.Header.Set(key, value)
		}
	}

	return stdReq, nil
}

func parseRequestURL(raw string) (*url.URL, error) {
	rawURL := strings.TrimSpace(raw)
	if rawURL == "" {
		return nil, errors.New("http: url is required")
	}

	u, err := url.Parse(rawURL)
	if err != nil {
		return nil, fmt.Errorf("http: parse url: %w", err)
	}
	if u.Scheme == "" {
		return nil, errors.New("http: url scheme is required")
	}
	if u.Host == "" {
		return nil, errors.New("http: url host is required")
	}

	u.Scheme = strings.ToLower(u.Scheme)
	u.Host = strings.ToLower(u.Host)

	return u, nil
}

func fromStdResponse(res *stdhttp.Response, p *Policies) (*Response, error) {
	if res == nil {
		return nil, errors.New("http: response is nil")
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

	data, err := io.ReadAll(io.LimitReader(body, limit+1))
	if err != nil {
		return nil, err
	}

	if int64(len(data)) > limit {
		return nil, fmt.Errorf("http: response body exceeds limit: %d > %d", len(data), limit)
	}

	return data, nil
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
