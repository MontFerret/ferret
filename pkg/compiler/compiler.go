package compiler

import (
	"errors"
	goruntime "runtime"

	"github.com/MontFerret/ferret/pkg/compiler/internal"

	"github.com/MontFerret/ferret/pkg/stdlib"
	"github.com/MontFerret/ferret/pkg/vm"

	"github.com/MontFerret/ferret/pkg/parser"
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

func (c *Compiler) Compile(query string) (program *vm.Program, err error) {
	if query == "" {
		return nil, internal.ErrEmptyQuery
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

	l := internal.NewVisitor(query)

	p.Visit(l)

	if l.Err != nil {
		return nil, l.Err
	}

	program = &vm.Program{}
	program.Bytecode = l.Emitter.Bytecode()
	program.Constants = l.Symbols.Constants()
	program.CatchTable = l.CatchTable.All()
	program.Registers = l.Registers.Size()
	program.Params = l.Symbols.Params()

	return program, err
}

func (c *Compiler) MustCompile(query string) *vm.Program {
	program, err := c.Compile(query)

	if err != nil {
		panic(err)
	}

	return program
}
