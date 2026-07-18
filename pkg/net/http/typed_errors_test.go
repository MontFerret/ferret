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
	_, err := toStdRequest(nil, &Request{URL: "https://example.com"}, newTestPolicy(t))
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
	if got := err.Error(); !strings.Contains(got, "invalid URL escape") {
		t.Fatalf("expected safe parse reason, got %q", got)
	} else if strings.Contains(got, "http://%") {
		t.Fatalf("URL parse error included the raw URL: %q", got)
	}
	assertNotPolicyError(t, err)
}

func TestURLParseErrorRedactsMalformedCredentials(t *testing.T) {
	const (
		rawURL = "https://user:malformed-secret@%"
		secret = "malformed-secret"
	)

	err := newTestPolicy(t).Eval(&Request{URL: rawURL})
	var parseErr *URLParseError
	if !errors.As(err, &parseErr) || parseErr.Err == nil {
		t.Fatalf("expected URLParseError with a cause, got %T: %v", err, err)
	}
	if strings.Contains(err.Error(), secret) || strings.Contains(err.Error(), rawURL) {
		t.Fatalf("URL parse error leaked malformed credentials: %v", err)
	}
}

func TestURLParseErrorOnPreparedRequestPath(t *testing.T) {
	tests := []struct {
		run  func() error
		name string
	}{
		{
			name: "policy eval",
			run: func() error {
				return newTestPolicy(t).Eval(&Request{URL: "http://%"})
			},
		},
		{
			name: "client do",
			run: func() error {
				_, err := newTestClient(t).Do(
					context.Background(),
					&Request{URL: "http://%"},
				)

				return err
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.run()
			if err == nil {
				t.Fatal("expected URL parse error")
			}

			var parseErr *URLParseError
			if !errors.As(err, &parseErr) || parseErr.Err == nil {
				t.Fatalf("expected URLParseError with a cause, got %T: %v", err, err)
			}
			assertNotPolicyError(t, err)
		})
	}
}

func TestInvalidMethodError(t *testing.T) {
	err := newTestPolicy(t).Eval(&Request{Method: "BAD METHOD", URL: "https://example.com"})
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

func TestURLValidationErrorsForMissingStructuralFields(t *testing.T) {
	tests := []struct {
		name   string
		rawURL string
		field  string
	}{
		{name: "url", field: "url"},
		{name: "scheme", rawURL: "example.com", field: "scheme"},
		{name: "host", rawURL: "http:///path", field: "host"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := parseRequestURL(tt.rawURL)
			if err == nil {
				t.Fatal("expected URL validation error")
			}

			var validationErr *URLValidationError
			if !errors.As(err, &validationErr) {
				t.Fatalf("expected URLValidationError, got %T: %v", err, err)
			}
			if validationErr.Field != tt.field || validationErr.Reason != "is required" {
				t.Fatalf("unexpected URL validation error: %#v", validationErr)
			}
			if got, want := err.Error(), "http: url "+tt.field+" is required"; got != want {
				t.Fatalf("expected %q, got %q", want, got)
			}
			assertNotPolicyError(t, err)

			err = newTestPolicy(t).Eval(&Request{URL: tt.rawURL})
			if !errors.As(err, &validationErr) || validationErr.Field != tt.field {
				t.Fatalf("expected Policy.Eval URLValidationError for %q, got %T: %v", tt.field, err, err)
			}
		})
	}
}

func TestRequestBodyLimitError(t *testing.T) {
	err := newTestPolicy(t, WithMaxRequestSize(3)).Eval(&Request{
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
	if limitErr.Size != 4 || limitErr.Limit != 3 {
		t.Fatalf("unexpected response body limit error: %#v", limitErr)
	}
	if got := err.Error(); !strings.Contains(got, "response body exceeds limit") {
		t.Fatalf("expected response limit message, got %q", got)
	}
	assertNotPolicyError(t, err)
}

func TestRedirectLimitErrorSurvivesURLErrorWrapping(t *testing.T) {
	roundTrips := 0
	policy := newTestPolicy(t, WithMaxRedirects(1))
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
	if roundTrips != 2 {
		t.Fatalf("expected initial request plus one followed redirect, got %d request(s)", roundTrips)
	}
	assertNotPolicyError(t, err)
}

func TestPolicyErrorClassificationRemainsDistinct(t *testing.T) {
	err := newTestPolicy(t).Eval(&Request{Method: stdhttp.MethodConnect, URL: "https://example.com"})
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
