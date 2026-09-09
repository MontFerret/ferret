package ferret

import (
	"io"

	"github.com/MontFerret/ferret/v2/pkg/engine"
	"github.com/MontFerret/ferret/v2/pkg/vm"
)

// SessionOption configures a native execution or debug session before resource acquisition.
type SessionOption = engine.SessionOption

// WithDebugFormat configures bounded debugger value formatting.
func WithDebugFormat(format DebugFormatOptions) SessionOption {
	return engine.WithDebugFormat(format)
}

// WithEnvironmentOptions is an advanced escape hatch that appends VM environment
// options to the created session. Prefer the root session options for ordinary
// embedding configuration.
func WithEnvironmentOptions(opts ...vm.EnvironmentOption) SessionOption {
	return engine.WithEnvironmentOptions(opts...)
}

// WithOutputContentType selects the output codec content type for session results.
func WithOutputContentType(value string) SessionOption {
	return engine.WithOutputContentType(value)
}

// WithSessionFSRoot selects the rooted filesystem used by one execution
// session. The session owns the replacement filesystem and inherits the
// engine's read-only policy.
func WithSessionFSRoot(value string) SessionOption {
	return engine.WithSessionFSRoot(value)
}

// WithSessionParams merges the provided parameter map into the session environment,
// overriding existing keys while preserving any other previously defined parameters.
func WithSessionParams(params map[string]any) SessionOption {
	return engine.WithSessionParams(params)
}

// WithSessionRuntimeParams merges the provided Params into the session environment,
// overriding existing keys while preserving any other previously defined parameters.
func WithSessionRuntimeParams(params Params) SessionOption {
	return engine.WithSessionRuntimeParams(params)
}

// WithSessionParam adds or overrides a single session parameter.
func WithSessionParam(name string, value any) SessionOption {
	return engine.WithSessionParam(name, value)
}

// WithSessionRuntimeParam adds or overrides a single session parameter using a pre-converted Value.
func WithSessionRuntimeParam(name string, value Value) SessionOption {
	return engine.WithSessionRuntimeParam(name, value)
}

// WithSessionLog sets the writer for logging output.
// The writer can be any io.Writer, such as os.Stdout or a file.
func WithSessionLog(writer io.Writer) SessionOption {
	return engine.WithSessionLog(writer)
}

// WithSessionLogLevel sets the logging level for the session.
// The logging level determines the severity of log messages that will be recorded.
func WithSessionLogLevel(lvl LogLevel) SessionOption {
	return engine.WithSessionLogLevel(lvl)
}

// WithSessionLogFields sets the fields to be included in log entries for the session.
// These fields can provide additional context for debugging and monitoring purposes.
func WithSessionLogFields(fields map[string]any) SessionOption {
	return engine.WithSessionLogFields(fields)
}
