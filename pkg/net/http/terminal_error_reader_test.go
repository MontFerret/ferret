package http

import "io"

type terminalErrorReader struct {
	err  error
	data []byte
}

func newTerminalErrorReader(data string, err error) *terminalErrorReader {
	return &terminalErrorReader{
		data: []byte(data),
		err:  err,
	}
}

func (r *terminalErrorReader) Read(p []byte) (int, error) {
	if len(r.data) == 0 {
		return 0, io.EOF
	}

	n := copy(p, r.data)
	r.data = r.data[n:]
	if len(r.data) == 0 {
		return n, r.err
	}

	return n, nil
}
