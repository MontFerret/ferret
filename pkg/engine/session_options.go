package engine

import (
	"fmt"
	"io"
	"strings"

	gooptions "github.com/ziflex/go-options"

	"github.com/MontFerret/ferret/v2/pkg/debugger"
	"github.com/MontFerret/ferret/v2/pkg/logging"
	"github.com/MontFerret/ferret/v2/pkg/runtime"
	"github.com/MontFerret/ferret/v2/pkg/vm"
)

// SessionOption configures a native execution or debug session before resource acquisition.
type SessionOption = gooptions.Option[sessionConfig]

// WithDebugFormat configures bounded debugger value formatting.
func WithDebugFormat(format debugger.FormatOptions) SessionOption {
	return func(session *sessionConfig) error {
		return gooptions.New(func(session *sessionConfig, format debugger.FormatOptions) {
			session.config.DebugFormat = format
		}).
			Value(format).
			Named("debug format").
			Validators(gooptions.Check(func(format debugger.FormatOptions) error {
				if format.MaxDepth <= 0 || format.MaxItems <= 0 || format.MaxBytes <= 0 {
					return fmt.Errorf("debug format limits must be positive")
				}

				return nil
			})).
			Build()(session)
	}
}

// WithEnvironmentOptions is an advanced escape hatch that appends VM environment
// options to the created session. Prefer the root session options for ordinary
// embedding configuration.
func WithEnvironmentOptions(opts ...vm.EnvironmentOption) SessionOption {
	return func(session *sessionConfig) error {
		if len(opts) == 0 {
			return nil
		}

		for _, opt := range opts {
			if opt == nil {
				continue
			}

			session.config.Environment = append(session.config.Environment, opt)
		}

		return nil
	}
}

// WithOutputContentType selects the output codec content type for session results.
func WithOutputContentType(value string) SessionOption {
	return gooptions.New(func(cfg *sessionConfig, value string) {
		cfg.config.OutputContentType = strings.TrimSpace(value)
	}).
		Value(value).
		Named("output content type").
		Validators(gooptions.NotBlank[string]()).
		Build()
}

// WithSessionFSRoot selects the rooted filesystem used by one execution
// session. The session owns the replacement filesystem and inherits the
// engine's read-only policy.
func WithSessionFSRoot(value string) SessionOption {
	return gooptions.New(func(o *sessionConfig, value string) {
		o.config.FSRoot = strings.TrimSpace(value)
	}).
		Value(value).
		Named("fs root").
		Validators(gooptions.NotBlank[string]()).
		Build()
}

// WithSessionParams merges the provided parameter map into the session environment,
// overriding existing keys while preserving any other previously defined parameters.
func WithSessionParams(params map[string]any) SessionOption {
	return func(opts *sessionConfig) error {
		return opts.config.SetParams(params)
	}
}

// WithSessionRuntimeParams merges the provided runtime.Params into the session environment,
// overriding existing keys while preserving any other previously defined parameters.
func WithSessionRuntimeParams(params runtime.Params) SessionOption {
	return func(s *sessionConfig) error {
		if len(params) == 0 {
			return nil
		}

		s.config.Environment = append(s.config.Environment, vm.WithParams(params))

		return nil
	}
}

// WithSessionParam adds or overrides a single session parameter.
func WithSessionParam(name string, value any) SessionOption {
	return func(opts *sessionConfig) error {
		return opts.config.SetParam(name, value)
	}
}

// WithSessionRuntimeParam adds or overrides a single session parameter using a pre-converted runtime.Value.
func WithSessionRuntimeParam(name string, value runtime.Value) SessionOption {
	return func(s *sessionConfig) error {
		s.config.Environment = append(s.config.Environment, vm.WithParam(name, value))

		return nil
	}
}

// WithSessionLog sets the writer for logging output.
// The writer can be any io.Writer, such as os.Stdout or a file.
func WithSessionLog(writer io.Writer) SessionOption {
	return func(opts *sessionConfig) error {
		if writer == nil {
			return fmt.Errorf("log writer cannot be nil")
		}

		opts.config.Logger = append(opts.config.Logger, logging.WithWriter(writer))

		return nil
	}
}

// WithSessionLogLevel sets the logging level for the session.
// The logging level determines the severity of log messages that will be recorded.
func WithSessionLogLevel(lvl logging.LogLevel) SessionOption {
	return func(opts *sessionConfig) error {
		if lvl < logging.TraceLevel || lvl > logging.Disabled {
			return fmt.Errorf("invalid log level: %v", lvl)
		}

		opts.config.Logger = append(opts.config.Logger, logging.WithLevel(lvl))

		return nil
	}
}

// WithSessionLogFields sets the fields to be included in log entries for the session.
// These fields can provide additional context for debugging and monitoring purposes.
func WithSessionLogFields(fields map[string]any) SessionOption {
	return func(opts *sessionConfig) error {
		if len(fields) == 0 {
			return nil
		}

		opts.config.Logger = append(opts.config.Logger, logging.WithFields(fields))

		return nil
	}
}
