package msgpack

import (
	"bytes"
	"context"
	"errors"
	"math"
	"testing"

	vmmsgpack "github.com/vmihailenco/msgpack/v5"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

type (
	lengthList struct {
		runtime.List
		length runtime.Int
	}

	lengthMap struct {
		runtime.Map
		length runtime.Int
	}
)

func (v lengthList) Length(context.Context) (runtime.Int, error) { return v.length, nil }
func (v lengthMap) Length(context.Context) (runtime.Int, error)  { return v.length, nil }

func TestMessagePackCollectionLength(t *testing.T) {
	maximum := min(uint64(math.MaxInt), uint64(math.MaxUint32))
	for _, length := range []runtime.Int{0, 1, runtime.Int(maximum)} {
		got, err := collectionLength(length)
		if err != nil || runtime.Int(got) != length {
			t.Errorf("collectionLength(%d) = %d, %v", length, got, err)
		}
	}

	for _, length := range []runtime.Int{-1, runtime.Int(maximum) + 1, math.MaxInt64} {
		for _, value := range []runtime.Value{
			lengthList{length: length}, lengthMap{length: length},
		} {
			var buffer bytes.Buffer
			// Nil embedded collections ensure traversal would panic if it were reached.
			err := (encoder{}).encodeValue(t.Context(), vmmsgpack.NewEncoder(&buffer), value)
			if !errors.Is(err, runtime.ErrRange) || buffer.Len() != 0 {
				t.Errorf("%T length %d: wrote %d bytes, error %v", value, length, buffer.Len(), err)
			}
		}
	}
}
