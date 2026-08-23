package logging

import (
	"context"

	"github.com/rs/zerolog"
)

type Logger = zerolog.Logger

func (l LogLevel) String() string {
	return zerolog.Level(l).String()
}

// New creates a logger from the provided options. It returns any failures
// reported while applying the options.
func New(opts ...Option) (zerolog.Logger, error) {
	cfg, err := applyOptions(opts...)
	if err != nil {
		return zerolog.Logger{}, err
	}

	c := zerolog.New(cfg.writer).With().Timestamp()

	for k, v := range cfg.fields {
		c = c.Interface(k, v)
	}

	return c.Logger().Level(zerolog.Level(cfg.level)), nil
}

// NewFrom derives a logger from base. It returns any failures reported while
// applying the options.
func NewFrom(base Logger, opts ...Option) (zerolog.Logger, error) {
	if len(opts) == 0 {
		return base, nil
	}

	cfg, err := applyOptions(opts...)
	if err != nil {
		return zerolog.Logger{}, err
	}

	c := base.With()

	for k, v := range cfg.fields {
		c = c.Interface(k, v)
	}

	l := c.Logger()

	if cfg.hasLevel {
		l = l.Level(zerolog.Level(cfg.level))
	}

	return l, nil
}

// With derives a logger from the one stored in ctx and returns a context that
// contains it. It returns any failures reported while applying the options.
func With(ctx context.Context, opts ...Option) (context.Context, error) {
	logger, err := NewFrom(From(ctx), opts...)
	if err != nil {
		return nil, err
	}

	return logger.WithContext(ctx), nil
}

func From(ctx context.Context) zerolog.Logger {
	return *zerolog.Ctx(ctx)
}
