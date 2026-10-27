package encoding_test

import "bytes"

// encodeReentrantWriter calls the same codec while its outer operation still
// owns its write buffer and state.
type encodeReentrantWriter struct {
	beforeWrite func() error
	bytes.Buffer
	entered bool
}

func (w *encodeReentrantWriter) Write(data []byte) (int, error) {
	if !w.entered {
		w.entered = true
		if err := w.beforeWrite(); err != nil {
			return 0, err
		}
	}

	return w.Buffer.Write(data)
}
