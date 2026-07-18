package http

import "fmt"

// PolicyError describes an HTTP request rejected by the access policy.
type PolicyError struct {
	// Target identifies whether the original request or a redirect was denied.
	Target PolicyTarget
	// Subject identifies the method, URL component, header, host, or address.
	Subject string
	// Reason explains why the policy denied the subject.
	Reason string
}

// Error returns the human-readable policy denial.
func (e *PolicyError) Error() string {
	if e == nil {
		return ErrPolicyDenied.Error()
	}

	return fmt.Sprintf("http: %s blocked by access policy: %s: %s", e.Target, e.Subject, e.Reason)
}

// Unwrap makes policy denials detectable with errors.Is.
func (e *PolicyError) Unwrap() error {
	return ErrPolicyDenied
}

func newPolicyError(target PolicyTarget, subject, reason string) error {
	return &PolicyError{
		Target:  target,
		Subject: subject,
		Reason:  reason,
	}
}
