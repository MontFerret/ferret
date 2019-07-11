package templates

import (
	"fmt"
	"github.com/MontFerret/ferret/pkg/drivers"
)

const setInnerTextTemplate = `
	(element, value) => {
		element.innerText = value;
	}
`

func SetInnerText() string {
	return setInnerTextTemplate
}

func SetInnerTextBySelector(selector, innerText string) string {
	return fmt.Sprintf(`
		const selector = "%s";
		const found = document.querySelector(selector)

		if (found == null) {
			throw new Error('%s');
		}

		found.innerText = "%s"
	`,
		selector,
		drivers.ErrNotFound,
		innerText,
	)
}
