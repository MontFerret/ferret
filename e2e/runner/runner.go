package runner

import (
	"context"
	"encoding/json"
	"github.com/MontFerret/ferret/pkg/compiler"
	"github.com/MontFerret/ferret/pkg/html"
	"github.com/MontFerret/ferret/pkg/html/cdp"
	"github.com/MontFerret/ferret/pkg/html/static"
	"github.com/MontFerret/ferret/pkg/runtime"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"
)

type (
	Settings struct {
		StaticServerAddress  string
		DynamicServerAddress string
		CDPAddress           string
		Dir                  string
	}

	Result struct {
		name     string
		duration time.Duration
		err      error
	}

	Summary struct {
		passed   int
		failed   int
		duration time.Duration
	}

	Runner struct {
		logger   zerolog.Logger
		settings Settings
	}
)

func New(logger zerolog.Logger, settings Settings) *Runner {
	return &Runner{
		logger,
		settings,
	}
}

func (r *Runner) Run() error {
	results, err := r.runQueries(r.settings.Dir)

	if err != nil {
		return err
	}

	sum := r.report(results)

	var event *zerolog.Event

	if sum.failed == 0 {
		event = r.logger.Info()
	} else {
		event = r.logger.Error()
	}

	event.
		Timestamp().
		Int("passed", sum.passed).
		Int("failed", sum.failed).
		Dur("time", sum.duration).
		Msg("Completed")

	if sum.failed > 0 {
		return errors.New("failed")
	}

	return nil
}

func (r *Runner) runQueries(dir string) ([]Result, error) {
	files, err := ioutil.ReadDir(dir)

	if err != nil {
		r.logger.Error().
			Timestamp().
			Err(err).
			Str("dir", dir).
			Msg("failed to read scripts directory")

		return nil, err
	}

	results := make([]Result, 0, len(files))

	c := compiler.New()

	if err := c.RegisterFunctions(Assertions()); err != nil {
		return nil, err
	}

	// read scripts
	for _, f := range files {
		fName := filepath.Join(dir, f.Name())
		b, err := ioutil.ReadFile(fName)

		if err != nil {
			results = append(results, Result{
				name: fName,
				err:  errors.Wrap(err, "failed to read script file"),
			})

			continue
		}

		results = append(results, r.runQuery(c, fName, string(b)))
	}

	return results, nil
}

func (r *Runner) runQuery(c *compiler.FqlCompiler, name, script string) Result {
	start := time.Now()

	p, err := c.Compile(script)

	if err != nil {
		return Result{
			name:     name,
			duration: time.Duration(0) * time.Millisecond,
			err:      errors.Wrap(err, "failed to compile query"),
		}
	}

	ctx := context.Background()
	ctx = html.WithDynamic(
		ctx,
		cdp.NewDriver(cdp.WithCDP(r.settings.CDPAddress)),
	)

	ctx = html.WithStatic(ctx, static.NewDriver())

	out, err := p.Run(
		ctx,
		runtime.WithLog(os.Stdout),
		runtime.WithParam("static", r.settings.StaticServerAddress),
		runtime.WithParam("dynamic", r.settings.DynamicServerAddress),
	)

	duration := time.Now().Sub(start)

	if err != nil {
		return Result{
			name:     name,
			duration: duration,
			err:      errors.Wrap(err, "failed to execute query"),
		}
	}

	var result string

	if err := json.Unmarshal(out, &result); err != nil {
		return Result{
			name:     name,
			duration: duration,
			err:      err,
		}
	}

	if result == "" {
		return Result{
			name:     name,
			duration: duration,
		}
	}

	return Result{
		name:     name,
		duration: duration,
		err:      errors.New(result),
	}
}

func (r *Runner) report(results []Result) Summary {
	var failed int
	var passed int
	var sumDuration time.Duration

	for _, res := range results {
		if res.err != nil {
			r.logger.Error().
				Timestamp().
				Err(res.err).
				Str("file", res.name).
				Dur("time", res.duration).
				Msg("Test failed")

			failed++
		} else {
			r.logger.Info().
				Timestamp().
				Str("file", res.name).
				Dur("time", res.duration).
				Msg("Test passed")

			passed++
		}

		sumDuration += res.duration
	}

	return Summary{
		passed:   passed,
		failed:   failed,
		duration: sumDuration,
	}
}
