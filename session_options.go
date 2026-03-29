package ferret

import (
	"fmt"
	"io"
	"strings"

	encodingjson "github.com/MontFerret/ferret/v2/pkg/encoding/json"
	"github.com/MontFerret/ferret/v2/pkg/logging"
	"github.com/MontFerret/ferret/v2/pkg/runtime"
	"github.com/MontFerret/ferret/v2/pkg/vm"
)

type (
	sessionOptions struct {
		logger            []logging.Option
		outputContentType string
		env               []vm.EnvironmentOption
	}

	SessionOption func(*sessionOptions) error
)

func newSessionOptions(setters []SessionOption) (*sessionOptions, error) {
	opts := &sessionOptions{
		outputContentType: encodingjson.ContentType,
	}

	for _, setter := range setters {
		if setter == nil {
			continue
		}

		if err := setter(opts); err != nil {
			return nil, err
		}
	}

	return opts, nil
}

func WithEnvironmentOptions(opts ...vm.EnvironmentOption) SessionOption {
	return func(session *sessionOptions) error {
		if session == nil {
			return nil
		}

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

func WithOutputContentType(contentType string) SessionOption {
	return func(session *sessionOptions) error {
		if session == nil {
			return nil
		}

		trimmed := strings.TrimSpace(contentType)
		if trimmed == "" {
			return fmt.Errorf("output content type cannot be empty")
		}

		session.outputContentType = trimmed
		return nil
	}
}

func WithSessionParams(params runtime.Params) SessionOption {
	return WithEnvironmentOptions(vm.WithParams(params))
}

func WithSessionParam(name string, value runtime.Value) SessionOption {
	return WithEnvironmentOptions(vm.WithParam(name, value))
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
func WithSessionLogLevel(lvl logging.LogLevel) SessionOption {
	return func(opts *sessionOptions) error {
		if lvl < logging.TraceLevel || lvl > logging.Disabled {
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
		if fields == nil {
			return fmt.Errorf("log fields cannot be nil")
		}

		opts.logger = append(opts.logger, logging.WithFields(fields))

		return nil
	}
}
