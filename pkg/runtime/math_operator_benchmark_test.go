package runtime

import (
	"context"
	"testing"
	"time"
)

var (
	arithmeticBenchmarkSink      Value
	arithmeticBenchmarkErrorSink error
)

func BenchmarkArithmeticAddInt(b *testing.B) {
	benchmarkBinaryArithmetic(b, NewInt(40), NewInt(2), Add)
}

func BenchmarkArithmeticAddString(b *testing.B) {
	benchmarkBinaryArithmetic(b, NewString("value="), NewInt(42), Add)
}

func BenchmarkArithmeticSubtractMixed(b *testing.B) {
	benchmarkBinaryArithmetic(b, NewInt(40), NewFloat(2.5), Subtract)
}

func BenchmarkArithmeticMultiplyFloat(b *testing.B) {
	benchmarkBinaryArithmetic(b, NewFloat(4.25), NewFloat(2.5), Multiply)
}

func BenchmarkArithmeticDivideMixed(b *testing.B) {
	benchmarkBinaryArithmetic(b, NewInt(40), NewFloat(2.5), Divide)
}

func BenchmarkArithmeticModuloFloat(b *testing.B) {
	benchmarkBinaryArithmetic(b, NewFloat(40.5), NewInt(3), Mod)
}

func BenchmarkArithmeticAddDuration(b *testing.B) {
	benchmarkBinaryArithmetic(b, NewDuration(time.Second), NewDuration(2*time.Millisecond), Add)
}

func BenchmarkArithmeticMultiplyDuration(b *testing.B) {
	benchmarkBinaryArithmetic(b, NewDuration(time.Second), NewFloat(2.5), Multiply)
}

func BenchmarkArithmeticDivideDuration(b *testing.B) {
	benchmarkBinaryArithmetic(b, NewDuration(time.Second), NewInt(2), Divide)
}

func BenchmarkArithmeticAddDateTime(b *testing.B) {
	benchmarkBinaryArithmetic(
		b,
		NewDateTime(time.Date(2026, time.August, 5, 12, 0, 0, 0, time.UTC)),
		NewDuration(time.Second),
		Add,
	)
}

func BenchmarkArithmeticAddHostBuiltin(b *testing.B) {
	host := &arithmeticBenchmarkHost{result: NewInt(42)}
	benchmarkBinaryArithmetic(b, host, NewInt(2), Add)
}

func BenchmarkArithmeticAddBuiltinHost(b *testing.B) {
	host := &arithmeticBenchmarkHost{result: NewInt(42)}
	benchmarkBinaryArithmetic(b, NewInt(40), host, Add)
}

func BenchmarkArithmeticAddHostHost(b *testing.B) {
	left := &arithmeticBenchmarkHost{result: NewInt(42)}
	right := &arithmeticBenchmarkHost{result: NewInt(42)}
	benchmarkBinaryArithmetic(b, left, right, Add)
}

func BenchmarkArithmeticAddUnsupportedHost(b *testing.B) {
	left := &arithmeticBenchmarkHost{unsupported: true}
	right := &arithmeticBenchmarkHost{unsupported: true}
	ctx := context.Background()

	if _, err := Add(ctx, left, right); err == nil {
		b.Fatal("expected unsupported host addition to fail")
	}

	b.ReportAllocs()
	b.ResetTimer()

	var err error
	for b.Loop() {
		_, err = Add(ctx, left, right)
	}

	arithmeticBenchmarkErrorSink = err
}

func benchmarkBinaryArithmetic(
	b *testing.B,
	left Value,
	right Value,
	operation func(context.Context, Value, Value) (Value, error),
) {
	b.Helper()
	ctx := context.Background()

	if _, err := operation(ctx, left, right); err != nil {
		b.Fatal(err)
	}

	b.ReportAllocs()
	b.ResetTimer()

	var result Value
	for b.Loop() {
		result, _ = operation(ctx, left, right)
	}

	arithmeticBenchmarkSink = result
}
