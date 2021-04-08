package main

import (
	"bufio"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/MontFerret/ferret"
	"github.com/MontFerret/ferret/pkg/drivers"
	"github.com/MontFerret/ferret/pkg/drivers/cdp"
	"github.com/MontFerret/ferret/pkg/drivers/http"
	"github.com/MontFerret/ferret/pkg/runtime"
	"github.com/MontFerret/ferret/pkg/runtime/core"
)

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
		pair := strings.SplitN(entry, ":", 2)

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
	conn = flag.String(
		"cdp",
		"",
		"set CDP address",
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

	var query string

	stat, _ := os.Stdin.Stat()

	if (stat.Mode() & os.ModeCharDevice) == 0 {
		// check whether the app is getting a query via standard input
		std := bufio.NewReader(os.Stdin)

		b, err := ioutil.ReadAll(std)

		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		query = string(b)
	} else if flag.NArg() > 0 {
		// backward compatibility
		content, err := ioutil.ReadFile(flag.Arg(0))

		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		query = string(content)
	} else {
		fmt.Println(flag.NArg())
		fmt.Println("Missed file")
		os.Exit(1)
	}

	p, err := params.ToMap()

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if err := execFile(query, p); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func execFile(query string, params map[string]interface{}) error {
	ctx := drivers.WithContext(
		context.Background(),
		http.NewDriver(),
		drivers.AsDefault(),
	)

	ctx = drivers.WithContext(
		ctx,
		cdp.NewDriver(cdp.WithAddress(*conn)),
	)

	i := ferret.New()
	out, err := i.Exec(ctx, query, runtime.WithParams(params))

	if err != nil {
		return err
	}

	fmt.Println(string(out))

	return nil
}
