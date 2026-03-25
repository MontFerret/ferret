package exec

import (
	"errors"
	"fmt"
	"strings"

	"github.com/MontFerret/ferret/v2/pkg/vm"
)

func IsRuntimeError(actual, expected any) error {
	err, ok := actual.(error)
	if !ok || err == nil {
		return errors.New("expected error, got nothing")
	}

	var rtErr *vm.RuntimeError
	if !errors.As(err, &rtErr) {
		return errors.New("expected a VM runtime error")
	}

	switch ex := expected.(type) {
	case *ExpectedRuntimeError:
		if ex.Message != "" && rtErr.Message != ex.Message {
			return fmt.Errorf("expected runtime error message '%s', got '%s'", ex.Message, rtErr.Message)
		}

		formatted := rtErr.Format()

		if ex.Format != "" && formatted != ex.Format {
			return fmt.Errorf("unexpected runtime error format\nexpected:\n%s\nactual:\n%s", ex.Format, formatted)
		}

		for _, needle := range ex.Contains {
			if !strings.Contains(formatted, needle) {
				return fmt.Errorf("expected formatted error to contain '%s'", needle)
			}
		}

		for _, needle := range ex.NotContains {
			if strings.Contains(formatted, needle) {
				return fmt.Errorf("expected formatted error to not contain '%s'", needle)
			}
		}
	case error:
		if ex.Error() != "" && rtErr.Message != ex.Error() {
			return fmt.Errorf("expected runtime error message '%s', got '%s'", ex.Error(), rtErr.Message)
		}
	case string:
		if ex != "" && rtErr.Message != ex {
			return fmt.Errorf("expected runtime error message '%s', got '%s'", ex, rtErr.Message)
		}
	}

	return nil
}
