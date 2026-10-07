package mock

import (
	"context"

	ferretencoding "github.com/MontFerret/ferret/v2/pkg/encoding"
	encodingjson "github.com/MontFerret/ferret/v2/pkg/encoding/json"
	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

type Text struct {
	value runtime.String
}

func NewText(value string) *Text {
	return &Text{value: runtime.NewString(value)}
}

func (t *Text) MarshalJSON() ([]byte, error) {
	return ferretencoding.EncodeBytes(context.Background(), encodingjson.Default, t.value)
}

func (t *Text) String() string {
	return t.value.String()
}

func (t *Text) Hash() uint64 {
	return t.value.Hash()
}

func (t *Text) Copy() runtime.Value {
	return t
}

func (t *Text) Iterate(ctx context.Context) (runtime.Iterator, error) {
	return runtime.NewArrayWith(t.value).Iterate(ctx)
}
