package runtime

import (
	"context"
	"fmt"
)

// MergeMapsInto copies source entries into dst, with later sources winning.
// Values are cloned when supported and otherwise copied according to Value.Copy.
// The caller owns dst; sources are borrowed. An error may leave dst partially updated.
func MergeMapsInto(ctx context.Context, dst Map, sources ...Map) error {
	return mergeMapsInto(ctx, dst, false, sources)
}

// MergeMapsDeepInto merges maps recursively, replacing other conflicting values.
// The caller must own dst and its nested maps exclusively. Incoming values are
// cloned or copied before insertion; sources are borrowed. Errors may leave dst
// partially updated. Arrays are replaced, never concatenated.
func MergeMapsDeepInto(ctx context.Context, dst Map, sources ...Map) error {
	return mergeMapsInto(ctx, dst, true, sources)
}

func mergeMapsInto(ctx context.Context, dst Map, deep bool, sources []Map) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	for _, src := range sources {
		err := src.ForEach(ctx, func(ctx context.Context, value, key Value) (Boolean, error) {
			if err := ctx.Err(); err != nil {
				return false, err
			}

			if err := mergeMapEntry(ctx, dst, key, value, deep); err != nil {
				return false, fmt.Errorf("key %q: %w", key.String(), err)
			}

			return true, nil
		})
		if err != nil {
			return err
		}
	}

	return nil
}

func mergeMapEntry(ctx context.Context, dst Map, key, value Value, deep bool) error {
	if incoming, ok := value.(Map); deep && ok {
		existing, found, err := dst.Lookup(ctx, key)
		if err != nil {
			return err
		}

		if nested, ok := existing.(Map); found && ok {
			return MergeMapsDeepInto(ctx, nested, incoming)
		}
	}

	copied, err := CloneOrCopy(ctx, value)
	if err != nil {
		return err
	}

	return dst.Set(ctx, key, copied)
}

// KeepMapKeys removes all but the selected string keys from dst.
// Missing and repeated keys have no effect. The caller owns dst; retained values
// are untouched. An error may leave dst partially updated.
func KeepMapKeys(ctx context.Context, dst Map, keys ...String) error {
	return filterMapKeys(ctx, dst, keys, true)
}

// OmitMapKeys removes selected string keys from dst.
// Missing and repeated keys have no effect. The caller owns dst; retained values
// are untouched. An error may leave dst partially updated.
func OmitMapKeys(ctx context.Context, dst Map, keys ...String) error {
	return filterMapKeys(ctx, dst, keys, false)
}

func filterMapKeys(ctx context.Context, dst Map, keys []String, keep bool) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	selected := make(map[String]struct{}, len(keys))
	for _, key := range keys {
		selected[key] = struct{}{}
	}

	current, err := dst.Keys(ctx)
	if err != nil {
		return err
	}

	// Keys may be a live host view. Finish reading it before the first removal.
	var removals []String
	err = current.ForEach(ctx, func(ctx context.Context, value Value, _ Int) (Boolean, error) {
		if err := ctx.Err(); err != nil {
			return false, err
		}

		key, err := CastString(value)
		if err != nil {
			return false, err
		}

		_, included := selected[key]
		if included != keep {
			removals = append(removals, key)
		}

		return true, nil
	})
	if err != nil {
		return err
	}

	for _, key := range removals {
		if err := ctx.Err(); err != nil {
			return err
		}

		if err := dst.RemoveKey(ctx, key); err != nil {
			return fmt.Errorf("key %q: %w", key, err)
		}
	}

	return nil
}
