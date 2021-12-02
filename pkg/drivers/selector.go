package drivers

import (
	"hash/fnv"

	"github.com/wI2L/jettison"

	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
)

type (
	QuerySelectorKind int

	QuerySelector struct {
		core.Value
		kind  QuerySelectorKind
		value values.String
	}
)

const (
	UnknownSelector QuerySelectorKind = iota
	CSSSelector
	XPathSelector
)

var (
	qsvStr = map[QuerySelectorKind]string{
		UnknownSelector: "unknown",
		CSSSelector:     "css",
		XPathSelector:   "xpath",
	}
)

func (v QuerySelectorKind) String() string {
	str, found := qsvStr[v]

	if found {
		return str
	}

	return qsvStr[UnknownSelector]
}

func NewCSSSelector(value values.String) QuerySelector {
	return QuerySelector{
		kind:  CSSSelector,
		value: value,
	}
}

func NewXPathSelector(value values.String) QuerySelector {
	return QuerySelector{
		kind:  XPathSelector,
		value: value,
	}
}

func (q QuerySelector) Kind() QuerySelectorKind {
	return q.kind
}

func (q QuerySelector) MarshalJSON() ([]byte, error) {
	return jettison.MarshalOpts(map[string]string{
		"kind":  q.kind.String(),
		"value": q.value.String(),
	}, jettison.NoHTMLEscaping())
}

func (q QuerySelector) Type() core.Type {
	return QuerySelectorType
}

func (q QuerySelector) String() string {
	return q.value.String()
}

func (q QuerySelector) Compare(other core.Value) int64 {
	if other.Type() != QuerySelectorType {
		return Compare(QuerySelectorType, other.Type())
	}

	otherQS := other.(*QuerySelector)

	if q.kind == otherQS.Kind() {
		return q.value.Compare(values.NewString(otherQS.String()))
	}

	if q.kind == CSSSelector {
		return -1
	}

	return 0
}

func (q QuerySelector) Unwrap() interface{} {
	return q.value
}

func (q QuerySelector) Hash() uint64 {
	h := fnv.New64a()

	h.Write([]byte(q.Type().String()))
	h.Write([]byte(":"))
	h.Write([]byte(q.kind.String()))
	h.Write([]byte(":"))
	h.Write([]byte(q.value.String()))

	return h.Sum64()
}

func (q QuerySelector) Copy() core.Value {
	return q
}
