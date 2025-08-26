package runtime

import (
	"context"
	"io"

	"github.com/rs/zerolog"
)

type (
	LogLevel int8

	LogOptions struct {
		Writer io.Writer
		Level  LogLevel
		Fields map[string]interface{}
	}
)

const (
	DebugLevel LogLevel = iota
	InfoLevel
	WarnLevel
	ErrorLevel
	FatalLevel
	PanicLevel
	NoLevel
	Disabled

	TraceLevel LogLevel = -1
)

func ParseLogLevel(input string) (LogLevel, error) {
	lvl, err := zerolog.ParseLevel(input)

	if err != nil {
		return NoLevel, err
	}

	return LogLevel(lvl), nil
}

func MustParseLogLevel(input string) LogLevel {
	lvl, err := zerolog.ParseLevel(input)

	if err != nil {
		panic(err)
	}

	return LogLevel(lvl)
}

func (l LogLevel) String() string {
	return zerolog.Level(l).String()
}

func LoggerWithContext(ctx context.Context, opts LogOptions) context.Context {
	c := zerolog.New(opts.Writer).With().Timestamp()

	for k, v := range opts.Fields {
		c = c.Interface(k, v)
	}

	logger := c.Logger().Level(zerolog.Level(opts.Level))

	return logger.WithContext(ctx)
}

func LoggerFromContext(ctx context.Context) zerolog.Logger {
	found := zerolog.Ctx(ctx)

	if found == nil {
		panic("logger is not set")
	}

	return *found
}

func LogWithName(ctx zerolog.Context, name string) zerolog.Context {
	return ctx.Str("component", name)
}
