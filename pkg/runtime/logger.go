package runtime

import (
	"context"
	"io"

	"github.com/rs/zerolog"
)

type (
	LogLevel int8

	LogSettings struct {
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

func NewLogger(opts LogSettings) zerolog.Logger {
	c := zerolog.New(opts.Writer).With().Timestamp()

	for k, v := range opts.Fields {
		c = c.Interface(k, v)
	}

	return c.Logger().Level(zerolog.Level(opts.Level))
}

func WithLogger(ctx context.Context, opts LogSettings) context.Context {
	return NewLogger(opts).WithContext(ctx)
}

func GetLogger(ctx context.Context) zerolog.Logger {
	found := zerolog.Ctx(ctx)

	if found == nil {
		panic("logger is not set")
	}

	return *found
}
