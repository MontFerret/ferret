package app

import (
	"context"
	"fmt"
	"github.com/MontFerret/ferret/pkg/compiler"
	"github.com/MontFerret/ferret/pkg/runtime"
	"io/ioutil"
	"os"
)

func Exec(pathToFile, cdpConn string) {
	query, err := ioutil.ReadFile(pathToFile)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
		return
	}

	ferret := compiler.New()

	prog, err := ferret.Compile(string(query))

	if err != nil {
		fmt.Println("Failed to compile the query")
		fmt.Println(err)
		os.Exit(1)
		return
	}

	timer := NewTimer()
	timer.Start()

	out, err := prog.Run(context.Background(), runtime.WithBrowser(cdpConn))

	timer.Stop()

	fmt.Println(timer.Print())

	if err != nil {
		fmt.Println("Failed to execute the query")
		fmt.Println(err)
		os.Exit(1)
		return
	}

	fmt.Println(string(out))
}
