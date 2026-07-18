package http

import "strings"

// MultiPolicyConfigurationError reports multiple malformed or contradictory
// policy options in deterministic order.
type MultiPolicyConfigurationError struct {
	// Errors contains the individual policy configuration failures.
	Errors []*PolicyConfigurationError
}

func newMultiPolicyConfigurationError(errors []*PolicyConfigurationError) error {
	filtered := make([]*PolicyConfigurationError, 0, len(errors))
	for _, err := range errors {
		if err != nil {
			filtered = append(filtered, err)
		}
	}

	switch len(filtered) {
	case 0:
		return nil
	case 1:
		return filtered[0]
	default:
		return &MultiPolicyConfigurationError{Errors: filtered}
	}
}

// Error returns the ordered policy configuration failures with the shared
// sentinel prefix rendered once.
func (e *MultiPolicyConfigurationError) Error() string {
	if e == nil {
		return ErrInvalidPolicyConfiguration.Error()
	}

	var builder strings.Builder
	builder.WriteString(ErrInvalidPolicyConfiguration.Error())

	for _, err := range e.Errors {
		if err == nil {
			continue
		}

		builder.WriteString("\n- ")
		builder.WriteString(err.sanitizedDetail())
	}

	return builder.String()
}

// Unwrap exposes the individual configuration failures in deterministic order.
func (e *MultiPolicyConfigurationError) Unwrap() []error {
	if e == nil {
		return nil
	}

	errors := make([]error, 0, len(e.Errors))
	for _, err := range e.Errors {
		if err != nil {
			errors = append(errors, err)
		}
	}

	return errors
}

// Is makes the aggregate directly identifiable as invalid policy configuration,
// including when it contains no individual failures.
func (e *MultiPolicyConfigurationError) Is(target error) bool {
	return target == ErrInvalidPolicyConfiguration
}
