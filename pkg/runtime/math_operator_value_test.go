package runtime_test

import (
	"context"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

var arithmeticCapabilityType = runtime.NewType(
	"runtime_test",
	"ArithmeticCapabilityValue",
	func(value runtime.Value) bool {
		_, ok := value.(*arithmeticCapabilityValue)

		return ok
	},
)

type (
	arithmeticCapabilityResponse struct {
		value runtime.Value
		err   error
	}

	arithmeticCapabilityCall struct {
		ctx     context.Context
		operand runtime.Value
		method  string
	}

	arithmeticCapabilityValue struct {
		responses map[string]arithmeticCapabilityResponse
		name      string
		calls     []arithmeticCapabilityCall
	}
)

var (
	_ runtime.Value          = (*arithmeticCapabilityValue)(nil)
	_ runtime.Additive       = (*arithmeticCapabilityValue)(nil)
	_ runtime.Multiplicative = (*arithmeticCapabilityValue)(nil)
	_ runtime.Arithmetic     = (*arithmeticCapabilityValue)(nil)
)

func newArithmeticCapabilityValue(name string) *arithmeticCapabilityValue {
	return &arithmeticCapabilityValue{
		responses: make(map[string]arithmeticCapabilityResponse),
		name:      name,
	}
}

func (v *arithmeticCapabilityValue) String() string {
	return v.name
}

func (v *arithmeticCapabilityValue) Hash() uint64 {
	return 1
}

func (v *arithmeticCapabilityValue) Copy() runtime.Value {
	return v
}

func (v *arithmeticCapabilityValue) Type() runtime.Type {
	return arithmeticCapabilityType
}

func (v *arithmeticCapabilityValue) Add(ctx context.Context, right runtime.Value) (runtime.Value, error) {
	return v.invoke(ctx, "Add", right)
}

func (v *arithmeticCapabilityValue) RightAdd(ctx context.Context, left runtime.Value) (runtime.Value, error) {
	return v.invoke(ctx, "RightAdd", left)
}

func (v *arithmeticCapabilityValue) Subtract(ctx context.Context, right runtime.Value) (runtime.Value, error) {
	return v.invoke(ctx, "Subtract", right)
}

func (v *arithmeticCapabilityValue) RightSubtract(ctx context.Context, left runtime.Value) (runtime.Value, error) {
	return v.invoke(ctx, "RightSubtract", left)
}

func (v *arithmeticCapabilityValue) Multiply(ctx context.Context, right runtime.Value) (runtime.Value, error) {
	return v.invoke(ctx, "Multiply", right)
}

func (v *arithmeticCapabilityValue) RightMultiply(ctx context.Context, left runtime.Value) (runtime.Value, error) {
	return v.invoke(ctx, "RightMultiply", left)
}

func (v *arithmeticCapabilityValue) Divide(ctx context.Context, right runtime.Value) (runtime.Value, error) {
	return v.invoke(ctx, "Divide", right)
}

func (v *arithmeticCapabilityValue) RightDivide(ctx context.Context, left runtime.Value) (runtime.Value, error) {
	return v.invoke(ctx, "RightDivide", left)
}

func (v *arithmeticCapabilityValue) Mod(ctx context.Context, right runtime.Value) (runtime.Value, error) {
	return v.invoke(ctx, "Mod", right)
}

func (v *arithmeticCapabilityValue) RightMod(ctx context.Context, left runtime.Value) (runtime.Value, error) {
	return v.invoke(ctx, "RightMod", left)
}

func (v *arithmeticCapabilityValue) invoke(
	ctx context.Context,
	method string,
	operand runtime.Value,
) (runtime.Value, error) {
	v.calls = append(v.calls, arithmeticCapabilityCall{
		ctx:     ctx,
		method:  method,
		operand: operand,
	})

	if response, ok := v.responses[method]; ok {
		return response.value, response.err
	}

	return runtime.NewString(method), nil
}
