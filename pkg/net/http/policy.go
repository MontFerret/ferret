package http

import (
	"bytes"
	"context"
	"fmt"
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
		configurationError    error
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
// Construction fails when an option is malformed or contradictory. The zero
// value of Policy is intentionally deny-all; embedders should call NewPolicy
// when they want the standard defaults.
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
	if p.configurationError != nil {
		return p.configurationError
	}

	if p.timeout < 0 {
		return newPolicyConfigurationError(
			"WithTimeout",
			p.timeout.String(),
			"must not be negative",
		)
	}

	if p.maxRedirects < 0 {
		return newPolicyConfigurationError(
			"WithMaxRedirects",
			fmt.Sprint(p.maxRedirects),
			"must not be negative",
		)
	}

	if p.maxRequestSize < 0 {
		return newPolicyConfigurationError(
			"WithMaxRequestSize",
			fmt.Sprint(p.maxRequestSize),
			"must not be negative",
		)
	}

	if p.maxResponseSize < 0 {
		return newPolicyConfigurationError(
			"WithMaxResponseSize",
			fmt.Sprint(p.maxResponseSize),
			"must not be negative",
		)
	}

	if p.maxResponseHeaderSize <= 0 {
		return newPolicyConfigurationError(
			"WithMaxResponseHeaderSize",
			fmt.Sprint(p.maxResponseHeaderSize),
			"must be positive",
		)
	}

	for _, scheme := range p.allowedSchemes {
		if err := validateConfiguredScheme(scheme); err != nil {
			return newPolicyConfigurationError("WithAllowedSchemes", scheme, err.Error())
		}
	}

	for _, method := range p.allowedMethods {
		if !isValidMethod(normalizeMethod(method)) {
			return newPolicyConfigurationError(
				"WithAllowedMethods",
				method,
				"must be a non-empty HTTP method token",
			)
		}
	}

	for _, host := range p.allowedHosts {
		if _, err := normalizeConfiguredHost(host); err != nil {
			return newPolicyConfigurationError("WithAllowedHosts", host, err.Error())
		}
	}

	for _, host := range p.blockedHosts {
		if _, err := normalizeConfiguredHost(host); err != nil {
			return newPolicyConfigurationError("WithBlockedHosts", host, err.Error())
		}
	}

	for _, header := range p.blockedRequestHeaders {
		if err := validateHeaderName(strings.TrimSpace(header)); err != nil {
			return newPolicyConfigurationError(
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

	return p.normalizeDefaultHeaders()
}

func (p *Policy) normalizeDefaultHeaders() error {
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
			return newPolicyConfigurationError(input.option, input.key, err.Reason)
		}

		canonicalKey := stdhttp.CanonicalHeaderKey(key)

		if isReservedRequestHeader(canonicalKey) {
			return newPolicyConfigurationError(
				input.option,
				canonicalKey,
				"request header is reserved for the transport",
			)
		}

		if err := validateHeaderValue(canonicalKey, input.value); err != nil {
			return newPolicyConfigurationError(input.option, canonicalKey, err.Reason)
		}

		if value, exists := headers[canonicalKey]; exists {
			if value != input.value {
				return newPolicyConfigurationError(
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
			return newPolicyConfigurationError(
				sources[header],
				header,
				"default header is also configured as blocked",
			)
		}
	}

	p.defaultHeaders = headers
	p.defaultHeaderInputs = nil

	return nil
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
	if p.configurationError != nil {
		return
	}

	p.configurationError = newPolicyConfigurationError(option, value, reason)
}

// Eval validates an outbound request against the policy.
func (p *Policy) Eval(req *Request) error {
	_, err := p.prepareRequest(context.Background(), req)

	return err
}

func (p *Policy) prepareRequest(ctx context.Context, req *Request) (*stdhttp.Request, error) {
	if req == nil {
		return nil, ErrNilRequest
	}

	if err := p.validateMethod(req.Method, PolicyTargetRequest); err != nil {
		return nil, err
	}

	rawURL := strings.TrimSpace(req.URL)
	if rawURL == "" {
		return nil, &URLValidationError{Field: "url", Reason: "is required"}
	}

	method := normalizeRequestMethod(req.Method)
	stdReq, err := stdhttp.NewRequestWithContext(ctx, method, rawURL, bytes.NewReader(req.Body))

	if err != nil {
		if ctx == nil {
			return nil, &RequestBuildError{Err: err}
		}

		return nil, &URLParseError{Err: err}
	}

	if stdReq.URL.Scheme == "" {
		return nil, &URLValidationError{Field: "scheme", Reason: "is required"}
	}

	if stdReq.URL.Host == "" {
		return nil, &URLValidationError{Field: "host", Reason: "is required"}
	}

	stdReq.URL.Scheme = asciiLower(stdReq.URL.Scheme)
	stdReq.URL.Host = asciiLower(stdReq.URL.Host)

	if err := p.validateURL(stdReq.URL, PolicyTargetRequest); err != nil {
		return nil, err
	}

	headers, err := p.copyValidatedRequestHeaders(req.Headers, PolicyTargetRequest)
	if err != nil {
		return nil, err
	}

	stdReq.Header = headers

	for key, value := range p.defaultHeaders {
		if _, exists := stdReq.Header[key]; !exists {
			stdReq.Header.Set(key, value)
		}
	}

	if p.maxRequestSize > 0 && int64(len(req.Body)) > p.maxRequestSize {
		return nil, &RequestBodyLimitError{
			Size:  int64(len(req.Body)),
			Limit: p.maxRequestSize,
		}
	}

	return stdReq, nil
}

func (p *Policy) validateMethod(method string, target PolicyTarget) error {
	normalized := normalizeRequestMethod(method)

	if !isValidMethod(normalized) {
		return &InvalidMethodError{Method: method}
	}

	if !containsValue(p.allowedMethods, normalized) {
		return newPolicyError(target, fmt.Sprintf("method %q", normalized), "method is not allowed")
	}

	return nil
}

func (p *Policy) validateRequestHeaders(headers map[string][]string, target PolicyTarget) error {
	_, err := p.copyValidatedRequestHeaders(headers, target)

	return err
}

func (p *Policy) copyValidatedRequestHeaders(
	headers map[string][]string,
	target PolicyTarget,
) (stdhttp.Header, error) {
	result := make(stdhttp.Header, len(headers))
	keys := make([]string, 0, len(headers))

	for key := range headers {
		keys = append(keys, key)
	}

	sort.Strings(keys)

	for _, key := range keys {
		values := headers[key]
		if err := validateHeaderName(key); err != nil {
			return nil, err
		}

		canonicalKey := stdhttp.CanonicalHeaderKey(key)
		if isReservedRequestHeader(canonicalKey) {
			return nil, newPolicyError(
				target,
				fmt.Sprintf("header %q", canonicalKey),
				"request header is reserved for the transport",
			)
		}

		if p.isBlockedHeader(canonicalKey) {
			return nil, newPolicyError(
				target,
				fmt.Sprintf("header %q", canonicalKey),
				"request header is not allowed",
			)
		}

		for _, value := range values {
			if err := validateHeaderValue(canonicalKey, value); err != nil {
				return nil, err
			}
		}

		result[canonicalKey] = append(result[canonicalKey], values...)
	}

	return result, nil
}

func (p *Policy) validateURL(u *url.URL, target PolicyTarget) error {
	if u == nil {
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
