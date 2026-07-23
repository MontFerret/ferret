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
// Policy.Eval validates an already constructed *net/http.Request without
// mutating it. Policy.Prepare adds any missing default headers to the request,
// then calls Eval; defaults added before a later validation failure remain on
// the request. Neither method accepts a separate context parameter. Attach a
// context when constructing the standard request or by calling its WithContext
// method.
//
// Eval accepts only strict outbound client-request state. RequestURI must be
// empty; Host, when set, must be canonically authority-equivalent to URL.Host;
// Close must be false; and TransferEncoding and Trailer must be empty. Reserved
// transport-controlled entries in Header are rejected. With a finite request
// body limit, a non-empty body must have a positive, known ContentLength; an
// unknown-length body is rejected before transport with RequestBodyLengthError.
//
// Eval and Prepare are request preflight checks. Embedders implementing a
// custom Client must integrate every policy stage exposed by Policy:
//
//   - call Prepare or Eval before the initial request;
//   - apply Timeout to the backend or request context;
//   - call CheckRedirect before following every redirect;
//   - call EvalConnection on the resolved address immediately before connect;
//   - apply MaxResponseHeaderSize before parsing response headers; and
//   - apply MaxResponseSize natively or materialize with ReadResponseBody.
//
// ReadResponseBody does not close its reader. The caller retains response-body
// ownership. Header limits cannot provide a memory bound when checked only
// after a backend has parsed the headers. Ferret's built-in Client supplies all
// of these controls; custom transports passed to NewWithClient or
// NewWithTransport retain the documented transport-level responsibilities.
//
// Policy.Eval accepting *net/http.Request instead of Ferret's *Request is an
// intentional source break. Call Prepare when migrating code that needs policy
// defaults, or Eval when defaults have already been applied.
//
// Policy construction is fallible, so callers should handle its error at
// startup:
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
// Apply defaults and validate a standard-library request before passing it to
// a custom client:
//
//	stdReq, err := nethttp.NewRequestWithContext(ctx, nethttp.MethodGet, targetURL, nil)
//	if err != nil {
//		return err
//	}
//	if err := policy.Prepare(stdReq); err != nil {
//		return err
//	}
//	response, err := customClient.Do(stdReq)
//
// When policy defaults are not needed, call policy.Eval(stdReq) instead to
// validate the request without modifying it.
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
