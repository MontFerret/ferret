package uapi

import (
	"context"

	apidebugger "github.com/MontFerret/api/debugger"
	apisource "github.com/MontFerret/api/source"
	"github.com/MontFerret/ferret/v2/pkg/debugger"
)

type debugSession struct {
	native *debugger.Session
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

func (s *debugSession) Pause(ctx context.Context) error {
	return wrapDiagnosticError(s.native.Pause(ctx))
}

func (s *debugSession) ReplaceBreakpoints(ctx context.Context, sourceName string, requests []apidebugger.BreakpointRequest) ([]apidebugger.Breakpoint, error) {
	values, err := s.native.ReplaceBreakpoints(ctx, sourceName, requests)

	return values, wrapDiagnosticError(err)
}

func (s *debugSession) SetBreakpoint(ctx context.Context, location apisource.Location) (apidebugger.Breakpoint, error) {
	value, err := s.native.SetBreakpoint(ctx, location)

	return value, wrapDiagnosticError(err)
}

func (s *debugSession) SetBreakpointAt(ctx context.Context, location apisource.Location, options apidebugger.BreakpointOptions) (apidebugger.Breakpoint, error) {
	value, err := s.native.SetBreakpointAt(ctx, location, options)

	return value, wrapDiagnosticError(err)
}

func (s *debugSession) DeleteBreakpoint(ctx context.Context, id apidebugger.BreakpointID) error {
	return wrapDiagnosticError(s.native.DeleteBreakpoint(ctx, id))
}

func (s *debugSession) Breakpoints(ctx context.Context) ([]apidebugger.Breakpoint, error) {
	values, err := s.native.Breakpoints(ctx)

	return values, wrapDiagnosticError(err)
}

func (s *debugSession) Frames(ctx context.Context) ([]apidebugger.Frame, error) {
	values, err := s.native.Frames(ctx)

	return values, wrapDiagnosticError(err)
}

func (s *debugSession) Locals(ctx context.Context) ([]apidebugger.Variable, error) {
	values, err := s.native.Locals(ctx)

	return values, wrapDiagnosticError(err)
}

func (s *debugSession) FrameLocals(ctx context.Context, frame int) ([]apidebugger.Variable, error) {
	values, err := s.native.FrameLocals(ctx, frame)

	return values, wrapDiagnosticError(err)
}

func (s *debugSession) Variables(ctx context.Context, reference apidebugger.ValueReference) ([]apidebugger.Variable, error) {
	values, err := s.native.Variables(ctx, reference)

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
	return wrapDiagnosticError(s.native.Close())
}
