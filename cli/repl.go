package cli

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/MontFerret/ferret/pkg/compiler"
	"github.com/MontFerret/ferret/pkg/runtime"
	"github.com/MontFerret/ferret/pkg/runtime/logging"
	"github.com/chzyer/readline"
)

func Repl(version string, opts Options) {
	ferret := compiler.New()

	fmt.Printf("Welcome to Ferret REPL %s\n", version)
	fmt.Println("Please use `exit` or `Ctrl-D` to exit this program.")

	rl, err := readline.NewEx(&readline.Config{
		Prompt:          "> ",
		InterruptPrompt: "^C",
		EOFPrompt:       "exit",
		AutoComplete:    NewAutoCompleter(ferret.RegisteredFunctions()),
	})

	if err != nil {
		panic(err)
	}

	defer rl.Close()

	var commands []string
	var multiline bool

	var timer *Timer

	if opts.ShowTime {
		timer = NewTimer()
	}

	l := NewLogger()

	ctx, err := opts.WithContext(context.Background())

	if err != nil {
		fmt.Println("Failed to register HTML drivers")
		fmt.Println(err)
		os.Exit(1)
		return
	}

	ctx, cancel := context.WithCancel(ctx)
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGHUP)

	exit := func() {
		cancel()
		l.Close()
	}

	go func() {
		for {
			<-c
			exit()
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

		if query == "exit" {
			exit()
			os.Exit(0)
			return
		}

		program, err := ferret.Compile(query)

		if err != nil {
			fmt.Println("Failed to parse the query")
			fmt.Println(err)
			continue
		}

		if opts.ShowTime {
			timer.Start()
		}

		out, err := program.Run(
			ctx,
			runtime.WithLog(l),
			runtime.WithLogLevel(logging.DebugLevel),
			runtime.WithParams(opts.Params),
		)

		if err != nil {
			fmt.Println("Failed to execute the query")
			fmt.Println(err)
			continue
		}

		fmt.Println(string(out))

		if opts.ShowTime {
			timer.Stop()
			fmt.Println(timer.Print())
		}
	}
}
