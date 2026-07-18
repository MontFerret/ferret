package http

import (
	"context"
	"errors"
	stdhttp "net/http"
	"net/url"
	"strings"
	"testing"
)

func TestRequestBuildErrorPreservesCause(t *testing.T) {
	_, err := toStdRequest(nil, &Request{URL: "https://example.com"}, NewPolicy())
	if err == nil {
		t.Fatal("expected request build error")
	}

	var buildErr *RequestBuildError
	if !errors.As(err, &buildErr) {
		t.Fatalf("expected RequestBuildError, got %T: %v", err, err)
	}
	if buildErr.Err == nil || !errors.Is(err, buildErr.Err) {
		t.Fatalf("expected wrapped request build cause, got %v", err)
	}
	if got, want := err.Error(), "http: build request: "+buildErr.Err.Error(); got != want {
		t.Fatalf("expected %q, got %q", want, got)
	}
	assertNotPolicyError(t, err)
}

func TestURLParseErrorPreservesCause(t *testing.T) {
	_, err := parseRequestURL("http://%")
	if err == nil {
		t.Fatal("expected URL parse error")
	}

	var parseErr *URLParseError
	if !errors.As(err, &parseErr) {
		t.Fatalf("expected URLParseError, got %T: %v", err, err)
	}
	if parseErr.Err == nil || !errors.Is(err, parseErr.Err) {
		t.Fatalf("expected wrapped URL parse cause, got %v", err)
	}
	if got, want := err.Error(), "http: parse url: "+parseErr.Err.Error(); got != want {
		t.Fatalf("expected %q, got %q", want, got)
	}
	assertNotPolicyError(t, err)
}

func TestInvalidMethodError(t *testing.T) {
	err := NewPolicy().Eval(&Request{Method: "BAD METHOD", URL: "https://example.com"})
	if err == nil {
		t.Fatal("expected invalid method error")
	}

	var methodErr *InvalidMethodError
	if !errors.As(err, &methodErr) {
		t.Fatalf("expected InvalidMethodError, got %T: %v", err, err)
	}
	if methodErr.Method != "BAD METHOD" {
		t.Fatalf("expected rejected method %q, got %q", "BAD METHOD", methodErr.Method)
	}
	if got, want := err.Error(), `http: invalid method "BAD METHOD"`; got != want {
		t.Fatalf("expected %q, got %q", want, got)
	}
	assertNotPolicyError(t, err)
}

func TestURLValidationErrorForMissingHost(t *testing.T) {
	tests := []struct {
		validate func() error
		name     string
	}{
		{
			name: "parsed request URL",
			validate: func() error {
				_, err := parseRequestURL("http:///path")
				return err
			},
		},
		{
			name: "policy URL validation",
			validate: func() error {
				return NewPolicy().validateURL(&url.URL{Scheme: "http"}, PolicyTargetRequest)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.validate()
			if err == nil {
				t.Fatal("expected URL validation error")
			}

			var validationErr *URLValidationError
			if !errors.As(err, &validationErr) {
				t.Fatalf("expected URLValidationError, got %T: %v", err, err)
			}
			if validationErr.Field != "host" || validationErr.Reason != "is required" {
				t.Fatalf("unexpected URL validation error: %#v", validationErr)
			}
			if got, want := err.Error(), "http: url host is required"; got != want {
				t.Fatalf("expected %q, got %q", want, got)
			}
			assertNotPolicyError(t, err)
		})
	}
}

func TestRequestBodyLimitError(t *testing.T) {
	err := NewPolicy(WithMaxRequestSize(3)).Eval(&Request{
		URL:  "https://example.com",
		Body: []byte("four"),
	})
	if err == nil {
		t.Fatal("expected request body limit error")
	}

	var limitErr *RequestBodyLimitError
	if !errors.As(err, &limitErr) {
		t.Fatalf("expected RequestBodyLimitError, got %T: %v", err, err)
	}
	if limitErr.Size != 4 || limitErr.Limit != 3 {
		t.Fatalf("unexpected request body limit error: %#v", limitErr)
	}
	if got, want := err.Error(), "http: request body exceeds limit: 4 > 3"; got != want {
		t.Fatalf("expected %q, got %q", want, got)
	}
	assertNotPolicyError(t, err)
}

func TestResponseBodyLimitError(t *testing.T) {
	_, err := readResponseBody(strings.NewReader("four"), 3)
	if err == nil {
		t.Fatal("expected response body limit error")
	}

	var limitErr *ResponseBodyLimitError
	if !errors.As(err, &limitErr) {
		t.Fatalf("expected ResponseBodyLimitError, got %T: %v", err, err)
	}
	if limitErr.Limit != 3 {
		t.Fatalf("expected response body limit 3, got %d", limitErr.Limit)
	}
	if got, want := err.Error(), "http: response body exceeds limit of 3 bytes"; got != want {
		t.Fatalf("expected %q, got %q", want, got)
	}
	assertNotPolicyError(t, err)
}

func TestRedirectLimitErrorSurvivesURLErrorWrapping(t *testing.T) {
	roundTrips := 0
	policy := NewPolicy(WithMaxRedirects(1))
	client := &defaultHTTPClient{
		policy: policy,
		client: stdhttp.Client{Transport: testRoundTripper(func(*stdhttp.Request) (*stdhttp.Response, error) {
			roundTrips++
			return responseWithBody(
				stdhttp.StatusFound,
				"",
				stdhttp.Header{"Location": {"/next"}},
			), nil
		})},
	}

	_, err := client.Do(context.Background(), &Request{URL: "https://example.com/start"})
	if err == nil {
		t.Fatal("expected redirect limit error")
	}

	var urlErr *url.Error
	if !errors.As(err, &urlErr) {
		t.Fatalf("expected surrounding url.Error, got %T: %v", err, err)
	}

	var limitErr *RedirectLimitError
	if !errors.As(err, &limitErr) {
		t.Fatalf("expected RedirectLimitError, got %T: %v", err, err)
	}
	if limitErr.Limit != 1 {
		t.Fatalf("expected redirect limit 1, got %d", limitErr.Limit)
	}
	if got, want := limitErr.Error(), "http: stopped after 1 redirect(s)"; got != want {
		t.Fatalf("expected %q, got %q", want, got)
	}
	if roundTrips != 1 {
		t.Fatalf("expected one request before redirect rejection, got %d", roundTrips)
	}
	assertNotPolicyError(t, err)
}

func TestPolicyErrorClassificationRemainsDistinct(t *testing.T) {
	err := NewPolicy().Eval(&Request{Method: stdhttp.MethodConnect, URL: "https://example.com"})
	if err == nil {
		t.Fatal("expected policy denial")
	}
	if !errors.Is(err, ErrPolicyDenied) {
		t.Fatalf("expected ErrPolicyDenied, got %v", err)
	}

	var policyErr *PolicyError
	if !errors.As(err, &policyErr) {
		t.Fatalf("expected PolicyError, got %T: %v", err, err)
	}
}

func assertNotPolicyError(t *testing.T, err error) {
	t.Helper()

	if errors.Is(err, ErrPolicyDenied) {
		t.Fatalf("expected dedicated non-policy error, got %v", err)
	}

	var policyErr *PolicyError
	if errors.As(err, &policyErr) {
		t.Fatalf("expected dedicated non-policy error, got PolicyError: %v", err)
	}
}
