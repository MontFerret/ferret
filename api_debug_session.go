package ferret

import (
	"context"

	apidebugger "github.com/MontFerret/api/debugger"
	"github.com/MontFerret/api/result"
	apisource "github.com/MontFerret/api/source"
	coredebugger "github.com/MontFerret/ferret/v2/pkg/debugger"
)

type apiDebugSession struct {
	session *DebugSession
}

var _ apidebugger.Session = (*apiDebugSession)(nil)

func newAPIDebugSession(session *DebugSession) *apiDebugSession {
	return &apiDebugSession{session: session}
}

func (s *apiDebugSession) Start(ctx context.Context) (*apidebugger.Event, error) {
	event, err := s.session.Start(ctx)

	return s.event(event), err
}

func (s *apiDebugSession) Continue(ctx context.Context) (*apidebugger.Event, error) {
	event, err := s.session.Continue(ctx)

	return s.event(event), err
}

func (s *apiDebugSession) Step(ctx context.Context) (*apidebugger.Event, error) {
	event, err := s.session.Step(ctx)

	return s.event(event), err
}

func (s *apiDebugSession) Next(ctx context.Context) (*apidebugger.Event, error) {
	event, err := s.session.Next(ctx)

	return s.event(event), err
}

func (s *apiDebugSession) Out(ctx context.Context) (*apidebugger.Event, error) {
	event, err := s.session.Out(ctx)

	return s.event(event), err
}

func (s *apiDebugSession) Pause() error {
	return s.session.Pause()
}

func (s *apiDebugSession) SetBreakpoint(location apisource.Location) (apidebugger.Breakpoint, error) {
	return s.SetBreakpointAt(location, apidebugger.BreakpointOptions{
		BindingMode: apidebugger.BreakpointBindNextExecutableInFile,
	})
}

func (s *apiDebugSession) SetBreakpointAt(
	location apisource.Location,
	opts apidebugger.BreakpointOptions,
) (apidebugger.Breakpoint, error) {
	breakpoint, err := s.session.SetBreakpointAt(
		coredebugger.SourceLocation{
			File:   location.File,
			Line:   location.Line,
			Column: location.Column,
		},
		coredebugger.BreakpointOptions{
			BindingMode: coredebugger.BreakpointBindingMode(opts.BindingMode),
		},
	)

	return s.breakpoint(breakpoint), err
}

func (s *apiDebugSession) DeleteBreakpoint(id apidebugger.BreakpointID) error {
	return s.session.DeleteBreakpoint(coredebugger.BreakpointID(id))
}

func (s *apiDebugSession) Breakpoints() []apidebugger.Breakpoint {
	breakpoints := s.session.Breakpoints()
	converted := make([]apidebugger.Breakpoint, len(breakpoints))

	for i := range breakpoints {
		converted[i] = s.breakpoint(breakpoints[i])
	}

	return converted
}

func (s *apiDebugSession) Frames() ([]apidebugger.Frame, error) {
	frames, err := s.session.Frames()
	if err != nil {
		return nil, err
	}

	converted := make([]apidebugger.Frame, len(frames))

	for i := range frames {
		converted[i] = apidebugger.Frame{
			Name:       frames[i].Name,
			Location:   s.location(frames[i].Location).Location,
			FunctionID: apidebugger.FunctionID(frames[i].FunctionID),
		}
	}

	return converted, nil
}

func (s *apiDebugSession) Locals() ([]apidebugger.Variable, error) {
	variables, err := s.session.Locals()
	if err != nil {
		return nil, err
	}

	return s.variables(variables), nil
}

func (s *apiDebugSession) FrameLocals(frame int) ([]apidebugger.Variable, error) {
	variables, err := s.session.FrameLocals(frame)
	if err != nil {
		return nil, err
	}

	return s.variables(variables), nil
}

func (s *apiDebugSession) Variables(
	reference apidebugger.ValueReference,
) ([]apidebugger.Variable, error) {
	variables, err := s.session.Variables(coredebugger.ValueReference(reference))
	if err != nil {
		return nil, err
	}

	return s.variables(variables), nil
}

func (s *apiDebugSession) Evaluate(ctx context.Context, expression string) (apidebugger.Value, error) {
	value, err := s.session.Evaluate(ctx, expression)

	return s.value(value), err
}

func (s *apiDebugSession) EvaluateFrame(
	ctx context.Context,
	frame int,
	expression string,
) (apidebugger.Value, error) {
	value, err := s.session.EvaluateFrame(ctx, frame, expression)

	return s.value(value), err
}

func (s *apiDebugSession) Close() error {
	return s.session.Close()
}

func (s *apiDebugSession) event(event *coredebugger.Event) *apidebugger.Event {
	if event == nil {
		return nil
	}

	converted := &apidebugger.Event{
		Error:            event.Error,
		Reason:           apidebugger.ReasonFromString(string(event.Reason)),
		HitBreakpointIDs: make([]apidebugger.BreakpointID, len(event.HitBreakpointIDs)),
		Location:         s.location(event.Location),
		Depth:            event.Depth,
	}

	for i := range event.HitBreakpointIDs {
		converted.HitBreakpointIDs[i] = apidebugger.BreakpointID(event.HitBreakpointIDs[i])
	}

	if event.Output != nil {
		converted.Output = &result.Output{
			ContentType: event.Output.ContentType,
			Content:     event.Output.Content,
		}
	}

	return converted
}

func (s *apiDebugSession) breakpoint(breakpoint coredebugger.Breakpoint) apidebugger.Breakpoint {
	return apidebugger.Breakpoint{
		Location: apisource.Range{
			Location: apisource.Location{
				Position: apisource.Position{
					Line:   breakpoint.Line,
					Column: breakpoint.Column,
				},
				File: breakpoint.File,
			},
		},
		RequestedLocation: apisource.Location{
			Position: apisource.Position{
				Line:   breakpoint.RequestedLine,
				Column: breakpoint.RequestedColumn,
			},
			File: breakpoint.File,
		},
		ID:          apidebugger.BreakpointID(breakpoint.ID),
		PointID:     apidebugger.PointID(breakpoint.PointID),
		FunctionID:  apidebugger.FunctionID(breakpoint.FunctionID),
		BindingMode: apidebugger.BreakpointBindingMode(breakpoint.BindingMode),
		Bound:       breakpoint.Bound,
	}
}

func (s *apiDebugSession) variables(variables []coredebugger.Variable) []apidebugger.Variable {
	converted := make([]apidebugger.Variable, len(variables))

	for i := range variables {
		converted[i] = apidebugger.Variable{
			Name:    variables[i].Name,
			Value:   s.value(variables[i].Value),
			Mutable: variables[i].Mutable,
			Param:   variables[i].Param,
		}
	}

	return converted
}

func (s *apiDebugSession) value(value coredebugger.Value) apidebugger.Value {
	return apidebugger.Value{
		Type:      value.Type,
		Display:   value.Display,
		Reference: apidebugger.ValueReference(value.Reference),
	}
}

func (s *apiDebugSession) location(location coredebugger.Location) apisource.Range {
	return apisource.Range{
		Location: apisource.Location{
			Position: apisource.Position{
				Line:   location.Line,
				Column: location.Column,
			},
			File: location.File,
		},
		Span: apisource.Span{
			Start: location.Span.Start,
			End:   location.Span.End,
		},
	}
}
