package encoding_test

import "bytes"

type countedScanner struct {
	*bytes.Reader
	reads       int
	readBytes   int
	unreadBytes int
	closed      bool
}

func (r *countedScanner) Read(dst []byte) (int, error) {
	r.reads++

	return r.Reader.Read(dst)
}

func (r *countedScanner) ReadByte() (byte, error) {
	r.readBytes++

	return r.Reader.ReadByte()
}

func (r *countedScanner) UnreadByte() error {
	r.unreadBytes++

	return r.Reader.UnreadByte()
}

func (r *countedScanner) Close() error {
	r.closed = true

	return nil
}
