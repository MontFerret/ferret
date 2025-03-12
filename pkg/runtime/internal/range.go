package internal

import (
	"context"
	"encoding/binary"
	"fmt"
	"hash/fnv"

	"github.com/wI2L/jettison"

	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values/types"
)

type Range struct {
	start int64
	end   int64
}

func NewRange(start, end int64) *Range {
	return &Range{start, end}
}

func (r *Range) Start() int64 {
	return r.start
}

func (r *Range) End() int64 {
	return r.end
}

func (r *Range) String() string {
	return fmt.Sprintf("%d..%d", r.start, r.end)
}

func (r *Range) Unwrap() interface{} {
	return []int64{r.start, r.end}
}

func (r *Range) Hash() uint64 {
	h := fnv.New64a()

	h.Write([]byte(types.Range.String()))
	h.Write([]byte(":"))

	startMultiplier := 1
	if r.start < 0 {
		h.Write([]byte("-"))
		startMultiplier = -1
	}

	bytes := make([]byte, 8)
	binary.LittleEndian.PutUint64(bytes, uint64(r.start*int64(startMultiplier)))
	h.Write(bytes)
	h.Write([]byte(".."))

	endMultiplier := 1
	if r.start < 0 {
		h.Write([]byte("-"))
		endMultiplier = -1
	}

	binary.LittleEndian.PutUint64(bytes, uint64(r.end*int64(endMultiplier)))
	h.Write(bytes)

	return h.Sum64()
}

func (r *Range) Copy() core.Value {
	return NewRange(r.start, r.end)
}

func (r *Range) Iterate(_ context.Context) (core.Iterator, error) {
	return NewRangeIterator(r), nil
}

func (r *Range) MarshalJSON() ([]byte, error) {
	start := r.start
	end := r.end

	var arr []int64

	if start <= end {
		arr = r.populateArray(start, end, r.calculateCapacity(start, end), true)
	} else {
		arr = r.populateArray(start, end, r.calculateCapacity(start, end), false)
	}

	return jettison.MarshalOpts(arr, jettison.NoHTMLEscaping())
}

func (r *Range) Compare(_ context.Context, other core.Value) (int64, error) {
	otherRange, ok := other.(*Range)

	if !ok {
		return types.Compare(types.Range, core.Reflect(other)), nil
	}

	if r.start == otherRange.start && r.end == otherRange.end {
		return 0, nil
	} else if r.start < otherRange.start || r.end < otherRange.end {
		return -1, nil
	} else {
		return 1, nil
	}
}

func (r *Range) calculateCapacity(start int64, end int64) int64 {
	var capacity int64
	if start <= end {
		if end < 0 {
			capacity = start + (end * -1) + 1
		} else {
			capacity = end - start + 1
		}
	} else {
		if start < 0 {
			capacity = end + (start * -1) + 1
		} else {
			capacity = start - end + 1
		}
	}
	return capacity
}

func (r *Range) populateArray(start int64, end int64, capacity int64, ascending bool) []int64 {
	arr := make([]int64, 0, capacity)

	if ascending {
		// start to end
		for i := start; i <= end; i++ {
			arr = append(arr, i)
		}
	} else {
		// end to start
		for i := start; i >= end; i-- {
			arr = append(arr, i)
		}
	}
	return arr
}
