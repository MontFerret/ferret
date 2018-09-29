package logging

import (
	"context"
	"github.com/gofrs/uuid"
	"github.com/rs/zerolog"
	"io"
)

type (
	Level uint8

	Options struct {
		Writer io.Writer
		Level  Level
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

func WithContext(ctx context.Context, opts *Options) context.Context {
	id, err := uuid.NewV4()

	if err != nil {
		panic(err)
	}

	logger := zerolog.New(opts.Writer).
		With().
		Str("id", id.String()).
		Logger()

	logger.WithLevel(zerolog.Level(opts.Level))

	return logger.WithContext(ctx)
}

func FromContext(ctx context.Context) *zerolog.Logger {
	return zerolog.Ctx(ctx)
}
