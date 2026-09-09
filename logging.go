package ferret

import "github.com/MontFerret/ferret/v2/pkg/logging"

// LogLevel controls the minimum severity emitted by Ferret logging.
type LogLevel = logging.LogLevel

const (
	// LogTrace enables the most detailed diagnostic logging.
	LogTrace LogLevel = logging.TraceLevel

	// LogDebug enables diagnostic logging intended for development and troubleshooting.
	LogDebug LogLevel = logging.DebugLevel

	// LogInfo enables informational logging.
	LogInfo LogLevel = logging.InfoLevel

	// LogWarn enables warning and higher-severity logging.
	LogWarn LogLevel = logging.WarnLevel

	// LogError enables error and higher-severity logging.
	LogError LogLevel = logging.ErrorLevel

	// LogFatal enables fatal and panic logging.
	LogFatal LogLevel = logging.FatalLevel

	// LogPanic enables only panic-level logging.
	LogPanic LogLevel = logging.PanicLevel

	// LogNone represents log events without an assigned severity.
	LogNone LogLevel = logging.NoLevel

	// LogDisabled disables Ferret logging.
	LogDisabled LogLevel = logging.Disabled
)

// ParseLogLevel parses a textual Ferret log level.
func ParseLogLevel(input string) (LogLevel, error) {
	return logging.ParseLogLevel(input)
}

// MustParseLogLevel parses a textual Ferret log level and panics if it is invalid.
func MustParseLogLevel(input string) LogLevel {
	return logging.MustParseLogLevel(input)
}
