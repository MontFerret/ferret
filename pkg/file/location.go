package file

type Location struct {
	line   int
	column int
	start  int
	end    int
}

func NewLocation(line, column, start, end int) *Location {
	return &Location{
		line:   line,
		column: column,
		start:  start,
		end:    end,
	}
}

func EmptyLocation() *Location {
	return &Location{
		line:   0,
		column: 0,
		start:  0,
		end:    0,
	}
}

func (l *Location) Line() int {
	return l.line
}

func (l *Location) Column() int {
	return l.column
}

func (l *Location) Start() int {
	return l.start
}

func (l *Location) End() int {
	return l.end
}
