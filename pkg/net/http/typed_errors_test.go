package http

import (
	"context"
	"errors"
	"io"
	stdhttp "net/http"
	"net/url"
	"strings"
	"testing"
)

func TestRequestBuildErrorPreservesCause(t *testing.T) {
	_, err := toStdRequest(nil, &Request{URL: "https://example.com"})
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
	_, err := toStdRequest(context.Background(), &Request{URL: "http://%"})
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

	_, err := toStdRequest(context.Background(), &Request{URL: rawURL})
	var parseErr *URLParseError
	if !errors.As(err, &parseErr) || parseErr.Err == nil {
		t.Fatalf("expected URLParseError with a cause, got %T: %v", err, err)
	}
	if strings.Contains(err.Error(), secret) || strings.Contains(err.Error(), rawURL) {
		t.Fatalf("URL parse error leaked malformed credentials: %v", err)
	}
}

func TestURLParseErrorOnConversionAndClientPaths(t *testing.T) {
	tests := []struct {
		run  func() error
		name string
	}{
		{
			name: "conversion",
			run: func() error {
				_, err := toStdRequest(context.Background(), &Request{URL: "http://%"})

				return err
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
	err := newTestPolicy(t).Eval(newTestPolicyRequest(t, "BAD METHOD", "https://example.com"))
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
			_, err := toStdRequest(context.Background(), &Request{URL: tt.rawURL})
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

			var policyReq *stdhttp.Request
			switch tt.field {
			case "url":
				policyReq = &stdhttp.Request{Method: stdhttp.MethodGet}
			default:
				policyReq = newTestPolicyGETRequest(t, tt.rawURL)
			}

			err = newTestPolicy(t).Eval(policyReq)
			if !errors.As(err, &validationErr) || validationErr.Field != tt.field {
				t.Fatalf("expected Policy.Eval URLValidationError for %q, got %T: %v", tt.field, err, err)
			}
		})
	}
}

func TestRequestBodyLimitError(t *testing.T) {
	req := newTestPolicyGETRequest(t, "https://example.com")
	req.Body = io.NopCloser(strings.NewReader("four"))
	req.ContentLength = 4
	err := newTestPolicy(t, WithMaxRequestSize(3)).Eval(req)
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

func TestPolicyEvalRequestBodyBoundariesAndUnknownLengths(t *testing.T) {
	tests := []struct {
		body          io.ReadCloser
		name          string
		contentLength int64
		wantUnknown   bool
		wantOversized bool
	}{
		{name: "nil body", contentLength: 0},
		{name: "explicit no body", body: stdhttp.NoBody, contentLength: 0},
		{name: "known below", body: io.NopCloser(strings.NewReader("12")), contentLength: 2},
		{name: "known at limit", body: io.NopCloser(strings.NewReader("123")), contentLength: 3},
		{name: "known above limit", body: io.NopCloser(strings.NewReader("1234")), contentLength: 4, wantOversized: true},
		{name: "unknown zero", body: io.NopCloser(strings.NewReader("body")), contentLength: 0, wantUnknown: true},
		{name: "unknown negative", body: io.NopCloser(strings.NewReader("body")), contentLength: -1, wantUnknown: true},
	}

	policy := newTestPolicy(t, WithMaxRequestSize(3))
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := newTestPolicyGETRequest(t, "https://example.com")
			req.Body = tt.body
			req.ContentLength = tt.contentLength

			err := policy.Eval(req)
			switch {
			case tt.wantUnknown:
				var lengthErr *RequestBodyLengthError
				if !errors.As(err, &lengthErr) {
					t.Fatalf("expected RequestBodyLengthError, got %T: %v", err, err)
				}
				if lengthErr.ContentLength != tt.contentLength || lengthErr.Limit != 3 {
					t.Fatalf("unexpected request body length error: %#v", lengthErr)
				}
				assertNotPolicyError(t, err)
			case tt.wantOversized:
				var limitErr *RequestBodyLimitError
				if !errors.As(err, &limitErr) {
					t.Fatalf("expected RequestBodyLimitError, got %T: %v", err, err)
				}
				if limitErr.Size != tt.contentLength || limitErr.Limit != 3 {
					t.Fatalf("unexpected request body limit error: %#v", limitErr)
				}
				assertNotPolicyError(t, err)
			default:
				if err != nil {
					t.Fatalf("expected body to pass: %v", err)
				}
			}
		})
	}
}

func TestPolicyEvalAllowsUnknownBodyLengthWhenUnlimited(t *testing.T) {
	policy := newTestPolicy(t, WithUnlimitedRequestSize())
	for _, contentLength := range []int64{-1, 0} {
		req := newTestPolicyGETRequest(t, "https://example.com")
		req.Body = io.NopCloser(strings.NewReader("unknown"))
		req.ContentLength = contentLength

		if err := policy.Eval(req); err != nil {
			t.Fatalf("expected unknown content length %d to be allowed: %v", contentLength, err)
		}
	}
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
	client := newDefaultHTTPClient(policy, stdhttp.Client{
		Transport: testRoundTripper(func(*stdhttp.Request) (*stdhttp.Response, error) {
			roundTrips++
			return responseWithBody(
				stdhttp.StatusFound,
				"",
				stdhttp.Header{"Location": {"/next"}},
			), nil
		})},
	)

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
	err := newTestPolicy(t).Eval(newTestPolicyRequest(
		t,
		stdhttp.MethodConnect,
		"https://example.com",
	))
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
