package http

import (
	"errors"
	"strings"
	"testing"
)

func TestMultiPolicyConfigurationErrorNilSafety(t *testing.T) {
	var nilErr *MultiPolicyConfigurationError
	if got := nilErr.Error(); got != ErrInvalidPolicyConfiguration.Error() {
		t.Fatalf("unexpected nil aggregate text: %q", got)
	}
	if got := nilErr.Unwrap(); got != nil {
		t.Fatalf("expected nil aggregate to unwrap to nil, got %#v", got)
	}
	if !nilErr.Is(ErrInvalidPolicyConfiguration) {
		t.Fatal("expected nil aggregate to match ErrInvalidPolicyConfiguration")
	}

	empty := &MultiPolicyConfigurationError{}
	if got := empty.Error(); got != ErrInvalidPolicyConfiguration.Error() {
		t.Fatalf("unexpected empty aggregate text: %q", got)
	}
	if !errors.Is(empty, ErrInvalidPolicyConfiguration) {
		t.Fatal("expected empty aggregate to match ErrInvalidPolicyConfiguration")
	}
}

func TestMultiPolicyConfigurationErrorFiltersNilChildren(t *testing.T) {
	first := newPolicyConfigurationError("WithTimeout", "-1ns", "must not be negative")
	second := newPolicyConfigurationError("WithAllowedSchemes", "://", "must be a valid URL scheme")
	err := newMultiPolicyConfigurationError([]*PolicyConfigurationError{nil, first, nil, second})

	multiErr := requireMultiPolicyConfigurationError(t, err)
	if len(multiErr.Errors) != 2 || multiErr.Errors[0] != first || multiErr.Errors[1] != second {
		t.Fatalf("unexpected filtered children: %#v", multiErr.Errors)
	}

	unwrapped := multiErr.Unwrap()
	if len(unwrapped) != 2 || unwrapped[0] != first || unwrapped[1] != second {
		t.Fatalf("unexpected unwrapped children: %#v", unwrapped)
	}

	rendered := multiErr.Error()
	if strings.Count(rendered, ErrInvalidPolicyConfiguration.Error()) != 1 ||
		!strings.Contains(rendered, first.sanitizedDetail()) ||
		!strings.Contains(rendered, second.sanitizedDetail()) {
		t.Fatalf("unexpected aggregate text: %q", rendered)
	}

	multiErr.Errors[0] = nil
	if got := multiErr.Unwrap(); len(got) != 1 || got[0] != second {
		t.Fatalf("expected public nil children to be ignored, got %#v", got)
	}
}

func TestNewMultiPolicyConfigurationErrorCollapsesSmallInputs(t *testing.T) {
	leaf := newPolicyConfigurationError("WithMaxRequestSize", "-1", "must not be negative")

	if err := newMultiPolicyConfigurationError(nil); err != nil {
		t.Fatalf("expected no errors to collapse to nil, got %v", err)
	}
	if err := newMultiPolicyConfigurationError([]*PolicyConfigurationError{nil, leaf, nil}); err != leaf {
		t.Fatalf("expected one non-nil error to remain a leaf, got %T: %v", err, err)
	}
}
