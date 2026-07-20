package http

import (
	"context"
	"errors"
	"io"
	"math"
	"math/rand"
	stdhttp "net/http"
	"net/url"
	"strings"
	"testing"
	"testing/iotest"
)

func TestReadResponseBodyBoundaries(t *testing.T) {
	tests := []struct {
		name      string
		body      string
		want      string
		limit     int64
		wantSize  int64
		wantLimit int64
	}{
		{name: "below", body: "12", limit: 3, want: "12"},
		{name: "at", body: "123", limit: 3, want: "123"},
		{name: "above", body: "1234", limit: 3, wantSize: 4, wantLimit: 3},
		{name: "unlimited", body: "1234", limit: 0, want: "1234"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, err := readResponseBody(strings.NewReader(tt.body), tt.limit)
			if tt.wantLimit == 0 {
				if err != nil {
					t.Fatalf("read response body: %v", err)
				}
				if got := string(body); got != tt.want {
					t.Fatalf("expected body %q, got %q", tt.want, got)
				}
				return
			}

			if body != nil {
				t.Fatalf("expected oversized body not to be returned, got %q", body)
			}
			limitErr := requireResponseBodyLimitError(t, err)
			if limitErr.Size != tt.wantSize || limitErr.Limit != tt.wantLimit {
				t.Fatalf("unexpected body limit error: %#v", limitErr)
			}
		})
	}
}

func TestReadResponseBodyMaxInt64DoesNotOverflow(t *testing.T) {
	body, err := readResponseBody(strings.NewReader("tiny"), math.MaxInt64)
	if err != nil {
		t.Fatalf("read with maximum limit: %v", err)
	}
	if got := string(body); got != "tiny" {
		t.Fatalf("expected body %q, got %q", "tiny", got)
	}
}

func TestReadResponseBodyStopsAfterLimitPlusOne(t *testing.T) {
	const (
		limit     = int64(3)
		bodySize  = 100
		wantReads = int(limit + 1)
	)

	source := strings.NewReader(strings.Repeat("x", bodySize))
	body, err := readResponseBody(source, limit)
	if body != nil {
		t.Fatalf("expected oversized body not to be returned, got %q", body)
	}

	limitErr := requireResponseBodyLimitError(t, err)
	if limitErr.Size != limit+1 || limitErr.Limit != limit {
		t.Fatalf("unexpected body limit error: %#v", limitErr)
	}
	if consumed := bodySize - source.Len(); consumed != wantReads {
		t.Fatalf("expected exactly %d bytes to be read, got %d", wantReads, consumed)
	}
}

func TestReadResponseBodyLimitTakesPrecedenceAfterExtraByte(t *testing.T) {
	readErr := errors.New("read failed after data")
	body, err := readResponseBody(newTerminalErrorReader("four", readErr), 3)
	if body != nil {
		t.Fatalf("expected oversized body not to be returned, got %q", body)
	}

	limitErr := requireResponseBodyLimitError(t, err)
	if limitErr.Size != 4 || limitErr.Limit != 3 {
		t.Fatalf("unexpected body limit error: %#v", limitErr)
	}
	if errors.Is(err, readErr) {
		t.Fatalf("expected observed size violation to take precedence, got %v", err)
	}
}

func TestSaturatedIncrement(t *testing.T) {
	if got := saturatedIncrement(3); got != 4 {
		t.Fatalf("expected ordinary increment to return 4, got %d", got)
	}
	if got := saturatedIncrement(math.MaxInt64); got != math.MaxInt64 {
		t.Fatalf("expected maximum observed size to saturate, got %d", got)
	}
}

func TestFromStdResponseAlwaysClosesBody(t *testing.T) {
	tests := []struct {
		reader  io.Reader
		name    string
		limit   int64
		wantErr bool
	}{
		{name: "success", reader: strings.NewReader("ok"), limit: 3},
		{name: "limit", reader: strings.NewReader("four"), limit: 3, wantErr: true},
		{name: "read error", reader: iotest.ErrReader(errors.New("read failed")), limit: 3, wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body := newTrackedResponseBody(tt.reader)
			policy := newTestPolicy(t, WithMaxResponseSize(tt.limit))
			_, err := fromStdResponse(&stdhttp.Response{
				StatusCode: stdhttp.StatusOK,
				Status:     "200 OK",
				Header:     make(stdhttp.Header),
				Body:       body,
			}, policy)
			if tt.wantErr && err == nil {
				t.Fatal("expected response read error")
			}
			if !tt.wantErr && err != nil {
				t.Fatalf("materialize response: %v", err)
			}
			if !body.closed.Load() {
				t.Fatal("expected response body to be closed")
			}
		})
	}
}

func TestFromStdResponseNilResponse(t *testing.T) {
	response, err := fromStdResponse(nil, newTestPolicy(t))
	if response != nil {
		t.Fatalf("expected no response, got %#v", response)
	}
	if !errors.Is(err, ErrNilResponse) {
		t.Fatalf("expected ErrNilResponse, got %v", err)
	}
}

func TestClientNilResponsePreservesErrNilResponse(t *testing.T) {
	client := &defaultHTTPClient{
		policy: newTestPolicy(t),
		client: stdhttp.Client{Transport: testRoundTripper(func(*stdhttp.Request) (*stdhttp.Response, error) {
			return nil, nil
		})},
	}

	response, err := client.Do(context.Background(), &Request{URL: "https://example.com"})
	if response != nil {
		t.Fatalf("expected no response, got %#v", response)
	}
	if !errors.Is(err, ErrNilResponse) {
		t.Fatalf("expected ErrNilResponse through client wrapper, got %v", err)
	}

	var urlErr *url.Error
	if !errors.As(err, &urlErr) {
		t.Fatalf("expected surrounding url.Error, got %T: %v", err, err)
	}
}

func TestClientEnforcesDefaultResponseBodyLimit(t *testing.T) {
	policy := newTestPolicy(t)
	client := &defaultHTTPClient{
		policy: policy,
		client: stdhttp.Client{Transport: testRoundTripper(func(*stdhttp.Request) (*stdhttp.Response, error) {
			return &stdhttp.Response{
				StatusCode: stdhttp.StatusOK,
				Status:     "200 OK",
				Header:     make(stdhttp.Header),
				Body: io.NopCloser(io.LimitReader(
					rand.New(rand.NewSource(1)),
					defaultMaxResponseSize+1,
				)),
			}, nil
		})},
	}

	response, err := client.Do(context.Background(), &Request{URL: "https://example.com"})
	if response != nil {
		t.Fatalf("expected no oversized response, got %#v", response)
	}
	limitErr := requireResponseBodyLimitError(t, err)
	if limitErr.Size != defaultMaxResponseSize+1 || limitErr.Limit != defaultMaxResponseSize {
		t.Fatalf("unexpected default response body limit error: %#v", limitErr)
	}
}

func requireResponseBodyLimitError(t *testing.T, err error) *ResponseBodyLimitError {
	t.Helper()

	if err == nil {
		t.Fatal("expected response body limit error")
	}

	var limitErr *ResponseBodyLimitError
	if !errors.As(err, &limitErr) {
		t.Fatalf("expected ResponseBodyLimitError, got %T: %v", err, err)
	}

	return limitErr
}
