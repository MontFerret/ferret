package main

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"os/signal"
	"path/filepath"
	rt "runtime"
	"runtime/pprof"
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

type (
	Timer struct {
		start time.Time
		end   time.Time
	}

	Profiler struct {
		labels []string
		timers map[string]*Timer
		allocs map[string]*rt.MemStats
		cpus   map[string]*bytes.Buffer
		heaps  map[string]*bytes.Buffer
	}
)

func NewProfiler() *Profiler {
	return &Profiler{
		labels: make([]string, 0, 10),
		timers: make(map[string]*Timer),
		allocs: make(map[string]*rt.MemStats),
		cpus:   make(map[string]*bytes.Buffer),
		heaps:  make(map[string]*bytes.Buffer),
	}
}

func (p *Profiler) StartTimer(label string) {
	timer := &Timer{
		start: time.Now(),
	}

	p.timers[label] = timer
	p.addLabel(label)
}

func (p *Profiler) StopTimer(label string) {
	timer, found := p.timers[label]

	if !found {
		panic(fmt.Sprintf("Timer not found: %s", label))
	}

	timer.end = time.Now()
}

func (p *Profiler) HeapSnapshot(label string) {
	heap := &bytes.Buffer{}

	err := pprof.WriteHeapProfile(heap)

	if err != nil {
		panic(err)
	}

	p.heaps[label] = heap
	p.addLabel(label)
}

func (p *Profiler) Allocations(label string) {
	stats := &rt.MemStats{}

	rt.ReadMemStats(stats)

	p.allocs[label] = stats
	p.addLabel(label)
}

func (p *Profiler) StartCPU(label string) {
	b := &bytes.Buffer{}

	if err := pprof.StartCPUProfile(b); err != nil {
		panic(err)
	}

	p.cpus[label] = b
	p.addLabel(label)
}

func (p *Profiler) StopCPU() {
	pprof.StopCPUProfile()
}

func (p *Profiler) Print(label string) {
	writer := &bytes.Buffer{}

	timer, found := p.timers[label]

	if found {
		fmt.Fprintln(writer, fmt.Sprintf("Time: %s", timer.end.Sub(timer.start)))
	}

	stats, found := p.allocs[label]

	if found {
		fmt.Fprintln(writer, fmt.Sprintf("Alloc: %s", byteCountDecimal(stats.Alloc)))
		fmt.Fprintln(writer, fmt.Sprintf("Frees: %s", byteCountDecimal(stats.Frees)))
		fmt.Fprintln(writer, fmt.Sprintf("Total Alloc: %s", byteCountDecimal(stats.TotalAlloc)))
		fmt.Fprintln(writer, fmt.Sprintf("Heap Alloc: %s", byteCountDecimal(stats.HeapAlloc)))
		fmt.Fprintln(writer, fmt.Sprintf("Heap Sys: %s", byteCountDecimal(stats.HeapSys)))
		fmt.Fprintln(writer, fmt.Sprintf("Heap Idle: %s", byteCountDecimal(stats.HeapIdle)))
		fmt.Fprintln(writer, fmt.Sprintf("Heap In Use: %s", byteCountDecimal(stats.HeapInuse)))
		fmt.Fprintln(writer, fmt.Sprintf("Heap Released: %s", byteCountDecimal(stats.HeapReleased)))
		fmt.Fprintln(writer, fmt.Sprintf("Heap Objects: %d", stats.HeapObjects))
	}

	//cpu, found := p.cpus[label]
	//
	//if found {
	//	fmt.Fprintln(writer, cpu.String())
	//}

	if writer.Len() > 0 {
		fmt.Println(fmt.Sprintf("%s:", label))
		fmt.Println("-----")
		fmt.Println(writer.String())
	}
}

func (p *Profiler) PrintAll() {
	for _, label := range p.labels {
		p.Print(label)
	}
}

func (p *Profiler) addLabel(label string) {
	var found bool

	for _, l := range p.labels {
		if l == label {
			found = true
			break
		}
	}

	if !found {
		p.labels = append(p.labels, label)
	}
}

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

	profiler = flag.Bool(
		"profiler",
		false,
		"enables CPU and Memory profiler",
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
	beforeExec := "Before Execution"
	exec := "Execution"
	afterExec := "After Execution"

	prof := NewProfiler()

	if *profiler {
		prof.Allocations(beforeExec)
		prof.StartCPU(exec)
		prof.StartTimer(exec)
	}

	out, err := engine.Exec(ctx, query, opts...)

	if *profiler {
		prof.Allocations(afterExec)
		prof.StopTimer(exec)
		prof.StopCPU()

		prof.PrintAll()

		if out != nil {
			fmt.Println(fmt.Sprintf("Output size: %s", byteCountDecimal(uint64(len(out)))))
			fmt.Println("")
		}
	}

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Println(string(out))

	return nil
}

func analyzeQuery(engine *ferret.Instance, query string) error {
	beforeCompilation := "Before Compilation"
	compilation := "Compilation"
	afterCompilation := "After Compilation"
	prof := NewProfiler()

	fullProf := *profiler

	if fullProf {
		prof.Allocations(beforeCompilation)
	}

	prof.StartTimer(compilation)

	engine.MustCompile(query)

	prof.StopTimer(compilation)

	if fullProf {
		prof.Allocations(afterCompilation)
	}

	prof.PrintAll()

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
