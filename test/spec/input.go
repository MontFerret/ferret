package spec

import (
	"github.com/MontFerret/ferret/v2/pkg/bytecode"
	"github.com/MontFerret/ferret/v2/pkg/compiler"
	"github.com/MontFerret/ferret/v2/pkg/source"
)

type (
	ProgramBuilder func(name string, c *compiler.Compiler) (*bytecode.Program, error)

	ProgramSource struct {
		Build ProgramBuilder
		Name  string
	}

	Input struct {
		Source     ProgramSource
		Expression string
	}
)

func NewExpressionInput(expression string) Input {
	return Input{
		Expression: expression,
	}
}

func NewProgramInput(program *bytecode.Program, name ...string) Input {
	var srcName string

	if len(name) > 0 {
		srcName = name[0]
	} else {
		if program.Source != nil {
			srcName = joinExpression(program.Source.Content())
		} else {
			srcName = joinBytecode(program.Bytecode)
		}
	}

	return Input{
		Source: ProgramSource{
			Name: srcName,
			Build: func(name string, c *compiler.Compiler) (*bytecode.Program, error) {
				return program, nil
			},
		},
	}
}

func NewProgramSourceInput(source ProgramSource) Input {
	return Input{
		Source: source,
	}
}

func (i Input) Merge(other Input) Input {
	if i.Expression == "" && other.Expression != "" {
		return Input{
			Expression: other.Expression,
		}
	}

	return Input{
		Source: other.Source,
	}
}

func (i Input) String() string {
	if i.Expression != "" {
		return joinExpression(i.Expression)
	}

	if i.Source.Name != "" {
		return i.Source.Name
	}

	return "Anonymous program"
}

func (i Input) ResolveProgram(name string, c *compiler.Compiler) (*bytecode.Program, error) {
	if i.Source.Build != nil {
		return i.Source.Build(name, c)
	}

	return c.Compile(source.New(name, i.Expression))
}
