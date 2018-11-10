package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/MontFerret/ferret/e2e/runner"
	"github.com/MontFerret/ferret/e2e/server"
	"github.com/rs/zerolog"
	"os"
	"path/filepath"
)

var (
	testsDir = flag.String(
		"tests",
		"./tests",
		"root directory with test scripts",
	)

	pagesDir = flag.String(
		"pages",
		"./pages",
		"root directory with test pages",
	)

	cdp = flag.String(
		"cdp",
		"http://0.0.0.0:9222",
		"address of remote Chrome instance",
	)
)

func main() {
	flag.Parse()

	logger := zerolog.New(zerolog.ConsoleWriter{Out: os.Stderr})

	staticPort := uint64(8080)
	static := server.New(server.Settings{
		Port: staticPort,
		Dir:  filepath.Join(*pagesDir, "static"),
	})

	dynamicPort := uint64(8081)
	dynamic := server.New(server.Settings{
		Port: dynamicPort,
		Dir:  filepath.Join(*pagesDir, "dynamic"),
	})

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		if err := static.Start(); err != nil {
			logger.Info().Timestamp().Msg("shutting down the static pages server")
		}
	}()

	go func() {
		if err := dynamic.Start(); err != nil {
			logger.Info().Timestamp().Msg("shutting down the dynamic pages server")
		}
	}()

	dirname := *testsDir

	if dirname == "" {
		d, err := filepath.Abs(filepath.Dir(os.Args[0]))

		if err != nil {
			logger.Fatal().Timestamp().Err(err).Msg("failed to get testsDir")

			return
		}

		dirname = d
	}

	r := runner.New(logger, runner.Settings{
		StaticServerAddress:  fmt.Sprintf("http://0.0.0.0:%d", staticPort),
		DynamicServerAddress: fmt.Sprintf("http://0.0.0.0:%d", dynamicPort),
		CDPAddress:           *cdp,
		Dir:                  *testsDir,
	})

	err := r.Run()

	if err := static.Stop(ctx); err != nil {
		logger.Fatal().Timestamp().Err(err).Msg("failed to stop the static pages server")
	}

	if err := dynamic.Stop(ctx); err != nil {
		logger.Fatal().Timestamp().Err(err).Msg("failed to stop the dynamic pages server")
	}

	if err != nil {
		os.Exit(1)
	}
}
