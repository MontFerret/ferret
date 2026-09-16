package debugger

import (
	"context"
	"sort"
	"strconv"

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
	params       []string
	format       FormatOptions
	nextValueRef ValueReference
}

func newInspector(config Config, src *source.Source, points *debugpoint.Index, format FormatOptions) inspector {
	return inspector{
		execution:    config.Execution,
		values:       config.Values,
		valueRefs:    make(map[ValueReference]runtime.Value),
		pointIndex:   points,
		params:       append([]string(nil), config.Params...),
		source:       src,
		format:       format,
		nextValueRef: 1,
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

func (s *inspector) frameLocals(frame int) ([]Variable, error) {
	locals, err := s.rawFrameLocals(frame)
	if err != nil {
		return nil, err
	}

	out := make([]Variable, 0, len(locals)+len(s.execution.Params()))

	for _, local := range locals {
		out = append(out, Variable{Name: local.Name, Mutable: local.Mutable, Value: s.debugValue(local.Value)})
	}

	params := s.execution.Params()
	names := append([]string(nil), s.params...)

	sort.Strings(names)

	for _, name := range names {
		if value, exists := params.Get(name); exists {
			out = append(out, Variable{Name: "@" + name, Param: true, Value: s.debugValue(value)})
		}
	}

	return out, nil
}

func (s *inspector) variables(reference ValueReference) ([]Variable, error) {
	if !reference.Valid() {
		return nil, runtime.Error(runtime.ErrInvalidArgument, "debug value reference must be positive")
	}

	value, exists := s.valueRefs[reference]
	if !exists {
		return nil, runtime.Errorf(runtime.ErrNotFound, "debug value reference %d", reference)
	}

	inspection, ok := s.expandableInspection(value)
	if !ok {
		return nil, runtime.Errorf(runtime.ErrInvalidOperation, "debug value reference %d is not expandable", reference)
	}

	if inspection.Kind == vm.DebugValueArray {
		out := make([]Variable, 0, len(inspection.Items))

		for i, item := range inspection.Items {
			out = append(out, Variable{
				Name:  strconv.Itoa(i),
				Value: s.debugValue(item.Value),
			})
		}

		return out, nil
	}

	items := append([]vm.DebugValueItem(nil), inspection.Items...)
	sort.Slice(items, func(i, j int) bool { return items[i].Key < items[j].Key })

	out := make([]Variable, 0, len(items))
	for _, item := range items {
		out = append(out, Variable{
			Name:  item.Key,
			Value: s.debugValue(item.Value),
		})
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
		locals: values,
		params: s.execution.Params(),
		values: s.values,
	})
	if err != nil {
		return Value{}, err
	}

	return s.debugValue(value), nil
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

func (s *inspector) debugValue(value runtime.Value) Value {
	info, _ := s.values.DebugInfo(value)
	typeName := info.TypeName

	if typeName == "" {
		typeName = s.values.TypeName(value)
		info.TypeName = typeName
	} else {
		typeName = boundedText(typeName, s.format.MaxBytes)
	}

	return Value{
		Reference: s.referenceForValue(value),
		Type:      typeName,
		Display:   formatValueWithInfo(value, info, s.values, s.format),
	}
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
	if len(s.valueRefs) == 0 {
		return
	}

	clear(s.valueRefs)
}

func (s *inspector) referenceForValue(value runtime.Value) ValueReference {
	if _, ok := s.expandableInspection(value); !ok {
		return 0
	}

	reference := s.nextValueRef
	s.nextValueRef++
	s.valueRefs[reference] = value

	return reference
}

func (s *inspector) expandableInspection(value runtime.Value) (vm.DebugValueInspection, bool) {
	inspection, ok := s.values.Inspect(value, s.format.MaxItems)
	if !ok || !inspection.Complete || inspection.Length <= 0 || inspection.Length > s.format.MaxItems {
		return vm.DebugValueInspection{}, false
	}

	if inspection.Kind != vm.DebugValueArray && inspection.Kind != vm.DebugValueObject {
		return vm.DebugValueInspection{}, false
	}

	return inspection, true
}
