package session

import (
	"context"
	"errors"

	"github.com/MontFerret/ferret/v2/pkg/encoding"
	"github.com/MontFerret/ferret/v2/pkg/internal/encodingownership"
	"github.com/MontFerret/ferret/v2/pkg/runtime"
	"github.com/MontFerret/ferret/v2/pkg/vm"
)

// Materialize encodes a VM result and adopts resources discovered by the encoder.
// The caller remains responsible for closing the result, including on failure.
func Materialize(ctx context.Context, registry *encoding.Registry, contentType string, res *vm.Result) (*encoding.Content, error) {
	codec, err := registry.Codec(contentType)
	if err != nil {
		return nil, err
	}

	var encodeErr error
	content, materializeErr := vm.Materialize[*encoding.Content](res, func(value runtime.Value) (vm.Materialized[*encoding.Content], error) {
		adopter := outputValueAdopter{ctx: ctx, result: res}
		encodeCtx := encodingownership.WithValueAdopter(ctx, adopter.Adopt)
		enc := codec.EncodeWith().PreHook(func(_ context.Context, value runtime.Value) error {
			res.AdoptValue(value)

			return nil
		}).Encoder()

		data, err := encoding.EncodeBytes(encodeCtx, enc, value)
		encodeErr = err
		if err != nil && len(data) == 0 {
			return vm.Materialized[*encoding.Content]{}, nil
		}

		// The collecting buffer owns its bytes. Keep the failure separate because vm.Materialize
		// intentionally discards a materializer's value when it returns an error.
		metadata := encoding.Metadata{ContentType: codec.ContentType()}
		if err == nil {
			metadata.Length = int64(len(data))
			metadata.LengthKnown = true
		}

		return vm.Materialized[*encoding.Content]{
			Value: &encoding.Content{
				Metadata: metadata,
				Data:     data,
			},
		}, nil
	})

	if materializeErr != nil {
		return content, errors.Join(materializeErr, encodeErr)
	}

	return content, encodeErr
}
