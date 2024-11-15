package compiler

import (
	"errors"

	goruntime "runtime"

	"github.com/MontFerret/ferret/pkg/parser"
	"github.com/MontFerret/ferret/pkg/runtime"
	"github.com/MontFerret/ferret/pkg/stdlib"
)

type Compiler struct {
	*NamespaceContainer
}

func New(setters ...Option) *Compiler {
	c := &Compiler{}
	c.NamespaceContainer = NewRootNamespace()

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

func (c *Compiler) Compile(query string) (program *runtime.Program, err error) {
	if query == "" {
		return nil, ErrEmptyQuery
	}

	defer func() {
		if r := recover(); r != nil {
			buf := make([]byte, 1024)
			n := goruntime.Stack(buf, false)
			stackTrace := string(buf[:n])

			// find out exactly what the error was and set err
			// Find out exactly what the error was and set err
			switch x := r.(type) {
			case string:
				err = errors.New(x + "\n" + stackTrace)
			case error:
				err = errors.New(x.Error() + "\n" + stackTrace)
			default:
				err = errors.New("unknown panic\n" + stackTrace)
			}

			program = nil
		}
	}()

	p := parser.New(query)
	p.AddErrorListener(newErrorListener())

	l := newVisitor(query)

	p.Visit(l)

	if l.err != nil {
		return nil, l.err
	}

	program = &runtime.Program{}
	program.Bytecode = l.emitter.instructions
	program.Constants = l.symbols.constants
	program.CatchTable = l.catchTable
	program.Registers = int(l.registers.nextRegister)

	return program, err
}

func (c *Compiler) MustCompile(query string) *runtime.Program {
	program, err := c.Compile(query)

	if err != nil {
		panic(err)
	}

	return program
}
