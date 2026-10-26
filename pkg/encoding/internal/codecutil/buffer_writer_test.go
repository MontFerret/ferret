package codecutil_test

import (
	"bytes"
	"context"
	"errors"
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/encoding/internal/codecutil"
)

func TestBufferWriterSharesMixedWriteAccounting(t *testing.T) {
	ctx, cancel := context.WithCancel(t.Context())
	defer cancel()

	var dst bytes.Buffer
	op, err := codecutil.NewOperation(ctx, &dst, "writer")
	if err != nil {
		t.Fatal(err)
	}

	writer := codecutil.BufferWriter{Writer: codecutil.Writer{Dst: &dst, Op: op}, Buffer: &dst}
	base := &writer.Writer
	if n, err := writer.Write(bytes.Repeat([]byte{'a'}, codecutil.ChunkSize-2)); n != codecutil.ChunkSize-2 || err != nil {
		t.Fatalf("initial write: %d, %v", n, err)
	}

	if err := base.WriteByte('b'); err != nil {
		t.Fatal(err)
	}

	cancel()
	if err := writer.WriteByte('c'); !errors.Is(err, context.Canceled) || dst.Len() != codecutil.ChunkSize-1 {
		t.Fatalf("mixed write missed checkpoint: length=%d error=%v", dst.Len(), err)
	}

	if n, err := writer.Write([]byte("later")); n != 0 || !errors.Is(err, context.Canceled) {
		t.Fatalf("buffer adapter forgot failure: %d, %v", n, err)
	}

	if n, err := base.Write([]byte("later")); n != 0 || !errors.Is(err, context.Canceled) || dst.Len() != codecutil.ChunkSize-1 {
		t.Fatalf("base writer forgot failure: length=%d write=%d error=%v", dst.Len(), n, err)
	}

	if !errors.Is(writer.Error(), context.Canceled) || !errors.Is(base.Error(), context.Canceled) {
		t.Fatalf("failure not shared: adapter=%v base=%v", writer.Error(), base.Error())
	}
}

func TestBufferWriterChunkDelegationRetainsAccounting(t *testing.T) {
	ctx, cancel := context.WithCancel(t.Context())
	defer cancel()

	var dst bytes.Buffer
	op, err := codecutil.NewOperation(ctx, &dst, "writer")
	if err != nil {
		t.Fatal(err)
	}

	writer := codecutil.BufferWriter{Writer: codecutil.Writer{Dst: &dst, Op: op}, Buffer: &dst}
	base := &writer.Writer
	first := bytes.Repeat([]byte{'a'}, codecutil.ChunkSize+3)
	if n, err := writer.Write(first); n != len(first) || err != nil {
		t.Fatalf("chunked write: %d, %v", n, err)
	}

	cancel()
	second := bytes.Repeat([]byte{'b'}, codecutil.ChunkSize-4)
	if n, err := base.Write(second); n != len(second) || err != nil {
		t.Fatalf("write before checkpoint: %d, %v", n, err)
	}

	if n, err := writer.Write([]byte{'c'}); n != 0 || !errors.Is(err, context.Canceled) || dst.Len() != 2*codecutil.ChunkSize-1 {
		t.Fatalf("chunk remainder lost: length=%d write=%d error=%v", dst.Len(), n, err)
	}

	if err := base.WriteByte('d'); !errors.Is(err, context.Canceled) {
		t.Fatalf("base byte write forgot failure: %v", err)
	}

	if !bytes.Equal(dst.Bytes(), append(first, second...)) {
		t.Fatal("chunked output changed")
	}
}
