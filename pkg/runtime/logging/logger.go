package logging

import (
	"context"
	"io"

	"github.com/rs/zerolog"
)

type (
	Level int8

	Options struct {
		Writer io.Writer
		Level  Level
		Fields map[string]interface{}
	}
)

const (
	DebugLevel Level = iota
	InfoLevel
	WarnLevel
	ErrorLevel
	FatalLevel
	PanicLevel
	NoLevel
	Disabled

	TraceLevel Level = -1
)

func ParseLevel(input string) (Level, error) {
	lvl, err := zerolog.ParseLevel(input)

	if err != nil {
		return NoLevel, err
	}

	return Level(lvl), nil
}

func MustParseLevel(input string) Level {
	lvl, err := zerolog.ParseLevel(input)

	if err != nil {
		panic(err)
	}

	return Level(lvl)
}

func (l Level) String() string {
	return zerolog.Level(l).String()
}

func WithContext(ctx context.Context, opts Options) context.Context {
	c := zerolog.New(opts.Writer).With().Timestamp()

	for k, v := range opts.Fields {
		c = c.Interface(k, v)
	}

	logger := c.Logger().Level(zerolog.Level(opts.Level))

	return logger.WithContext(ctx)
}

func FromContext(ctx context.Context) zerolog.Logger {
	found := zerolog.Ctx(ctx)

	if found == nil {
		panic("logger is not set")
	}

	return *found
}

func WithName(ctx zerolog.Context, name string) zerolog.Context {
	return ctx.Str("component", name)
}
