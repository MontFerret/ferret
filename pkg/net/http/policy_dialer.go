package http

import (
	"context"
	"net"
	"net/netip"
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
	p := d.policy
	if p == nil {
		p = &Policy{}
	}

	host, _, err := net.SplitHostPort(address)
	if err != nil {
		return p.EvalConnection(netip.Addr{})
	}

	addr, ok := parseIPAddress(host)
	if !ok {
		return p.EvalConnection(netip.Addr{})
	}

	return p.EvalConnection(addr)
}
