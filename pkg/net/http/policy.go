package http

import (
	"fmt"
	"io"
	stdhttp "net/http"
	"net/netip"
	"net/url"
	"sort"
	"strings"
	"time"
)

type (
	// Policy describes HTTP request policy and validation behavior. Its zero
	// value is a deny-all policy; use NewPolicy for Ferret's secure defaults.
	Policy struct {
		configurationErrors   []*PolicyConfigurationError
		defaultHeaders        map[string]string
		allowedSchemes        []string
		allowedMethods        []string
		allowedHosts          []string
		blockedHosts          []string
		blockedRequestHeaders []string
		defaultHeaderInputs   []defaultHeaderInput
		maxResponseSize       int64
		maxResponseHeaderSize int64
		maxRedirects          int
		maxRequestSize        int64
		timeout               time.Duration
		followRedirects       bool
		allowLocalhost        bool
		allowPrivateNetworks  bool
		allowLinkLocal        bool
	}

	defaultHeaderInput struct {
		option string
		key    string
		value  string
	}

	// PolicyOption configures a Policy during NewPolicy construction.
	PolicyOption func(*Policy)
)

// NewPolicy builds a reusable policy with Ferret's secure HTTP defaults.
// Construction returns a PolicyConfigurationError for one failure or a
// MultiPolicyConfigurationError for multiple failures. The zero value of
// Policy is intentionally deny-all; embedders should call NewPolicy when they
// want the standard defaults.
func NewPolicy(options ...PolicyOption) (*Policy, error) {
	p := &Policy{
		followRedirects: true,
		allowedSchemes:  []string{"http", "https"},
		allowedMethods: []string{
			stdhttp.MethodGet,
			stdhttp.MethodHead,
			stdhttp.MethodPost,
			stdhttp.MethodPut,
			stdhttp.MethodPatch,
			stdhttp.MethodDelete,
			stdhttp.MethodOptions,
		},
		timeout:               defaultTimeout,
		maxRedirects:          defaultMaxRedirects,
		maxRequestSize:        defaultMaxRequestSize,
		maxResponseSize:       defaultMaxResponseSize,
		maxResponseHeaderSize: defaultMaxResponseHeaderSize,
	}

	for _, option := range options {
		if option == nil {
			continue
		}

		option(p)
	}

	if err := p.validateConfiguration(); err != nil {
		return nil, err
	}

	return p, nil
}

func (p *Policy) validateConfiguration() error {
	if p.timeout < 0 {
		p.setConfigurationErrorIfMissing(
			"WithTimeout",
			p.timeout.String(),
			"must not be negative",
		)
	}

	if p.maxRedirects < 0 {
		p.setConfigurationErrorIfMissing(
			"WithMaxRedirects",
			fmt.Sprint(p.maxRedirects),
			"must not be negative",
		)
	}

	if p.maxRequestSize < 0 {
		p.setConfigurationErrorIfMissing(
			"WithMaxRequestSize",
			fmt.Sprint(p.maxRequestSize),
			"must not be negative",
		)
	}

	if p.maxResponseSize < 0 {
		p.setConfigurationErrorIfMissing(
			"WithMaxResponseSize",
			fmt.Sprint(p.maxResponseSize),
			"must not be negative",
		)
	}

	if p.maxResponseHeaderSize <= 0 {
		p.setConfigurationErrorIfMissing(
			"WithMaxResponseHeaderSize",
			fmt.Sprint(p.maxResponseHeaderSize),
			"must be positive",
		)
	}

	for _, scheme := range p.allowedSchemes {
		if err := validateConfiguredScheme(scheme); err != nil {
			p.setConfigurationErrorIfMissing("WithAllowedSchemes", scheme, err.Error())
		}
	}

	for _, method := range p.allowedMethods {
		if !isValidMethod(normalizeMethod(method)) {
			p.setConfigurationErrorIfMissing(
				"WithAllowedMethods",
				method,
				"must be a non-empty HTTP method token",
			)
		}
	}

	for _, host := range p.allowedHosts {
		if _, err := normalizeConfiguredHost(host); err != nil {
			p.setConfigurationErrorIfMissing("WithAllowedHosts", host, err.Error())
		}
	}

	for _, host := range p.blockedHosts {
		if _, err := normalizeConfiguredHost(host); err != nil {
			p.setConfigurationErrorIfMissing("WithBlockedHosts", host, err.Error())
		}
	}

	for _, header := range p.blockedRequestHeaders {
		if err := validateHeaderName(strings.TrimSpace(header)); err != nil {
			p.setConfigurationErrorIfMissing(
				"WithBlockedRequestHeaders",
				header,
				err.Reason,
			)
		}
	}

	p.allowedSchemes = normalizeValues(p.allowedSchemes)
	p.allowedMethods = normalizeMethods(p.allowedMethods)
	p.allowedHosts = normalizeHosts(p.allowedHosts)
	p.blockedHosts = normalizeHosts(p.blockedHosts)
	p.blockedRequestHeaders = normalizeHeaders(p.blockedRequestHeaders)

	p.normalizeDefaultHeaders()

	return newMultiPolicyConfigurationError(p.configurationErrors)
}

func (p *Policy) normalizeDefaultHeaders() {
	inputs := append([]defaultHeaderInput(nil), p.defaultHeaderInputs...)

	if len(p.defaultHeaders) > 0 {
		keys := make([]string, 0, len(p.defaultHeaders))

		for key := range p.defaultHeaders {
			keys = append(keys, key)
		}

		sort.Strings(keys)

		for _, key := range keys {
			inputs = append(inputs, defaultHeaderInput{
				option: "WithDefaultHeaders",
				key:    key,
				value:  p.defaultHeaders[key],
			})
		}
	}

	headers := make(map[string]string, len(inputs))
	sources := make(map[string]string, len(inputs))

	for _, input := range inputs {
		key := strings.TrimSpace(input.key)

		if err := validateHeaderName(key); err != nil {
			p.setConfigurationErrorIfMissing(input.option, input.key, err.Reason)

			continue
		}

		canonicalKey := stdhttp.CanonicalHeaderKey(key)

		if isReservedRequestHeader(canonicalKey) {
			p.setConfigurationErrorIfMissing(
				input.option,
				canonicalKey,
				"request header is reserved for the transport",
			)

			continue
		}

		if err := validateHeaderValue(canonicalKey, input.value); err != nil {
			p.setConfigurationErrorIfMissing(input.option, canonicalKey, err.Reason)

			continue
		}

		if value, exists := headers[canonicalKey]; exists {
			if value != input.value {
				p.setConfigurationError(
					input.option,
					canonicalKey,
					"conflicts with another default for the same header",
				)
			}

			continue
		}

		headers[canonicalKey] = input.value
		sources[canonicalKey] = input.option
	}

	keys := make([]string, 0, len(headers))

	for header := range headers {
		keys = append(keys, header)
	}

	sort.Strings(keys)

	for _, header := range keys {
		if p.isBlockedHeader(header) {
			p.setConfigurationError(
				sources[header],
				header,
				"default header is also configured as blocked",
			)
		}
	}

	p.defaultHeaders = headers
	p.defaultHeaderInputs = nil
}

func (p *Policy) addDefaultHeader(option, key, value string) {
	key = strings.TrimSpace(key)
	p.defaultHeaderInputs = append(p.defaultHeaderInputs, defaultHeaderInput{
		option: option,
		key:    key,
		value:  value,
	})

	if err := validateHeaderName(key); err != nil {
		p.setConfigurationError(option, key, err.Reason)

		return
	}

	canonicalKey := stdhttp.CanonicalHeaderKey(key)

	if isReservedRequestHeader(canonicalKey) {
		p.setConfigurationError(
			option,
			canonicalKey,
			"request header is reserved for the transport",
		)

		return
	}

	if err := validateHeaderValue(canonicalKey, value); err != nil {
		p.setConfigurationError(option, canonicalKey, err.Reason)
	}
}

func (p *Policy) setConfigurationError(option, value, reason string) {
	p.configurationErrors = append(
		p.configurationErrors,
		newPolicyConfigurationError(option, value, reason),
	)
}

func (p *Policy) setConfigurationErrorIfMissing(option, value, reason string) {
	for _, err := range p.configurationErrors {
		if err != nil && err.Option == option && err.Value == value && err.Reason == reason {
			return
		}
	}

	p.setConfigurationError(option, value, reason)
}

// Eval validates an outbound standard-library request against the policy
// without modifying it.
func (p *Policy) Eval(req *stdhttp.Request) error {
	return p.eval(req, PolicyTargetRequest)
}

// Prepare adds configured default headers that are missing from req, then
// validates the effective outbound request. Defaults added before a validation
// failure remain on req; no other request fields are modified.
func (p *Policy) Prepare(req *stdhttp.Request) error {
	if err := p.applyDefaults(req); err != nil {
		return err
	}

	return p.Eval(req)
}

// CheckRedirect validates a redirect using the callback contract expected by
// net/http.Client. When redirects are disabled it returns net/http.ErrUseLastResponse
// so callers can retain the redirect response.
func (p *Policy) CheckRedirect(req *stdhttp.Request, via []*stdhttp.Request) error {
	if req == nil {
		return ErrNilRequest
	}

	if !p.followRedirects {
		return stdhttp.ErrUseLastResponse
	}

	limit := p.maxRedirects
	if len(via) > limit {
		return &RedirectLimitError{Limit: limit}
	}

	return p.eval(req, PolicyTargetRedirect)
}

// EvalConnection validates an already-resolved concrete destination address.
// Callers must invoke it immediately before connecting so DNS rebinding cannot
// bypass request-time host validation.
func (p *Policy) EvalConnection(addr netip.Addr) error {
	subject := "destination address"

	if addr.IsValid() {
		subject = addressSubject(addr)
	}

	return p.validateAddress(PolicyTargetConnection, subject, addr)
}

// Timeout returns the overall request timeout. Zero means no policy timeout.
func (p *Policy) Timeout() time.Duration {
	return p.timeout
}

// MaxResponseSize returns the materialized response-body limit in bytes. Zero
// means response bodies are unlimited.
func (p *Policy) MaxResponseSize() int64 {
	return p.maxResponseSize
}

// MaxResponseHeaderSize returns the response-header limit in bytes. Backends
// must apply it before parsing response headers for the limit to bound memory.
func (p *Policy) MaxResponseHeaderSize() int64 {
	return p.maxResponseHeaderSize
}

// EvalResponseSize validates an observed or materialized response-body size.
// A negative size represents an unknown length and must instead be enforced
// while reading, for example with ReadResponseBody.
func (p *Policy) EvalResponseSize(size int64) error {
	limit := p.maxResponseSize

	if size < 0 || limit <= 0 || size <= limit {
		return nil
	}

	return &ResponseBodyLimitError{
		Size:  size,
		Limit: limit,
	}
}

// ReadResponseBody materializes a response body while enforcing the configured
// size limit. It does not close body; ownership remains with the caller.
func (p *Policy) ReadResponseBody(body io.Reader) ([]byte, error) {
	if body == nil {
		return nil, nil
	}

	limit := p.maxResponseSize
	if limit <= 0 {
		return io.ReadAll(body)
	}

	data, err := io.ReadAll(io.LimitReader(body, saturatedIncrement(limit)))

	if sizeErr := p.EvalResponseSize(int64(len(data))); sizeErr != nil {
		return nil, sizeErr
	}

	if err != nil {
		return nil, err
	}

	return data, nil
}

func (p *Policy) applyDefaults(req *stdhttp.Request) error {
	if req == nil {
		return ErrNilRequest
	}

	if len(p.defaultHeaders) == 0 {
		return nil
	}

	if req.Header == nil {
		req.Header = make(stdhttp.Header, len(p.defaultHeaders))
	}

	keys := make([]string, 0, len(p.defaultHeaders))
	for key := range p.defaultHeaders {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	for _, key := range keys {
		present := false
		for requestKey := range req.Header {
			if strings.EqualFold(requestKey, key) {
				present = true

				break
			}
		}

		if !present {
			req.Header[key] = []string{p.defaultHeaders[key]}
		}
	}

	return nil
}

func (p *Policy) eval(req *stdhttp.Request, target PolicyTarget) error {
	if req == nil {
		return ErrNilRequest
	}

	if err := p.validateMethod(req.Method, target); err != nil {
		return err
	}

	if err := p.validateURL(req.URL, target); err != nil {
		return err
	}

	if err := p.validateOutboundRequest(req, target); err != nil {
		return err
	}

	if err := p.validateRequestHeaders(req.Header, target); err != nil {
		return err
	}

	return p.validateRequestBody(req)
}

func (p *Policy) validateMethod(method string, target PolicyTarget) error {
	effectiveMethod := method
	if effectiveMethod == "" {
		effectiveMethod = stdhttp.MethodGet
	}

	if !isValidMethod(effectiveMethod) {
		return &InvalidMethodError{Method: method}
	}

	if !containsValue(p.allowedMethods, effectiveMethod) {
		return newPolicyError(
			target,
			fmt.Sprintf("method %q", effectiveMethod),
			"method is not allowed",
		)
	}

	return nil
}

func (p *Policy) validateOutboundRequest(req *stdhttp.Request, target PolicyTarget) error {
	if req.RequestURI != "" {
		return newPolicyError(
			target,
			"request URI",
			"must be empty for outbound requests",
		)
	}

	if req.Host != "" && !requestAuthoritiesEquivalent(req.Host, req.URL.Host) {
		return p.reservedRequestHeaderError(target, "Host")
	}

	if req.Close {
		return p.reservedRequestHeaderError(target, "Connection")
	}

	if len(req.TransferEncoding) > 0 {
		return p.reservedRequestHeaderError(target, "Transfer-Encoding")
	}

	if len(req.Trailer) > 0 {
		return p.reservedRequestHeaderError(target, "Trailer")
	}

	return nil
}

func (p *Policy) reservedRequestHeaderError(target PolicyTarget, header string) error {
	return newPolicyError(
		target,
		fmt.Sprintf("header %q", stdhttp.CanonicalHeaderKey(header)),
		"request header is reserved for the transport",
	)
}

func (p *Policy) validateRequestHeaders(headers stdhttp.Header, target PolicyTarget) error {
	var keyBuffer [8]string
	keys := keyBuffer[:0]

	if len(headers) > len(keyBuffer) {
		keys = make([]string, 0, len(headers))
	}

	for key := range headers {
		keys = append(keys, key)
	}

	sort.Strings(keys)

	for _, key := range keys {
		values := headers[key]
		if err := validateHeaderName(key); err != nil {
			return err
		}

		canonicalKey := stdhttp.CanonicalHeaderKey(key)
		if isReservedRequestHeader(canonicalKey) {
			return newPolicyError(
				target,
				fmt.Sprintf("header %q", canonicalKey),
				"request header is reserved for the transport",
			)
		}

		if p.isBlockedHeader(canonicalKey) {
			return newPolicyError(
				target,
				fmt.Sprintf("header %q", canonicalKey),
				"request header is not allowed",
			)
		}

		for _, value := range values {
			if err := validateHeaderValue(canonicalKey, value); err != nil {
				return err
			}
		}
	}

	return nil
}

func (p *Policy) validateRequestBody(req *stdhttp.Request) error {
	if p.maxRequestSize <= 0 || req.Body == nil || req.Body == stdhttp.NoBody {
		return nil
	}

	if req.ContentLength > p.maxRequestSize {
		return &RequestBodyLimitError{
			Size:  req.ContentLength,
			Limit: p.maxRequestSize,
		}
	}

	if req.ContentLength <= 0 {
		return &RequestBodyLengthError{
			ContentLength: req.ContentLength,
			Limit:         p.maxRequestSize,
		}
	}

	return nil
}

func (p *Policy) validateURL(u *url.URL, target PolicyTarget) error {
	if u == nil || *u == (url.URL{}) {
		return &URLValidationError{Field: "url", Reason: "is required"}
	}

	if u.Scheme == "" {
		return &URLValidationError{Field: "scheme", Reason: "is required"}
	}

	if u.Host == "" {
		return &URLValidationError{Field: "host", Reason: "is required"}
	}

	if u.User != nil {
		return newPolicyError(target, "URL credentials", "URL user information is not allowed")
	}

	scheme := asciiLower(u.Scheme)
	if !containsValue(p.allowedSchemes, scheme) {
		return newPolicyError(target, fmt.Sprintf("scheme %q", u.Scheme), "scheme is not allowed")
	}

	rawHostname := u.Hostname()

	if !isASCII(rawHostname) {
		return newPolicyError(
			target,
			fmt.Sprintf("host %q", rawHostname),
			"internationalized hostnames must use ASCII/punycode",
		)
	}

	hostname := canonicalHostKey(rawHostname)

	if containsHost(p.blockedHosts, u) {
		return newPolicyError(target, fmt.Sprintf("host %q", hostname), "host is blocked")
	}

	if len(p.allowedHosts) > 0 && !containsHost(p.allowedHosts, u) {
		return newPolicyError(target, fmt.Sprintf("host %q", hostname), "host is not allowed")
	}

	if isLocalhostName(hostname) && !p.allowLocalhost {
		return newPolicyError(target, fmt.Sprintf("host %q", hostname), "localhost is not allowed")
	}

	if addr, ok := parseIPAddress(hostname); ok {
		return p.validateAddress(target, addressSubject(addr), addr)
	}

	return nil
}

func (p *Policy) validateAddress(target PolicyTarget, subject string, addr netip.Addr) error {
	if !addr.IsValid() {
		return newPolicyError(target, subject, "invalid address is not allowed")
	}

	addr = addr.Unmap()
	if embedded, ok := wellKnownNAT64IPv4(addr); ok {
		return p.validateAddress(target, subject, embedded)
	}

	if addr.IsLoopback() {
		if p.allowLocalhost {
			return nil
		}

		return newPolicyError(target, subject, "localhost is not allowed")
	}

	if addr.IsPrivate() || carrierGradeNAT.Contains(addr) {
		if p.allowPrivateNetworks {
			return nil
		}

		return newPolicyError(target, subject, "private networks are not allowed")
	}

	if addr.IsLinkLocalUnicast() {
		if p.allowLinkLocal {
			return nil
		}

		return newPolicyError(target, subject, "link-local addresses are not allowed")
	}

	if addr.IsUnspecified() || ipv4CurrentNetwork.Contains(addr) {
		return newPolicyError(target, subject, "unspecified addresses are not allowed")
	}

	if addr.IsMulticast() {
		return newPolicyError(target, subject, "multicast addresses are not allowed")
	}

	if ipv4Reserved.Contains(addr) || ipv6SiteLocal.Contains(addr) {
		return newPolicyError(target, subject, "reserved addresses are not allowed")
	}

	if isAlwaysBlockedAddress(addr) || !addr.IsGlobalUnicast() {
		return newPolicyError(target, subject, "non-public addresses are not allowed")
	}

	return nil
}

func (p *Policy) isBlockedHeader(key string) bool {
	return containsValue(p.blockedRequestHeaders, stdhttp.CanonicalHeaderKey(key))
}
