package source

type (
	Position struct {
		Line   int `json:"line"`
		Column int `json:"column"`
	}

	Span struct {
		Start int `json:"start"`
		End   int `json:"end"`
	}

	Location struct {
		File string `json:"file"`
		Position
	}

	Range struct {
		Location
		Span Span `json:"span"`
	}
)
