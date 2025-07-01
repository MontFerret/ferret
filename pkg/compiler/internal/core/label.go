package core

type (
	Label int

	labelRef struct {
		pos   int
		field int
	}

	labelDef struct {
		addr int
	}
)

const InvalidLabel Label = -1
