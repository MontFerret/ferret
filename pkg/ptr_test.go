package pkg

import (
	"testing"

	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
)

type Pair struct {
	Key   any
	Value any
}

type Item interface {
	GetKey() any
	GetValue() any
}

func (p *Pair) GetKey() any {
	return p.Key
}

func (p *Pair) GetValue() any {
	return p.Value
}

func BenchmarkPtr(b *testing.B) {
	stack1 := make([]*Pair, b.N)
	stack2 := make([]*Pair, b.N)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		stack1[i] = &Pair{Key: i, Value: i}

		if i > 0 {
			stack2[i] = stack1[i-1]
		}
	}
}

func BenchmarkVal(b *testing.B) {
	stack1 := make([]Pair, b.N)
	stack2 := make([]Pair, b.N)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		stack1[i] = Pair{Key: i, Value: i}

		if i > 0 {
			stack2[i] = stack1[i-1]
		}
	}
}

func BenchmarkInt(b *testing.B) {
	stack1 := make([]Item, b.N)
	stack2 := make([]Item, b.N)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		stack1[i] = &Pair{Key: i, Value: i}

		if i > 0 {
			stack2[i] = stack1[i-1]
		}
	}
}

func BenchmarkAny(b *testing.B) {
	stack1 := make([]any, b.N)
	stack2 := make([]any, b.N)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		stack1[i] = Pair{Key: i, Value: i}

		if i > 0 {
			stack2[i] = stack1[i-1]
		}
	}
}

func BenchmarkValue(b *testing.B) {
	stack1 := make([]core.Value, b.N)
	stack2 := make([]core.Value, b.N)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		stack1[i] = &values.Tuple{
			First:  values.NewInt(i),
			Second: values.NewInt(i),
		}

		if i > 0 {
			stack2[i] = stack1[i-1]
		}
	}
}

func BenchmarkBox(b *testing.B) {
	stack1 := make([]core.Value, b.N)
	stack2 := make([]core.Value, b.N)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		stack1[i] = &values.Boxed{
			Value: Pair{
				Key:   values.Int(i),
				Value: values.Int(i),
			},
		}

		if i > 0 {
			stack2[i] = stack1[i-1]
		}
	}
}
