package sdk

import (
	"reflect"
	"strings"
)

const (
	runtimeTagName = "ferret"
	jsonTagName    = "json"
)

// Tag returns the binding tag name for the provided struct field.
// It also falls back to "json" tag if "ferret" tag is not present.
// It returns false when the field has no runtime tag or is explicitly ignored.
func Tag(field reflect.StructField) (string, bool) {
	raw, ok := field.Tag.Lookup(runtimeTagName)
	if ok {
		return parseTag(raw, field.Name)
	}

	raw, ok = field.Tag.Lookup(jsonTagName)
	if ok {
		return parseTag(raw, field.Name)
	}

	return "", false
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
