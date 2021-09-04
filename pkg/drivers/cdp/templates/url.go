package templates

import "github.com/MontFerret/ferret/pkg/drivers/cdp/eval"

const getURL = `() => window.location.toString()`

func GetURL() *eval.Function {
	return eval.F(getURL)
}
