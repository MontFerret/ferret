package templates

import (
	"fmt"
	"github.com/MontFerret/ferret/pkg/drivers/cdp/eval"
	"github.com/MontFerret/ferret/pkg/runtime/values"
	"github.com/mafredri/cdp/protocol/runtime"
)

const xpath = `(el, expression, resType) => {
	const unwrap = (item) => {
		return item.nodeType != 2 ? item : item.nodeValue;
	};
	const out = document.evaluate(
		expression,
		el,
		null,
		resType == null ? XPathResult.ANY_TYPE : resType
	);
	let result;

	switch (out.resultType) {
		case XPathResult.UNORDERED_NODE_ITERATOR_TYPE:
		case XPathResult.ORDERED_NODE_ITERATOR_TYPE: {
			result = [];
			let item;

			while ((item = out.iterateNext())) {
				result.push(unwrap(item));
			}

			break;
		}
		case XPathResult.UNORDERED_NODE_SNAPSHOT_TYPE:
		case XPathResult.ORDERED_NODE_SNAPSHOT_TYPE: {
			result = [];

			for (let i = 0; i < out.snapshotLength; i++) {
				const item = out.snapshotItem(i);

				if (item != null) {
					result.push(unwrap(item));
				}
			}
			break;
		}
		case XPathResult.NUMBER_TYPE: {
			result = out.numberValue;
			break;
		}
		case XPathResult.STRING_TYPE: {
			result = out.stringValue;
			break;
		}
		case XPathResult.BOOLEAN_TYPE: {
			result = out.booleanValue;
			break;
		}
		case XPathResult.ANY_UNORDERED_NODE_TYPE:
		case XPathResult.FIRST_ORDERED_NODE_TYPE: {
			const node = out.singleNodeValue;
			
			if (node != null) {
				result = unwrap(node);
			}
			
			break;
		}
		default: {
			break;
		}
	}

	return result;
}
`

var (
	xpathAsElementFragment = fmt.Sprintf(`
const xpath = %s;
const found = xpath(el, selector, XPathResult.FIRST_ORDERED_NODE_TYPE);
`, xpath)

	xpathAsElementArrayFragment = fmt.Sprintf(`
const xpath = %s;
const found = xpath(el, selector, XPathResult.ORDERED_NODE_ITERATOR_TYPE);
`, xpath)
)

func XPath(id runtime.RemoteObjectID, expression values.String) *eval.Function {
	return eval.F(xpath).WithArgRef(id).WithArgValue(expression)
}
