package logging

import (
	"context"
	"github.com/rs/zerolog"
)

type Level uint8

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

func From(ctx context.Context) *zerolog.Logger {
	return zerolog.Ctx(ctx)
}
