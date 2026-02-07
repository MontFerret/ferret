package base

import (
	"context"
	j "encoding/json"

	"github.com/MontFerret/ferret/pkg/file"

	"github.com/MontFerret/ferret/pkg/compiler"
	"github.com/MontFerret/ferret/pkg/vm"
)

func Compile(expression string) (*vm.Program, error) {
	c := compiler.New(compiler.WithOptimizationLevel(compiler.O0))

	return c.Compile(file.NewSource("", expression))
}

func Run(p *vm.Program, opts ...vm.EnvironmentOption) ([]byte, error) {
	instance := vm.New(p)
	env := vm.NewEnvironment(opts)

	type Salt struct{}

	out, err := instance.Run(context.WithValue(context.Background(), "test-salt", &Salt{}), env)

	if err != nil {
		return nil, err
	}

	return out.MarshalJSON()
}

func Exec(p *vm.Program, raw bool, opts ...vm.EnvironmentOption) (any, error) {
	out, err := Run(p, opts...)

	if err != nil {
		return 0, err
	}

	if raw {
		return string(out), nil
	}

	var res any

	err = j.Unmarshal(out, &res)

	if err != nil {
		return nil, err
	}

	return res, err
}
