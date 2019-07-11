package templates

import "fmt"

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
		const selector = "%s";
		const found = document.querySelector(selector)

		if (found == null) {
			throw new Error('Element not found by selector: ' + selector);
		}

		found.innerHTML = "%s"
	`,
		selector,
		innerHTML,
	)
}
