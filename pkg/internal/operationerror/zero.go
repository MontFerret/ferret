package operationerror

import (
	"errors"
	"fmt"
)

var (
	// ErrDivisionByZero identifies a failed runtime division without owning its user-facing diagnostic.
	ErrDivisionByZero = errors.New("division by zero")
	// ErrModuloByZero identifies a failed runtime modulo without owning its user-facing diagnostic.
	ErrModuloByZero = errors.New("modulo by zero")
)

// DivisionByZero preserves the operation's public cause while adding division identity.
func DivisionByZero(cause error) error {
	return fmt.Errorf("%w: %w", cause, ErrDivisionByZero)
}

// ModuloByZero preserves the operation's public cause while adding modulo identity.
func ModuloByZero(cause error) error {
	return fmt.Errorf("%w: %w", cause, ErrModuloByZero)
}
