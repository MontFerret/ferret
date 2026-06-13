package vm

import (
	"context"
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/bytecode"
	"github.com/MontFerret/ferret/v2/pkg/internal/debugpoint"
)

func BenchmarkDebugSourcePointDispatch(b *testing.B) {
	point := bytecode.DebugPoint{ID: 0, PC: 7, FunctionID: -1}
	execution := &debugExecution{points: debugpoint.New([]bytecode.DebugPoint{point})}
	control := debugControl{
		owner: execution,
		mode:  DebugResumeContinue,
	}
	state := sourcePointState{pc: point.PC, pointID: point.ID}
	ctx := context.Background()

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		action, err := control.onSourcePoint(ctx, state)
		if err != nil {
			b.Fatal(err)
		}
		if action != sourcePointContinue {
			b.Fatalf("unexpected source point action: %d", action)
		}
	}
}
