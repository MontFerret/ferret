package templates

import (
	"fmt"
	"github.com/MontFerret/ferret/pkg/drivers"
)

func Blur() string {
	return `
		(el) => {
			el.blur()
		}
	`
}

func BlurBySelector(selector string) string {
	return fmt.Sprintf(`
		(parent) => {
			const el = parent.querySelector('%s');

			if (el == null) {
				throw new Error('%s')
			}

			el.blur();
		}
`, selector, drivers.ErrNotFound)
}
