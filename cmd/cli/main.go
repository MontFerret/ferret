package main

import (
	"flag"
	"fmt"
	"github.com/MontFerret/ferret/cmd/cli/app"
	"github.com/MontFerret/ferret/pkg/browser"
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
		"Chrome DevTools Protocol address",
	)

	launchBrowser = flag.Bool(
		"cdp-launch",
		false,
		"launch Chrome",
	)

	noUserData = flag.Bool(
		"no-user-data",
		false,
		"do not create a separate location for browser sessions",
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

	// no files to execute
	// run REPL
	if flag.NArg() == 0 {
		app.Repl(Version, cdpConn)
	} else {
		app.Exec(flag.Arg(0), cdpConn)
	}
}
