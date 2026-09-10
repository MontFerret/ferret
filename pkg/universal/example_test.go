package universal_test

import (
	"context"
	"fmt"

	"github.com/MontFerret/api"

	"github.com/MontFerret/ferret/v2"
	"github.com/MontFerret/ferret/v2/pkg/universal"
)

func ExampleWrap() {
	native, err := ferret.New()
	if err != nil {
		panic(err)
	}

	defer native.Close()

	var portable api.Runtime = universal.Wrap(native)
	defer portable.Close()

	output, err := portable.Run(context.Background(), api.NewAnonymousSource("RETURN @value + 1"), api.WithParam("value", 41))
	if err != nil {
		panic(err)
	}

	fmt.Println(string(output.Content))
	// Output: 42
}

func ExampleNew() {
	portable, err := universal.New(ferret.WithParam("value", 41))
	if err != nil {
		panic(err)
	}

	defer portable.Close()

	output, err := portable.Run(context.Background(), api.NewAnonymousSource("RETURN @value + 1"))
	if err != nil {
		panic(err)
	}

	fmt.Println(string(output.Content))
	// Output: 42
}
