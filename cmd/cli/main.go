package main

import (
	"bufio"
	"flag"
	"fmt"
	"github.com/MontFerret/ferret/cmd/cli/app"
	"github.com/MontFerret/ferret/pkg/browser"
	"io/ioutil"
	"os"
)

var Version string

var (
	help = flag.Bool(
		"help",
		false,
		"show this list",
	)

	version = flag.Bool(
		"version",
		false,
		"show REPL version",
	)

	conn = flag.String(
		"cdp",
		"http://127.0.0.1:9222",
		"set CDP address",
	)

	launchBrowser = flag.Bool(
		"cdp-launch",
		false,
		"launch Chrome",
	)
)

func main() {
	flag.Parse()

	if *help {
		flag.PrintDefaults()
		os.Exit(0)
		return
	}

	if *version {
		fmt.Println(Version)
		os.Exit(0)
		return
	}

	cdpConn := *conn

	if cdpConn == "" && *launchBrowser {
		opts := make([]browser.Option, 0, 2)

		//if *noUserData {
		//	opts = append(opts, browser.WithoutUserDataDir())
		//}

		// TODO: Make it optional.
		opts = append(opts, browser.WithoutUserDataDir())

		// we need to launch Chrome instance
		b, err := browser.Launch(opts...)

		if err != nil {
			fmt.Println(fmt.Sprintf("Failed to launch browser:"))
			fmt.Println(err)
			os.Exit(1)
		}

		cdpConn = b.DebuggingAddress()

		defer b.Close()
	}

	opts := app.Options{
		Cdp: cdpConn,
	}

	stat, _ := os.Stdin.Stat()

	if (stat.Mode() & os.ModeCharDevice) == 0 {
		// check whether the app is getting a query via standard input
		std := bufio.NewReader(os.Stdin)

		b, err := ioutil.ReadAll(std)

		if err != nil {
			fmt.Println(err)
			os.Exit(1)
			return
		}

		app.Exec(string(b), opts)
		return
	}

	// filename was passed
	if flag.NArg() > 0 {
		app.ExecFile(flag.Arg(0), opts)
		return
	}

	// nothing was passed, run REPL
	app.Repl(Version, opts)
}
