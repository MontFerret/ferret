package sdk_test

import (
	"context"
	"fmt"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
	"github.com/MontFerret/ferret/v2/pkg/sdk"
)

func ExampleDecodeValue() {
	type options struct {
		Temperature *float64 `ferret:"temperature"`
		Name        string   `ferret:"name"`
	}

	input := runtime.NewObjectWith(map[string]runtime.Value{
		"NAME":        runtime.NewString("gpt"),
		"temperature": runtime.NewFloat(0),
	})

	decoded, err := sdk.DecodeValue[options](
		context.Background(),
		input,
		sdk.RequireType(runtime.TypeMap),
		sdk.OnlyFields("name", "temperature"),
		sdk.DisallowUnknownFields(),
		sdk.DisallowNoneValues(),
	)
	if err != nil {
		panic(err)
	}

	fmt.Println(decoded.Name, *decoded.Temperature)
	// Output: gpt 0
}
