package optimization

type Level int

const (
	None Level = iota
	Basic
	Full
)

var (
	levelNames = map[Level]string{
		None:  "none",
		Basic: "basic",
		Full:  "full",
	}
)

func (l Level) String() string {
	if name, ok := levelNames[l]; ok {
		return name
	}

	return "unknown"
}
