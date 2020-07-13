package templates

const domReadyTemplate = `
if (document.readyState === 'complete') {
	return true;
}

return null;
`

func DOMReady() string {
	return domReadyTemplate
}
