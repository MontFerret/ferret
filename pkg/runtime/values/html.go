package values

import "github.com/MontFerret/ferret/pkg/runtime/core"

type (
	HtmlNode interface {
		core.Value

		NodeType() Int

		NodeName() String

		Length() Int

		InnerText() String

		InnerHtml() String

		Value() core.Value

		GetAttributes() core.Value

		GetAttribute(name String) core.Value

		GetChildNodes() core.Value

		GetChildNode(idx Int) core.Value

		QuerySelector(selector String) core.Value

		QuerySelectorAll(selector String) core.Value
	}

	HtmlDocument interface {
		HtmlNode

		Url() core.Value
	}
)
