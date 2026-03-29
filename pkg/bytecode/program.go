package bytecode

import (
	"fmt"

	"github.com/goccy/go-json"

	"github.com/MontFerret/ferret/v2/pkg/source"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

const Version = 1

type (
	// Catch stores an inclusive instruction range [start, end] and an optional recovery jump target.
	Catch [3]int

	Metadata struct {
		Labels                 map[int]string  `json:"labels"`
		CompilerVersion        string          `json:"compilerVersion"`
		AggregatePlans         []AggregatePlan `json:"aggregatePlans"`
		AggregateSelectorSlots []int           `json:"aggregateSelectorSlots,omitempty"`
		MatchFailTargets       []int           `json:"matchFailTargets,omitempty"`
		DebugSpans             []source.Span   `json:"debugSpans"`
		OptimizationLevel      int             `json:"optimizationLevel"`
	}

	Program struct {
		Source     *source.Source
		Functions  Functions
		Bytecode   []Instruction
		Constants  []runtime.Value
		CatchTable []Catch
		Params     []string
		Metadata   Metadata
		ISAVersion int
		Registers  int
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
		ISAVersion: p.ISAVersion,
		Source:     p.Source,
		Registers:  p.Registers,
		Bytecode:   p.Bytecode,
		Constants:  constants,
		CatchTable: p.CatchTable,
		Params:     p.Params,
		Functions:  p.Functions,
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
	p.ISAVersion = decoded.ISAVersion
	p.Registers = decoded.Registers
	p.Bytecode = decoded.Bytecode
	p.Constants = constants
	p.CatchTable = decoded.CatchTable
	p.Params = decoded.Params
	p.Functions = decoded.Functions
	p.Metadata = decoded.Metadata

	return nil
}
