package runtime

import (
	"context"
	"encoding/binary"
	"fmt"
	"hash/fnv"
	"math"

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

func (r *Range) Length(_ context.Context) (Int, error) {
	var distance uint64
	if r.start <= r.end {
		distance = uint64(r.end) - uint64(r.start)
	} else {
		distance = uint64(r.start) - uint64(r.end)
	}

	if distance >= uint64(math.MaxInt64) {
		return ZeroInt, Error(ErrRange, "range length exceeds runtime.Int capacity")
	}

	return Int(distance + 1), nil
}

func (r *Range) MarshalJSON() ([]byte, error) {
	start := r.start
	end := r.end

	capacity, err := r.Length(context.Background())
	if err != nil {
		return nil, err
	}

	arr := r.populateArray(start, capacity, start <= end)

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

func (r *Range) populateArray(start, capacity Int, ascending bool) []Int {
	arr := make([]Int, 0, capacity)

	for offset := Int(0); offset < capacity; offset++ {
		if ascending {
			arr = append(arr, start+offset)
		} else {
			arr = append(arr, start-offset)
		}
	}

	return arr
}
