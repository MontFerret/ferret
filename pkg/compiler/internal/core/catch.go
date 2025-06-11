package core

import (
	"github.com/MontFerret/ferret/pkg/vm"
)

type CatchStack struct {
	entries []vm.Catch
}

func NewCatchStack() *CatchStack {
	return &CatchStack{
		entries: make([]vm.Catch, 0),
	}
}

func (cs *CatchStack) Push(start, end, jump int) {
	cs.entries = append(cs.entries, vm.Catch{start, end, jump})
}

func (cs *CatchStack) Pop() {
	if len(cs.entries) > 0 {
		cs.entries = cs.entries[:len(cs.entries)-1]
	}
}

func (cs *CatchStack) Find(pos int) (vm.Catch, bool) {
	for _, c := range cs.entries {
		if pos >= c[0] && pos <= c[1] {
			return c, true
		}
	}

	return vm.Catch{}, false
}

func (cs *CatchStack) Clear() {
	cs.entries = cs.entries[:0]
}

func (cs *CatchStack) Len() int {
	return len(cs.entries)
}

func (cs *CatchStack) All() []vm.Catch {
	return cs.entries
}
