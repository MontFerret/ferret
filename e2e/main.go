package main

import (
	"context"
	"flag"
	"fmt"
	"net"
	"os"
	"os/signal"
	"path/filepath"

	"github.com/MontFerret/ferret/e2e/runner"
	"github.com/MontFerret/ferret/e2e/server"

	"github.com/rs/zerolog"
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

func getOutboundIP() (net.IP, error) {
	conn, err := net.Dial("udp", "8.8.8.8:80")

	if err != nil {
		return nil, err
	}

	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)

	return localAddr.IP, nil
}

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

	if *testsDir == "" {
		_, err := filepath.Abs(filepath.Dir(os.Args[0]))

		if err != nil {
			logger.Fatal().Timestamp().Err(err).Msg("failed to get testsDir")

			return
		}
	}

	var ipAddr string

	// we need it in those cases when a Chrome instance is running inside a container
	// and it needs an external IP to get access to our static web server
	outIP, err := getOutboundIP()

	if err != nil {
		ipAddr = "0.0.0.0"
		logger.Warn().Err(err).Msg("Failed to get outbound IP address")
	} else {
		ipAddr = outIP.String()
	}

	r := runner.New(logger, runner.Settings{
		StaticServerAddress:  fmt.Sprintf("http://%s:%d", ipAddr, staticPort),
		DynamicServerAddress: fmt.Sprintf("http://%s:%d", ipAddr, dynamicPort),
		CDPAddress:           *cdp,
		Dir:                  *testsDir,
		Filter:               *filter,
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

	err = r.Run(ctx)

	if err != nil {
		os.Exit(1)
	}

	os.Exit(0)
}
