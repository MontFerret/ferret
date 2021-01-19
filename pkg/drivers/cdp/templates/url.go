package templates

const getUrl = `() => document.location.toString()`

func GetURL() string {
	return getUrl
}