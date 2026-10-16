package msgpack_test

import vmmsgpack "github.com/vmihailenco/msgpack/v5"

type (
	ignoringWriterValue struct {
		marshalerValue
		failure error
	}

	errorWriter struct {
		failure error
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
