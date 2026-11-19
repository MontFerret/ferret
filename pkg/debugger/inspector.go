package debugger

import (
	"context"
	"slices"

	apidebugger "github.com/MontFerret/api/debugger"

	"github.com/MontFerret/ferret/v2/pkg/bytecode"
	"github.com/MontFerret/ferret/v2/pkg/internal/debugpoint"
	"github.com/MontFerret/ferret/v2/pkg/runtime"
	"github.com/MontFerret/ferret/v2/pkg/source"
	"github.com/MontFerret/ferret/v2/pkg/vm"
)

// inspector borrows paused execution and owns its presentation and references.
// The session command lock protects every operation, including invalidation.
type inspector struct {
	execution    vm.DebugExecution
	values       vm.DebugValueAccess
	valueRefs    map[ValueReference]runtime.Value
	pointIndex   *debugpoint.Index
	source       *source.Source
	occurrences  map[occurrence]ValueReference
	lifecycle    *sessionLifecycle
	params       []string
	format       FormatOptions
	lastValueRef ValueReference
	handleLimit  int
}

func newInspector(lifecycle *sessionLifecycle, config Config, src *source.Source, points *debugpoint.Index, format FormatOptions) inspector {
	params := slices.Clone(config.Params)
	slices.Sort(params)
	params = slices.Compact(params)

	return inspector{
		execution:   config.Execution,
		values:      config.Values,
		valueRefs:   make(map[ValueReference]runtime.Value),
		pointIndex:  points,
		params:      params,
		source:      src,
		format:      format,
		handleLimit: 1 << 16,
		occurrences: make(map[occurrence]ValueReference),
		lifecycle:   lifecycle,
	}
}

func (s *inspector) frames() ([]Frame, error) {
	frames, err := s.execution.Frames()
	if err != nil {
		return nil, err
	}

	out := make([]Frame, 0, len(frames))

	for _, frame := range frames {
		out = append(out, Frame{
			Name:       frame.Name,
			FunctionID: apidebugger.FunctionID(frame.FunctionID),
			Location:   s.locationForPC(frame.PC, frame.FunctionID),
		})
	}

	return out, nil
}

func (s *inspector) frameLocals(ctx context.Context, frame int) ([]Variable, error) {
	locals, err := s.rawFrameLocals(frame)
	if err != nil {
		return nil, err
	}

	batch := inspectionBatch{owner: s}
	out := make([]Variable, 0, len(locals)+len(s.params))

	names := make(map[string]struct{}, len(locals))
	for _, local := range locals {
		if _, duplicate := names[local.Name]; duplicate {
			return nil, runtime.Errorf(runtime.ErrUnexpected, "duplicate visible local %q in frame %d", local.Name, frame)
		}
		names[local.Name] = struct{}{}

		value, err := batch.describe(ctx, occurrence{frame: frame, kind: localBinding, name: local.Name}, local.Value)
		if err != nil {
			return nil, err
		}

		out = append(out, Variable{Name: local.Name, Mutable: local.Mutable, Value: value})
	}

	if len(s.params) > 0 {
		params := s.execution.Params()

		for _, name := range s.params {
			if raw, ok := params.Get(name); ok {
				value, err := batch.describe(ctx, occurrence{frame: frame, kind: paramBinding, name: name}, raw)
				if err != nil {
					return nil, err
				}

				out = append(out, Variable{Name: "@" + name, Param: true, Value: value})
			}
		}
	}

	if err := batch.publish(ctx); err != nil {
		return nil, err
	}
	return out, nil
}

func (s *inspector) variables(ctx context.Context, reference ValueReference, page VariablePage) ([]Variable, error) {
	page, err := page.Normalize()
	if err != nil {
		return nil, err
	}

	raw, ok := s.valueRefs[reference]
	if !ok || !reference.Valid() {
		return nil, apidebugger.ErrInvalidValueReference
	}

	// Close can interrupt a cooperative provider without canceling execution on
	// request cancellation. Only one command can own activeCancel at a time.
	readCtx, cancel := context.WithCancelCause(ctx)
	id := s.lifecycle.setActiveCancel(cancel)
	defer s.lifecycle.clearActiveCancel(id)
	defer cancel(nil)

	items, err := s.values.ReadChildren(readCtx, raw, page)
	if err != nil {
		return nil, err
	}

	if int64(len(items)) > page.Count {
		return nil, runtime.Error(runtime.ErrUnexpected, "debug value access exceeded requested count")
	}

	batch := inspectionBatch{owner: s}
	out := make([]Variable, 0, len(items))
	for _, item := range items {
		kind := indexedChild
		if item.Kind == VariableFilterNamed {
			kind = namedChild
		}

		value, err := batch.describe(readCtx, occurrence{parent: reference, kind: kind, position: item.Position}, item.Value)
		if err != nil {
			return nil, err
		}

		out = append(out, Variable{Name: item.Name, Value: value})
	}

	if err := batch.publish(ctx); err != nil {
		return nil, err
	}
	return out, nil
}

func (s *inspector) evaluateFrame(ctx context.Context, frame int, expression string) (Value, error) {
	locals, err := s.rawFrameLocals(frame)
	if err != nil {
		return Value{}, err
	}

	values := make(map[string]runtime.Value, len(locals))

	for _, local := range locals {
		values[local.Name] = local.Value
	}

	value, err := evaluateExpression(ctx, expression, evalScope{
		locals:   values,
		params:   s.execution.Params(),
		values:   s.values,
		maxBytes: s.format.MaxBytes,
	})
	if err != nil {
		return Value{}, err
	}

	batch := inspectionBatch{owner: s}
	result, err := batch.describe(ctx, occurrence{}, value)
	if err != nil {
		return Value{}, err
	}

	if err := batch.publish(ctx); err != nil {
		return Value{}, err
	}

	return result, nil
}

func (s *inspector) rawFrameLocals(frame int) ([]vm.DebugLocal, error) {
	if inspector, ok := s.execution.(vm.DebugFrameInspector); ok {
		return inspector.FrameLocals(frame)
	}

	if frame != 0 {
		return nil, runtime.Error(runtime.ErrInvalidOperation, "debug execution does not support caller frame inspection")
	}

	return s.execution.Locals()
}

func (s *inspector) locationForPC(pc int, functionID bytecode.FunctionID) source.Location {
	var point *bytecode.DebugPoint

	if functionID == bytecode.NoFunction || functionID.Valid() {
		point = s.pointIndex.NearestBeforeOrAtInFunction(functionID, pc)
	} else {
		point = s.pointIndex.NearestBeforeOrAt(pc)
	}

	if point == nil {
		return source.Location{SourceName: s.source.Name()}
	}

	return s.source.LocationAt(point.Span)
}

func (s *inspector) resetValueReferences() {
	clear(s.valueRefs)
	clear(s.occurrences)
}
