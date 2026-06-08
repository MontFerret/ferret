package vm

import "github.com/MontFerret/ferret/v2/pkg/runtime"

func assertNumericOperand(value runtime.Value) error {
	return runtime.AssertNumber(value)
}

func assertNumericOperands(left, right runtime.Value) error {
	if err := runtime.AssertNumber(left); err != nil {
		return err
	}

	return runtime.AssertNumber(right)
}
