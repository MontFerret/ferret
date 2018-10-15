package values

import "github.com/MontFerret/ferret/pkg/runtime/core"

type (
	HTMLNode interface {
		core.Value

		NodeType() Int

		NodeName() String

		Length() Int

		InnerText() String

		InnerHTML() String

		Value() core.Value

		GetAttributes() core.Value

		GetAttribute(name String) core.Value

		GetChildNodes() core.Value

		GetChildNode(idx Int) core.Value

		QuerySelector(selector String) core.Value

		QuerySelectorAll(selector String) core.Value

		InnerHTMLBySelector(selector String) String

		InnerHTMLBySelectorAll(selector String) *Array

		InnerTextBySelector(selector String) String

		InnerTextBySelectorAll(selector String) *Array

		CountBySelector(selector String) Int
	}

	HTMLDocument interface {
		HTMLNode

		URL() core.Value
	}
)
