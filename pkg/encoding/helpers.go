package encoding

import (
	"mime"
	"strings"
)

func normalizeContentType(contentType string) (string, error) {
	trimmed := strings.TrimSpace(contentType)
	if trimmed == "" {
		return "", ErrEmptyContentType
	}

	mediaType, _, err := mime.ParseMediaType(trimmed)
	if err == nil && mediaType != "" {
		trimmed = mediaType
	}

	normalized := strings.ToLower(strings.TrimSpace(trimmed))
	if normalized == "" {
		return "", ErrEmptyContentType
	}

	return normalized, nil
}
