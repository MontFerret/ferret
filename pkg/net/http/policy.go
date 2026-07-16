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
	// Policy configures HTTP client behavior.
	Policy func(*Policies)

	// Policies describes HTTP request policy and validation behavior.
	Policies struct {
		DefaultHeaders        map[string]string
		AllowedSchemes        []string
		AllowedHosts          []string
		BlockedHosts          []string
		BlockedRequestHeaders []string
		Timeout               time.Duration
		MaxRequestSize        int64
		MaxResponseSize       int64
		MaxRedirects          int
		FollowRedirects       bool
		AllowLocalhost        bool
		AllowPrivateNetworks  bool
		AllowLinkLocal        bool
	}
)

// NewPolicies builds a policy set with Ferret's default HTTP policy values.
// By default, requests are limited to public destinations; localhost, private,
// and link-local destinations require their corresponding explicit opt-in.
func NewPolicies(setters ...Policy) Policies {
	p := Policies{
		FollowRedirects: true,
		AllowedSchemes:  []string{"http", "https"},
	}

	for _, setter := range setters {
		if setter == nil {
			continue
		}

		setter(&p)
	}

	p.AllowedSchemes = normalizeValues(p.AllowedSchemes)
	p.AllowedHosts = normalizeHosts(p.AllowedHosts)
	p.BlockedHosts = normalizeHosts(p.BlockedHosts)
	p.BlockedRequestHeaders = normalizeHeaders(p.BlockedRequestHeaders)

	if len(p.DefaultHeaders) > 0 {
		headers := make(map[string]string, len(p.DefaultHeaders))
		for key, value := range p.DefaultHeaders {
			key = strings.TrimSpace(key)
			if key == "" {
				continue
			}

			headers[stdhttp.CanonicalHeaderKey(key)] = value
		}
		p.DefaultHeaders = headers
	}

	return p
}

// Eval validates an outbound request against the policy.
func (p Policies) Eval(req *Request) error {
	if req == nil {
		return ErrNilRequest
	}

	method := strings.TrimSpace(req.Method)
	if method == "" {
		method = stdhttp.MethodGet
	}

	if !isValidMethod(method) {
		return fmt.Errorf("http: invalid method %q", req.Method)
	}

	u, err := parseRequestURL(req.URL)
	if err != nil {
		return err
	}

	if err := p.validateURL(u, policyTargetRequest); err != nil {
		return err
	}

	if p.MaxRequestSize > 0 && int64(len(req.Body)) > p.MaxRequestSize {
		return fmt.Errorf(
			"http: request body exceeds limit: %d > %d",
			len(req.Body),
			p.MaxRequestSize,
		)
	}

	return nil
}

func (p *Policies) validateURL(u *url.URL, target policyTarget) error {
	scheme := strings.ToLower(u.Scheme)
	if !containsValue(p.AllowedSchemes, scheme) {
		return newPolicyError(target, fmt.Sprintf("scheme %q", u.Scheme), "scheme is not allowed")
	}

	hostname := canonicalHostKey(u.Hostname())
	if hostname == "" {
		return fmt.Errorf("http: url host is required")
	}

	if containsHost(p.BlockedHosts, u) {
		return newPolicyError(target, fmt.Sprintf("host %q", hostname), "host is blocked")
	}

	if len(p.AllowedHosts) > 0 && !containsHost(p.AllowedHosts, u) {
		return newPolicyError(target, fmt.Sprintf("host %q", hostname), "host is not allowed")
	}

	if isLocalhostName(hostname) && !p.AllowLocalhost {
		return newPolicyError(target, fmt.Sprintf("host %q", hostname), "localhost is not allowed")
	}

	if addr, ok := parseIPAddress(hostname); ok {
		return p.validateAddress(target, addressSubject(addr), addr)
	}

	return nil
}

func (p *Policies) validateAddress(target policyTarget, subject string, addr netip.Addr) error {
	if !addr.IsValid() {
		return newPolicyError(target, subject, "invalid address is not allowed")
	}

	addr = addr.Unmap()
	if embedded, ok := wellKnownNAT64IPv4(addr); ok {
		return p.validateAddress(target, subject, embedded)
	}

	if addr.IsLoopback() {
		if p.AllowLocalhost {
			return nil
		}

		return newPolicyError(target, subject, "localhost is not allowed")
	}

	if addr.IsPrivate() || carrierGradeNAT.Contains(addr) {
		if p.AllowPrivateNetworks {
			return nil
		}

		return newPolicyError(target, subject, "private networks are not allowed")
	}

	if addr.IsLinkLocalUnicast() {
		if p.AllowLinkLocal {
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

func (p *Policies) isBlockedHeader(key string) bool {
	return containsValue(p.BlockedRequestHeaders, stdhttp.CanonicalHeaderKey(key))
}
