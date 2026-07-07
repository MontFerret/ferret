package http

// Response contains the materialized HTTP response.
type Response struct {
	Headers    Headers `json:"headers"`
	Status     string  `json:"status"`
	Body       []byte  `json:"body,omitempty"`
	StatusCode int     `json:"statusCode"`
}
