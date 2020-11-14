package templates

const getPreviousElementSibling = "(el) => el.previousElementSibling"
const getNextElementSibling = "(el) => el.nextElementSibling"

func GetPreviousElementSibling() string {
	return getPreviousElementSibling
}

func GetNextElementSibling() string {
	return getNextElementSibling
}
