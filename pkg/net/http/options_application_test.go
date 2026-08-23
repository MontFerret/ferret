package http

import (
	"errors"
	stdhttp "net/http"
	"testing"
	"time"

	"github.com/ziflex/go-options"
)

func TestPolicyOptionsApplyInOrder(t *testing.T) {
	policy, err := NewPolicy(
		WithTimeout(time.Second),
		WithTimeout(2*time.Second),
		WithMaxRedirects(2),
		WithMaxRedirects(3),
		WithAllowedSchemes("http"),
		WithAllowedSchemes("https"),
		WithAllowLocalhost(false),
		WithAllowLocalhost(true),
	)
	if err != nil {
		t.Fatalf("NewPolicy() error = %v", err)
	}

	if policy.timeout != 2*time.Second {
		t.Fatalf("timeout = %s, want %s", policy.timeout, 2*time.Second)
	}
	if policy.maxRedirects != 3 {
		t.Fatalf("max redirects = %d, want 3", policy.maxRedirects)
	}
	if len(policy.allowedSchemes) != 1 || policy.allowedSchemes[0] != "https" {
		t.Fatalf("allowed schemes = %v, want [https]", policy.allowedSchemes)
	}
	if !policy.allowLocalhost {
		t.Fatal("allow localhost = false, want true")
	}
}

func TestPolicyConstructorsReturnJoinedOptionErrors(t *testing.T) {
	first := errors.New("first option failed")
	second := errors.New("second option failed")
	setters := []PolicyOption{
		func(_ *Policy) error {
			return first
		},
		nil,
		func(_ *Policy) error {
			return second
		},
	}

	tests := []struct {
		construct func(...PolicyOption) (bool, error)
		name      string
	}{
		{
			name: "NewPolicy",
			construct: func(setters ...PolicyOption) (bool, error) {
				policy, err := NewPolicy(setters...)

				return policy != nil, err
			},
		},
		{
			name: "New",
			construct: func(setters ...PolicyOption) (bool, error) {
				client, err := New(setters...)

				return client != nil, err
			},
		},
		{
			name: "NewWithTransport",
			construct: func(setters ...PolicyOption) (bool, error) {
				client, err := NewWithTransport(stdhttp.DefaultTransport, setters...)

				return client != nil, err
			},
		},
		{
			name: "NewWithClient",
			construct: func(setters ...PolicyOption) (bool, error) {
				client, err := NewWithClient(&stdhttp.Client{}, setters...)

				return client != nil, err
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			constructed, err := tt.construct(setters...)
			if constructed {
				t.Fatal("constructor returned a value, want nil")
			}
			if !errors.Is(err, first) || !errors.Is(err, second) {
				t.Fatalf("constructor error = %v, want both option errors", err)
			}
			if errors.Is(err, ErrInvalidPolicyConfiguration) {
				t.Fatalf("constructor error = %v, want non-policy application errors", err)
			}
		})
	}
}

func TestNewPolicyClassifiesMixedOptionErrors(t *testing.T) {
	want := errors.New("option failed")
	policy, err := NewPolicy(
		WithMaxResponseSize(-1),
		func(_ *Policy) error {
			return want
		},
	)
	if policy != nil {
		t.Fatalf("NewPolicy() policy = %#v, want nil", policy)
	}
	if !errors.Is(err, want) {
		t.Fatalf("NewPolicy() error = %v, want option error", err)
	}
	if !errors.Is(err, ErrInvalidPolicyConfiguration) {
		t.Fatalf("NewPolicy() error = %v, want policy validation classification", err)
	}

	var validationErr options.ValidationError
	if !errors.As(err, &validationErr) {
		t.Fatalf("NewPolicy() error = %T, want options.ValidationError", err)
	}
	if validationErr.Field != "max response size" || validationErr.Value != "-1" {
		t.Fatalf("NewPolicy() validation error = %+v", validationErr)
	}
}
