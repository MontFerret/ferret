package diagnostics

import (
	"errors"
	"strings"
)

func FormatMessage(str string) string {
	// Capitalize first letter and ensure it ends with a period.
	if str == "" {
		return ""
	}

	str = strings.TrimSpace(str)
	str = strings.ToLower(str)

	if len(str) > 0 {
		str = strings.ToUpper(str[:1]) + str[1:]
	}

	return str
}

func Unwrap(err error) (error, error) {
	if err == nil {
		return nil, nil
	}

	wrapped := errors.Unwrap(err)

	if wrapped == nil {
		return nil, err
	}

	originalMsg := err.Error()
	wrappedMsg := wrapped.Error()

	// Filter where the wrapped message ends in the original
	if idx := strings.Index(originalMsg, wrappedMsg); idx != -1 {
		// Extract everything after the wrapped message
		rest := originalMsg[idx+len(wrappedMsg):]
		rest = strings.TrimPrefix(rest, ": ")

		if rest != "" {
			return wrapped, errors.New(rest)
		}
	}

	return nil, errors.New(originalMsg)
}
