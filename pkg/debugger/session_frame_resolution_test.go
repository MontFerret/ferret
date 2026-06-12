package debugger

import (
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/bytecode"
	"github.com/MontFerret/ferret/v2/pkg/source"
	"github.com/MontFerret/ferret/v2/pkg/vm"
)

func TestSessionFramesDoNotResolveLocationsAcrossFunctions(t *testing.T) {
	src := source.New("frames.fql", "RETURN 1\nRETURN 2\nRETURN 3")
	points := []bytecode.DebugPoint{
		{PC: 2, Span: source.Span{Start: 0, End: 8}, FunctionID: 0},
		{PC: 4, Span: source.Span{Start: 9, End: 17}, FunctionID: 1},
		{PC: 6, Span: source.Span{Start: 18, End: 26}, FunctionID: 0},
	}
	session, err := NewSession(Config{
		Execution: &fakeExecution{
			frames: []vm.DebugFrame{{Name: "first", FunctionID: 0, PC: 5}},
			status: vm.DebugExecutionPaused,
		},
		Values:      &fakeValueAccess{inner: vm.NewDebugValueAccess()},
		Services:    &fakeSessionServices{},
		Source:      src,
		DebugPoints: points,
	})
	if err != nil {
		t.Fatal(err)
	}
	defer session.Close()

	frames, err := session.Frames()
	if err != nil {
		t.Fatal(err)
	}
	if len(frames) != 1 || frames[0].Location.Line != 1 {
		t.Fatalf("resolved frame through another function's point: %#v", frames)
	}
}
