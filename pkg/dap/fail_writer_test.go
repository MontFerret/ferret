package dap

import (
	"bytes"
	"errors"
)

type failAfterWriter struct {
	Err       error
	Buffer    bytes.Buffer
	FailAfter int
	Writes    int
}

func (w *failAfterWriter) Write(data []byte) (int, error) {
	if w.Writes >= w.FailAfter {
		if w.Err == nil {
			return 0, errors.New("write failed")
		}

		return 0, w.Err
	}

	w.Writes++
	return w.Buffer.Write(data)
}
