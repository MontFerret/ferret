package templates

import (
	"fmt"
	"github.com/MontFerret/ferret/pkg/drivers/cdp/eval"
)

const getChildren = "(el) => Array.from(el.children)"
const getChildrenCount = "(el) => el.children.length"

func GetChildren() string {
	return getChildren
}

func GetChildrenCount() string {
	return getChildrenCount
}

func GetChildByIndex(idx int64) string {
	return fmt.Sprintf(`
		(el) => el.children[%s]
`, eval.ParamInt(idx))
}
