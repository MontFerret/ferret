package core

import (
	"github.com/MontFerret/ferret/v2/pkg/bytecode"
)

type (
	CatchJumpMode int

	CatchStack struct {
		entries []bytecode.Catch
	}
)

const (
	CatchJumpModeNone CatchJumpMode = iota
	CatchJumpModeEnd
)

func NewCatchStack() *CatchStack {
	return &CatchStack{
		entries: make([]bytecode.Catch, 0),
	}
}

func (cs *CatchStack) Push(start, end, jump int) {
	cs.entries = append(cs.entries, bytecode.Catch{start, end, jump})
}

func (cs *CatchStack) Pop() {
	if len(cs.entries) > 0 {
		cs.entries = cs.entries[:len(cs.entries)-1]
	}
}

func (cs *CatchStack) Find(pos int) (bytecode.Catch, bool) {
	for _, c := range cs.entries {
		if pos >= c[0] && pos <= c[1] {
			return c, true
		}
	}

	return bytecode.Catch{}, false
}

func (cs *CatchStack) Clear() {
	cs.entries = cs.entries[:0]
}

func (cs *CatchStack) Len() int {
	return len(cs.entries)
}

func (cs *CatchStack) All() []bytecode.Catch {
	return cs.entries
}
