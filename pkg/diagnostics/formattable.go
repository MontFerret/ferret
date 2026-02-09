package diagnostics

type (
	Formattable interface {
		Format() string
	}

	FormattableError interface {
		error
		Formattable
	}
)

func Format(err error) string {
	if formattable, ok := err.(Formattable); ok {
		return formattable.Format()
	}

	return err.Error()
}
