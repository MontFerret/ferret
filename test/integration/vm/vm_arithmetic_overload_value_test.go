package vm_test

import (
	"context"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

var vmArithmeticOverloadType = runtime.NewType(
	"vm_test",
	"ArithmeticOverload",
	func(value runtime.Value) bool {
		_, ok := value.(*vmArithmeticOverloadValue)

		return ok
	},
)

type (
	vmArithmeticOverloadResponse struct {
		value runtime.Value
		err   error
	}

	vmArithmeticOverloadValue struct {
		responses map[string]vmArithmeticOverloadResponse
		name      string
	}
)

func newVMArithmeticOverloadValue(name string) *vmArithmeticOverloadValue {
	return &vmArithmeticOverloadValue{
		responses: make(map[string]vmArithmeticOverloadResponse),
		name:      name,
	}
}

func (v *vmArithmeticOverloadValue) String() string {
	return v.name
}

func (v *vmArithmeticOverloadValue) Hash() uint64 {
	return 1
}

func (v *vmArithmeticOverloadValue) Copy() runtime.Value {
	return v
}

func (v *vmArithmeticOverloadValue) Type() runtime.Type {
	return vmArithmeticOverloadType
}

func (v *vmArithmeticOverloadValue) Add(ctx context.Context, right runtime.Value) (runtime.Value, error) {
	return v.invoke(ctx, "Add", right)
}

func (v *vmArithmeticOverloadValue) RightAdd(ctx context.Context, left runtime.Value) (runtime.Value, error) {
	return v.invoke(ctx, "RightAdd", left)
}

func (v *vmArithmeticOverloadValue) Subtract(ctx context.Context, right runtime.Value) (runtime.Value, error) {
	return v.invoke(ctx, "Subtract", right)
}

func (v *vmArithmeticOverloadValue) RightSubtract(ctx context.Context, left runtime.Value) (runtime.Value, error) {
	return v.invoke(ctx, "RightSubtract", left)
}

func (v *vmArithmeticOverloadValue) Multiply(ctx context.Context, right runtime.Value) (runtime.Value, error) {
	return v.invoke(ctx, "Multiply", right)
}

func (v *vmArithmeticOverloadValue) RightMultiply(ctx context.Context, left runtime.Value) (runtime.Value, error) {
	return v.invoke(ctx, "RightMultiply", left)
}

func (v *vmArithmeticOverloadValue) Divide(ctx context.Context, right runtime.Value) (runtime.Value, error) {
	return v.invoke(ctx, "Divide", right)
}

func (v *vmArithmeticOverloadValue) RightDivide(ctx context.Context, left runtime.Value) (runtime.Value, error) {
	return v.invoke(ctx, "RightDivide", left)
}

func (v *vmArithmeticOverloadValue) Mod(ctx context.Context, right runtime.Value) (runtime.Value, error) {
	return v.invoke(ctx, "Mod", right)
}

func (v *vmArithmeticOverloadValue) RightMod(ctx context.Context, left runtime.Value) (runtime.Value, error) {
	return v.invoke(ctx, "RightMod", left)
}

func (v *vmArithmeticOverloadValue) invoke(
	ctx context.Context,
	method string,
	_ runtime.Value,
) (runtime.Value, error) {
	if err := ctx.Err(); err != nil {
		return runtime.None, err
	}

	if response, ok := v.responses[method]; ok {
		return response.value, response.err
	}

	return runtime.NewString(method), nil
}
