package main

import (
	"context"
	"fmt"
	"os"
	"strings"

	ferret "github.com/MontFerret/ferret/v2"
	"github.com/MontFerret/ferret/v2/pkg/module"
	"github.com/MontFerret/ferret/v2/pkg/runtime"
	"github.com/MontFerret/ferret/v2/pkg/sdk"
	"github.com/MontFerret/ferret/v2/pkg/source"
)

type formatOptions struct {
	Prefix string `ferret:"prefix"`
}

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run() error {
	engine, err := ferret.New(ferret.WithModules(newExampleModule()))
	if err != nil {
		return err
	}
	defer func() { _ = engine.Close() }()

	output, err := engine.Run(
		context.Background(),
		source.NewAnonymous(`RETURN EXAMPLE::FORMAT("ferret", { prefix: "hello " })`),
	)
	if err != nil {
		return err
	}

	fmt.Println(string(output.Content))
	return nil
}

func newExampleModule() module.Module {
	return sdk.NewModule("example", func(bootstrap module.Bootstrap) error {
		return sdk.RegisterFunctions(
			bootstrap.Host().Library().Namespace("EXAMPLE"),
			sdk.Func("UPPER", sdk.Bind1(upper)),
			sdk.Func("FORMAT", runtime.Function(format)),
		)
	})
}

func upper(_ context.Context, value runtime.String) (runtime.String, error) {
	return runtime.NewString(strings.ToUpper(value.String())), nil
}

func format(ctx context.Context, args ...runtime.Value) (runtime.Value, error) {
	if err := runtime.ValidateArgs(args, 1, 2); err != nil {
		return runtime.None, err
	}

	value, err := runtime.CastArgAt[runtime.String](args, 0)
	if err != nil {
		return runtime.None, err
	}

	options, err := sdk.DecodeArgOr(
		ctx,
		args,
		1,
		formatOptions{},
		sdk.DisallowUnknownFields(),
	)
	if err != nil {
		return runtime.None, err
	}

	return runtime.NewString(options.Prefix + strings.ToUpper(value.String())), nil
}
