package main

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
	"os/signal"
	"path/filepath"
	rt "runtime"
	"runtime/pprof"
	"strings"
	"time"

	"github.com/MontFerret/ferret/pkg/compiler"
	"github.com/MontFerret/ferret/pkg/file"
	"github.com/MontFerret/ferret/pkg/runtime"

	"github.com/rs/zerolog"

	"github.com/MontFerret/ferret"
	//"github.com/MontFerret/ferret/pkg/drivers/cdp"
	//"github.com/MontFerret/ferret/pkg/drivers/http"
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
		fmt.Fprintf(writer, "Time: %s\n", timer.end.Sub(timer.start))
	}

	stats, found := p.allocs[label]

	if found {
		fmt.Fprintf(writer, "Alloc: %s\n", byteCountDecimal(stats.Alloc))
		fmt.Fprintf(writer, "Frees: %s\n", byteCountDecimal(stats.Frees))
		fmt.Fprintf(writer, "Total Alloc: %s\n", byteCountDecimal(stats.TotalAlloc))
		fmt.Fprintf(writer, "Heap Alloc: %s\n", byteCountDecimal(stats.HeapAlloc))
		fmt.Fprintf(writer, "Heap Sys: %s\n", byteCountDecimal(stats.HeapSys))
		fmt.Fprintf(writer, "Heap Idle: %s\n", byteCountDecimal(stats.HeapIdle))
		fmt.Fprintf(writer, "Heap In Use: %s\n", byteCountDecimal(stats.HeapInuse))
		fmt.Fprintf(writer, "Heap Released: %s\n", byteCountDecimal(stats.HeapReleased))
		fmt.Fprintf(writer, "Heap Objects: %d\n", stats.HeapObjects)
	}

	//cpu, found := p.cpus[label]
	//
	//if found {
	//	fmt.Fprintln(writer, cpu.String())
	//}

	if writer.Len() > 0 {
		fmt.Printf("%s:\n", label)
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

func (p *Params) ToMap() (map[string]runtime.Value, error) {
	res := make(map[string]runtime.Value)

	for _, entry := range *p {
		pair := strings.SplitN(entry, ":", 2)

		if len(pair) < 2 {
			return nil, runtime.Error(runtime.ErrInvalidArgument, entry)
		}

		var value interface{}
		key := pair[0]

		err := json.Unmarshal([]byte(pair[1]), &value)

		if err != nil {
			fmt.Println(pair[1])
			return nil, err
		}

		res[key] = ferret.ToValue(value)
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
		runtime.ErrorLevel.String(),
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
		Level(zerolog.Level(runtime.MustParseLogLevel(*logLevel))).
		With().
		Timestamp().
		Logger()

	stat, _ := os.Stdin.Stat()

	var query string
	var files []string

	if (stat.Mode() & os.ModeCharDevice) == 0 {
		// check whether the app is getting a query via standard input
		std := bufio.NewReader(os.Stdin)

		b, err := io.ReadAll(std)

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

	var err error

	p, err := params.ToMap()

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	engine := ferret.New()
	//_ = engine.Drivers().Register(http.NewDriver())
	//_ = engine.Drivers().Register(cdp.NewDriver(cdp.WithAddress(*conn)))

	opts := []ferret.SessionOption{
		ferret.WithParams(p),
		ferret.WithLog(console),
		ferret.WithLogLevel(ferret.MustParseLogLevel(*logLevel)),
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
		err = runQuery(ctx, engine, opts, file.NewSource("stdin", query))
	} else {
		err = execFiles(ctx, engine, opts, files)
	}

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func execFiles(ctx context.Context, engine *ferret.Engine, opts []ferret.SessionOption, files []string) error {
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

			fileInfos, err := os.ReadDir(path)

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

		src, err := file.Read(path)

		if err != nil {
			log.Debug().Err(err).Msg("failed to read content")

			errList = append(errList, err)
			continue
		}

		log.Debug().Msg("successfully read file")
		log.Debug().Msg("executing file...")
		err = runQuery(ctx, engine, opts, src)

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

		return runtime.Errors(errList...)
	}

	return nil
}

func runQuery(ctx context.Context, engine *ferret.Engine, opts []ferret.SessionOption, query *file.Source) error {
	if !(*dryRun) {
		return execQuery(ctx, engine, opts, query)
	}

	return analyzeQuery(query)
}

func execQuery(ctx context.Context, engine *ferret.Engine, opts []ferret.SessionOption, query *file.Source) error {
	plan, err := engine.Compile(query)

	if err != nil {
		return err
	}

	sess := plan.NewSession(opts...)
	defer sess.Close()

	beforeExec := "Before Execution"
	exec := "Execution"
	afterExec := "After Execution"

	// TODO: We need add an option to separate compilation and execution phases
	prof := NewProfiler()

	if *profiler {
		prof.Allocations(beforeExec)
		prof.StartCPU(exec)
		prof.StartTimer(exec)
	}

	res, err := sess.Run(ctx)

	if *profiler {
		prof.Allocations(afterExec)
		prof.StopTimer(exec)
		prof.StopCPU()

		prof.PrintAll()
	}

	if err == nil {
		defer res.Close()

		if size, e := printResult(ctx, res); e != nil {
			err = e
		} else if *profiler {
			fmt.Printf("Output size: %s\n", byteCountDecimal(size))
			fmt.Println("")
		}
	}

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	return nil
}

type ResultPrinter struct {
	out  io.Writer
	size uint64
}

func (r *ResultPrinter) Write(p []byte) (n int, err error) {
	r.size += uint64(len(p))
	return r.out.Write(p)
}

func printResult(ctx context.Context, res ferret.Result) (uint64, error) {
	defer res.Close()

	isScalar := ferret.IsScalar(res)
	printer := &ResultPrinter{
		out: os.Stdout,
	}

	if !isScalar {
		printer.Write([]byte("["))
	}

	var counter int

	err := ferret.ForEach(ctx, res, func(val runtime.Value) error {
		j, err := val.MarshalJSON()

		if err != nil {
			return err
		}

		if counter > 0 && !isScalar {
			printer.Write([]byte(","))
		}

		_, err = printer.Write(j)

		counter++

		return err
	})

	if err != nil {
		return printer.size, err
	}

	if !isScalar {
		printer.Write([]byte("]"))
	}
	os.Stdout.Write([]byte("\n"))

	return printer.size, nil
}

func analyzeQuery(query *file.Source) error {
	beforeCompilation := "Before Compilation"
	compilation := "Compilation"
	afterCompilation := "After Compilation"
	prof := NewProfiler()

	c := compiler.New()

	fullProf := *profiler

	if fullProf {
		prof.Allocations(beforeCompilation)
	}

	prof.StartTimer(compilation)

	c.MustCompile(query)

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
