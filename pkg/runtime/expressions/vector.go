package expressions

import (
	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
	"math"
)

// Vector TODO: Finish and add unit tests
type Vector struct {
	slices   map[int][]core.Value
	capacity int
	current  int
}

func NewVector(capacity int) *Vector {
	el := new(Vector)
	el.slices = make(map[int][]core.Value)
	el.capacity = capacity
	el.current = -1

	return el
}

func (vec *Vector) Push(value core.Value) *Vector {
	slice := vec.getCurrentSlice()
	slice[len(slice)] = value

	return vec
}

func (vec *Vector) Get(index int) core.Value {
	var sliceIndex int
	var itemIndex int

	// if it's within a capacity
	if index < vec.capacity {
		sliceIndex = 0
		itemIndex = index
	} else {
		sliceIndex = int(math.Floor(float64(index / vec.capacity)))
		itemIndex = index % vec.capacity
	}

	// if out of range
	if sliceIndex >= len(vec.slices) {
		return values.None
	}

	return vec.slices[sliceIndex][itemIndex]
}

func (vec *Vector) ForEach(predicate func(item core.Value, index int)) {
	if len(vec.slices) == 0 {
		return
	}

	lastIndex := vec.capacity * len(vec.slices)

	for i := 0; i <= lastIndex; i++ {
		itemIndex := i % vec.capacity
		sliceIndex := int(math.Floor(float64(i / vec.capacity)))

		predicate(vec.slices[sliceIndex][itemIndex], i)
	}
}

func (vec *Vector) ToSlice() []core.Value {
	out := make([]core.Value, vec.capacity*len(vec.slices))

	if len(vec.slices) == 0 {
		return out
	}

	lastIndex := vec.capacity * len(vec.slices)

	for i := 0; i <= lastIndex; i++ {
		itemIndex := i % vec.capacity
		sliceIndex := int(math.Floor(float64(i / vec.capacity)))

		out[i] = vec.slices[sliceIndex][itemIndex]
	}

	return out
}

func (vec *Vector) getCurrentSlice() []core.Value {
	if vec.current < 0 {
		return vec.newSlice()
	}

	current := vec.slices[vec.current]

	if len(current) < vec.capacity {
		return current
	}

	return vec.newSlice()
}

func (vec *Vector) newSlice() []core.Value {
	slice := make([]core.Value, 0, vec.capacity)

	vec.current++
	vec.slices[vec.current] = slice

	return slice
}
