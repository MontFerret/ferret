package http

import (
	"errors"
	"math"
	stdhttp "net/http"
	"net/netip"
	"testing"
	"time"
)

func TestPolicyEnforcementAccessors(t *testing.T) {
	defaults := newTestPolicy(t)
	if defaults.Timeout() != defaultTimeout {
		t.Fatalf("expected timeout %s, got %s", defaultTimeout, defaults.Timeout())
	}
	if defaults.MaxResponseSize() != defaultMaxResponseSize {
		t.Fatalf(
			"expected response limit %d, got %d",
			defaultMaxResponseSize,
			defaults.MaxResponseSize(),
		)
	}
	if defaults.MaxResponseHeaderSize() != defaultMaxResponseHeaderSize {
		t.Fatalf(
			"expected response header limit %d, got %d",
			defaultMaxResponseHeaderSize,
			defaults.MaxResponseHeaderSize(),
		)
	}

	custom := newTestPolicy(
		t,
		WithTimeout(time.Second),
		WithUnlimitedResponseSize(),
		WithMaxResponseHeaderSize(2048),
	)
	if custom.Timeout() != time.Second {
		t.Fatalf("expected timeout %s, got %s", time.Second, custom.Timeout())
	}
	if custom.MaxResponseSize() != 0 {
		t.Fatalf("expected unlimited response size, got %d", custom.MaxResponseSize())
	}
	if custom.MaxResponseHeaderSize() != 2048 {
		t.Fatalf("expected response header limit 2048, got %d", custom.MaxResponseHeaderSize())
	}

	noTimeout := newTestPolicy(t, WithNoTimeout())
	if noTimeout.Timeout() != 0 {
		t.Fatalf("expected disabled timeout, got %s", noTimeout.Timeout())
	}
}

func TestPolicyCheckRedirect(t *testing.T) {
	redirect := newTestPolicyGETRequest(t, "https://example.com/next")

	t.Run("disabled", func(t *testing.T) {
		policy := newTestPolicy(t, WithFollowRedirects(false))
		if err := policy.CheckRedirect(redirect, nil); !errors.Is(err, stdhttp.ErrUseLastResponse) {
			t.Fatalf("expected ErrUseLastResponse, got %v", err)
		}
	})

	t.Run("exact limit", func(t *testing.T) {
		policy := newTestPolicy(t, WithMaxRedirects(2))
		if err := policy.CheckRedirect(
			redirect,
			[]*stdhttp.Request{{}, {}},
		); err != nil {
			t.Fatalf("expected redirect at configured limit to pass: %v", err)
		}
	})

	t.Run("over limit", func(t *testing.T) {
		policy := newTestPolicy(t, WithMaxRedirects(2))
		err := policy.CheckRedirect(
			redirect,
			[]*stdhttp.Request{{}, {}, {}},
		)
		var limitErr *RedirectLimitError
		if !errors.As(err, &limitErr) || limitErr.Limit != 2 {
			t.Fatalf("expected redirect limit error for 2 redirects, got %#v", err)
		}
	})

	t.Run("blocked target", func(t *testing.T) {
		policy := newTestPolicy(t)
		err := policy.CheckRedirect(
			newTestPolicyGETRequest(t, "http://127.0.0.1/private"),
			[]*stdhttp.Request{{}},
		)
		requirePolicyError(t, err, PolicyTargetRedirect)
	})

	t.Run("nil request", func(t *testing.T) {
		tests := []struct {
			name   string
			policy *Policy
			via    []*stdhttp.Request
		}{
			{
				name:   "enabled",
				policy: newTestPolicy(t),
				via:    []*stdhttp.Request{{}},
			},
			{
				name:   "disabled",
				policy: newTestPolicy(t, WithFollowRedirects(false)),
			},
			{
				name:   "over limit",
				policy: newTestPolicy(t, WithMaxRedirects(0)),
				via:    []*stdhttp.Request{{}},
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				err := tt.policy.CheckRedirect(nil, tt.via)
				if !errors.Is(err, ErrNilRequest) {
					t.Fatalf("expected ErrNilRequest, got %v", err)
				}
			})
		}
	})
}

func TestPolicyEvalConnection(t *testing.T) {
	tests := []struct {
		address netip.Addr
		name    string
		wantErr bool
	}{
		{name: "public", address: netip.MustParseAddr("93.184.216.34")},
		{name: "invalid", address: netip.Addr{}, wantErr: true},
		{name: "loopback", address: netip.MustParseAddr("127.0.0.1"), wantErr: true},
		{name: "private", address: netip.MustParseAddr("10.0.0.1"), wantErr: true},
		{name: "link local", address: netip.MustParseAddr("169.254.169.254"), wantErr: true},
		{name: "nat64 private", address: netip.MustParseAddr("64:ff9b::10.0.0.1"), wantErr: true},
	}

	policy := newTestPolicy(t)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := policy.EvalConnection(tt.address)
			if !tt.wantErr {
				if err != nil {
					t.Fatalf("expected concrete address to pass: %v", err)
				}
				return
			}

			requirePolicyError(t, err, PolicyTargetConnection)
		})
	}
}

func TestPolicyEvalResponseSize(t *testing.T) {
	policy := newTestPolicy(t, WithMaxResponseSize(3))
	for _, size := range []int64{-1, 0, 2, 3} {
		if err := policy.EvalResponseSize(size); err != nil {
			t.Fatalf("expected response size %d to pass: %v", size, err)
		}
	}

	err := policy.EvalResponseSize(4)
	limitErr := requireResponseBodyLimitError(t, err)
	if limitErr.Size != 4 || limitErr.Limit != 3 {
		t.Fatalf("unexpected response limit error: %#v", limitErr)
	}

	unlimited := newTestPolicy(t, WithUnlimitedResponseSize())
	if err := unlimited.EvalResponseSize(math.MaxInt64); err != nil {
		t.Fatalf("expected unlimited response size to pass: %v", err)
	}
}

func TestPolicyReadResponseBodyNilAndOwnership(t *testing.T) {
	policy := newTestPolicy(t, WithMaxResponseSize(3))
	body, err := policy.ReadResponseBody(nil)
	if err != nil {
		t.Fatalf("read nil body: %v", err)
	}
	if body != nil {
		t.Fatalf("expected nil materialized body, got %q", body)
	}

	tracked := newTrackedResponseBody(newTerminalErrorReader("ok", nil))
	body, err = policy.ReadResponseBody(tracked)
	if err != nil {
		t.Fatalf("read tracked body: %v", err)
	}
	if string(body) != "ok" {
		t.Fatalf("expected body %q, got %q", "ok", body)
	}
	if tracked.closed.Load() {
		t.Fatal("expected caller to retain response-body ownership")
	}
}
