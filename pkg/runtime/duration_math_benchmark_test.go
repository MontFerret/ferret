package runtime

import (
	"context"
	"testing"
	"time"
)

var durationScaleBenchmarkSink Value

func BenchmarkScaleDurationNumericStringMultiply(b *testing.B) {
	ctx := context.Background()
	duration := NewDuration(5 * time.Second)
	scalar := NewString("2.5")

	if _, err := scaleDuration(ctx, duration, scalar, false); err != nil {
		b.Fatal(err)
	}

	b.ReportAllocs()
	b.ResetTimer()

	var result Value
	for b.Loop() {
		result, _ = scaleDuration(ctx, duration, scalar, false)
	}

	durationScaleBenchmarkSink = result
}

func BenchmarkScaleDurationNumericStringDivide(b *testing.B) {
	ctx := context.Background()
	duration := NewDuration(5 * time.Second)
	scalar := NewString("2.5")

	if _, err := scaleDuration(ctx, duration, scalar, true); err != nil {
		b.Fatal(err)
	}

	b.ReportAllocs()
	b.ResetTimer()

	var result Value
	for b.Loop() {
		result, _ = scaleDuration(ctx, duration, scalar, true)
	}

	durationScaleBenchmarkSink = result
}
