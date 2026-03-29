package logging

import "github.com/rs/zerolog"

type LogLevel int8

const (
	DebugLevel LogLevel = iota
	InfoLevel
	WarnLevel
	ErrorLevel
	FatalLevel
	PanicLevel
	NoLevel
	Disabled

	TraceLevel LogLevel = -1
)

func ParseLogLevel(input string) (LogLevel, error) {
	lvl, err := zerolog.ParseLevel(input)

	if err != nil {
		return NoLevel, err
	}

	return LogLevel(lvl), nil
}

func MustParseLogLevel(input string) LogLevel {
	lvl, err := zerolog.ParseLevel(input)

	if err != nil {
		panic(err)
	}

	return LogLevel(lvl)
}
