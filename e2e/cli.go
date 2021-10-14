package main

import (
	"bufio"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"os/signal"
	"path/filepath"
	rt "runtime"
	"strings"
	"time"

	"github.com/rs/zerolog"

	"github.com/MontFerret/ferret"
	"github.com/MontFerret/ferret/pkg/drivers/cdp"
	"github.com/MontFerret/ferret/pkg/drivers/http"
	"github.com/MontFerret/ferret/pkg/runtime"
	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/logging"
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

	dryRun = flag.Bool(
		"dry-run",
		false,
		"compiles a given query, but does not execute",
	)

	logLevel = flag.String(
		"log-level",
		logging.ErrorLevel.String(),
		"log level",
	)
)

var logger zerolog.Logger

func main() {
	var params Params

	flag.Var(
		&params,
		"param",
		`query parameter (--param=foo:\"bar\", --param=id:1)`,
	)

	flag.Parse()

	console := zerolog.ConsoleWriter{
		Out:        os.Stderr,
		TimeFormat: "15:04:05.999",
	}
	logger = zerolog.New(console).
		Level(zerolog.Level(logging.MustParseLevel(*logLevel))).
		With().
		Timestamp().
		Logger()

	stat, _ := os.Stdin.Stat()

	var query string
	var files []string

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
		files = flag.Args()
	} else {
		fmt.Println(flag.NArg())
		fmt.Println("File or input stream are required")
		os.Exit(1)
	}

	p, err := params.ToMap()

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	engine := ferret.New()
	_ = engine.Drivers().Register(http.NewDriver())
	_ = engine.Drivers().Register(cdp.NewDriver(cdp.WithAddress(*conn)))

	opts := []runtime.Option{
		runtime.WithParams(p),
		runtime.WithLog(console),
		runtime.WithLogLevel(logging.MustParseLevel(*logLevel)),
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	signal.Notify(c, os.Kill)

	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		for {
			<-c
			cancel()
		}
	}()

	if query != "" {
		err = runQuery(ctx, engine, opts, query)
	} else {
		err = execFiles(ctx, engine, opts, files)
	}

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func execFiles(ctx context.Context, engine *ferret.Instance, opts []runtime.Option, files []string) error {
	errList := make([]error, 0, len(files))

	for _, path := range files {
		log := logger.With().Str("path", path).Logger()
		log.Debug().Msg("checking path...")

		info, err := os.Stat(path)

		if err != nil {
			log.Debug().Err(err).Msg("failed to get path info")

			errList = append(errList, err)
			continue
		}

		if info.IsDir() {
			log.Debug().Msg("path points to a directory. retrieving list of files...")

			fileInfos, err := ioutil.ReadDir(path)

			if err != nil {
				log.Debug().Err(err).Msg("failed to retrieve list of files")

				errList = append(errList, err)
				continue
			}

			log.Debug().Int("size", len(fileInfos)).Msg("retrieved list of files. starting to iterate...")

			dirFiles := make([]string, 0, len(fileInfos))

			for _, info := range fileInfos {
				if filepath.Ext(info.Name()) == ".fql" {
					dirFiles = append(dirFiles, filepath.Join(path, info.Name()))
				}
			}

			if len(dirFiles) > 0 {
				if err := execFiles(ctx, engine, opts, dirFiles); err != nil {
					log.Debug().Err(err).Msg("failed to execute files")

					errList = append(errList, err)
				} else {
					log.Debug().Int("size", len(fileInfos)).Err(err).Msg("successfully executed files")
				}
			} else {
				log.Debug().Int("size", len(fileInfos)).Err(err).Msg("no FQL files found")
			}

			continue
		}

		log.Debug().Msg("path points to a file. starting to read content")

		out, err := ioutil.ReadFile(path)

		if err != nil {
			log.Debug().Err(err).Msg("failed to read content")

			errList = append(errList, err)
			continue
		}

		log.Debug().Msg("successfully read file")
		log.Debug().Msg("executing file...")
		err = runQuery(ctx, engine, opts, string(out))

		if err != nil {
			log.Debug().Err(err).Msg("failed to execute file")

			errList = append(errList, err)
			continue
		}

		log.Debug().Msg("successfully executed file")
	}

	if len(errList) > 0 {
		if len(errList) == len(files) {
			logger.Debug().Errs("errors", errList).Msg("failed to execute file(s)")
		} else {
			logger.Debug().Errs("errors", errList).Msg("executed with errors")
		}

		return core.Errors(errList...)
	}

	return nil
}

func runQuery(ctx context.Context, engine *ferret.Instance, opts []runtime.Option, query string) error {
	if !(*dryRun) {
		return execQuery(ctx, engine, opts, query)
	}

	return analyzeQuery(engine, query)
}

func execQuery(ctx context.Context, engine *ferret.Instance, opts []runtime.Option, query string) error {
	out, err := engine.Exec(ctx, query, opts...)

	if err != nil {
		return err
	}

	fmt.Println(string(out))

	return nil
}

func analyzeQuery(engine *ferret.Instance, query string) error {
	memBefore := &rt.MemStats{}
	rt.ReadMemStats(memBefore)

	timeBefore := time.Now()
	out, err := engine.Compile(query)

	if err != nil {
		return err
	}

	timeAfter := time.Since(timeBefore)
	memAfter := &rt.MemStats{}
	rt.ReadMemStats(memAfter)

	fmt.Println(out.Source())
	fmt.Println(fmt.Sprintf(`Time: %s`, timeAfter))
	fmt.Println(fmt.Sprintf(`Memory before: %s`, byteCountDecimal(memBefore.Alloc)))
	fmt.Println(fmt.Sprintf(`Memory after: %s`, byteCountDecimal(memAfter.Alloc)))

	return nil
}

func byteCountDecimal(b uint64) string {
	const unit = 1000

	if b < unit {
		return fmt.Sprintf("%d B", b)
	}

	div, exp := int64(unit), 0
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(b)/float64(div), "kMGTPE"[exp])
}
