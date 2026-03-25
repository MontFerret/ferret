package assert

type (
	Unary func(actual any) error

	Binary func(actual, expected any) error
)
