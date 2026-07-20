package http

import (
	"context"
	"errors"
	"io"
	stdhttp "net/http"
	"reflect"
	"testing"
)

func TestToStdRequestIsPolicyIndependentAndDoesNotMutateDTO(t *testing.T) {
	type contextKey struct{}

	ctx, cancel := context.WithCancel(context.WithValue(context.Background(), contextKey{}, "caller"))
	defer cancel()

	req := &Request{
		Method: stdhttp.MethodConnect,
		URL:    "HTTP://127.0.0.1/internal",
		Headers: Headers{
			"Authorization": {"Bearer test-token"},
			"X-Values":      {"one", "two"},
			"x-values":      {"three"},
		},
		Body: []byte("body"),
	}
	wantMethod := req.Method
	wantURL := req.URL
	wantHeaders := Headers{
		"Authorization": {"Bearer test-token"},
		"X-Values":      {"one", "two"},
		"x-values":      {"three"},
	}
	wantBody := append([]byte(nil), req.Body...)

	stdReq, err := toStdRequest(ctx, req)
	if err != nil {
		t.Fatalf("convert request: %v", err)
	}
	if stdReq.Method != stdhttp.MethodConnect {
		t.Fatalf("expected method %q, got %q", stdhttp.MethodConnect, stdReq.Method)
	}
	if stdReq.URL.Hostname() != "127.0.0.1" {
		t.Fatalf("expected loopback URL to be converted without policy evaluation, got %q", stdReq.URL)
	}
	if got := stdReq.Header.Get("Authorization"); got != "Bearer test-token" {
		t.Fatalf("expected policy-sensitive header to be converted, got %q", got)
	}
	if got := stdReq.Header.Values("X-Values"); !reflect.DeepEqual(got, []string{"one", "two", "three"}) {
		t.Fatalf("expected deterministic case-equivalent header merge, got %v", got)
	}
	if stdReq.Host != stdReq.URL.Host || stdReq.RequestURI != "" || stdReq.Close ||
		len(stdReq.TransferEncoding) != 0 || len(stdReq.Trailer) != 0 {
		t.Fatalf("conversion produced non-client request state: %#v", stdReq)
	}
	if got := stdReq.Context().Value(contextKey{}); got != "caller" {
		t.Fatalf("expected caller context value, got %v", got)
	}
	if stdReq.Context() != ctx {
		t.Fatal("conversion replaced the caller-owned context")
	}

	body, err := io.ReadAll(stdReq.Body)
	if err != nil {
		t.Fatalf("read converted body: %v", err)
	}
	if err := stdReq.Body.Close(); err != nil {
		t.Fatalf("close converted body: %v", err)
	}
	if string(body) != "body" || stdReq.ContentLength != int64(len(body)) {
		t.Fatalf("unexpected converted body=%q content-length=%d", body, stdReq.ContentLength)
	}
	if stdReq.GetBody == nil {
		t.Fatal("expected non-empty DTO body to have replay metadata")
	}
	replay, err := stdReq.GetBody()
	if err != nil {
		t.Fatalf("get replay body: %v", err)
	}
	replayedBody, err := io.ReadAll(replay)
	if err != nil {
		t.Fatalf("read replay body: %v", err)
	}
	if err := replay.Close(); err != nil {
		t.Fatalf("close replay body: %v", err)
	}
	if !reflect.DeepEqual(replayedBody, req.Body) {
		t.Fatalf("expected replay body %q, got %q", req.Body, replayedBody)
	}

	if req.Method != wantMethod || req.URL != wantURL ||
		!reflect.DeepEqual(req.Headers, wantHeaders) || !reflect.DeepEqual(req.Body, wantBody) {
		t.Fatalf("conversion mutated DTO: %#v", req)
	}

	stdReq.Header["X-Values"][0] = "changed"
	if req.Headers["X-Values"][0] != "one" {
		t.Fatalf("converted headers alias DTO values: %v", req.Headers["X-Values"])
	}

	cancel()
	<-stdReq.Context().Done()
	if !errors.Is(stdReq.Context().Err(), context.Canceled) {
		t.Fatalf("expected caller cancellation to propagate, got %v", stdReq.Context().Err())
	}
}

func TestToStdRequestNormalizesDTOMethod(t *testing.T) {
	tests := []struct {
		name   string
		method string
		want   string
	}{
		{name: "empty", want: stdhttp.MethodGet},
		{name: "blank", method: " \t ", want: stdhttp.MethodGet},
		{name: "trimmed lowercase", method: " post ", want: stdhttp.MethodPost},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := toStdRequest(context.Background(), &Request{
				Method: tt.method,
				URL:    "https://example.com",
			})
			if err != nil {
				t.Fatalf("convert request: %v", err)
			}
			if req.Method != tt.want {
				t.Fatalf("expected method %q, got %q", tt.want, req.Method)
			}
		})
	}
}

func TestToStdRequestRejectsNilDTO(t *testing.T) {
	req, err := toStdRequest(context.Background(), nil)
	if req != nil {
		t.Fatalf("expected no converted request, got %#v", req)
	}
	if !errors.Is(err, ErrNilRequest) {
		t.Fatalf("expected ErrNilRequest, got %T: %v", err, err)
	}
}

func TestToStdRequestDoesNotMutateDTOOnError(t *testing.T) {
	req := &Request{
		Method: "BAD METHOD",
		URL:    "HTTP://EXAMPLE.COM/path",
		Headers: Headers{
			"X-Test": {"one", "two"},
		},
		Body: []byte("body"),
	}
	want := &Request{
		Method: "BAD METHOD",
		URL:    "HTTP://EXAMPLE.COM/path",
		Headers: Headers{
			"X-Test": {"one", "two"},
		},
		Body: []byte("body"),
	}

	_, err := toStdRequest(context.Background(), req)
	var methodErr *InvalidMethodError
	if !errors.As(err, &methodErr) {
		t.Fatalf("expected InvalidMethodError, got %T: %v", err, err)
	}
	if !reflect.DeepEqual(req, want) {
		t.Fatalf("failed conversion mutated DTO: want=%#v got=%#v", want, req)
	}
}
