package values

import (
	"context"
	"encoding/binary"
	"fmt"
	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values/types"
	"github.com/wI2L/jettison"
	"hash/fnv"
)

type Range struct {
	start uint64
	end   uint64
}

func NewRange(start, end uint64) *Range {
	return &Range{start, end}
}
func (r *Range) Start() uint64 {
	return r.start
}

func (r *Range) End() uint64 {
	return r.end
}

func (r *Range) MarshalJSON() ([]byte, error) {
	arr := make([]uint64, r.end-r.start+1)

	for i := r.start; i <= r.end; i++ {
		arr[i-r.start] = i
	}

	return jettison.MarshalOpts(arr, jettison.NoHTMLEscaping())
}

func (r *Range) String() string {
	return fmt.Sprintf("%d..%d", r.start, r.end)
}

func (r *Range) Compare(other core.Value) int64 {
	otherRange, ok := other.(*Range)

	if !ok {
		return types.Compare(types.Range, core.Reflect(other))
	}

	if r.start == otherRange.start && r.end == otherRange.end {
		return 0
	} else if r.start < otherRange.start || r.end < otherRange.end {
		return -1
	} else {
		return 1
	}
}

func (r *Range) Unwrap() interface{} {
	return []uint64{r.start, r.end}
}

func (r *Range) Hash() uint64 {
	h := fnv.New64a()

	h.Write([]byte(types.Range.String()))
	h.Write([]byte(":"))
	bytes := make([]byte, 8)
	binary.LittleEndian.PutUint64(bytes, r.start)
	h.Write(bytes)
	h.Write([]byte(".."))
	binary.LittleEndian.PutUint64(bytes, r.end)
	h.Write(bytes)

	return h.Sum64()
}

func (r *Range) Copy() core.Value {
	return NewRange(r.start, r.end)
}

func (r *Range) Iterate(_ context.Context) (core.Iterator, error) {
	return NewRangeIterator(r), nil
}
