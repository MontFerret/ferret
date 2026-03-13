package runtime

import (
	"context"
	"encoding/binary"
	"fmt"
	"hash/fnv"

	"github.com/wI2L/jettison"
)

type Range struct {
	start Int
	end   Int
}

func NewRange(start, end Int) *Range {
	return &Range{start, end}
}

func (r *Range) Start() Int {
	return r.start
}

func (r *Range) End() Int {
	return r.end
}

func (r *Range) String() string {
	return fmt.Sprintf("%d..%d", r.start, r.end)
}

func (r *Range) Hash() uint64 {
	h := fnv.New64a()

	h.Write([]byte("range"))
	h.Write([]byte(":"))

	var startMultiplier Int
	startMultiplier = 1
	if r.start < 0 {
		h.Write([]byte("-"))
		startMultiplier = -1
	}

	bytes := make([]byte, 8)
	binary.LittleEndian.PutUint64(bytes, uint64(r.start*startMultiplier))
	h.Write(bytes)
	h.Write([]byte(".."))

	var endMultiplier Int
	endMultiplier = 1

	if r.start < 0 {
		h.Write([]byte("-"))
		endMultiplier = -1
	}

	binary.LittleEndian.PutUint64(bytes, uint64(r.end*endMultiplier))
	h.Write(bytes)

	return h.Sum64()
}

func (r *Range) Copy() Value {
	return NewRange(r.start, r.end)
}

func (r *Range) Iterate(_ context.Context) (Iterator, error) {
	return NewRangeIterator(r), nil
}

func (r *Range) MarshalJSON() ([]byte, error) {
	start := r.start
	end := r.end

	var arr []Int

	if start <= end {
		arr = r.populateArray(start, end, r.calculateCapacity(start, end), true)
	} else {
		arr = r.populateArray(start, end, r.calculateCapacity(start, end), false)
	}

	return jettison.MarshalOpts(arr, jettison.NoHTMLEscaping())
}

func (r *Range) Compare(_ context.Context, other Value) (int, error) {
	otherRange, ok := other.(*Range)

	if !ok {
		return CompareTypes(r, other), nil
	}

	if r.start == otherRange.start && r.end == otherRange.end {
		return 0, nil
	}

	if r.start < otherRange.start || r.end < otherRange.end {
		return -1, nil
	}

	return 1, nil
}

func (r *Range) calculateCapacity(start, end Int) Int {
	var capacity Int
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

func (r *Range) populateArray(start, end, capacity Int, ascending bool) []Int {
	arr := make([]Int, 0, capacity)

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
