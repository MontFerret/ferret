// Package http provides the sandboxed HTTP client used by Ferret host
// integrations.
//
// The built-in client permits HTTP and HTTPS requests to public destinations
// and validates original URLs, redirect destinations, DNS results, and the
// concrete address selected for each connection. Ambient proxy configuration
// is disabled. Localhost, private, and link-local destinations require explicit
// policy opt-ins.
//
// The default method allowlist is GET, HEAD, POST, PUT, PATCH, DELETE, and
// OPTIONS. URL credentials and transport-reserved headers are rejected,
// response headers are limited to 1 MiB, pooled connections are bounded, and
// requests have a 30-second overall timeout. Internationalized hostnames must
// be supplied in ASCII/punycode form. Host policy entries use exact normalized
// matching; they do not implicitly match subdomains.
package http
