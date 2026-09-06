package encoding

import (
	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

// RegisterLib registers serialization and escaping operations in the encoding namespace.
func RegisterLib(ns runtime.Namespace) {
	ns.Namespace("encoding").Function().A1().Add("json_parse", JSONParse).Add("json_stringify", JSONStringify).
		Add("query_escape", QueryEscape).Add("query_unescape", QueryUnescape).
		Add("base64_encode", Base64Encode).Add("base64_decode", Base64Decode).
		Add("html_escape", HTMLEscape).Add("html_unescape", HTMLUnescape)
}
