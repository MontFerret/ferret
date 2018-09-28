package cli

import (
	"context"
	"fmt"
	"github.com/MontFerret/ferret/pkg/compiler"
	"github.com/MontFerret/ferret/pkg/runtime"
	"github.com/MontFerret/ferret/pkg/runtime/logging"
	"github.com/chzyer/readline"
	"os"
	"os/signal"
	"strings"
	"syscall"
)

func Repl(version string, opts Options) {
	ferret := compiler.New()

	fmt.Printf("Welcome to Ferret REPL %s\n", version)
	fmt.Println("Please use `exit` or `Ctrl-D` to exit this program.")

	rl, err := readline.NewEx(&readline.Config{
		Prompt:          "> ",
		InterruptPrompt: "^C",
		EOFPrompt:       "exit",
	})

	if err != nil {
		panic(err)
	}

	defer rl.Close()

	var commands []string
	var multiline bool

	timer := NewTimer()

	l := NewLogger()

	ctx, cancel := context.WithCancel(context.Background())
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGHUP)

	go func() {
		for {
			<-c
			cancel()
			l.Close()
		}
	}()

	for {
		line, err := rl.Readline()

		if err != nil {
			break
		}

		line = strings.TrimSpace(line)

		if len(line) == 0 {
			continue
		}

		if strings.HasPrefix(line, "%") {
			line = line[1:]

			multiline = !multiline
		}

		if multiline {
			commands = append(commands, line)
			continue
		}

		commands = append(commands, line)
		query := strings.TrimSpace(strings.Join(commands, "\n"))
		commands = make([]string, 0, 10)

		if query == "" {
			continue
		}

		program, err := ferret.Compile(query)

		if err != nil {
			fmt.Println("Failed to parse the query")
			fmt.Println(err)
			continue
		}

		timer.Start()

		out, err := program.Run(
			ctx,
			runtime.WithBrowser(opts.Cdp),
			runtime.WithLog(l),
			runtime.WithLogLevel(logging.DebugLevel),
		)

		timer.Stop()
		fmt.Println(timer.Print())

		if err != nil {
			fmt.Println("Failed to execute the query")
			fmt.Println(err)
			continue
		}

		fmt.Println(string(out))
	}
}
