package internal

import "strings"

type testConfiguration struct {
	printWidth     uint64
	tabWidth       uint64
	singleQuote    bool
	bracketSpacing bool
}

func defaultTestConfiguration() *testConfiguration {
	return &testConfiguration{
		printWidth:     80,
		tabWidth:       4,
		bracketSpacing: true,
	}
}

func (c *testConfiguration) PrintWidth() uint64 {
	return c.printWidth
}

func (c *testConfiguration) TabWidth() uint64 {
	return c.tabWidth
}

func (c *testConfiguration) SingleQuote() bool {
	return c.singleQuote
}

func (c *testConfiguration) BracketSpacing() bool {
	return c.bracketSpacing
}

func (c *testConfiguration) FormatKeyword(value string) string {
	return strings.ToLower(value)
}
