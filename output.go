package ferret

import (
	"github.com/MontFerret/ferret/v2/pkg/encoding"
	"github.com/MontFerret/ferret/v2/pkg/runtime"
	"github.com/MontFerret/ferret/v2/pkg/vm"
)

type Output struct {
	ContentType string
	Content     []byte
}

func newOutput(registry *encoding.Registry, contentType string, res *vm.Result) (*Output, error) {
	codec, err := registry.Codec(contentType)
	if err != nil {
		return nil, err
	}

	return vm.Materialize[*Output](res, func(value runtime.Value) (vm.Materialized[*Output], error) {
		enc := codec.EncodeWith().PreHook(func(value runtime.Value) error {
			res.AdoptValue(value)
			return nil
		}).Encoder()

		data, err := enc.Encode(value)

		if err != nil {
			return vm.Materialized[*Output]{}, err
		}

		return vm.Materialized[*Output]{
			Value: &Output{
				ContentType: codec.ContentType(),
				Content:     data,
			},
		}, nil

	})
}
