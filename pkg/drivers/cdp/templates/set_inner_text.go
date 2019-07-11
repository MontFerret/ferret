package templates

import "fmt"

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
			throw new Error('Element not found by selector: ' + selector);
		}

		found.innerText = "%s"
	`,
		selector,
		innerText,
	)
}
