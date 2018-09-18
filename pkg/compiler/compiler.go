package compiler

import (
	"github.com/MontFerret/ferret/pkg/parser"
	"github.com/MontFerret/ferret/pkg/runtime"
	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/stdlib"
	"github.com/pkg/errors"
)

type FqlCompiler struct {
	funcs map[string]core.Function
}

func New() *FqlCompiler {
	return &FqlCompiler{
		stdlib.NewLib(),
	}
}

func (c *FqlCompiler) RegisterFunction(name string, fun core.Function) error {
	_, exists := c.funcs[name]

	if exists {
		return errors.Errorf("function already exists: %s", name)
	}

	c.funcs[name] = fun

	return nil
}

func (c *FqlCompiler) Compile(query string) (program *runtime.Program, err error) {
	if query == "" {
		return nil, ErrEmptyQuery
	}

	p := parser.New(query)
	p.AddErrorListener(&errorListener{})

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

	l := newVisitor(c.funcs)

	res := p.Visit(l).(*result)

	if res.Ok() {
		program = res.Data().(*runtime.Program)
	} else {
		err = res.Error()
	}

	return program, err
}
