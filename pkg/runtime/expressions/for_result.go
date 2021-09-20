package expressions

import (
	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
)

type ForResult struct {
	itemList    *values.Array
	hashTable   map[uint64]bool
	distinct    bool
	spread      bool
	passThrough bool
}

func NewForResult(capacity int) *ForResult {
	res := new(ForResult)
	res.itemList = values.NewArray(capacity)

	return res
}

func (f *ForResult) Distinct(distinct bool) *ForResult {
	f.distinct = distinct

	if f.distinct {
		f.hashTable = make(map[uint64]bool)
	} else {
		f.hashTable = nil
	}

	return f
}

func (f *ForResult) Spread(spread bool) *ForResult {
	f.spread = spread

	return f
}

func (f *ForResult) PassThrough(passThrough bool) *ForResult {
	f.passThrough = passThrough

	return f
}

func (f *ForResult) Push(value core.Value) {
	if f.passThrough {
		return
	}

	if f.distinct {
		// We need to check whether the value already exists in the result set
		hash := value.Hash()

		// if already exists
		// we skip it
		if f.hashTable[hash] {
			return
		}

		f.hashTable[hash] = true
	}

	if !f.spread {
		f.itemList.Push(value)

		return
	}

	elements, ok := value.(*values.Array)

	if !ok {
		f.itemList.Push(value)

		return
	}

	elements.ForEach(func(i core.Value, _ int) bool {
		f.Push(i)

		return true
	})
}

func (f *ForResult) ToArray() *values.Array {
	return f.itemList
}
