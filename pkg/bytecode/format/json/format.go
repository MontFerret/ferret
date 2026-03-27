package json

import (
	"fmt"

	"github.com/goccy/go-json"

	"github.com/MontFerret/ferret/v2/pkg/bytecode"
	"github.com/MontFerret/ferret/v2/pkg/bytecode/internal/persist"
)

type Format struct{}

// Default is the built-in JSON bytecode payload format.
var Default = Format{}

func (Format) Name() string {
	return "json"
}

func (Format) Marshal(program *bytecode.Program) ([]byte, error) {
	frame, err := persist.FromProgram(program)
	if err != nil {
		return nil, err
	}

	data, err := json.Marshal(frame)
	if err != nil {
		return nil, fmt.Errorf("bytecode json format: marshal payload: %w", err)
	}

	return data, nil
}

func (Format) Unmarshal(data []byte) (*bytecode.Program, error) {
	var frame persist.ProgramFrame
	if err := json.Unmarshal(data, &frame); err != nil {
		return nil, fmt.Errorf("bytecode json format: unmarshal payload: %w", err)
	}

	return persist.ToProgram(frame)
}
