package bytecode

import (
	"fmt"

	"github.com/goccy/go-json"

	"github.com/MontFerret/ferret/v2/pkg/file"
	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

type (
	Catch [3]int

	Metadata struct {
		AggregatePlans []AggregatePlan `json:"aggregatePlans"`
		DebugSpans     []file.Span     `json:"debugSpans"`
		Functions      map[string]int  `json:"functions"`
		Labels         map[int]string  `json:"labels"`
		UDFs           []UDF           `json:"udfs,omitempty"`
	}

	Program struct {
		Source     *file.Source
		Registers  int
		Bytecode   []Instruction
		Constants  []runtime.Value
		CatchTable []Catch
		Params     []string
		Metadata   Metadata
	}
)

func (p *Program) MarshalJSON() ([]byte, error) {
	if p == nil {
		return []byte("null"), nil
	}

	constants := make([]constantJSON, len(p.Constants))

	for i, value := range p.Constants {
		encoded, err := encodeConstant(value)

		if err != nil {
			return nil, fmt.Errorf("bytecode.Program: encode constant %d: %w", i, err)
		}

		constants[i] = encoded
	}

	payload := programJSON{
		Source:     p.Source,
		Registers:  p.Registers,
		Bytecode:   p.Bytecode,
		Constants:  constants,
		CatchTable: p.CatchTable,
		Params:     p.Params,
		Metadata:   p.Metadata,
	}

	return json.Marshal(payload)
}

func (p *Program) UnmarshalJSON(data []byte) error {
	if p == nil {
		return fmt.Errorf("bytecode.Program: UnmarshalJSON on nil pointer")
	}

	var decoded programJSON
	if err := json.Unmarshal(data, &decoded); err != nil {
		return err
	}

	constants := make([]runtime.Value, len(decoded.Constants))

	for i, value := range decoded.Constants {
		decodedValue, err := decodeConstant(value)

		if err != nil {
			return fmt.Errorf("bytecode.Program: decode constant %d: %w", i, err)
		}

		constants[i] = decodedValue
	}

	p.Source = decoded.Source
	p.Registers = decoded.Registers
	p.Bytecode = decoded.Bytecode
	p.Constants = constants
	p.CatchTable = decoded.CatchTable
	p.Params = decoded.Params
	p.Metadata = decoded.Metadata

	return nil
}
