package http

import "fmt"

// PolicyConfigurationError reports a malformed or contradictory policy option.
type PolicyConfigurationError struct {
	// Option identifies the option that supplied the invalid configuration.
	Option string
	// Value identifies the invalid non-secret input.
	Value string
	// Reason explains why the configuration is invalid.
	Reason string
}

func newPolicyConfigurationError(option, value, reason string) *PolicyConfigurationError {
	return &PolicyConfigurationError{
		Option: option,
		Value:  value,
		Reason: reason,
	}
}

// Error returns the human-readable policy configuration failure.
func (e *PolicyConfigurationError) Error() string {
	detail := e.sanitizedDetail()
	if detail == "" {
		return ErrInvalidPolicyConfiguration.Error()
	}

	return fmt.Sprintf("%s: %s", ErrInvalidPolicyConfiguration, detail)
}

func (e *PolicyConfigurationError) sanitizedDetail() string {
	if e == nil {
		return ""
	}

	if e.Value == "" {
		return fmt.Sprintf("%s: %s", e.Option, e.Reason)
	}

	return fmt.Sprintf("%s value %q: %s", e.Option, e.Value, e.Reason)
}

// Unwrap makes configuration failures detectable with errors.Is.
func (e *PolicyConfigurationError) Unwrap() error {
	return ErrInvalidPolicyConfiguration
}
