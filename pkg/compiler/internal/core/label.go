package core

type (
	Label struct {
		id   labelID
		name string
		addr int
	}

	labelID int

	labelRef struct {
		pos   int
		field int
	}
)

func (l Label) String() string {
	return l.name
}
