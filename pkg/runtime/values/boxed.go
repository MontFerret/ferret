package values

import (
	"fmt"
	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values/types"
	"github.com/wI2L/jettison"
	"hash/fnv"
)

// Boxed represents an arbitrary value that can be boxed as a runtime Value.
type Boxed struct {
	value any
}

func NewBoxedValue(value any) *Boxed {
	return &Boxed{value}
}

func (b *Boxed) MarshalJSON() ([]byte, error) {
	return jettison.MarshalOpts(b.value, jettison.NoHTMLEscaping())
}

func (b *Boxed) String() string {
	return fmt.Sprintf("%v", b.value)
}

func (b *Boxed) Unwrap() interface{} {
	return b.value
}

func (b *Boxed) Hash() uint64 {
	h := fnv.New64a()

	h.Write([]byte(types.Boxed.String()))
	h.Write([]byte(":"))
	h.Write([]byte(fmt.Sprintf("%v", b.value)))

	return h.Sum64()
}

func (b *Boxed) Copy() core.Value {
	return NewBoxedValue(b.value)
}
