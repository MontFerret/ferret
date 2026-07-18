package http

import (
	"context"
	"net"
	"syscall"
	"time"
)

type policyDialer struct {
	policy *Policy
	dialer net.Dialer
}

func newPolicyDialer(policy *Policy) *policyDialer {
	d := &policyDialer{
		policy: policy,
	}
	d.dialer = net.Dialer{
		Timeout:        30 * time.Second,
		KeepAlive:      30 * time.Second,
		ControlContext: d.controlContext,
	}

	return d
}

func (d *policyDialer) DialContext(ctx context.Context, network, address string) (net.Conn, error) {
	return d.dialer.DialContext(ctx, network, address)
}

func (d *policyDialer) controlContext(
	_ context.Context,
	_ string,
	address string,
	_ syscall.RawConn,
) error {
	host, _, err := net.SplitHostPort(address)
	if err != nil {
		return newPolicyError(
			PolicyTargetConnection,
			"destination address",
			"invalid address is not allowed",
		)
	}

	addr, ok := parseIPAddress(host)
	if !ok {
		return newPolicyError(
			PolicyTargetConnection,
			"destination address",
			"invalid address is not allowed",
		)
	}

	p := d.policy
	if p == nil {
		p = &Policy{}
	}

	return p.validateAddress(PolicyTargetConnection, addressSubject(addr), addr)
}
