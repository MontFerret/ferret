package assert

import (
	"errors"
	"fmt"

	"github.com/MontFerret/ferret/v2/pkg/diagnostics"
	"github.com/goccy/go-json"
	"github.com/google/go-cmp/cmp"
)

func Nil(actual any) error {
	if actual == nil {
		return nil
	}

	return fmt.Errorf("actual value is not nil: %v", actual)
}

func NotNil(actual any) error {
	if actual != nil {
		return nil
	}

	return fmt.Errorf("actual value is nil")
}

func Equal(actual, expected any) error {
	if numericEqual(actual, expected) {
		return nil
	}

	if !cmp.Equal(actual, expected) {
		return fmt.Errorf("actual value is not equal to expected value: %s", cmp.Diff(actual, expected))
	}

	return nil
}

func Error(actual any) error {
	if actual == nil {
		return errors.New("actual value is nil")
	}

	_, ok := actual.(error)

	if !ok {
		return fmt.Errorf("actual value is not an error: %T", actual)
	}

	return nil
}

func HaveSameItems(actual, expected any) error {
	actualArr := actual.([]any)
	expectedArr := expected.([]any)

	actualSize := len(actualArr)
	expectedSize := len(expectedArr)

	if actualSize != expectedSize {
		return fmt.Errorf("expected %d items, got %d", expectedSize, actualSize)
	}

	for i := 0; i < actualSize; i++ {
		if err := Equal(actualArr[i], expectedArr[i]); err != nil {
			return fmt.Errorf("actual item %d is different from expected item %d: %w", i, i, err)
		}
	}

	return nil
}

// EqualJSON compares actual and expected by JSON-encoding both values.
// This normalizes numeric types (e.g., int vs float64) and map ordering.
func EqualJSON(actual, expected any) error {
	actualJSON, err := canonicalJSON(actual)
	if err != nil {
		return fmt.Errorf("actual value is not JSON-comparable: %v", err)
	}

	expectedJSON, err := canonicalJSON(expected)
	if err != nil {
		return fmt.Errorf("expected value is not JSON-comparable: %v", err)
	}

	return Equal(actualJSON, expectedJSON)
}

func numericEqual(actual, expected any) bool {
	actualFloat, actualOK := toFloat64(actual)
	expectedFloat, expectedOK := toFloat64(expected)

	return actualOK && expectedOK && actualFloat == expectedFloat
}

func toFloat64(value any) (float64, bool) {
	switch v := value.(type) {
	case int:
		return float64(v), true
	case int8:
		return float64(v), true
	case int16:
		return float64(v), true
	case int32:
		return float64(v), true
	case int64:
		return float64(v), true
	case uint:
		return float64(v), true
	case uint8:
		return float64(v), true
	case uint16:
		return float64(v), true
	case uint32:
		return float64(v), true
	case uint64:
		return float64(v), true
	case float32:
		return float64(v), true
	case float64:
		return v, true
	case json.Number:
		f, err := v.Float64()
		return f, err == nil
	default:
		return 0, false
	}
}

func canonicalJSON(value any) (string, error) {
	data, err := toJSONBytes(value)
	if err != nil {
		return "", err
	}

	var normalized any
	if err := json.Unmarshal(data, &normalized); err != nil {
		return "", err
	}

	canonical, err := json.Marshal(normalized)
	if err != nil {
		return "", err
	}

	return string(canonical), nil
}

func toJSONBytes(value any) ([]byte, error) {
	switch v := value.(type) {
	case json.RawMessage:
		return v, nil
	case []byte:
		if json.Valid(v) {
			return v, nil
		}
		return json.Marshal(string(v))
	case string:
		if json.Valid([]byte(v)) {
			return []byte(v), nil
		}
		return json.Marshal(v)
	default:
		return json.Marshal(v)
	}
}

func DiagnosticError(a, e any) error {
	if a == nil {
		return errors.New("expected a compilation error")
	}

	if e == nil {
		return errors.New("expected a compilation error")
	}

	actual, ok := a.(*diagnostics.Diagnostic)

	if !ok {
		return fmt.Errorf("actual value is not diagnostics.Diagnostic: %T", a)
	}

	expected, ok := CastExpectedError(e)

	if !ok {
		return fmt.Errorf("expected value is not *ExpectedError: %T", e)
	}

	if expected.Kind != "" && actual.Kind != expected.Kind {
		return fmt.Errorf("expected error kind %s, got %s", expected.Kind, actual.Kind)
	}

	if expected.Message != "" && actual.Message != expected.Message {
		return fmt.Errorf("expected error message '%s', got '%s'", expected.Message, actual.Message)
	}

	if expected.Hint != "" && actual.Hint != expected.Hint {
		return fmt.Errorf("expected error hint '%s', got '%s'", expected.Hint, actual.Hint)
	}

	if expected.Format != "" {
		actualFormat := actual.Format()

		if !cmp.Equal(actualFormat, expected.Format) {
			return fmt.Errorf("expected error format:\n%s\ngot:\n%s\n\nDiff:\n%s", expected.Format, actualFormat, cmp.Diff(expected.Format, actualFormat))
		}
	}

	return nil
}

func DiagnosticErrors(a, e any) error {
	if a == nil {
		return errors.New("expected a multi compilation error")
	}

	if e == nil {
		return errors.New("expected a multi compilation error")
	}

	actual, ok := a.(*diagnostics.Diagnostics[*diagnostics.Diagnostic])

	if !ok {
		return fmt.Errorf("actual value is not diagnostics.Diagnostics[*diagnostics.Diagnostic]: %T", a)
	}

	expected, ok := CastExpectedMultiError(e)

	if !ok {
		return fmt.Errorf("expected value is not *ExpectedMultiError: %T", e)
	}

	if expected.Number > 0 && actual.Size() != expected.Number {
		return fmt.Errorf("expected %d errors, got %d", expected.Number, actual.Size())
	}

	if len(expected.Errors) > 0 {
		for i, err := range actual.Errors() {
			if i >= len(expected.Errors) {
				break
			}

			msg := DiagnosticError(err, expected.Errors[i])

			if msg != nil {
				return msg
			}
		}
	}

	return nil
}
