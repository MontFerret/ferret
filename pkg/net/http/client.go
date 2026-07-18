package http

import (
	"context"
	"errors"
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
		policy *Policy
		client stdhttp.Client
	}
)

// New constructs an HTTP client with the provided policy options.
func New(options ...PolicyOption) Client {
	policy := NewPolicy(options...)
	dialer := newPolicyDialer(policy)

	return &defaultHTTPClient{
		policy: policy,
		client: stdhttp.Client{
			Transport: newPolicyTransport(dialer, policy.maxResponseHeaderSize),
			Timeout:   policy.timeout,
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
		p = NewPolicy()
	}

	stdReq, err := toStdRequest(ctx, req, p)
	if err != nil {
		return nil, err
	}

	client := d.client
	client.Timeout = p.timeout
	client.CheckRedirect = d.checkRedirect

	res, err := client.Do(stdReq)
	if err != nil {
		var policyErr *PolicyError
		if errors.As(err, &policyErr) {
			return nil, policyErr
		}

		return nil, err
	}

	return fromStdResponse(res, p)
}

func (d *defaultHTTPClient) CloseIdleConnections() {
	d.client.CloseIdleConnections()
}

func (d *defaultHTTPClient) checkRedirect(req *stdhttp.Request, via []*stdhttp.Request) error {
	p := d.policy

	if p == nil {
		p = NewPolicy()
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

	if err := p.validateMethod(req.Method, PolicyTargetRedirect); err != nil {
		return err
	}

	if err := p.validateURL(req.URL, PolicyTargetRedirect); err != nil {
		return err
	}

	return p.validateRequestHeaders(req.Header, PolicyTargetRedirect)
}
