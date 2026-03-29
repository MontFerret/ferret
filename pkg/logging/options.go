package logging

import (
	"io"
)

type (
	loggerOptions struct {
		writer   io.Writer
		fields   map[string]any
		level    LogLevel
		hasLevel bool
	}

	Option func(*loggerOptions)
)

func newOptions(opts ...Option) *loggerOptions {
	o := &loggerOptions{
		writer: io.Discard,
		level:  ErrorLevel,
	}

	for _, opt := range opts {
		opt(o)
	}

	return o
}

func WithWriter(writer io.Writer) Option {
	return func(s *loggerOptions) {
		if writer == nil {
			return
		}

		s.writer = writer
	}
}

func WithField(key string, value any) Option {
	return func(s *loggerOptions) {
		if s.fields == nil {
			s.fields = make(map[string]any)
		}

		s.fields[key] = value
	}
}

func WithFields(fields map[string]any) Option {
	return func(s *loggerOptions) {
		if len(fields) == 0 {
			return
		}

		if s.fields == nil {
			s.fields = make(map[string]any, len(fields))
		}

		for k, v := range fields {
			s.fields[k] = v
		}
	}
}

func WithLevel(level LogLevel) Option {
	return func(s *loggerOptions) {
		if level < TraceLevel || level > Disabled {
			return
		}

		s.level = level
		s.hasLevel = true
	}
}
