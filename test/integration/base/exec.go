package base

import (
	"context"
	j "encoding/json"

	"github.com/MontFerret/ferret/v2/pkg/bytecode"
	ferretencoding "github.com/MontFerret/ferret/v2/pkg/encoding"
	encodingjson "github.com/MontFerret/ferret/v2/pkg/encoding/json"
	"github.com/MontFerret/ferret/v2/pkg/file"

	"github.com/MontFerret/ferret/v2/pkg/compiler"
	"github.com/MontFerret/ferret/v2/pkg/vm"
)

type contextKey string

const testSaltKey contextKey = "test-salt"

func Compile(expression string) (*bytecode.Program, error) {
	c := compiler.New(compiler.WithOptimizationLevel(compiler.O0))

	return c.Compile(file.NewSource("", expression))
}

func Run(p *bytecode.Program, opts ...vm.EnvironmentOption) ([]byte, error) {
	instance, err := vm.New(p)
	if err != nil {
		return nil, err
	}

	env, err := vm.NewEnvironment(opts)
	if err != nil {
		return nil, err
	}

	type Salt struct{}

	ctx := context.WithValue(context.Background(), testSaltKey, &Salt{})
	ctx = ferretencoding.WithRegistry(ctx, ferretencoding.NewRegistry())

	out, err := instance.Run(ctx, env)

	if err != nil {
		return nil, err
	}

	return encodingjson.Default.Encode(out)
}

func Exec(p *bytecode.Program, raw bool, opts ...vm.EnvironmentOption) (any, error) {
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
