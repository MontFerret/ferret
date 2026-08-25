package ferret

import (
	"errors"
	"fmt"
	"io"
	"strings"

	gooptions "github.com/ziflex/go-options"

	"github.com/MontFerret/api"
	"github.com/MontFerret/ferret/v2/pkg/debugger"
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
		debugFormat       debugger.FormatOptions
	}

	nativeSessionOption = gooptions.Option[sessionOptions]
)

func defaultSessionOptions() sessionOptions {
	return sessionOptions{
		outputContentType: encodingjson.ContentType,
		debugFormat:       debugger.DefaultFormatOptions(),
	}
}

func newSessionOptions(setters []SessionOption) (sessionOptions, error) {
	opts := defaultSessionOptions()
	if len(setters) == 0 {
		return opts, nil
	}

	var failures []error

	for _, setter := range setters {
		if setter == nil {
			continue
		}

		target := newSessionOptionTarget(1)
		if err := setter(target); err != nil {
			failures = append(failures, err)
		}

		updated, err := gooptions.ApplyTo(opts, target.setters...)
		if err != nil {
			failures = append(failures, err)
			continue
		}

		opts = updated
	}

	return opts, errors.Join(failures...)
}

func wrapSessionOption(setter nativeSessionOption) SessionOption {
	return func(options SessionOptions) error {
		target, ok := options.(*sessionOptionTarget)
		if !ok {
			return fmt.Errorf("Ferret session option cannot be applied to %T", options)
		}

		target.setters = append(target.setters, setter)

		return nil
	}
}

// WithParam sets a session parameter for the execution.
func WithParam(name string, value any) SessionOption {
	return api.WithParam(name, value)
}

// WithParams merges the provided parameter map into the session environment,
// overriding existing keys while preserving any other previously defined parameters.
func WithParams(params map[string]any) SessionOption {
	return api.WithParams(params)
}

// WithDebugFormat configures bounded debugger value formatting.
func WithDebugFormat(format DebugFormatOptions) SessionOption {
	return wrapSessionOption(gooptions.New(func(session *sessionOptions, format DebugFormatOptions) {
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
		Build())
}

// WithEnvironmentOptions appends VM environment options to the created session.
func WithEnvironmentOptions(opts ...vm.EnvironmentOption) SessionOption {
	return wrapSessionOption(environmentOptionsOption(opts...))
}

// WithOutputContentType selects the output codec content type for session results.
func WithOutputContentType(contentType string) SessionOption {
	return api.WithOutputContentType(contentType)
}

func outputContentTypeOption(contentType string) nativeSessionOption {
	return gooptions.New(func(session *sessionOptions, contentType string) {
		session.outputContentType = strings.TrimSpace(contentType)
	}).
		Value(contentType).
		Named("output content type").
		Validators(gooptions.NotBlank[string]()).
		Build()
}

func sessionParamsOption(params map[string]any) nativeSessionOption {
	return func(s *sessionOptions) error {
		if len(params) == 0 {
			return nil
		}

		rtp, err := runtime.NewParamsFrom(params)

		if err != nil {
			return fmt.Errorf("failed to convert params to runtime.Params: %w", err)
		}

		return environmentOptionsOption(vm.WithParams(rtp))(s)
	}
}

// WithSessionRuntimeParams merges the provided runtime.Params into the session environment,
// overriding existing keys while preserving any other previously defined parameters.
func WithSessionRuntimeParams(params runtime.Params) SessionOption {
	return wrapSessionOption(func(s *sessionOptions) error {
		if len(params) == 0 {
			return nil
		}

		return environmentOptionsOption(vm.WithParams(params))(s)
	})
}

func sessionParamOption(name string, value any) nativeSessionOption {
	return func(s *sessionOptions) error {
		if name == "" {
			return fmt.Errorf("param name cannot be empty")
		}

		if value == nil {
			return fmt.Errorf("param value cannot be nil")
		}

		rtp, err := runtime.NewParamsFrom(map[string]any{name: value})
		if err != nil {
			return fmt.Errorf("failed to convert param to runtime.Params: %w", err)
		}

		return environmentOptionsOption(vm.WithParams(rtp))(s)
	}
}

// WithSessionRuntimeParam adds or overrides a single session parameter using a pre-converted runtime.Value.
func WithSessionRuntimeParam(name string, value runtime.Value) SessionOption {
	return wrapSessionOption(func(s *sessionOptions) error {
		if name == "" {
			return fmt.Errorf("param name cannot be empty")
		}

		if value == nil {
			return fmt.Errorf("param value cannot be nil")
		}

		return environmentOptionsOption(vm.WithParam(name, value))(s)
	})
}

// WithSessionLog sets the writer for logging output.
// The writer can be any io.Writer, such as os.Stdout or a file.
func WithSessionLog(writer io.Writer) SessionOption {
	return wrapSessionOption(func(opts *sessionOptions) error {
		if writer == nil {
			return fmt.Errorf("log writer cannot be nil")
		}

		opts.logger = append(opts.logger, logging.WithWriter(writer))

		return nil
	})
}

// WithSessionLogLevel sets the logging level for the session.
// The logging level determines the severity of log messages that will be recorded.
func WithSessionLogLevel(lvl logging.LogLevel) SessionOption {
	return wrapSessionOption(func(opts *sessionOptions) error {
		if lvl < logging.TraceLevel || lvl > logging.Disabled {
			return fmt.Errorf("invalid log level: %v", lvl)
		}

		opts.logger = append(opts.logger, logging.WithLevel(lvl))

		return nil
	})
}

// WithSessionLogFields sets the fields to be included in log entries for the session.
// These fields can provide additional context for debugging and monitoring purposes.
func WithSessionLogFields(fields map[string]any) SessionOption {
	return wrapSessionOption(func(opts *sessionOptions) error {
		if len(fields) == 0 {
			return nil
		}

		opts.logger = append(opts.logger, logging.WithFields(fields))

		return nil
	})
}

func environmentOptionsOption(opts ...vm.EnvironmentOption) nativeSessionOption {
	return func(session *sessionOptions) error {
		if session == nil || len(opts) == 0 {
			return nil
		}

		for _, opt := range opts {
			if opt != nil {
				session.env = append(session.env, opt)
			}
		}

		return nil
	}
}
