package runtime

import (
	"context"
	"github.com/MontFerret/ferret/pkg/runtime/logging"
	"github.com/gofrs/uuid"
	"github.com/rs/zerolog"
	"io"
	"os"
)

type (
	Options struct {
		proxy     string
		cdp       string
		variables map[string]interface{}
		logWriter io.Writer
		logLevel  zerolog.Level
	}

	Option func(*Options)
)

func newOptions() *Options {
	return &Options{
		cdp:       "http://0.0.0.0:9222",
		variables: make(map[string]interface{}),
		logWriter: os.Stdout,
		logLevel:  zerolog.ErrorLevel,
	}
}

func WithParam(name string, value interface{}) Option {
	return func(options *Options) {
		options.variables[name] = value
	}
}

func WithBrowser(address string) Option {
	return func(options *Options) {
		options.cdp = address
	}
}

func WithProxy(address string) Option {
	return func(options *Options) {
		// TODO: add implementation
		options.proxy = address
	}
}

func WithLog(writer io.Writer) Option {
	return func(options *Options) {
		options.logWriter = writer
	}
}

func WithLogLevel(lvl logging.Level) Option {
	return func(options *Options) {
		options.logLevel = zerolog.Level(lvl)
	}
}

func (opts *Options) withContext(parent context.Context) context.Context {
	ctx := context.WithValue(
		parent,
		"variables",
		opts.variables,
	)

	id, err := uuid.NewV4()

	if err != nil {
		panic(err)
	}

	logger := zerolog.New(opts.logWriter).
		With().
		Str("id", id.String()).
		Logger()
	logger.WithLevel(opts.logLevel)

	ctx = logger.WithContext(ctx)

	return ctx
}
