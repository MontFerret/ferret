package runner

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/MontFerret/ferret/pkg/compiler"
	"github.com/MontFerret/ferret/pkg/drivers"
	"github.com/MontFerret/ferret/pkg/drivers/cdp"
	"github.com/MontFerret/ferret/pkg/drivers/http"
	"github.com/MontFerret/ferret/pkg/runtime"

	"github.com/gobwas/glob"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
)

type (
	Settings struct {
		StaticServerAddress  string
		DynamicServerAddress string
		CDPAddress           string
		Dir                  string
		Filter               string
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

func (r *Runner) Run(ctx context.Context) error {
	ctx = drivers.WithContext(
		ctx,
		cdp.NewDriver(cdp.WithAddress(r.settings.CDPAddress)),
	)

	ctx = drivers.WithContext(
		ctx,
		cdp.NewDriver(cdp.WithAddress(r.settings.CDPAddress),
			cdp.WithCustomName("cdp_headers"),
			cdp.WithHeader("Single_header", []string{"single_header_value"}),
			cdp.WithHeaders(drivers.HTTPHeaders{
				"Multi_set_header":  []string{"multi_set_header_value"},
				"Multi_set_header2": []string{"multi_set_header2_value"},
			}),
		),
	)

	ctx = drivers.WithContext(
		ctx,
		cdp.NewDriver(cdp.WithAddress(r.settings.CDPAddress),
			cdp.WithCustomName("cdp_cookies"),
			cdp.WithCookie(drivers.HTTPCookie{
				Name:     "single_cookie",
				Value:    "single_cookie_value",
				Path:     "/",
				MaxAge:   0,
				Secure:   false,
				HTTPOnly: false,
				SameSite: 0,
			}),
			cdp.WithCookies([]drivers.HTTPCookie{
				{
					Name:     "multi_set_cookie",
					Value:    "multi_set_cookie_value",
					Path:     "/",
					MaxAge:   0,
					Secure:   false,
					HTTPOnly: false,
					SameSite: 0,
				},
				{
					Name:     "multi_set_cookie2",
					Value:    "multi_set_cookie2_value",
					Path:     "/",
					MaxAge:   0,
					Secure:   false,
					HTTPOnly: false,
					SameSite: 0,
				},
			}),
		),
	)

	ctx = drivers.WithContext(
		ctx,
		http.NewDriver(),
		drivers.AsDefault(),
	)

	ctx = drivers.WithContext(
		ctx,
		http.NewDriver(
			http.WithCustomName("http_headers"),
			http.WithHeader("Single_header", []string{"single_header_value"}),
			http.WithHeaders(drivers.HTTPHeaders{
				"Multi_set_header":  []string{"multi_set_header_value"},
				"Multi_set_header2": []string{"multi_set_header2_value"},
			}),
		),
	)

	ctx = drivers.WithContext(
		ctx,
		http.NewDriver(
			http.WithCustomName("http_cookies"),
			http.WithCookie(drivers.HTTPCookie{
				Name:     "single_cookie",
				Value:    "single_cookie_value",
				Path:     "/",
				MaxAge:   0,
				Secure:   false,
				HTTPOnly: false,
				SameSite: 0,
			}),
			http.WithCookies([]drivers.HTTPCookie{
				{
					Name:     "multi_set_cookie",
					Value:    "multi_set_cookie_value",
					Path:     "/",
					MaxAge:   0,
					Secure:   false,
					HTTPOnly: false,
					SameSite: 0,
				},
				{
					Name:     "multi_set_cookie2",
					Value:    "multi_set_cookie2_value",
					Path:     "/",
					MaxAge:   0,
					Secure:   false,
					HTTPOnly: false,
					SameSite: 0,
				},
			}),
		),
	)

	results, err := r.runQueries(ctx, r.settings.Dir)

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
		Str("duration", sum.duration.String()).
		Msg("Completed")

	if sum.failed > 0 {
		return errors.New("failed")
	}

	return nil
}

func (r *Runner) runQueries(ctx context.Context, dir string) ([]Result, error) {
	results := make([]Result, 0, 50)

	c := compiler.New()

	// backward compatible
	if err := Assertions(c); err != nil {
		return nil, err
	}

	ns := c.Namespace("T")

	if err := Assertions(ns); err != nil {
		return nil, err
	}

	if err := HTTPHelpers(ns.Namespace("HTTP")); err != nil {
		return nil, err
	}

	var filter glob.Glob
	var useFilter bool

	if r.settings.Filter != "" {
		f, err := glob.Compile(r.settings.Filter)

		if err != nil {
			return nil, err
		}

		filter = f
		useFilter = true
	}

	err := r.traverseDir(ctx, dir, func(name string) error {
		if useFilter {
			if !filter.Match(name) {
				return nil
			}
		}

		b, err := ioutil.ReadFile(name)

		if err != nil {
			results = append(results, Result{
				name: name,
				err:  errors.Wrap(err, "failed to read script file"),
			})

			return nil
		}

		r.logger.Info().Timestamp().Str("name", name).Msg("Running test")

		select {
		case <-ctx.Done():
			return context.Canceled
		default:
			result := r.runQuery(ctx, c, name, string(b))

			if result.err == nil {
				r.logger.Info().
					Timestamp().
					Str("file", result.name).
					Str("duration", result.duration.String()).
					Msg("Test passed")
			} else {
				r.logger.Error().
					Timestamp().
					Err(result.err).
					Str("file", result.name).
					Str("duration", result.duration.String()).
					Msg("Test failed")
			}

			results = append(results, result)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return results, nil
}

func (r *Runner) runQuery(ctx context.Context, c *compiler.FqlCompiler, name, script string) Result {
	start := time.Now()

	p, err := c.Compile(script)

	if err != nil {
		return Result{
			name:     name,
			duration: time.Duration(0) * time.Millisecond,
			err:      errors.Wrap(err, "failed to compile query"),
		}
	}

	mustFail := r.mustFail(name)

	out, err := p.Run(
		ctx,
		runtime.WithLog(zerolog.ConsoleWriter{Out: os.Stdout}),
		runtime.WithParam("static", r.settings.StaticServerAddress),
		runtime.WithParam("dynamic", r.settings.DynamicServerAddress),
	)

	duration := time.Since(start)

	if err != nil {
		if mustFail {
			return Result{
				name:     name,
				duration: duration,
			}
		}

		return Result{
			name:     name,
			duration: duration,
			err:      errors.Wrap(err, "failed to execute query"),
		}
	}

	if mustFail {
		return Result{
			name:     name,
			duration: duration,
			err:      errors.New("expected to fail"),
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

			failed++
		} else {

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

func (r *Runner) traverseDir(ctx context.Context, dir string, iteratee func(name string) error) error {
	files, err := ioutil.ReadDir(dir)

	if err != nil {
		r.logger.Error().
			Timestamp().
			Err(err).
			Str("dir", dir).
			Msg("failed to read scripts directory")

		return err
	}

	for _, file := range files {
		name := filepath.Join(dir, file.Name())

		if file.IsDir() {
			if err := r.traverseDir(ctx, name, iteratee); err != nil {
				return err
			}

			continue
		}

		if err := iteratee(name); err != nil {
			return err
		}
	}

	return nil
}

func (r *Runner) mustFail(name string) bool {
	return strings.HasSuffix(name, ".fail.fql")
}
