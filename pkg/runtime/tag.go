package runtime

import (
	"reflect"
	"strings"
)

const TagName = "ferret"

// Tag returns the binding tag name for the provided struct field.
// It returns false when the field has no runtime tag or is explicitly ignored.
func Tag(field reflect.StructField) (string, bool) {
	raw := field.Tag.Get(TagName)

	if raw == "" {
		return "", false
	}

	return parseTag(raw, field.Name)
}

func parseTag(tag, fallback string) (string, bool) {
	name := strings.Split(tag, ",")[0]

	if name == "-" {
		return "", false
	}

	if name == "" {
		return fallback, true
	}

	return name, true
}
