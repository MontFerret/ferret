package msgpack

import (
	"math"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

// MessagePack collection headers are uint32, but the encoder accepts native int.
func collectionLength(length runtime.Int) (int, error) {
	size, ok := runtime.ToNativeInt(length)
	if !ok || length < 0 || uint64(length) > math.MaxUint32 {
		return 0, runtime.Error(runtime.ErrRange, "collection length exceeds MessagePack capacity")
	}

	return size, nil
}
