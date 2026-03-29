package logging

import (
	"context"

	"github.com/rs/zerolog"
)

type Logger = zerolog.Logger

func (l LogLevel) String() string {
	return zerolog.Level(l).String()
}

func New(opts ...Option) zerolog.Logger {
	o := newOptions(opts...)
	c := zerolog.New(o.writer).With().Timestamp()

	for k, v := range o.fields {
		c = c.Interface(k, v)
	}

	return c.Logger().Level(zerolog.Level(o.level))
}

func NewFrom(base Logger, opts ...Option) zerolog.Logger {
	if len(opts) == 0 {
		return base
	}

	o := newOptions(opts...)
	c := base.With()

	for k, v := range o.fields {
		c = c.Interface(k, v)
	}

	l := c.Logger()

	if o.hasLevel {
		l = l.Level(zerolog.Level(o.level))
	}

	return l
}

func With(ctx context.Context, opts ...Option) context.Context {
	return NewFrom(From(ctx), opts...).WithContext(ctx)
}

func From(ctx context.Context) zerolog.Logger {
	return *zerolog.Ctx(ctx)
}
