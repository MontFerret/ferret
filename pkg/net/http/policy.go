package http

import (
	"fmt"
	stdhttp "net/http"
	"net/netip"
	"net/url"
	"strings"
	"time"
)

type (
	// Policy describes HTTP request policy and validation behavior.
	Policy struct {
		defaultHeaders        map[string]string
		allowedSchemes        []string
		allowedMethods        []string
		allowedHosts          []string
		blockedHosts          []string
		blockedRequestHeaders []string
		timeout               time.Duration
		maxRequestSize        int64
		maxResponseSize       int64
		maxResponseHeaderSize int64
		maxRedirects          int
		followRedirects       bool
		allowLocalhost        bool
		allowPrivateNetworks  bool
		allowLinkLocal        bool
	}

	// PolicyOption configures a Policy.
	PolicyOption func(*Policy)
)

// NewPolicy builds a policy with Ferret's secure HTTP defaults. Requests are
// limited to public destinations; GET, HEAD, POST, PUT, PATCH, DELETE, and
// OPTIONS; ASCII/punycode hosts; 1 MiB response headers; and a 30-second
// overall timeout. Localhost, private, and link-local destinations require
// their corresponding opt-in.
func NewPolicy(options ...PolicyOption) *Policy {
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
		maxResponseHeaderSize: defaultMaxResponseHeaderSize,
	}

	for _, option := range options {
		if option == nil {
			continue
		}

		option(p)
	}

	p.allowedSchemes = normalizeValues(p.allowedSchemes)
	p.allowedMethods = normalizeMethods(p.allowedMethods)
	p.allowedHosts = normalizeHosts(p.allowedHosts)
	p.blockedHosts = normalizeHosts(p.blockedHosts)
	p.blockedRequestHeaders = normalizeHeaders(p.blockedRequestHeaders)

	if len(p.defaultHeaders) > 0 {
		headers := make(map[string]string, len(p.defaultHeaders))

		for key, value := range p.defaultHeaders {
			key = strings.TrimSpace(key)

			if key == "" {
				continue
			}

			headers[stdhttp.CanonicalHeaderKey(key)] = value
		}

		p.defaultHeaders = headers
	}

	return p
}

// Eval validates an outbound request against the policy.
func (p *Policy) Eval(req *Request) error {
	if req == nil {
		return ErrNilRequest
	}

	if err := p.validateMethod(req.Method, PolicyTargetRequest); err != nil {
		return err
	}

	u, err := parseRequestURL(req.URL)
	if err != nil {
		return err
	}

	if err := p.validateURL(u, PolicyTargetRequest); err != nil {
		return err
	}

	if err := p.validateRequestHeaders(req.Headers, PolicyTargetRequest); err != nil {
		return err
	}

	if p.maxRequestSize > 0 && int64(len(req.Body)) > p.maxRequestSize {
		return &RequestBodyLimitError{
			Size:  int64(len(req.Body)),
			Limit: p.maxRequestSize,
		}
	}

	return nil
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
	for key := range headers {
		key = strings.TrimSpace(key)

		if key == "" {
			continue
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
	}

	return nil
}

func (p *Policy) validateURL(u *url.URL, target PolicyTarget) error {
	if u.User != nil {
		return newPolicyError(target, "URL credentials", "URL user information is not allowed")
	}

	scheme := asciiLower(u.Scheme)
	if !containsValue(p.allowedSchemes, scheme) {
		return newPolicyError(target, fmt.Sprintf("scheme %q", u.Scheme), "scheme is not allowed")
	}

	rawHostname := u.Hostname()
	if rawHostname == "" {
		return &URLValidationError{Field: "host", Reason: "is required"}
	}

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
