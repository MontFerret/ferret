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

	port = flag.Uint64(
		"port",
		8080,
		"server port",
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

	s := server.New(server.Settings{
		Port: *port,
		Dir:  *pagesDir,
	})

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		if err := s.Start(); err != nil {
			logger.Info().Timestamp().Msg("shutting down the server")
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
		ServerAddress: fmt.Sprintf("http://0.0.0.0:%d", *port),
		CDPAddress:    *cdp,
		Dir:           *testsDir,
	})

	err := r.Run()

	if err := s.Stop(ctx); err != nil {
		logger.Fatal().Timestamp().Err(err).Msg("failed to stop server")
	}

	if err != nil {
		os.Exit(1)
	}
}
