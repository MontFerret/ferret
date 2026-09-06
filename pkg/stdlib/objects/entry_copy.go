package objects

import (
	"context"
	"fmt"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

func setCopiedEntry(ctx context.Context, dst runtime.KeyWritable, key runtime.String, value runtime.Value) error {
	copied, err := runtime.CloneOrCopy(ctx, value)
	if err != nil {
		return fmt.Errorf("key %q: %w", key, err)
	}

	if err := dst.Set(ctx, key, copied); err != nil {
		return fmt.Errorf("key %q: %w", key, err)
	}

	return nil
}
