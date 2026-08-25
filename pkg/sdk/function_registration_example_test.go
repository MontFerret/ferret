package sdk_test

import (
	"context"
	"fmt"

	apisource "github.com/MontFerret/api/source"
	ferret "github.com/MontFerret/ferret/v2"
	"github.com/MontFerret/ferret/v2/pkg/module"
	"github.com/MontFerret/ferret/v2/pkg/runtime"
	"github.com/MontFerret/ferret/v2/pkg/sdk"
)

func ExampleRegisterFunctions() {
	definitionCount := 0
	logicalNameCount := 0
	duplicateRejected := false

	mod := sdk.NewModule("overloads", func(bootstrap module.Bootstrap) error {
		library := bootstrap.Host().Library()
		before, err := library.Build()
		if err != nil {
			return err
		}
		ns := library.Namespace("example")
		one := runtime.Function1(func(context.Context, runtime.Value) (runtime.Value, error) {
			return runtime.NewString("fixed-one"), nil
		})

		if err := sdk.RegisterFunctions(ns,
			sdk.Func("pick", one),
			sdk.Func("pick", runtime.Function2(func(context.Context, runtime.Value, runtime.Value) (runtime.Value, error) {
				return runtime.NewString("fixed-two"), nil
			})),
			sdk.Func("pick", runtime.Function(func(_ context.Context, args ...runtime.Value) (runtime.Value, error) {
				return runtime.NewString(fmt.Sprintf("variadic-%d", len(args))), nil
			})),
		); err != nil {
			return err
		}

		duplicateRejected = sdk.RegisterFunctions(ns, sdk.Func("PICK", one)) != nil
		functions, err := library.Build()
		if err != nil {
			return err
		}
		definitionCount = functions.Size() - before.Size()
		logicalNameCount = len(functions.List()) - len(before.List())

		return nil
	})

	engine, err := ferret.New(ferret.WithModules(mod))
	if err != nil {
		panic(err)
	}
	defer func() { _ = engine.Close() }()

	for _, query := range []string{
		`return example::pick(1)`,
		`return example::pick(1, 2)`,
		`return example::pick(1, 2, 3, 4, 5)`,
	} {
		output, runErr := engine.Run(context.Background(), apisource.File{
			Name:    "anonymous",
			Content: query,
		})
		if runErr != nil {
			panic(runErr)
		}
		fmt.Println(string(output.Content))
	}

	fmt.Println("duplicate rejected:", duplicateRejected)
	fmt.Printf("definitions=%d logical names=%d\n", definitionCount, logicalNameCount)

	// Output:
	// "fixed-one"
	// "fixed-two"
	// "variadic-5"
	// duplicate rejected: true
	// definitions=3 logical names=1
}
