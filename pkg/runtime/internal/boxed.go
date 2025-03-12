package internal

import (
	"fmt"
	"hash/fnv"

	"github.com/wI2L/jettison"

	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values/types"
)

// Boxed represents an arbitrary Value that can be boxed as a runtime Value.
type Boxed struct {
	Value any
}

func NewBoxedValue(value any) *Boxed {
	return &Boxed{value}
}

func (b *Boxed) MarshalJSON() ([]byte, error) {
	return jettison.MarshalOpts(b.Value, jettison.NoHTMLEscaping())
}

func (b *Boxed) String() string {
	return fmt.Sprintf("%v", b.Value)
}

func (b *Boxed) Unwrap() any {
	return b.Value
}

func (b *Boxed) Hash() uint64 {
	h := fnv.New64a()

	h.Write([]byte(types.Boxed.String()))
	h.Write([]byte(":"))
	h.Write([]byte(fmt.Sprintf("%v", b.Value)))

	return h.Sum64()
}

func (b *Boxed) Copy() core.Value {
	return NewBoxedValue(b.Value)
}
