package runtime

import "context"

var arithmeticBenchmarkHostType = NewType(
	"runtime",
	"ArithmeticBenchmarkHost",
	func(value Value) bool {
		_, ok := value.(*arithmeticBenchmarkHost)

		return ok
	},
)

type arithmeticBenchmarkHost struct {
	result      Value
	unsupported bool
}

func (v *arithmeticBenchmarkHost) String() string {
	return "host"
}

func (v *arithmeticBenchmarkHost) Hash() uint64 {
	return 1
}

func (v *arithmeticBenchmarkHost) Copy() Value {
	return v
}

func (v *arithmeticBenchmarkHost) Type() Type {
	return arithmeticBenchmarkHostType
}

func (v *arithmeticBenchmarkHost) Add(_ context.Context, _ Value) (Value, error) {
	if v.unsupported {
		return None, ErrUnsupportedOperands
	}

	return v.result, nil
}

func (v *arithmeticBenchmarkHost) RightAdd(_ context.Context, _ Value) (Value, error) {
	if v.unsupported {
		return None, ErrUnsupportedOperands
	}

	return v.result, nil
}
