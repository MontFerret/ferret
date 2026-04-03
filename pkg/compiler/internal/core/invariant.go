package core

import "fmt"

// InvariantViolation marks a broken compiler-internal assumption.
// These panics should indicate bugs in compiler state management rather than user input errors.
type InvariantViolation struct {
	message string
}

func (e *InvariantViolation) Error() string {
	return e.message
}

func PanicInvariant(message string) {
	panic(&InvariantViolation{
		message: "compiler invariant violated: " + message,
	})
}

func PanicInvariantf(format string, args ...any) {
	PanicInvariant(fmt.Sprintf(format, args...))
}
