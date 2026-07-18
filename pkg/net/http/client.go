package http

import (
	"context"
	"errors"
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

// New constructs an HTTP client with the provided policy options. Invalid
// configuration returns a PolicyConfigurationError for one failure or a
// MultiPolicyConfigurationError for multiple failures. Both match
// ErrInvalidPolicyConfiguration with errors.Is and expose details through
// errors.As.
func New(options ...PolicyOption) (Client, error) {
	policy, err := NewPolicy(options...)
	if err != nil {
		return nil, err
	}

	dialer := newPolicyDialer(policy)
	transport := newResponseValidatingTransport(
		newPolicyTransport(dialer, policy.maxResponseHeaderSize),
	)

	return &defaultHTTPClient{
		policy: policy,
		client: stdhttp.Client{
			Transport: transport,
			Timeout:   policy.timeout,
		},
	}, nil
}

func (d *defaultHTTPClient) Do(ctx context.Context, req *Request) (*Response, error) {
	if ctx == nil {
		ctx = context.Background()
	}

	p := d.policy
	if p == nil {
		p = &Policy{}
	}

	stdReq, err := toStdRequest(ctx, req)
	if err != nil {
		return nil, err
	}
	if err := p.Prepare(stdReq); err != nil {
		return nil, err
	}

	client := d.client
	client.Transport = newResponseValidatingTransport(client.Transport)
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
		p = &Policy{}
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
