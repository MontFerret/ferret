package main

import (
	"bufio"
	"bytes"
	"context"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

	godap "github.com/google/go-dap"
)

func TestServeDAPKeepsStdoutPureAndRoutesTrace(t *testing.T) {
	input := dapInitializeDisconnectInput(t)

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	if err := serveDAP(context.Background(), input, &stdout, &stderr, true, ""); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(stderr.String(), "recv ") || !strings.Contains(stderr.String(), "send ") {
		t.Fatalf("expected DAP trace on stderr, got %q", stderr.String())
	}

	reader := bufio.NewReader(bytes.NewReader(stdout.Bytes()))
	if _, ok := readDAPMessage(t, reader).(*godap.InitializeResponse); !ok {
		t.Fatal("expected initialize response")
	}
	if _, ok := readDAPMessage(t, reader).(*godap.DisconnectResponse); !ok {
		t.Fatal("expected disconnect response")
	}
	if _, ok := readDAPMessage(t, reader).(*godap.TerminatedEvent); !ok {
		t.Fatal("expected terminated event")
	}
	if _, err := godap.ReadProtocolMessage(reader); err != io.EOF {
		t.Fatalf("stdout contained non-DAP data: %v", err)
	}
}

func TestServeDAPRoutesTraceToLogFile(t *testing.T) {
	input := dapInitializeDisconnectInput(t)
	logPath := filepath.Join(t.TempDir(), "dap.log")

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	if err := serveDAP(context.Background(), input, &stdout, &stderr, true, logPath); err != nil {
		t.Fatal(err)
	}
	if stderr.Len() != 0 {
		t.Fatalf("expected no trace on stderr, got %q", stderr.String())
	}

	logged, err := os.ReadFile(logPath)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(logged), "recv ") || !strings.Contains(string(logged), "send ") {
		t.Fatalf("expected DAP trace in log file, got %q", logged)
	}

	reader := bufio.NewReader(bytes.NewReader(stdout.Bytes()))
	if _, ok := readDAPMessage(t, reader).(*godap.InitializeResponse); !ok {
		t.Fatal("expected initialize response")
	}
	if _, ok := readDAPMessage(t, reader).(*godap.DisconnectResponse); !ok {
		t.Fatal("expected disconnect response")
	}
	if _, ok := readDAPMessage(t, reader).(*godap.TerminatedEvent); !ok {
		t.Fatal("expected terminated event")
	}
	if _, err := godap.ReadProtocolMessage(reader); err != io.EOF {
		t.Fatalf("stdout contained non-DAP data: %v", err)
	}
}

func dapInitializeDisconnectInput(t *testing.T) *bytes.Buffer {
	t.Helper()

	var input bytes.Buffer
	if err := godap.WriteProtocolMessage(&input, &godap.InitializeRequest{
		Request: godap.Request{
			ProtocolMessage: godap.ProtocolMessage{Seq: 1, Type: "request"},
			Command:         "initialize",
		},
		Arguments: godap.InitializeRequestArguments{
			AdapterID:       "ferret",
			LinesStartAt1:   true,
			ColumnsStartAt1: true,
		},
	}); err != nil {
		t.Fatal(err)
	}
	if err := godap.WriteProtocolMessage(&input, &godap.DisconnectRequest{
		Request: godap.Request{
			ProtocolMessage: godap.ProtocolMessage{Seq: 2, Type: "request"},
			Command:         "disconnect",
		},
	}); err != nil {
		t.Fatal(err)
	}

	return &input
}

func readDAPMessage(t *testing.T, reader *bufio.Reader) godap.Message {
	t.Helper()

	message, err := godap.ReadProtocolMessage(reader)
	if err != nil {
		t.Fatal(err)
	}

	return message
}
