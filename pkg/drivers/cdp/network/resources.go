package network

import "github.com/mafredri/cdp/protocol/network"

var (
	resourceTypeMapping = map[string]network.ResourceType{
		"document":           network.ResourceTypeDocument,
		"stylesheet":         network.ResourceTypeStylesheet,
		"css":                network.ResourceTypeStylesheet,
		"image":              network.ResourceTypeImage,
		"media":              network.ResourceTypeMedia,
		"font":               network.ResourceTypeFont,
		"script":             network.ResourceTypeScript,
		"js":                 network.ResourceTypeScript,
		"texttrack":          network.ResourceTypeTextTrack,
		"xhr":                network.ResourceTypeXHR,
		"ajax":               network.ResourceTypeXHR,
		"fetch":              network.ResourceTypeFetch,
		"eventsource":        network.ResourceTypeEventSource,
		"websocket":          network.ResourceTypeWebSocket,
		"manifest":           network.ResourceTypeManifest,
		"sxg":                network.ResourceTypeSignedExchange,
		"ping":               network.ResourceTypePing,
		"cspViolationReport": network.ResourceTypeCSPViolationReport,
		"other":              network.ResourceTypeOther,
	}
)

func toResourceType(alias string) network.ResourceType {
	rt, found := resourceTypeMapping[alias]

	if found {
		return rt
	}

	return network.ResourceTypeNotSet
}
