package mem

import (
	"regexp"

	"github.com/MontFerret/ferret/pkg/runtime"
)

type (
	Cache struct {
		FuncHash  uint64
		Functions map[int]*CachedFunction
		Regexps   map[int]*regexp.Regexp
	}

	CachedFunction struct {
		Fn0 runtime.Function0
		Fn1 runtime.Function1
		Fn2 runtime.Function2
		Fn3 runtime.Function3
		Fn4 runtime.Function4
		FnV runtime.Function
	}
)

func NewCache() *Cache {
	return &Cache{
		Functions: make(map[int]*CachedFunction),
		Regexps:   make(map[int]*regexp.Regexp),
	}
}
