package logging

import (
	"context"
	"io"

	"github.com/rs/zerolog"
)

type (
	Level uint8

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
)

func ParseLevel(input string) (Level, error) {
	lvl, err := zerolog.ParseLevel(input)

	if err != nil {
		return NoLevel, err
	}

	return Level(lvl), nil
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

func FromContext(ctx context.Context) *zerolog.Logger {
	return zerolog.Ctx(ctx)
}
