package vm

import (
	"context"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
	"github.com/MontFerret/ferret/v2/pkg/vm/internal/data"
)

// Only these in-memory object implementations have debugger-safe accessors.
func debugObject(value runtime.Value) runtime.Map {
	switch v := value.(type) {
	case *runtime.Object:
		if v != nil {
			return v
		}
	case *data.FastObject:
		if v != nil {
			return v
		}
	}

	return nil
}

func inspectDebugMap(value runtime.Map, maxItems int) DebugValueInspection {
	length, _ := value.Length(context.Background())
	out := DebugValueInspection{Kind: DebugValueObject, Length: int(length)}

	if maxItems <= 0 || out.Length > maxItems {
		return out
	}

	keys, _ := value.Keys(context.Background())
	out.Items = make([]DebugValueItem, 0, out.Length)

	for i := runtime.Int(0); i < length; i++ {
		key, _ := keys.At(context.Background(), i)
		name, ok := key.(runtime.String)
		if !ok {
			return out
		}

		item, _ := value.Get(context.Background(), name)
		out.Items = append(out.Items, DebugValueItem{Key: name.String(), Value: item})
	}

	out.Complete = true

	return out
}
