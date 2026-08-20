package runtime_test

import (
	"context"
	"fmt"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

func ExampleAddable() {
	host := exampleAddableValue{label: "host"}

	left, _ := runtime.Add(context.Background(), host, runtime.NewInt(2))
	right, _ := runtime.Add(context.Background(), runtime.NewInt(2), host)

	fmt.Println(left)
	fmt.Println(right)
	// Output:
	// host+2
	// 2+host
}
