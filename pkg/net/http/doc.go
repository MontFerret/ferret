// Package http provides Ferret's sandboxed HTTP client and its reusable access
// policy.
//
// NewPolicy validates all options and applies secure defaults: HTTP and HTTPS
// to public destinations; GET, HEAD, POST, PUT, PATCH, DELETE, and OPTIONS; a
// 30-second timeout; 10 followed redirects; 16 MiB request and response bodies;
// and 1 MiB response headers. Localhost, private, and link-local destinations
// require explicit opt-ins. The zero value of Policy is deny-all.
//
// Numeric option values below zero are invalid, zero restores the documented
// default, and positive values override it. WithNoTimeout,
// WithUnlimitedRequestSize, and WithUnlimitedResponseSize are the only ways to
// disable their corresponding limits. Response headers always remain bounded.
//
// Policies reject URL credentials, malformed headers, and transport-controlled
// request headers including Host and Content-Length. Internationalized hostnames
// must be supplied in explicit ASCII/punycode form. Host rules use exact
// normalized matching and never implicitly include subdomains. Original URLs,
// redirects, DNS results, and concrete connection addresses are validated;
// dial-time denials use PolicyTargetConnection. Ambient proxy configuration is
// disabled.
//
// Embedders may reuse Policy.Eval with their own Client implementation. Policy
// construction is fallible, so callers should handle its error at startup:
//
//	policy, err := http.NewPolicy(
//		http.WithAllowedHosts("api.example.com"),
//		http.WithMaxResponseSize(32 << 20),
//	)
//	if err != nil {
//		if errors.Is(err, http.ErrInvalidPolicyConfiguration) {
//			var issue *http.PolicyConfigurationError
//			if errors.As(err, &issue) {
//				log.Printf("invalid %s: %s", issue.Option, issue.Reason)
//			}
//		}
//		return err
//	}
//
// A single configuration failure is returned as PolicyConfigurationError;
// multiple failures are returned as MultiPolicyConfigurationError. The
// aggregate unwraps its individual failures, so errors.As can inspect either
// the aggregate or its first PolicyConfigurationError. Its Errors field
// provides all failures in deterministic validation order.
//
// Structural, configuration, limit, and policy failures are typed. Callers
// should inspect them with errors.Is and errors.As; error strings are intended
// for safe diagnostics and omit header values.
//
//	var policyErr *http.PolicyError
//	if errors.As(err, &policyErr) {
//		log.Printf(
//			"HTTP policy denied %s %q: %s",
//			policyErr.Target,
//			policyErr.Subject,
//			policyErr.Reason,
//		)
//	}
package http
