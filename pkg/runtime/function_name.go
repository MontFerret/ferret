package runtime

import "strings"

// NamespaceSeparator separates namespace and function segments in registered names.
const NamespaceSeparator = "::"

// CanonicalRegisteredName returns the canonical identity of a registered name.
// It lowercases ASCII A-Z bytes in every qualified segment and leaves all other
// bytes unchanged. When name is already canonical, it returns name without allocating.
// This is the stable identity contract for runtime registries, extensions, and tooling.
func CanonicalRegisteredName(name string) string {
	for i := 0; i < len(name); i++ {
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

// HasTerminalFunctionName reports whether name contains a non-empty function
// segment after its final namespace separator. It performs no other validation.
// This is the stable validation contract for extension and tooling registration paths.
func HasTerminalFunctionName(name string) bool {
	terminalStart := 0
	if separator := strings.LastIndex(name, NamespaceSeparator); separator >= 0 {
		terminalStart = separator + len(NamespaceSeparator)
	}

	return terminalStart < len(name)
}
