package ferret

import "github.com/MontFerret/ferret/v2/pkg/engine"

// LogLevel controls the minimum severity emitted by Ferret logging.
type LogLevel = engine.LogLevel

const (
	// LogTrace enables the most detailed diagnostic logging.
	LogTrace LogLevel = engine.LogTrace

	// LogDebug enables diagnostic logging intended for development and troubleshooting.
	LogDebug LogLevel = engine.LogDebug

	// LogInfo enables informational logging.
	LogInfo LogLevel = engine.LogInfo

	// LogWarn enables warning and higher-severity logging.
	LogWarn LogLevel = engine.LogWarn

	// LogError enables error and higher-severity logging.
	LogError LogLevel = engine.LogError

	// LogFatal enables fatal and panic logging.
	LogFatal LogLevel = engine.LogFatal

	// LogPanic enables only panic-level logging.
	LogPanic LogLevel = engine.LogPanic

	// LogNone represents log events without an assigned severity.
	LogNone LogLevel = engine.LogNone

	// LogDisabled disables Ferret logging.
	LogDisabled LogLevel = engine.LogDisabled
)

// ParseLogLevel parses a textual Ferret log level.
func ParseLogLevel(input string) (LogLevel, error) {
	return engine.ParseLogLevel(input)
}

// MustParseLogLevel parses a textual Ferret log level and panics if it is invalid.
func MustParseLogLevel(input string) LogLevel {
	return engine.MustParseLogLevel(input)
}
