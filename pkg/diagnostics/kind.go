package diagnostics

type Kind string

const (
	Unknown         Kind = ""
	Unsupported     Kind = "Unsupported"
	UnexpectedError Kind = "UnexpectedError"
	TypeError       Kind = "TypeError"
)

func (ek Kind) String() string {
	return string(ek)
}
