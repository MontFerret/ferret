package vm_test

import (
	"errors"
	"fmt"
	"strings"

	"github.com/MontFerret/ferret/v2/pkg/vm"
)

type ExpectedRuntimeError struct {
	Message     string
	Format      string
	Contains    []string
	NotContains []string
}

func ShouldBeRuntimeError(actual any, expected ...any) string {
	err, ok := actual.(error)
	if !ok || err == nil {
		return "expected a runtime error"
	}

	var rtErr *vm.RuntimeError
	if !errors.As(err, &rtErr) {
		return "expected a VM runtime error"
	}

	if len(expected) == 0 || expected[0] == nil {
		return "expected runtime error metadata"
	}

	switch ex := expected[0].(type) {
	case *ExpectedRuntimeError:
		if ex.Message != "" && rtErr.Message != ex.Message {
			return fmt.Sprintf("expected runtime error message '%s', got '%s'", ex.Message, rtErr.Message)
		}

		formatted := rtErr.Format()

		if ex.Format != "" && formatted != ex.Format {
			return fmt.Sprintf("unexpected runtime error format\nexpected:\n%s\nactual:\n%s", ex.Format, formatted)
		}

		for _, needle := range ex.Contains {
			if !strings.Contains(formatted, needle) {
				return fmt.Sprintf("expected formatted error to contain '%s'", needle)
			}
		}

		for _, needle := range ex.NotContains {
			if strings.Contains(formatted, needle) {
				return fmt.Sprintf("expected formatted error to not contain '%s'", needle)
			}
		}
	case error:
		if ex.Error() != "" && rtErr.Message != ex.Error() {
			return fmt.Sprintf("expected runtime error message '%s', got '%s'", ex.Error(), rtErr.Message)
		}
	case string:
		if ex != "" && rtErr.Message != ex {
			return fmt.Sprintf("expected runtime error message '%s', got '%s'", ex, rtErr.Message)
		}
	}

	return ""
}
