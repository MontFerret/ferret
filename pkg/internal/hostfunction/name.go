package hostfunction

import "strings"

const namespaceSeparator = "::"

// CanonicalName returns the canonical host-function identity while preserving
// every namespace segment exactly as supplied.
func CanonicalName(name string) string {
	terminalStart := 0
	if separator := strings.LastIndex(name, namespaceSeparator); separator >= 0 {
		terminalStart = separator + len(namespaceSeparator)
	}

	for i := terminalStart; i < len(name); i++ {
		if name[i] < 'A' || name[i] > 'Z' {
			continue
		}

		canonical := []byte(name)
		for j := i; j < len(canonical); j++ {
			if canonical[j] >= 'A' && canonical[j] <= 'Z' {
				canonical[j] += 'a' - 'A'
			}
		}

		return string(canonical)
	}

	return name
}

// HasTerminalName reports whether name contains a non-empty terminal function segment.
func HasTerminalName(name string) bool {
	terminalStart := 0
	if separator := strings.LastIndex(name, namespaceSeparator); separator >= 0 {
		terminalStart = separator + len(namespaceSeparator)
	}

	return terminalStart < len(name)
}
