package core

// InvariantViolation marks a broken compiler-internal assumption.
// These panics should indicate bugs in compiler state management rather than user input errors.
type InvariantViolation struct {
	message string
}

func (e *InvariantViolation) Error() string {
	return e.message
}
