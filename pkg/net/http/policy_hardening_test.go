package http

import (
	"errors"
	stdhttp "net/http"
	"strings"
	"testing"
)

func TestPolicyAllowedMethods(t *testing.T) {
	for _, method := range []string{
		"",
		stdhttp.MethodGet,
		stdhttp.MethodHead,
		stdhttp.MethodPost,
		stdhttp.MethodPut,
		stdhttp.MethodPatch,
		stdhttp.MethodDelete,
		stdhttp.MethodOptions,
	} {
		t.Run("default_"+strings.TrimSpace(method), func(t *testing.T) {
			err := newTestPolicy(t).Eval(newTestPolicyRequest(t, method, "https://example.com"))
			if err != nil {
				t.Fatalf("expected method %q to be allowed, got %v", method, err)
			}
		})
	}

	for _, method := range []string{stdhttp.MethodConnect, stdhttp.MethodTrace, "CUSTOM-METHOD"} {
		t.Run("denied_"+method, func(t *testing.T) {
			err := newTestPolicy(t).Eval(newTestPolicyRequest(t, method, "https://example.com"))
			policyErr := requirePolicyError(t, err, PolicyTargetRequest)
			if policyErr.Subject != `method "`+method+`"` || policyErr.Reason != "method is not allowed" {
				t.Fatalf("unexpected method policy error: %#v", policyErr)
			}
		})
	}

	err := newTestPolicy(t, WithAllowedMethods("SET")).Eval(
		newTestPolicyRequest(t, "ſET", "https://example.com"),
	)
	if err == nil || errors.Is(err, ErrPolicyDenied) {
		t.Fatalf("expected non-ASCII method to remain a syntax error, got %v", err)
	}
}

func TestPolicyConfiguredAllowedMethodsReplaceDefaults(t *testing.T) {
	policy := newTestPolicy(t, WithAllowedMethods(" custom-method ", "connect", "CUSTOM-METHOD"))

	for _, method := range []string{"CUSTOM-METHOD", stdhttp.MethodConnect} {
		if err := policy.Eval(newTestPolicyRequest(t, method, "https://example.com")); err != nil {
			t.Fatalf("expected configured method %q to be allowed, got %v", method, err)
		}
	}

	for _, method := range []string{"", stdhttp.MethodGet, stdhttp.MethodTrace, "custom-method"} {
		err := policy.Eval(newTestPolicyRequest(t, method, "https://example.com"))
		requirePolicyError(t, err, PolicyTargetRequest)
	}

	err := policy.Eval(newTestPolicyRequest(t, " CONNECT ", "https://example.com"))
	var methodErr *InvalidMethodError
	if !errors.As(err, &methodErr) {
		t.Fatalf("expected exact standard method validation, got %T: %v", err, err)
	}
}

func TestPolicyRejectsURLCredentials(t *testing.T) {
	for _, rawURL := range []string{
		"https://user@example.com",
		"https://user:password@example.com",
		"https://@example.com",
		"https://:password@example.com",
		"https://user%40name:password@example.com",
	} {
		t.Run(rawURL, func(t *testing.T) {
			err := newTestPolicy(t, WithBlockedRequestHeaders("Authorization")).Eval(
				newTestPolicyGETRequest(t, rawURL),
			)
			policyErr := requirePolicyError(t, err, PolicyTargetRequest)
			if policyErr.Subject != "URL credentials" || policyErr.Reason != "URL user information is not allowed" {
				t.Fatalf("unexpected URL credential policy error: %#v", policyErr)
			}
			if strings.Contains(err.Error(), "password") {
				t.Fatalf("policy error leaked URL credentials: %v", err)
			}
		})
	}

	req := newTestPolicyGETRequest(t, "https://example.com")
	req.Header = stdhttp.Header{"Authorization": {"Bearer token"}}
	err := newTestPolicy(t).Eval(req)
	if err != nil {
		t.Fatalf("expected explicit Authorization to remain supported, got %v", err)
	}
}

func TestPolicyRequiresASCIIPunycodeHosts(t *testing.T) {
	for _, rawURL := range []string{"https://éxample.com", "https://K.example"} {
		err := newTestPolicy(t).Eval(newTestPolicyGETRequest(t, rawURL))
		policyErr := requirePolicyError(t, err, PolicyTargetRequest)
		if policyErr.Reason != "internationalized hostnames must use ASCII/punycode" {
			t.Fatalf("unexpected non-ASCII host reason: %q", policyErr.Reason)
		}
	}

	if err := newTestPolicy(t).Eval(newTestPolicyGETRequest(t, "https://xn--xample-9ua.com")); err != nil {
		t.Fatalf("expected punycode hostname to be allowed, got %v", err)
	}
}

func TestPolicyHostMatchingIsExactAndPortAware(t *testing.T) {
	policy := newTestPolicy(t, WithAllowedHosts("example.com"))
	for _, rawURL := range []string{"https://example.com", "https://example.com:8443"} {
		if err := policy.Eval(newTestPolicyGETRequest(t, rawURL)); err != nil {
			t.Fatalf("expected hostname-only rule to allow %q, got %v", rawURL, err)
		}
	}

	err := policy.Eval(newTestPolicyGETRequest(t, "https://api.example.com"))
	requirePolicyError(t, err, PolicyTargetRequest)

	portPolicy := newTestPolicy(t, WithAllowedHosts("example.com:8443"))
	if err := portPolicy.Eval(newTestPolicyGETRequest(t, "https://example.com:8443")); err != nil {
		t.Fatalf("expected matching port-specific rule to pass, got %v", err)
	}
	requirePolicyError(
		t,
		portPolicy.Eval(newTestPolicyGETRequest(t, "https://example.com:443")),
		PolicyTargetRequest,
	)
}

func TestPolicyRejectsBlockedAndReservedRequestHeaders(t *testing.T) {
	req := newTestPolicyGETRequest(t, "https://example.com")
	req.Header = stdhttp.Header{"AUTHORIZATION": nil}
	err := newTestPolicy(t, WithBlockedRequestHeaders("authorization")).Eval(req)
	policyErr := requirePolicyError(t, err, PolicyTargetRequest)
	if policyErr.Subject != `header "Authorization"` || policyErr.Reason != "request header is not allowed" {
		t.Fatalf("unexpected blocked-header error: %#v", policyErr)
	}

	reserved := []string{
		"Connection",
		"Content-Length",
		"Host",
		"Proxy-Connection",
		"Keep-Alive",
		"TE",
		"Trailer",
		"Transfer-Encoding",
		"Upgrade",
		"Proxy-Authorization",
		"Proxy-Authenticate",
	}
	for _, header := range reserved {
		t.Run(header, func(t *testing.T) {
			req := newTestPolicyGETRequest(t, "https://example.com")
			req.Header = stdhttp.Header{strings.ToLower(header): nil}
			err := newTestPolicy(t).Eval(req)
			policyErr := requirePolicyError(t, err, PolicyTargetRequest)
			if policyErr.Reason != "request header is reserved for the transport" {
				t.Fatalf("unexpected reserved-header error: %#v", policyErr)
			}
		})
	}
}

func TestPolicyErrorsRemainDistinctFromValidationErrors(t *testing.T) {
	err := newTestPolicy(t).Eval(newTestPolicyRequest(t, "BAD METHOD", "https://example.com"))
	if err == nil {
		t.Fatal("expected invalid method error")
	}
	if errors.Is(err, ErrPolicyDenied) {
		t.Fatalf("expected syntax validation error, got policy denial: %v", err)
	}
}

func requirePolicyError(t *testing.T, err error, target PolicyTarget) *PolicyError {
	t.Helper()

	if err == nil {
		t.Fatal("expected policy error")
	}
	if !errors.Is(err, ErrPolicyDenied) {
		t.Fatalf("expected ErrPolicyDenied, got %v", err)
	}

	var policyErr *PolicyError
	if !errors.As(err, &policyErr) {
		t.Fatalf("expected PolicyError, got %T: %v", err, err)
	}
	if policyErr.Target != target {
		t.Fatalf("expected policy target %q, got %q", target, policyErr.Target)
	}

	return policyErr
}
