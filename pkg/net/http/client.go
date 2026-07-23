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

	return newDefaultHTTPClient(policy, stdhttp.Client{
		Transport: newPolicyTransport(dialer, policy.MaxResponseHeaderSize()),
	}), nil
}

// NewWithClient constructs a policy-aware client from a standard-library
// client. It snapshots the supplied client's fields without mutating it.
// Policy timeout and redirect settings take precedence over the corresponding
// client fields.
//
// A nil Transport is replaced with Ferret's policy-aware transport. A non-nil
// Transport is preserved and remains responsible for proxy behavior, DNS and
// concrete-address enforcement, and response-header limits. The transport and
// cookie jar remain shared with the supplied client; closing idle connections
// through the returned Client affects the shared transport pool.
func NewWithClient(client *stdhttp.Client, options ...PolicyOption) (Client, error) {
	if client == nil {
		return nil, ErrNilClient
	}

	policy, err := NewPolicy(options...)
	if err != nil {
		return nil, err
	}

	stdClient := *client
	transport, isStandardTransport := stdClient.Transport.(*stdhttp.Transport)

	if stdClient.Transport == nil || (isStandardTransport && transport == nil) {
		dialer := newPolicyDialer(policy)
		stdClient.Transport = newPolicyTransport(dialer, policy.MaxResponseHeaderSize())
	}

	return newDefaultHTTPClient(policy, stdClient), nil
}

// newDefaultHTTPClient snapshots and normalizes a standard-library client so
// the stored client can be reused safely across concurrent requests.
func newDefaultHTTPClient(policy *Policy, client stdhttp.Client) *defaultHTTPClient {
	result := &defaultHTTPClient{
		policy: policy,
		client: client,
	}
	result.client.Transport = newNilResponseGuardTransport(result.client.Transport)
	result.client.Timeout = policy.Timeout()
	result.client.CheckRedirect = policy.CheckRedirect

	return result
}

func (d *defaultHTTPClient) Do(ctx context.Context, req *Request) (*Response, error) {
	if ctx == nil {
		ctx = context.Background()
	}

	stdReq, err := toStdRequest(ctx, req)
	if err != nil {
		return nil, err
	}

	if err := d.policy.Prepare(stdReq); err != nil {
		return nil, err
	}

	res, err := d.client.Do(stdReq)
	if err != nil {
		var policyErr *PolicyError
		if errors.As(err, &policyErr) {
			return nil, policyErr
		}

		return nil, err
	}

	return fromStdResponse(res, d.policy)
}

func (d *defaultHTTPClient) CloseIdleConnections() {
	d.client.CloseIdleConnections()
}
