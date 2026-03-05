package vm

import (
	"context"

	"github.com/MontFerret/ferret/v2/pkg/bytecode"
	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

const (
	queryDescriptorKind    = "kind"
	queryDescriptorPayload = "payload"
	queryDescriptorOptions = "options"

	errQueryFormatUnexpected = "unexpected query format"
)

var (
	queryDescriptorKeyKind    = runtime.NewString(queryDescriptorKind)
	queryDescriptorKeyPayload = runtime.NewString(queryDescriptorPayload)
	queryDescriptorKeyOptions = runtime.NewString(queryDescriptorOptions)
)

func (vm *VM) applyQuery(ctx context.Context, reg []runtime.Value, src1 bytecode.Operand, constants []runtime.Value, src2 bytecode.Operand, dst bytecode.Operand) error {
	var src runtime.Value
	if src1.IsConstant() {
		src = constants[src1.Constant()]
	} else {
		src = reg[src1]
	}

	var arg runtime.Value

	if src2.IsConstant() {
		arg = constants[src2.Constant()]
	} else {
		arg = reg[src2]
	}

	var query runtime.Query

	switch value := arg.(type) {
	case runtime.ObjectLike:
		kind, err := value.Get(ctx, queryDescriptorKeyKind)
		if err != nil {
			if err := vm.setOrTryCatch(dst, runtime.None, err); err != nil {
				return err
			}

			break
		}

		payload, err := value.Get(ctx, queryDescriptorKeyPayload)
		if err != nil {
			if err := vm.setOrTryCatch(dst, runtime.None, err); err != nil {
				return err
			}

			break
		}

		options, err := value.Get(ctx, queryDescriptorKeyOptions)
		if err != nil {
			if err := vm.setOrTryCatch(dst, runtime.None, err); err != nil {
				return err
			}

			break
		}

		query = runtime.Query{
			Kind:    runtime.CastOr[runtime.String](kind, runtime.EmptyString),
			Payload: runtime.CastOr[runtime.String](payload, runtime.EmptyString),
			Options: options,
		}
	case *runtime.Array:
		length, err := value.Length(ctx)
		if err != nil {
			if err := vm.setOrTryCatch(dst, runtime.None, err); err != nil {
				return err
			}

			break
		}

		if length != 3 {
			if err := vm.setOrTryCatch(dst, runtime.None, runtime.Error(runtime.ErrInvalidOperation, errQueryFormatUnexpected)); err != nil {
				return err
			}

			break
		}

		kindVal, err := value.At(ctx, runtime.NewInt(0))
		if err != nil {
			if err := vm.setOrTryCatch(dst, runtime.None, err); err != nil {
				return err
			}

			break
		}

		payloadVal, err := value.At(ctx, runtime.NewInt(1))
		if err != nil {
			if err := vm.setOrTryCatch(dst, runtime.None, err); err != nil {
				return err
			}

			break
		}

		var optionsVal runtime.Value = runtime.None
		if length > 2 {
			optionsVal, err = value.At(ctx, runtime.NewInt(2))
			if err != nil {
				if err := vm.setOrTryCatch(dst, runtime.None, err); err != nil {
					return err
				}

				break
			}
		}

		kind, err := runtime.CastString(kindVal)
		if err != nil {
			if err := vm.setOrTryCatch(dst, runtime.None, runtime.TypeErrorOf(kindVal, runtime.TypeString)); err != nil {
				return err
			}

			break
		}

		payload := runtime.EmptyString
		if payloadVal != runtime.None {
			payload, err = runtime.CastString(payloadVal)
			if err != nil {
				if err := vm.setOrTryCatch(dst, runtime.None, runtime.TypeErrorOf(payloadVal, runtime.TypeString, runtime.TypeNone)); err != nil {
					return err
				}

				break
			}
		}

		query = runtime.Query{
			Kind:    kind,
			Payload: payload,
			Options: optionsVal,
		}
	default:
		// TODO: Give a more specific error message here
		if err := vm.setOrTryCatch(dst, runtime.None, runtime.Error(runtime.ErrInvalidOperation, errQueryFormatUnexpected)); err != nil {
			return err
		}

		return nil
	}

	if queryable, ok := src.(runtime.Queryable); ok {
		res, err := queryable.Query(ctx, query)
		if err == nil && res == nil {
			res = runtime.NewArray(0)
		}

		return vm.setOrTryCatch(dst, res, err)
	}

	if list, ok := src.(runtime.List); ok {
		out := runtime.NewArray(0)

		err := runtime.ForEach(ctx, list, func(ctx context.Context, value, _ runtime.Value) (runtime.Boolean, error) {
			queryable, ok := value.(runtime.Queryable)
			if !ok {
				return runtime.False, runtime.TypeErrorOf(value, runtime.TypeQueryable)
			}

			res, err := queryable.Query(ctx, query)
			if err != nil {
				return runtime.False, err
			}
			if res == nil {
				return runtime.True, nil
			}

			if err := runtime.ForEach(ctx, res, func(ctx context.Context, item, _ runtime.Value) (runtime.Boolean, error) {
				if err := out.Append(ctx, item); err != nil {
					return runtime.False, err
				}

				return runtime.True, nil
			}); err != nil {
				return runtime.False, err
			}

			return runtime.True, nil
		})

		return vm.setOrTryCatch(dst, out, err)
	}

	return vm.setOrTryCatch(dst, runtime.None, runtime.TypeErrorOf(src, runtime.TypeQueryable, runtime.TypeList))
}
