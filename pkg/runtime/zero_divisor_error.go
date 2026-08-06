package runtime

import "fmt"

func divisionByZeroError() error {
	return fmt.Errorf("%w: %w", ErrInvalidOperation, ErrDivisionByZero)
}

func moduloByZeroError() error {
	return fmt.Errorf("%w: %w", ErrInvalidOperation, ErrModuloByZero)
}
