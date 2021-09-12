package drivers

import (
	"hash/fnv"

	"github.com/wI2L/jettison"

	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
)

type (
	QuerySelectorVariant int

	QuerySelector struct {
		core.Value
		variant QuerySelectorVariant
		value   values.String
	}
)

const (
	UnknownSelector QuerySelectorVariant = iota
	CSSSelector
	XPathSelector
)

var (
	qsvStr = map[QuerySelectorVariant]string{
		UnknownSelector: "unknown",
		CSSSelector:     "css",
		XPathSelector:   "xpath",
	}
)

func (v QuerySelectorVariant) String() string {
	str, found := qsvStr[v]

	if found {
		return str
	}

	return qsvStr[UnknownSelector]
}

func NewCSSSelector(value values.String) QuerySelector {
	return QuerySelector{
		variant: CSSSelector,
		value:   value,
	}
}

func NewXPathSelector(value values.String) QuerySelector {
	return QuerySelector{
		variant: XPathSelector,
		value:   value,
	}
}

func (q QuerySelector) Variant() QuerySelectorVariant {
	return q.variant
}

func (q QuerySelector) MarshalJSON() ([]byte, error) {
	return jettison.MarshalOpts(map[string]string{
		"variant": q.variant.String(),
		"value":   q.value.String(),
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

	if q.variant == otherQS.Variant() {
		return q.value.Compare(values.NewString(otherQS.String()))
	}

	if q.variant == CSSSelector {
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
	h.Write([]byte(q.variant.String()))
	h.Write([]byte(":"))
	h.Write([]byte(q.value.String()))

	return h.Sum64()
}

func (q QuerySelector) Copy() core.Value {
	return q
}
