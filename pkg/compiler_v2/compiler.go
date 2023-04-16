package compiler_v2

import (
	"errors"

	"github.com/MontFerret/ferret/pkg/parser"
	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime_v2"
	"github.com/MontFerret/ferret/pkg/stdlib"
)

type Compiler struct {
	*NamespaceContainer
}

func New(setters ...Option) *Compiler {
	c := &Compiler{}
	c.NamespaceContainer = newRootNamespace()
	c.funcs = core.NewFunctions()

	opts := &Options{}

	for _, setter := range setters {
		setter(opts)
	}

	if !opts.noStdlib {
		if err := stdlib.RegisterLib(c.NamespaceContainer); err != nil {
			panic(err)
		}
	}

	return c
}

func (c *Compiler) Compile(query string) (program *runtime_v2.Program, err error) {
	if query == "" {
		return nil, ErrEmptyQuery
	}

	defer func() {
		if r := recover(); r != nil {
			// find out exactly what the error was and set err
			switch x := r.(type) {
			case string:
				err = errors.New(x)
			case error:
				err = x
			default:
				err = errors.New("unknown panic")
			}

			program = nil
		}
	}()

	p := parser.New(query)
	p.AddErrorListener(newErrorListener())

	l := newVisitor(query, c.funcs)

	p.Visit(l)

	if l.err != nil {
		return nil, l.err
	}

	program = &runtime_v2.Program{}
	program.Bytecode = l.bytecode
	program.Arguments = l.arguments
	program.Constants = l.constants
	program.Locations = l.locations

	return program, err
}

func (c *Compiler) MustCompile(query string) *runtime_v2.Program {
	program, err := c.Compile(query)

	if err != nil {
		panic(err)
	}

	return program
}
