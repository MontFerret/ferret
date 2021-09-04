package templates

import "github.com/MontFerret/ferret/pkg/drivers/cdp/eval"

const getURL = `() => return window.location.toString()`

func GetURL() *eval.Function {
	return eval.F(getURL)
}
