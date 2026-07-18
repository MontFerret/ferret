package http

import (
	"context"
	"fmt"
	stdhttp "net/http"
)

type (
	// Client executes HTTP requests for Ferret host integrations.
	Client interface {
		Do(ctx context.Context, req *Request) (*Response, error)
	}

	// IdleConnectionCloser is an optional capability implemented by clients that
	// can release pooled idle connections. Standalone callers may type-assert a
	// Client to this interface when deterministic cleanup is needed. Calls are
	// safe to repeat and do not interrupt active requests.
	IdleConnectionCloser interface {
		CloseIdleConnections()
	}

	defaultHTTPClient struct {
		policy    *Policies
		transport stdhttp.Client
	}
)

// New constructs an HTTP client with the provided policies.
func New(setters ...Policy) Client {
	policies := NewPolicies(setters...)
	dialer := newPolicyDialer(policies)

	return &defaultHTTPClient{
		policy: policies,
		transport: stdhttp.Client{
			Transport: newPolicyTransport(dialer),
		},
	}
}

func (d *defaultHTTPClient) Do(ctx context.Context, req *Request) (*Response, error) {
	if ctx == nil {
		ctx = context.Background()
	}

	if req == nil {
		return nil, ErrNilRequest
	}

	p := d.policy
	if p == nil {
		policies := NewPolicies()
		p = policies
	}

	stdReq, err := toStdRequest(ctx, req, p)
	if err != nil {
		return nil, err
	}

	client := d.transport

	if p.timeout > 0 {
		client.Timeout = p.timeout
	}

	client.CheckRedirect = d.checkRedirect

	res, err := client.Do(stdReq)
	if err != nil {
		return nil, err
	}

	return fromStdResponse(res, p)
}

func (d *defaultHTTPClient) CloseIdleConnections() {
	d.transport.CloseIdleConnections()
}

func (d *defaultHTTPClient) checkRedirect(req *stdhttp.Request, via []*stdhttp.Request) error {
	p := d.policy

	if p == nil {
		policies := NewPolicies()
		p = policies
	}

	if !p.followRedirects {
		return stdhttp.ErrUseLastResponse
	}

	limit := p.maxRedirects

	if limit == 0 {
		limit = 10
	}

	if len(via) >= limit {
		return fmt.Errorf("http: stopped after %d redirect(s)", limit)
	}

	return p.validateURL(req.URL, policyTargetRedirect)
}
