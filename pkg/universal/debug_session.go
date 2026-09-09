package universal

import (
	"context"
	"sync"

	apidebugger "github.com/MontFerret/api/debugger"
	apisource "github.com/MontFerret/api/source"
	"github.com/MontFerret/ferret/v2/pkg/debugger"
)

type debugSession struct {
	closeErr  error
	native    *debugger.Session
	closeOnce sync.Once
}

var _ apidebugger.Session = (*debugSession)(nil)

func (s *debugSession) Start(ctx context.Context) (*apidebugger.Event, error) {
	event, err := s.native.Start(ctx)

	return convertEvent(event), wrapDiagnosticError(err)
}

func (s *debugSession) Continue(ctx context.Context) (*apidebugger.Event, error) {
	event, err := s.native.Continue(ctx)

	return convertEvent(event), wrapDiagnosticError(err)
}

func (s *debugSession) StepIn(ctx context.Context) (*apidebugger.Event, error) {
	event, err := s.native.StepIn(ctx)

	return convertEvent(event), wrapDiagnosticError(err)
}

func (s *debugSession) StepOver(ctx context.Context) (*apidebugger.Event, error) {
	event, err := s.native.StepOver(ctx)

	return convertEvent(event), wrapDiagnosticError(err)
}

func (s *debugSession) StepOut(ctx context.Context) (*apidebugger.Event, error) {
	event, err := s.native.StepOut(ctx)

	return convertEvent(event), wrapDiagnosticError(err)
}

func (s *debugSession) Pause() error {
	return wrapDiagnosticError(s.native.Pause())
}

func (s *debugSession) SetBreakpoint(location apisource.Location) (apidebugger.Breakpoint, error) {
	value, err := s.native.SetBreakpoint(location)

	return value, wrapDiagnosticError(err)
}

func (s *debugSession) SetBreakpointAt(location apisource.Location, options apidebugger.BreakpointOptions) (apidebugger.Breakpoint, error) {
	value, err := s.native.SetBreakpointAt(location, options)

	return value, wrapDiagnosticError(err)
}

func (s *debugSession) DeleteBreakpoint(id apidebugger.BreakpointID) error {
	return wrapDiagnosticError(s.native.DeleteBreakpoint(id))
}

func (s *debugSession) Breakpoints() []apidebugger.Breakpoint { return s.native.Breakpoints() }

func (s *debugSession) Frames() ([]apidebugger.Frame, error) {
	values, err := s.native.Frames()

	return values, wrapDiagnosticError(err)
}

func (s *debugSession) Locals() ([]apidebugger.Variable, error) {
	values, err := s.native.Locals()

	return values, wrapDiagnosticError(err)
}

func (s *debugSession) FrameLocals(frame int) ([]apidebugger.Variable, error) {
	values, err := s.native.FrameLocals(frame)

	return values, wrapDiagnosticError(err)
}

func (s *debugSession) Variables(reference apidebugger.ValueReference) ([]apidebugger.Variable, error) {
	values, err := s.native.Variables(reference)

	return values, wrapDiagnosticError(err)
}

func (s *debugSession) Evaluate(ctx context.Context, expression string) (apidebugger.Value, error) {
	value, err := s.native.Evaluate(ctx, expression)

	return value, wrapDiagnosticError(err)
}

func (s *debugSession) EvaluateFrame(
	ctx context.Context,
	frame int,
	expression string,
) (apidebugger.Value, error) {
	value, err := s.native.EvaluateFrame(ctx, frame, expression)

	return value, wrapDiagnosticError(err)
}

func (s *debugSession) Close() error {
	s.closeOnce.Do(func() { s.closeErr = wrapDiagnosticError(s.native.Close()) })

	return s.closeErr
}
