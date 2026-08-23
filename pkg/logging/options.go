package logging

import (
	"io"

	"github.com/ziflex/go-options"
)

type (
	config struct {
		writer   io.Writer
		fields   map[string]any
		level    LogLevel
		hasLevel bool
	}

	Option = options.Option[config]
)

func defaultConfig() config {
	return config{
		writer: io.Discard,
		level:  ErrorLevel,
	}
}

func applyOptions(setters ...Option) (config, error) {
	return options.ApplyTo(defaultConfig(), setters...)
}

func WithWriter(writer io.Writer) Option {
	return func(config *config) error {
		if writer == nil {
			return nil
		}

		config.writer = writer

		return nil
	}
}

func WithField(key string, value any) Option {
	return func(config *config) error {
		if config.fields == nil {
			config.fields = make(map[string]any)
		}

		config.fields[key] = value

		return nil
	}
}

func WithFields(fields map[string]any) Option {
	return func(config *config) error {
		if len(fields) == 0 {
			return nil
		}

		if config.fields == nil {
			config.fields = make(map[string]any, len(fields))
		}

		for k, v := range fields {
			config.fields[k] = v
		}

		return nil
	}
}

func WithLevel(level LogLevel) Option {
	return options.New[config, LogLevel](func(config *config, level LogLevel) {
		config.level = level
		config.hasLevel = true
	}).
		Value(level).
		Validators(options.OneOf(TraceLevel, DebugLevel, InfoLevel, WarnLevel, ErrorLevel, FatalLevel, PanicLevel, Disabled)).
		Build()
}
