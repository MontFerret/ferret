package http

// PolicyTarget identifies which outbound request stage was rejected by the
// access policy.
type PolicyTarget string

const (
	// PolicyTargetRequest identifies the original outbound request.
	PolicyTargetRequest PolicyTarget = "request"

	// PolicyTargetRedirect identifies a redirect destination.
	PolicyTargetRedirect PolicyTarget = "redirect destination"

	// PolicyTargetConnection identifies the concrete address selected while
	// establishing an outbound connection.
	PolicyTargetConnection PolicyTarget = "connection destination"
)
