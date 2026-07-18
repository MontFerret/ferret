package http

import (
	"context"
	"net"
	"sync/atomic"
	"testing"
)

func newLoopbackResolver(tb testing.TB) (*net.Resolver, *atomic.Int64) {
	tb.Helper()

	server, err := net.ListenPacket("udp4", "127.0.0.1:0")
	if err != nil {
		tb.Fatalf("listen for test DNS: %v", err)
	}
	tb.Cleanup(func() { _ = server.Close() })

	queries := &atomic.Int64{}
	go serveLoopbackDNS(server, queries)

	resolver := &net.Resolver{
		PreferGo: true,
		Dial: func(ctx context.Context, _, _ string) (net.Conn, error) {
			var dialer net.Dialer

			return dialer.DialContext(ctx, "udp4", server.LocalAddr().String())
		},
	}

	return resolver, queries
}

func serveLoopbackDNS(server net.PacketConn, queries *atomic.Int64) {
	buffer := make([]byte, 1500)
	for {
		n, peer, err := server.ReadFrom(buffer)
		if err != nil {
			return
		}

		response := loopbackDNSResponse(buffer[:n])
		if response == nil {
			continue
		}

		queries.Add(1)
		_, _ = server.WriteTo(response, peer)
	}
}

func loopbackDNSResponse(query []byte) []byte {
	const dnsHeaderSize = 12
	if len(query) < dnsHeaderSize || query[4] != 0 || query[5] != 1 {
		return nil
	}

	offset := dnsHeaderSize
	for {
		if offset >= len(query) {
			return nil
		}

		labelSize := int(query[offset])
		offset++
		if labelSize == 0 {
			break
		}
		if labelSize > 63 || offset+labelSize > len(query) {
			return nil
		}

		offset += labelSize
	}

	if offset+4 > len(query) {
		return nil
	}

	questionEnd := offset + 4
	queryType := uint16(query[offset])<<8 | uint16(query[offset+1])
	response := append([]byte(nil), query[:questionEnd]...)
	response[2] = 0x80 | query[2]&0x01
	response[3] = 0x80
	response[6], response[7] = 0, 0
	response[8], response[9] = 0, 0
	response[10], response[11] = 0, 0

	var answer []byte
	switch queryType {
	case 1:
		answer = []byte{
			0xc0, 0x0c,
			0x00, 0x01,
			0x00, 0x01,
			0x00, 0x00, 0x00, 0x00,
			0x00, 0x04,
			127, 0, 0, 1,
		}
	case 28:
		answer = []byte{
			0xc0, 0x0c,
			0x00, 0x1c,
			0x00, 0x01,
			0x00, 0x00, 0x00, 0x00,
			0x00, 0x10,
			0, 0, 0, 0, 0, 0, 0, 0,
			0, 0, 0, 0, 0, 0, 0, 1,
		}
	default:
		return response
	}

	response[6], response[7] = 0, 1

	return append(response, answer...)
}
