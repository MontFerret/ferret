package base

import (
	"context"
	j "encoding/json"
	"errors"

	"github.com/MontFerret/ferret/v2/pkg/bytecode"
	ferretencoding "github.com/MontFerret/ferret/v2/pkg/encoding"
	encodingjson "github.com/MontFerret/ferret/v2/pkg/encoding/json"
	encodingmsgpack "github.com/MontFerret/ferret/v2/pkg/encoding/msgpack"
	"github.com/MontFerret/ferret/v2/pkg/file"

	"github.com/MontFerret/ferret/v2/pkg/compiler"
	"github.com/MontFerret/ferret/v2/pkg/runtime"
	"github.com/MontFerret/ferret/v2/pkg/vm"
)

type contextKey string

const testSaltKey contextKey = "test-salt"

func Compile(expression string) (*bytecode.Program, error) {
	c := compiler.New(compiler.WithOptimizationLevel(compiler.O0))

	return c.Compile(file.NewSource("", expression))
}

func newTestContext() context.Context {
	type Salt struct{}

	ctx := context.WithValue(context.Background(), testSaltKey, &Salt{})

	return ferretencoding.WithRegistry(ctx, ferretencoding.NewRegistry(encodingjson.Default, encodingmsgpack.Default))
}

func materializeJSONResult(out *vm.Result) ([]byte, error) {
	data, materializeErr := vm.Materialize[[]byte](out, func(value runtime.Value) (vm.Materialized[[]byte], error) {
		enc := encodingjson.Default.EncodeWith().PreHook(func(value runtime.Value) error {
			out.AdoptValue(value)
			return nil
		}).Encoder()

		data, err := enc.Encode(value)
		if err != nil {
			return vm.Materialized[[]byte]{}, err
		}

		return vm.Materialized[[]byte]{Value: data}, nil
	})

	return data, errors.Join(materializeErr, out.Close())
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

	ctx := newTestContext()

	out, err := instance.Run(ctx, env)

	if err != nil {
		return nil, err
	}

	return materializeJSONResult(out)
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
