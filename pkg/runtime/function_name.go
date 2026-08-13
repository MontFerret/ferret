package runtime

import "strings"

// NamespaceSeparator separates namespace and function segments in registered names.
const NamespaceSeparator = "::"

// NormalizeRegisteredName returns the case-insensitive lookup key for a registered name.
// It lowercases ASCII A-Z bytes in every qualified segment and leaves all other
// bytes unchanged. The result is the stable identity key for runtime registries,
// extensions, and tooling; it is not a presentation format. When name is already
// normalized, it returns name without allocating.
func NormalizeRegisteredName(name string) string {
	for i := 0; i < len(name); i++ {
		if name[i] < 'A' || name[i] > 'Z' {
			continue
		}

		normalized := []byte(name)
		for j := i; j < len(normalized); j++ {
			if normalized[j] >= 'A' && normalized[j] <= 'Z' {
				normalized[j] += 'a' - 'A'
			}
		}

		return string(normalized)
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
