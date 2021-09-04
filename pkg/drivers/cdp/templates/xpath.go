package templates

import (
	"github.com/MontFerret/ferret/pkg/drivers/cdp/eval"
	"github.com/MontFerret/ferret/pkg/runtime/values"
	"github.com/mafredri/cdp/protocol/runtime"
)

const xpath = `(el, expression) => {
	const unwrap = (item) => {
		return item.nodeType != 2 ? item : item.nodeValue;
	};
	const out = document.evaluate(
		expression,
		el,
		null,
		XPathResult.ANY_TYPE
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
			result = unwrap(out.singleNodeValue);
			break;
		}
		default: {
			break;
		}
	}

	return result;
}
`

func XPath(id runtime.RemoteObjectID, expression values.String) *eval.Function {
	return eval.F(xpath).WithArgRef(id).WithArgValue(expression)
}
