// Package codecutil owns the built-in codecs' bounded I/O, cancellation
// checkpoints, and acquired iterator cleanup. It does not change the public
// reader/writer contract.
package codecutil

import (
	"bytes"
	"context"
	"io"
	"strings"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

// ChunkSize bounds payload writes and reader requests between checkpoints.
const ChunkSize = 4096

type (
	// Operation samples traversal cancellation every 256 values or tokens. I/O and
	// fallback boundaries can check directly without adding ctx.Err to hot loops.
	Operation struct {
		ctx   context.Context
		done  <-chan struct{}
		steps uint8
	}

	// Writer detects short writes and stops forwarding after the first failure.
	Writer struct {
		Dst     io.Writer
		err     error
		Op      Operation
		pending int
		byte    [1]byte
	}

	// Reader observes non-EOF errors even when a decoder accepts accompanying data.
	Reader struct {
		Src io.Reader
		err error
		Op  Operation
	}

	scannerReader struct {
		*Reader
		scanner io.ByteScanner
	}
)

func NewOperation(ctx context.Context, input any, name string) (Operation, error) {
	if ctx == nil {
		return Operation{}, runtime.Error(runtime.ErrInvalidArgument, "codec context is nil")
	}

	if input == nil {
		return Operation{}, runtime.Errorf(runtime.ErrInvalidArgument, "codec %s is nil", name)
	}

	op := Operation{ctx: ctx, done: ctx.Done()}

	return op, op.Check()
}

func (op *Operation) Check() error {
	if op.done == nil {
		return nil
	}

	select {
	case <-op.done:
		return op.ctx.Err()
	default:
		return nil
	}
}

func (op *Operation) Step() error {
	if op.done == nil {
		return nil
	}

	op.steps++
	if op.steps == 0 {
		return op.Check()
	}

	return nil
}

func (w *Writer) Write(p []byte) (int, error) {
	if w.err != nil {
		return 0, w.err
	}

	if len(p) <= ChunkSize {
		if err := w.checkWrite(len(p)); err != nil {
			return 0, err
		}

		n, err := w.Dst.Write(p)
		if err == nil && n != len(p) {
			err = io.ErrShortWrite
		}

		w.err = err

		return n, err
	}

	written := 0

	for len(p) > 0 {
		size := min(len(p), ChunkSize)
		n, err := w.Write(p[:size])
		written += n

		if err != nil {
			return written, err
		}

		p = p[size:]
	}

	return written, nil
}

func (w *Writer) WriteByte(value byte) error {
	if w.err != nil {
		return w.err
	}

	if buffer, ok := w.Dst.(*bytes.Buffer); ok {
		if err := w.checkWrite(1); err != nil {
			return err
		}

		return buffer.WriteByte(value)
	}

	w.byte[0] = value
	_, err := w.Write(w.byte[:])

	return err
}

func (w *Writer) Error() error {
	return w.err
}

func (r *Reader) Read(p []byte) (int, error) {
	if err := r.Op.Check(); err != nil {
		return 0, err
	}

	n, err := r.Src.Read(p[:min(len(p), ChunkSize)])
	r.recordError(err)

	return n, err
}

// Input preserves ByteScanner support, avoiding an extra MessagePack buffer for
// byte readers. Scanner reads use sampled cancellation rather than per-byte checks.
func (r *Reader) Input() io.Reader {
	// These standard readers have no non-EOF failures to observe. Preserve their
	// native fast path when the operation has no cancellation channel.
	if r.Op.done == nil {
		switch r.Src.(type) {
		case *bytes.Reader, *strings.Reader:
			return r.Src
		}
	}

	if scanner, ok := r.Src.(io.ByteScanner); ok {
		return &scannerReader{Reader: r, scanner: scanner}
	}

	return r
}

func (r *Reader) Error() error {
	return r.err
}

func (r *scannerReader) ReadByte() (byte, error) {
	if err := r.Op.Step(); err != nil {
		return 0, err
	}

	value, err := r.scanner.ReadByte()
	r.recordError(err)

	return value, err
}

func (r *scannerReader) UnreadByte() error {
	return r.scanner.UnreadByte()
}

func (r *Reader) recordError(err error) {
	// Only bare EOF is normal completion. An error joined with EOF can carry
	// another I/O cause even if the decoder consumes all accompanying data.
	if r.err == nil && err != nil && err != io.EOF {
		r.err = err
	}
}

func (w *Writer) checkWrite(size int) error {
	if w.Op.done == nil {
		return nil
	}

	w.pending += size
	if w.pending < ChunkSize {
		return nil
	}

	w.pending = 0
	w.err = w.Op.Check()

	return w.err
}
