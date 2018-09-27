package common

import (
	"github.com/MontFerret/ferret/pkg/runtime/core"
	"sync"
)

type SyncValue struct {
	sync.Mutex
	value core.Value
}

func NewSyncValue(init core.Value) *SyncValue {
	val := new(SyncValue)
	val.value = init

	return val
}

func (sv *SyncValue) Get() core.Value {
	sv.Lock()
	defer sv.Unlock()

	return sv.value
}

func (sv *SyncValue) Set(val core.Value) {
	sv.Lock()
	defer sv.Unlock()

	sv.value = val
}
