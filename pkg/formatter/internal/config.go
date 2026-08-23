package internal

// Config contains resolved formatter configuration.
type Config struct {
	PrintWidth     uint64
	TabWidth       uint64
	SingleQuote    bool
	BracketSpacing bool
	CaseMode       CaseMode
}

// DefaultConfig returns the default formatter configuration.
func DefaultConfig() Config {
	return Config{
		PrintWidth:     80,
		TabWidth:       4,
		SingleQuote:    false,
		BracketSpacing: true,
		CaseMode:       CaseModeLower,
	}
}
