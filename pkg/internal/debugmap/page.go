package debugmap

import (
	"errors"
	"slices"
)

// Read sorts only keys and looks up values only in the selected interval.
func Read[V any](values map[string]V, bounds Bounds, start, end int, visit func(string, V) error) error {
	if !bounds.Pageable(len(values)) || start < 0 || end < start || end > len(values) {
		return errors.ErrUnsupported
	}

	keys := make([]string, 0, len(values))

	for key := range values {
		keys = append(keys, key)
	}

	slices.Sort(keys)

	for _, key := range keys[start:end] {
		if err := visit(key, values[key]); err != nil {
			return err
		}
	}

	return nil
}
