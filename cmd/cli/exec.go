package cli

import (
	"context"
	"fmt"
	"github.com/MontFerret/ferret/pkg/compiler"
	"github.com/MontFerret/ferret/pkg/runtime"
	"io/ioutil"
	"os"
)

func ExecFile(pathToFile string, opts Options) {
	query, err := ioutil.ReadFile(pathToFile)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
		return
	}

	Exec(string(query), opts)
}

func Exec(query string, opts Options) {
	ferret := compiler.New()

	prog, err := ferret.Compile(query)

	if err != nil {
		fmt.Println("Failed to compile the query")
		fmt.Println(err)
		os.Exit(1)
		return
	}

	out, err := prog.Run(
		context.Background(),
		runtime.WithBrowser(opts.Cdp),
	)

	if err != nil {
		fmt.Println("Failed to execute the query")
		fmt.Println(err)
		os.Exit(1)
		return
	}

	fmt.Println(string(out))
}
