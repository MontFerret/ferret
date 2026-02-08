package diagnostics

import "errors"

type Formattable interface {
	Format() string
}

func FormatError2(err error) string {
	if formattable, ok := err.(Formattable); ok {
		return formattable.Format()
	}

	if errors.Unwrap(err) != nil {
	}

	return err.Error()
}
