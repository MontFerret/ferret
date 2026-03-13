package core

type (
	Label struct {
		name string
		id   labelID
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
