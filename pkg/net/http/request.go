package http

// Request describes an outbound HTTP request.
type Request struct {
	Method  string  `json:"method"`
	URL     string  `json:"url"`
	Headers Headers `json:"headers"`
	Body    []byte  `json:"body,omitempty"`
}
