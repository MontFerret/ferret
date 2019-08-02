package templates

import (
	"fmt"
	"github.com/MontFerret/ferret/pkg/drivers"
)

const setInnerHTMLTemplate = `
	(element, value) => {
		element.innerHTML = value;
	}
`

func SetInnerHTML() string {
	return setInnerHTMLTemplate
}

func SetInnerHTMLBySelector(selector, innerHTML string) string {
	return fmt.Sprintf(`
		(el) => {
			const selector = '%s';
			const found = el.querySelector(selector);
	
			if (found == null) {
				throw new Error('%s');
			}
	
			found.innerHTML = '%s';
		}
	`,
		selector,
		drivers.ErrNotFound,
		innerHTML,
	)
}
