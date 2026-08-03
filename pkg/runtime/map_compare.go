package runtime

import (
	"context"
	"sort"
)

func sortedMapKeys(ctx context.Context, value Map) ([]string, error) {
	keys, err := value.Keys(ctx)
	if err != nil {
		return nil, err
	}

	length, err := keys.Length(ctx)
	if err != nil {
		return nil, err
	}

	result := make([]string, 0, length)
	for idx := ZeroInt; idx < length; idx++ {
		key, err := keys.At(ctx, idx)
		if err != nil {
			return nil, err
		}

		stringKey, err := CastString(key)
		if err != nil {
			return nil, err
		}

		result = append(result, string(stringKey))
	}

	sort.Strings(result)

	return result, nil
}
