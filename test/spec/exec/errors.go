package exec

type ExpectedRuntimeError struct {
	Message     string
	Format      string
	Contains    []string
	NotContains []string
}
