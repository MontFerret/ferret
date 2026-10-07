package uapi_test

import (
	"context"
	"errors"
	"fmt"

	"github.com/MontFerret/api"

	"github.com/MontFerret/ferret/v2"
	"github.com/MontFerret/ferret/v2/uapi"
)

func ExampleWrap() {
	native, err := ferret.New()
	if err != nil {
		panic(err)
	}

	defer native.Close()

	var portable api.Runtime = uapi.Wrap(native, "dev")
	defer portable.Close()

	handle, err := portable.Run(context.Background(), api.NewAnonymousSource("RETURN @value + 1"), api.WithParam("value", 41))
	if err != nil {
		panic(err)
	}

	defer handle.Close()
	output, err := handle.Collect(context.Background())
	if err != nil {
		panic(err)
	}

	fmt.Println(string(output.Data))
	// Output: 42
}

func ExampleNew() {
	portable, err := uapi.New("dev", ferret.WithParam("value", 41))
	if err != nil {
		panic(err)
	}

	defer portable.Close()

	handle, err := portable.Run(context.Background(), api.NewAnonymousSource("RETURN @value + 1"))
	if err != nil {
		panic(err)
	}

	defer handle.Close()
	output, err := handle.Collect(context.Background())
	if err != nil {
		panic(err)
	}

	fmt.Println(string(output.Data))
	// Output: 42
}

func ExampleRuntime_Run() {
	failure := errors.New("after-run hook failed")
	portable, err := uapi.New("dev", ferret.WithAfterRunHook(func(context.Context, error) error {
		return failure
	}))
	if err != nil {
		panic(err)
	}

	defer portable.Close()

	handle, err := portable.Run(context.Background(), api.NewAnonymousSource("RETURN 42"))
	if err != nil {
		panic(err)
	}

	defer handle.Close()
	output, err := handle.Collect(context.Background())
	if output != nil {
		fmt.Printf("output: %s\n", output.Data)
	}

	if err != nil {
		fmt.Printf("after-run failure: %t\n", errors.Is(err, failure))
	}

	// Output:
	// output: 42
	// after-run failure: true
}
