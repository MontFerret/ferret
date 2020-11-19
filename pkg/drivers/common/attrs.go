package common

import "strings"

var Attributes = []string{
	"abbr",
	"accept",
	"accept-charset",
	"accesskey",
	"action",
	"allowfullscreen",
	"allowpaymentrequest",
	"allowusermedia",
	"alt",
	"as",
	"async",
	"autocomplete",
	"autofocus",
	"autoplay",
	"challenge",
	"charset",
	"checked",
	"cite",
	"class",
	"color",
	"cols",
	"colspan",
	"command",
	"content",
	"contenteditable",
	"contextmenu",
	"controls",
	"coords",
	"crossorigin",
	"data",
	"datetime",
	"default",
	"defer",
	"dir",
	"dirname",
	"disabled",
	"download",
	"draggable",
	"dropzone",
	"enctype",
	"for",
	"form",
	"formaction",
	"formenctype",
	"formmethod",
	"formnovalidate",
	"formtarget",
	"headers",
	"height",
	"hidden",
	"high",
	"href",
	"hreflang",
	"http-equiv",
	"icon",
	"id",
	"inputmode",
	"integrity",
	"is",
	"ismap",
	"itemid",
	"itemprop",
	"itemref",
	"itemscope",
	"itemtype",
	"keytype",
	"kind",
	"label",
	"lang",
	"list",
	"loop",
	"low",
	"manifest",
	"max",
	"maxlength",
	"media",
	"mediagroup",
	"method",
	"min",
	"minlength",
	"multiple",
	"muted",
	"name",
	"nomodule",
	"nonce",
	"novalidate",
	"open",
	"optimum",
	"pattern",
	"ping",
	"placeholder",
	"playsinline",
	"poster",
	"preload",
	"radiogroup",
	"readonly",
	"referrerpolicy",
	"rel",
	"required",
	"reversed",
	"rows",
	"rowspan",
	"sandbox",
	"spellcheck",
	"scope",
	"scoped",
	"seamless",
	"selected",
	"shape",
	"size",
	"sizes",
	"sortable",
	"sorted",
	"slot",
	"span",
	"spellcheck",
	"src",
	"srcdoc",
	"srclang",
	"srcset",
	"start",
	"step",
	"style",
	"tabindex",
	"target",
	"title",
	"translate",
	"type",
	"typemustmatch",
	"updateviacache",
	"usemap",
	"value",
	"width",
	"workertype",
	"wrap",
}

const (
	ATTR_NAME_STYLE = "style"
)

var attrMap = make(map[string]bool)

func init() {
	for _, attr := range Attributes {
		attrMap[attr] = true
	}
}

func IsAttribute(name string) bool {
	_, isDefault := attrMap[name]

	if isDefault {
		return true
	}

	return strings.HasPrefix(name, "data-") ||
		strings.HasPrefix(name, "aria-")
}
