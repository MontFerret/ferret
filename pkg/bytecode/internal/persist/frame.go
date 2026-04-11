package persist

import (
	"fmt"
	"sort"
	"strconv"
	"time"

	"github.com/MontFerret/ferret/v2/pkg/bytecode"
	"github.com/MontFerret/ferret/v2/pkg/runtime"
	"github.com/MontFerret/ferret/v2/pkg/source"
)

const (
	constantTypeNone     = "none"
	constantTypeBool     = "bool"
	constantTypeInt      = "int"
	constantTypeFloat    = "float"
	constantTypeString   = "string"
	constantTypeBinary   = "binary"
	constantTypeDateTime = "datetime"
)

type (
	ProgramFrame struct {
		ISAVersion *int                `json:"isaVersion" msgpack:"isaVersion"`
		Registers  *int                `json:"registers" msgpack:"registers"`
		Bytecode   *[]InstructionFrame `json:"bytecode" msgpack:"bytecode"`
		Source     *SourceFrame        `json:"source,omitempty" msgpack:"source,omitempty"`
		Functions  FunctionsFrame      `json:"functions" msgpack:"functions"`
		Constants  []ConstantFrame     `json:"constants,omitempty" msgpack:"constants,omitempty"`
		CatchTable []CatchFrame        `json:"catchTable,omitempty" msgpack:"catchTable,omitempty"`
		Params     []string            `json:"params,omitempty" msgpack:"params,omitempty"`
		Metadata   MetadataFrame       `json:"metadata" msgpack:"metadata"`
	}

	InstructionFrame struct {
		Opcode   uint8    `json:"opcode" msgpack:"opcode"`
		Operands [3]int64 `json:"operands" msgpack:"operands"`
	}

	ConstantFrame struct {
		Bool     *bool    `json:"bool,omitempty" msgpack:"bool,omitempty"`
		Int      *int64   `json:"int,omitempty" msgpack:"int,omitempty"`
		Float    *float64 `json:"float,omitempty" msgpack:"float,omitempty"`
		String   *string  `json:"string,omitempty" msgpack:"string,omitempty"`
		DateTime *string  `json:"datetime,omitempty" msgpack:"datetime,omitempty"`
		Type     string   `json:"type" msgpack:"type"`
		Binary   []byte   `json:"binary,omitempty" msgpack:"binary,omitempty"`
	}

	CatchFrame struct {
		StartPC   int `json:"startPC" msgpack:"startPC"`
		EndPC     int `json:"endPC" msgpack:"endPC"`
		HandlerPC int `json:"handlerPC" msgpack:"handlerPC"`
	}

	FunctionsFrame struct {
		Host        []HostFunctionFrame `json:"host,omitempty" msgpack:"host,omitempty"`
		UserDefined []UDFFrame          `json:"userDefined,omitempty" msgpack:"userDefined,omitempty"`
	}

	HostFunctionFrame struct {
		Name  string `json:"name" msgpack:"name"`
		Arity int    `json:"arity" msgpack:"arity"`
	}

	UDFFrame struct {
		Name        string `json:"name" msgpack:"name"`
		DisplayName string `json:"displayName,omitempty" msgpack:"displayName,omitempty"`
		Entry       int    `json:"entry" msgpack:"entry"`
		Registers   int    `json:"registers" msgpack:"registers"`
		Params      int    `json:"params" msgpack:"params"`
	}

	MetadataFrame struct {
		Labels                 []LabelFrame         `json:"labels,omitempty" msgpack:"labels,omitempty"`
		CompilerVersion        string               `json:"compilerVersion,omitempty" msgpack:"compilerVersion,omitempty"`
		AggregatePlans         []AggregatePlanFrame `json:"aggregatePlans,omitempty" msgpack:"aggregatePlans,omitempty"`
		AggregateSelectorSlots []int                `json:"aggregateSelectorSlots,omitempty" msgpack:"aggregateSelectorSlots,omitempty"`
		CallArgumentSpans      [][]SpanFrame        `json:"callArgumentSpans,omitempty" msgpack:"callArgumentSpans,omitempty"`
		MatchFailTargets       []int                `json:"matchFailTargets,omitempty" msgpack:"matchFailTargets,omitempty"`
		DebugSpans             []SpanFrame          `json:"debugSpans,omitempty" msgpack:"debugSpans,omitempty"`
		OptimizationLevel      int                  `json:"optimizationLevel" msgpack:"optimizationLevel"`
	}

	LabelFrame struct {
		Name string `json:"name" msgpack:"name"`
		PC   int    `json:"pc" msgpack:"pc"`
	}

	AggregatePlanFrame struct {
		Keys             []string `json:"keys,omitempty" msgpack:"keys,omitempty"`
		Kinds            []int    `json:"kinds,omitempty" msgpack:"kinds,omitempty"`
		TrackGroupValues bool     `json:"trackGroupValues,omitempty" msgpack:"trackGroupValues,omitempty"`
	}

	SpanFrame struct {
		Start int `json:"start" msgpack:"start"`
		End   int `json:"end" msgpack:"end"`
	}

	SourceFrame struct {
		Name string `json:"name,omitempty" msgpack:"name,omitempty"`
		Text string `json:"text" msgpack:"text"`
	}
)

func FromProgram(program *bytecode.Program) (ProgramFrame, error) {
	if err := bytecode.ValidateProgram(program); err != nil {
		return ProgramFrame{}, err
	}

	constants := make([]ConstantFrame, len(program.Constants))
	for i, value := range program.Constants {
		encoded, err := encodeConstant(value)
		if err != nil {
			return ProgramFrame{}, fmt.Errorf("encode constant %d: %w", i, err)
		}

		constants[i] = encoded
	}

	instructions := make([]InstructionFrame, len(program.Bytecode))
	for i, inst := range program.Bytecode {
		instructions[i] = InstructionFrame{
			Opcode: uint8(inst.Opcode),
			Operands: [3]int64{
				int64(inst.Operands[0]),
				int64(inst.Operands[1]),
				int64(inst.Operands[2]),
			},
		}
	}

	catches := make([]CatchFrame, len(program.CatchTable))
	for i, entry := range program.CatchTable {
		catches[i] = CatchFrame{
			StartPC:   entry[0],
			EndPC:     entry[1],
			HandlerPC: entry[2],
		}
	}

	host := make([]HostFunctionFrame, 0, len(program.Functions.Host))
	if len(program.Functions.Host) > 0 {
		names := make([]string, 0, len(program.Functions.Host))
		for name := range program.Functions.Host {
			names = append(names, name)
		}

		sort.Strings(names)

		for _, name := range names {
			host = append(host, HostFunctionFrame{
				Name:  name,
				Arity: program.Functions.Host[name],
			})
		}
	}

	udfs := make([]UDFFrame, len(program.Functions.UserDefined))
	for i, udf := range program.Functions.UserDefined {
		udfs[i] = UDFFrame{
			Name:        udf.Name,
			DisplayName: udf.DisplayName,
			Entry:       udf.Entry,
			Registers:   udf.Registers,
			Params:      udf.Params,
		}
	}

	aggregatePlans := make([]AggregatePlanFrame, len(program.Metadata.AggregatePlans))
	for i, plan := range program.Metadata.AggregatePlans {
		keys := make([]string, len(plan.Keys))
		for j, key := range plan.Keys {
			keys[j] = key.String()
		}

		kinds := make([]int, len(plan.Kinds))
		for j, kind := range plan.Kinds {
			kinds[j] = int(kind)
		}

		aggregatePlans[i] = AggregatePlanFrame{
			Keys:             keys,
			Kinds:            kinds,
			TrackGroupValues: plan.TrackGroupValues,
		}
	}

	debugSpans := make([]SpanFrame, len(program.Metadata.DebugSpans))
	for i, span := range program.Metadata.DebugSpans {
		debugSpans[i] = SpanFrame{
			Start: span.Start,
			End:   span.End,
		}
	}

	callArgumentSpans := make([][]SpanFrame, len(program.Metadata.CallArgumentSpans))
	for i, spans := range program.Metadata.CallArgumentSpans {
		if len(spans) == 0 {
			continue
		}

		callArgumentSpans[i] = make([]SpanFrame, len(spans))
		for j, span := range spans {
			callArgumentSpans[i][j] = SpanFrame{
				Start: span.Start,
				End:   span.End,
			}
		}
	}

	labels := make([]LabelFrame, 0, len(program.Metadata.Labels))
	if len(program.Metadata.Labels) > 0 {
		pcs := make([]int, 0, len(program.Metadata.Labels))
		for pc := range program.Metadata.Labels {
			pcs = append(pcs, pc)
		}

		sort.Ints(pcs)

		for _, pc := range pcs {
			labels = append(labels, LabelFrame{
				PC:   pc,
				Name: program.Metadata.Labels[pc],
			})
		}
	}

	var source *SourceFrame
	if program.Source != nil {
		source = &SourceFrame{
			Name: program.Source.Name(),
			Text: program.Source.Content(),
		}
	}

	isaVersion := program.ISAVersion
	registers := program.Registers
	bytecodeCopy := instructions

	return ProgramFrame{
		ISAVersion: &isaVersion,
		Registers:  &registers,
		Bytecode:   &bytecodeCopy,
		Constants:  constants,
		CatchTable: catches,
		Params:     append([]string(nil), program.Params...),
		Functions: FunctionsFrame{
			Host:        host,
			UserDefined: udfs,
		},
		Metadata: MetadataFrame{
			Labels:                 labels,
			CompilerVersion:        program.Metadata.CompilerVersion,
			AggregatePlans:         aggregatePlans,
			AggregateSelectorSlots: append([]int(nil), program.Metadata.AggregateSelectorSlots...),
			CallArgumentSpans:      callArgumentSpans,
			MatchFailTargets:       append([]int(nil), program.Metadata.MatchFailTargets...),
			DebugSpans:             debugSpans,
			OptimizationLevel:      program.Metadata.OptimizationLevel,
		},
		Source: source,
	}, nil
}

func ToProgram(frame ProgramFrame) (*bytecode.Program, error) {
	if frame.ISAVersion == nil {
		return nil, fmt.Errorf("%w: missing isaVersion field", bytecode.ErrInvalidProgram)
	}

	if frame.Registers == nil {
		return nil, fmt.Errorf("%w: missing registers field", bytecode.ErrInvalidProgram)
	}

	if frame.Bytecode == nil {
		return nil, fmt.Errorf("%w: missing bytecode field", bytecode.ErrInvalidProgram)
	}

	constants := make([]runtime.Value, len(frame.Constants))
	for i, value := range frame.Constants {
		decoded, err := decodeConstant(value)
		if err != nil {
			return nil, fmt.Errorf("decode constant %d: %w", i, err)
		}

		constants[i] = decoded
	}

	instructions := make([]bytecode.Instruction, len(*frame.Bytecode))
	for i, inst := range *frame.Bytecode {
		var operands [3]bytecode.Operand

		for operandIdx, value := range inst.Operands {
			decoded, err := decodeInstructionOperand(value)
			if err != nil {
				return nil, fmt.Errorf("decode instruction %d operand %d: %w", i, operandIdx, err)
			}

			operands[operandIdx] = decoded
		}

		instructions[i] = bytecode.Instruction{
			Opcode:   bytecode.Opcode(inst.Opcode),
			Operands: operands,
		}
	}

	catches := make([]bytecode.Catch, len(frame.CatchTable))
	for i, entry := range frame.CatchTable {
		catches[i] = bytecode.Catch{entry.StartPC, entry.EndPC, entry.HandlerPC}
	}

	var host map[string]int
	if len(frame.Functions.Host) > 0 {
		host = make(map[string]int, len(frame.Functions.Host))
	}

	for _, entry := range frame.Functions.Host {
		if _, exists := host[entry.Name]; exists {
			return nil, fmt.Errorf("%w: duplicate host function %q", bytecode.ErrInvalidProgram, entry.Name)
		}

		host[entry.Name] = entry.Arity
	}

	udfs := make([]bytecode.UDF, len(frame.Functions.UserDefined))
	for i, udf := range frame.Functions.UserDefined {
		udfs[i] = bytecode.UDF{
			Name:        udf.Name,
			DisplayName: udf.DisplayName,
			Entry:       udf.Entry,
			Registers:   udf.Registers,
			Params:      udf.Params,
		}
	}

	aggregatePlans := make([]bytecode.AggregatePlan, len(frame.Metadata.AggregatePlans))
	for i, plan := range frame.Metadata.AggregatePlans {
		keys := make([]runtime.String, len(plan.Keys))
		for j, key := range plan.Keys {
			keys[j] = runtime.NewString(key)
		}

		kinds := make([]bytecode.AggregateKind, len(plan.Kinds))
		for j, kind := range plan.Kinds {
			kinds[j] = bytecode.AggregateKind(kind)
		}

		aggregatePlans[i] = bytecode.NewAggregatePlan(keys, kinds, plan.TrackGroupValues)
	}

	debugSpans := make([]source.Span, len(frame.Metadata.DebugSpans))
	for i, span := range frame.Metadata.DebugSpans {
		debugSpans[i] = source.Span{
			Start: span.Start,
			End:   span.End,
		}
	}

	callArgumentSpans := make([][]source.Span, len(frame.Metadata.CallArgumentSpans))
	for i, spans := range frame.Metadata.CallArgumentSpans {
		if len(spans) == 0 {
			continue
		}

		callArgumentSpans[i] = make([]source.Span, len(spans))
		for j, span := range spans {
			callArgumentSpans[i][j] = source.Span{
				Start: span.Start,
				End:   span.End,
			}
		}
	}

	var labels map[int]string
	if len(frame.Metadata.Labels) > 0 {
		labels = make(map[int]string, len(frame.Metadata.Labels))
	}

	for _, label := range frame.Metadata.Labels {
		if _, exists := labels[label.PC]; exists {
			return nil, fmt.Errorf("%w: duplicate label pc %d", bytecode.ErrInvalidProgram, label.PC)
		}

		labels[label.PC] = label.Name
	}

	var src *source.Source
	if frame.Source != nil {
		src = source.New(frame.Source.Name, frame.Source.Text)
	}

	program := &bytecode.Program{
		Source: src,
		Functions: bytecode.Functions{
			Host:        host,
			UserDefined: udfs,
		},
		Bytecode:   instructions,
		Constants:  constants,
		CatchTable: catches,
		Params:     append([]string(nil), frame.Params...),
		Metadata: bytecode.Metadata{
			Labels:                 labels,
			CompilerVersion:        frame.Metadata.CompilerVersion,
			AggregatePlans:         aggregatePlans,
			AggregateSelectorSlots: append([]int(nil), frame.Metadata.AggregateSelectorSlots...),
			CallArgumentSpans:      callArgumentSpans,
			MatchFailTargets:       append([]int(nil), frame.Metadata.MatchFailTargets...),
			DebugSpans:             debugSpans,
			OptimizationLevel:      frame.Metadata.OptimizationLevel,
		},
		ISAVersion: *frame.ISAVersion,
		Registers:  *frame.Registers,
	}

	if err := bytecode.ValidateProgram(program); err != nil {
		return nil, err
	}

	return program, nil
}

func decodeInstructionOperand(value int64) (bytecode.Operand, error) {
	if err := validateInstructionOperandForBitSize(value, strconv.IntSize); err != nil {
		return 0, err
	}

	return bytecode.Operand(value), nil
}

func validateInstructionOperandForBitSize(value int64, bitSize int) error {
	minValue, maxValue := operandRangeForBitSize(bitSize)
	if value < minValue || value > maxValue {
		return fmt.Errorf("%w: operand value %d overflows %d-bit operand range [%d,%d]", bytecode.ErrInvalidProgram, value, bitSize, minValue, maxValue)
	}

	return nil
}

func operandRangeForBitSize(bitSize int) (int64, int64) {
	maxValue := int64((uint64(1) << (bitSize - 1)) - 1)

	return -maxValue - 1, maxValue
}

func encodeConstant(value runtime.Value) (ConstantFrame, error) {
	if value == nil || value == runtime.None {
		return ConstantFrame{Type: constantTypeNone}, nil
	}

	switch v := value.(type) {
	case runtime.Boolean:
		value := bool(v)
		return ConstantFrame{Type: constantTypeBool, Bool: &value}, nil
	case runtime.Int:
		value := int64(v)
		return ConstantFrame{Type: constantTypeInt, Int: &value}, nil
	case runtime.Float:
		value := float64(v)
		return ConstantFrame{Type: constantTypeFloat, Float: &value}, nil
	case runtime.String:
		value := string(v)
		return ConstantFrame{Type: constantTypeString, String: &value}, nil
	case runtime.Binary:
		value := append([]byte(nil), []byte(v)...)
		return ConstantFrame{Type: constantTypeBinary, Binary: value}, nil
	case runtime.DateTime:
		value := v.Time.Format(time.RFC3339Nano)
		return ConstantFrame{Type: constantTypeDateTime, DateTime: &value}, nil
	default:
		return ConstantFrame{}, fmt.Errorf("%w: unsupported constant type %T", bytecode.ErrInvalidConstant, value)
	}
}

func decodeConstant(value ConstantFrame) (runtime.Value, error) {
	switch value.Type {
	case constantTypeNone:
		return runtime.None, nil
	case constantTypeBool:
		if value.Bool == nil {
			return nil, fmt.Errorf("%w: bool constant missing value", bytecode.ErrInvalidConstant)
		}

		return runtime.NewBoolean(*value.Bool), nil
	case constantTypeInt:
		if value.Int == nil {
			return nil, fmt.Errorf("%w: int constant missing value", bytecode.ErrInvalidConstant)
		}

		return runtime.NewInt64(*value.Int), nil
	case constantTypeFloat:
		if value.Float == nil {
			return nil, fmt.Errorf("%w: float constant missing value", bytecode.ErrInvalidConstant)
		}

		return runtime.NewFloat(*value.Float), nil
	case constantTypeString:
		if value.String == nil {
			return nil, fmt.Errorf("%w: string constant missing value", bytecode.ErrInvalidConstant)
		}

		return runtime.NewString(*value.String), nil
	case constantTypeBinary:
		return runtime.NewBinary(append([]byte(nil), value.Binary...)), nil
	case constantTypeDateTime:
		if value.DateTime == nil {
			return nil, fmt.Errorf("%w: datetime constant missing value", bytecode.ErrInvalidConstant)
		}

		parsed, err := time.Parse(time.RFC3339Nano, *value.DateTime)
		if err != nil {
			return nil, fmt.Errorf("%w: invalid datetime literal %q: %v", bytecode.ErrInvalidConstant, *value.DateTime, err)
		}

		return runtime.NewDateTime(parsed), nil
	default:
		return nil, fmt.Errorf("%w: unsupported constant type %q", bytecode.ErrInvalidConstant, value.Type)
	}
}
