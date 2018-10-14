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
	dir = flag.String(
		"dir",
		"",
		"root directory with test scripts",
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
	})

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Start server
	go func() {
		if err := s.Start(); err != nil {
			logger.Info().Msg("shutting down the server")
		}
	}()

	dirname := *dir

	if dirname == "" {
		d, err := filepath.Abs(filepath.Dir(os.Args[0]))

		if err != nil {
			logger.Fatal().Err(err).Msg("failed to get dir")

			return
		}

		dirname = d
	}

	r := runner.New(logger, runner.Settings{
		ServerAddress: fmt.Sprintf("http://0.0.0.0:%d", *port),
		CDPAddress:    *cdp,
		Dir:           *dir,
	})

	err := r.Run()

	if err := s.Stop(ctx); err != nil {
		logger.Fatal().Err(err).Msg("failed to stop server")
	}

	if err != nil {
		logger.Fatal().Err(err).Msg("failed to run tests")
	}
}
