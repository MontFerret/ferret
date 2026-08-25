package ferret

import apisource "github.com/MontFerret/api/source"

// NewSource creates a named FQL source file.
func NewSource(name, content string) Source {
	return apisource.New(name, content)
}

// NewAnonymousSource creates an FQL source file named "anonymous".
func NewAnonymousSource(content string) Source {
	return apisource.NewAnonymous(content)
}
