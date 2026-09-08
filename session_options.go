package ferret

import (
	"fmt"
	"io"

	gooptions "github.com/ziflex/go-options"

	"github.com/MontFerret/ferret/v2/pkg/logging"
	"github.com/MontFerret/ferret/v2/pkg/vm"
)

// SessionOption configures a native execution or debug session before resource acquisition.
type SessionOption = gooptions.Option[sessionOptions]

// WithDebugFormat configures bounded debugger value formatting.
func WithDebugFormat(format DebugFormatOptions) SessionOption {
	return func(session *sessionOptions) error {
		return gooptions.New(func(session *sessionOptions, format DebugFormatOptions) {
			session.debugFormat = format
		}).
			Value(format).
			Named("debug format").
			Validators(gooptions.Check(func(format DebugFormatOptions) error {
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
	return func(session *sessionOptions) error {
		if len(opts) == 0 {
			return nil
		}

		for _, opt := range opts {
			if opt == nil {
				continue
			}

			session.env = append(session.env, opt)
		}

		return nil
	}
}

// WithOutputContentType selects the output codec content type for session results.
func WithOutputContentType(contentType string) SessionOption {
	return func(opts *sessionOptions) error {
		return opts.setOutputContentType(contentType)
	}
}

// WithSessionFSRoot selects the rooted filesystem used by one execution
// session. The session owns the replacement filesystem and inherits the
// engine's read-only policy.
func WithSessionFSRoot(root string) SessionOption {
	return func(opts *sessionOptions) error {
		return opts.setFSRoot(root)
	}
}

// WithSessionParams merges the provided parameter map into the session environment,
// overriding existing keys while preserving any other previously defined parameters.
func WithSessionParams(params map[string]any) SessionOption {
	return func(opts *sessionOptions) error {
		return opts.setParams(params)
	}
}

// WithSessionRuntimeParams merges the provided Params into the session environment,
// overriding existing keys while preserving any other previously defined parameters.
func WithSessionRuntimeParams(params Params) SessionOption {
	return func(s *sessionOptions) error {
		if len(params) == 0 {
			return nil
		}

		return WithEnvironmentOptions(vm.WithParams(params))(s)
	}
}

// WithSessionParam adds or overrides a single session parameter.
func WithSessionParam(name string, value any) SessionOption {
	return func(opts *sessionOptions) error {
		return opts.setParam(name, value)
	}
}

// WithSessionRuntimeParam adds or overrides a single session parameter using a pre-converted Value.
func WithSessionRuntimeParam(name string, value Value) SessionOption {
	return func(s *sessionOptions) error {
		if name == "" {
			return fmt.Errorf("param name cannot be empty")
		}

		if value == nil {
			return fmt.Errorf("param value cannot be nil")
		}

		return WithEnvironmentOptions(vm.WithParam(name, value))(s)
	}
}

// WithSessionLog sets the writer for logging output.
// The writer can be any io.Writer, such as os.Stdout or a file.
func WithSessionLog(writer io.Writer) SessionOption {
	return func(opts *sessionOptions) error {
		if writer == nil {
			return fmt.Errorf("log writer cannot be nil")
		}

		opts.logger = append(opts.logger, logging.WithWriter(writer))

		return nil
	}
}

// WithSessionLogLevel sets the logging level for the session.
// The logging level determines the severity of log messages that will be recorded.
func WithSessionLogLevel(lvl LogLevel) SessionOption {
	return func(opts *sessionOptions) error {
		if lvl < LogTrace || lvl > LogDisabled {
			return fmt.Errorf("invalid log level: %v", lvl)
		}

		opts.logger = append(opts.logger, logging.WithLevel(lvl))

		return nil
	}
}

// WithSessionLogFields sets the fields to be included in log entries for the session.
// These fields can provide additional context for debugging and monitoring purposes.
func WithSessionLogFields(fields map[string]any) SessionOption {
	return func(opts *sessionOptions) error {
		if len(fields) == 0 {
			return nil
		}

		opts.logger = append(opts.logger, logging.WithFields(fields))

		return nil
	}
}
