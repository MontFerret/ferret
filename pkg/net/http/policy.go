package http

import (
	"fmt"
	"net"
	stdhttp "net/http"
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
	}
)

// NewPolicies builds a policy set with Ferret's default HTTP policy values.
func NewPolicies(setters ...Policy) Policies {
	p := Policies{
		FollowRedirects:      true,
		AllowedSchemes:       []string{"http", "https"},
		AllowLocalhost:       true,
		AllowPrivateNetworks: true,
	}

	for _, setter := range setters {
		if setter == nil {
			continue
		}

		setter(&p)
	}

	p.AllowedSchemes = normalizeValues(p.AllowedSchemes)
	p.AllowedHosts = normalizeValues(p.AllowedHosts)
	p.BlockedHosts = normalizeValues(p.BlockedHosts)
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

	if err := p.validateURL(u); err != nil {
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

func (p *Policies) validateURL(u *url.URL) error {
	scheme := strings.ToLower(u.Scheme)
	if !containsValue(p.AllowedSchemes, scheme) {
		return fmt.Errorf("http: scheme %q is not allowed", u.Scheme)
	}

	hostname := strings.ToLower(u.Hostname())
	if hostname == "" {
		return fmt.Errorf("http: url host is required")
	}

	if containsHost(p.BlockedHosts, u) {
		return fmt.Errorf("http: host %q is blocked", hostname)
	}

	if len(p.AllowedHosts) > 0 && !containsHost(p.AllowedHosts, u) {
		return fmt.Errorf("http: host %q is not allowed", hostname)
	}

	ip := net.ParseIP(hostname)
	if isLocalhost(hostname, ip) && !p.AllowLocalhost {
		return fmt.Errorf("http: localhost is not allowed")
	}

	if ip != nil && ip.IsPrivate() && !p.AllowPrivateNetworks {
		return fmt.Errorf("http: private network address %q is not allowed", hostname)
	}

	return nil
}

func (p *Policies) isBlockedHeader(key string) bool {
	return containsValue(p.BlockedRequestHeaders, stdhttp.CanonicalHeaderKey(key))
}

func (p *Policies) checkRedirect(req *stdhttp.Request, via []*stdhttp.Request) error {
	if !p.FollowRedirects {
		return stdhttp.ErrUseLastResponse
	}

	if p.MaxRedirects > 0 && len(via) >= p.MaxRedirects {
		return fmt.Errorf("http: stopped after %d redirect(s)", p.MaxRedirects)
	}

	return p.validateURL(req.URL)
}
