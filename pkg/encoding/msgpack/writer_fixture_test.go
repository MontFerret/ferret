package msgpack_test

import (
	"bytes"
	"context"
	"fmt"

	vmmsgpack "github.com/vmihailenco/msgpack/v5"

	"github.com/MontFerret/ferret/v2/pkg/encoding"
	ferretmsgpack "github.com/MontFerret/ferret/v2/pkg/encoding/msgpack"
	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

type (
	ignoringWriterValue struct {
		marshalerValue
		failure error
	}

	errorWriter struct {
		failure error
	}

	reentrantArrayEncoder struct {
		*runtime.Array
		ctx context.Context
	}

	opaqueBufferWriter struct {
		bytes.Buffer
	}

	stringSettingsValue struct {
		marshalerValue
		dictionary bool
	}

	cancelingWriterValue struct {
		marshalerValue
		cancel  context.CancelFunc
		failure error
	}

	customIntegerValue struct {
		runtime.Int
	}

	mapCancelWriter struct {
		cancel context.CancelFunc
		bytes.Buffer
	}

	payloadBoundaryWriter struct {
		failure error
		cancel  context.CancelFunc
		bytes.Buffer
		stopAt          int
		writesAfterStop int
		stopped         bool
	}
)

func (v *ignoringWriterValue) EncodeMsgpack(enc *vmmsgpack.Encoder) error {
	enc.UseCompactInts(true)
	_, _ = enc.Writer().Write([]byte{1, 2, 3})

	return v.failure
}

func (w errorWriter) Write([]byte) (int, error) {
	return 0, w.failure
}

func (v *reentrantArrayEncoder) EncodeMsgpack(enc *vmmsgpack.Encoder) error {
	if err := enc.EncodeArrayLen(2); err != nil {
		return err
	}

	if err := enc.EncodeInt(1); err != nil {
		return err
	}

	data, err := encoding.EncodeBytes(v.ctx, ferretmsgpack.Default, runtime.Int(7))
	if err != nil {
		return err
	}

	if !bytes.Equal(data, []byte{7}) {
		return fmt.Errorf("nested encoding: %x", data)
	}

	return enc.EncodeInt(2)
}

func (v *stringSettingsValue) EncodeMsgpack(enc *vmmsgpack.Encoder) error {
	if v.dictionary {
		enc.ResetDict(enc.Writer(), map[string]int{"repeat": 0})
	} else {
		enc.UseInternedStrings(true)
	}

	return enc.EncodeNil()
}

func (v *cancelingWriterValue) EncodeMsgpack(enc *vmmsgpack.Encoder) error {
	v.cancel()
	_, _ = enc.Writer().Write(make([]byte, 4096))
	_, _ = enc.Writer().Write([]byte{1})

	return v.failure
}

func (v *customIntegerValue) EncodeMsgpack(enc *vmmsgpack.Encoder) error {
	return enc.EncodeString("custom integer")
}

func (w *mapCancelWriter) Write(data []byte) (int, error) {
	n, err := w.Buffer.Write(data)
	if w.cancel != nil {
		cancel := w.cancel
		w.cancel = nil
		cancel()
	}

	return n, err
}

func (w *payloadBoundaryWriter) Write(data []byte) (int, error) {
	if w.stopped {
		w.writesAfterStop++

		return 0, w.failure
	}

	size := min(len(data), max(0, w.stopAt-w.Len()))
	n, _ := w.Buffer.Write(data[:size])
	if size < len(data) {
		w.stopped = true

		return n, w.failure
	}

	if w.cancel != nil && w.Len() == w.stopAt {
		w.stopped = true
		w.cancel()
	}

	return n, nil
}
