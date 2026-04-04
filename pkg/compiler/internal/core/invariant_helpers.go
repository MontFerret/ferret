package core

import "fmt"

func PanicInvariant(message string) {
	panic(&InvariantViolation{
		message: "compiler invariant violated: " + message,
	})
}

func PanicInvariantf(format string, args ...any) {
	PanicInvariant(fmt.Sprintf(format, args...))
}
