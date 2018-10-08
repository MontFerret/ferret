package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/MontFerret/ferret/cli"
	"github.com/MontFerret/ferret/cli/browser"
	"github.com/MontFerret/ferret/pkg/runtime/core"
	"io/ioutil"
	"os"
	"strings"
)

var Version string

type Params []string

func (p *Params) String() string {
	return "[" + strings.Join(*p, ",") + "]"
}

func (p *Params) Set(value string) error {
	*p = append(*p, value)
	return nil
}

func (p *Params) ToMap() (map[string]interface{}, error) {
	res := make(map[string]interface{})

	for _, entry := range *p {
		pair := strings.Split(entry, ":")

		if len(pair) < 2 {
			return nil, core.Error(core.ErrInvalidArgument, entry)
		}

		var value interface{}
		key := pair[0]

		err := json.Unmarshal([]byte(pair[1]), &value)

		if err != nil {
			fmt.Println(pair[1])
			return nil, err
		}

		res[key] = value
	}

	return res, nil
}

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
		"http://0.0.0.0:9222",
		"set CDP address",
	)

	launchBrowser = flag.Bool(
		"cdp-launch",
		false,
		"launch Chrome",
	)

	proxyAddress = flag.String(
		"proxy",
		"",
		"address of proxy server to use (only applicable for static pages)",
	)
)

func main() {
	var params Params

	flag.Var(
		&params,
		"param",
		`query parameter (--param=foo:\"bar\", --param=id:1)`,
	)

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

	p, err := params.ToMap()

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	opts := cli.Options{
		Cdp:    cdpConn,
		Params: p,
		Proxy:  *proxyAddress,
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

		cli.Exec(string(b), opts)
		return
	}

	// filename was passed
	if flag.NArg() > 0 {
		cli.ExecFile(flag.Arg(0), opts)
		return
	}

	// nothing was passed, run REPL
	cli.Repl(Version, opts)
}
