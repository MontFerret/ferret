package runner

import (
	"context"
	"encoding/json"
	"github.com/MontFerret/ferret/pkg/compiler"
	"github.com/MontFerret/ferret/pkg/runtime"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"io/ioutil"
	"path/filepath"
)

type (
	Settings struct {
		ServerAddress string
		CDPAddress    string
		Dir           string
	}

	Result struct {
		name string
		err  error
	}

	Stats struct {
		passed int
		failed int
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
	static, err := r.runQueries(filepath.Join(r.settings.Dir, "static"))

	if err != nil {
		return err
	}

	ss := r.report(static)

	dynamic, err := r.runQueries(filepath.Join(r.settings.Dir, "dynamic"))

	if err != nil {
		r.stat(ss)
		return err
	}

	ds := r.report(dynamic)

	r.stat(Stats{
		passed: ss.passed + ds.passed,
		failed: ss.failed + ds.failed,
	})

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
	c.RegisterFunctions(Assertions())

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
	p, err := c.Compile(script)

	if err != nil {
		return Result{
			name: name,
			err:  errors.Wrap(err, "failed to compile query"),
		}
	}

	out, err := p.Run(
		context.Background(),
		runtime.WithBrowser(r.settings.CDPAddress),
		runtime.WithParam("url", r.settings.ServerAddress),
	)

	if err != nil {
		return Result{
			name: name,
			err:  errors.Wrap(err, "failed to execute query"),
		}
	}

	var result string

	json.Unmarshal(out, &result)

	if result == "" {
		return Result{
			name: name,
		}
	}

	return Result{
		name: name,
		err:  errors.New(result),
	}
}

func (r *Runner) report(results []Result) Stats {
	s := Stats{}

	for _, res := range results {
		if res.err != nil {
			r.logger.Error().
				Timestamp().
				Err(res.err).
				Str("file", res.name).
				Msg("Test failed")

			s.failed++
		} else {
			r.logger.Info().
				Timestamp().
				Str("file", res.name).
				Msg("Test passed")

			s.passed++
		}
	}

	return s
}

func (r *Runner) stat(stats Stats) {
	r.logger.Info().
		Timestamp().
		Int("passed", stats.passed).
		Int("failed", stats.failed).
		Msg("Done")
}
