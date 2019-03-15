package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/MontFerret/ferret/e2e/runner"
	"github.com/MontFerret/ferret/e2e/server"
	"github.com/rs/zerolog"
	"os"
	"os/signal"
	"path/filepath"
	"regexp"
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

	filter = flag.String(
		"filter",
		"",
		"regexp expression to filter out tests",
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

	var filterR *regexp.Regexp

	if *filter != "" {
		r, err := regexp.Compile(*filter)

		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}

		filterR = r
	}

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
		Filter:               filterR,
	})

	ctx, cancel := context.WithCancel(context.Background())
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	go func() {
		for {
			<-c
			cancel()
		}
	}()

	err := r.Run(ctx)

	if err != nil {
		os.Exit(1)
	}

	os.Exit(0)
}
