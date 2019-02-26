package cli

import (
	"context"
	"fmt"
	"github.com/MontFerret/ferret/pkg/compiler"
	"github.com/MontFerret/ferret/pkg/runtime"
	"github.com/MontFerret/ferret/pkg/runtime/logging"
	"io/ioutil"
	"os"
	"os/signal"
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

	l := NewLogger()

	ctx, cancel := opts.WithContext(context.Background())

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	go func() {
		for {
			<-c
			cancel()
			l.Close()
		}
	}()

	var timer *Timer

	if opts.ShowTime {
		timer = NewTimer()
		timer.Start()
	}

	out, err := prog.Run(
		ctx,
		runtime.WithLog(l),
		runtime.WithLogLevel(logging.DebugLevel),
		runtime.WithParams(opts.Params),
	)

	if opts.ShowTime {
		timer.Stop()
	}

	if err != nil {
		fmt.Println("Failed to execute the query")
		fmt.Println(err)
		os.Exit(1)
		return
	}

	fmt.Println(string(out))

	if opts.ShowTime {
		fmt.Println(timer.Print())
	}
}
