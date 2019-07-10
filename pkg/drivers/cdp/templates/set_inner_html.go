package templates

const setInnerHTMLTemplate = `
	(element, value) => {
		element.innerHTML = value;
	}
`

func SetInnerHTML() string {
	return setInnerHTMLTemplate
}
