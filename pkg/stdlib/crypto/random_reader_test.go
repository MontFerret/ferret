package crypto

import "context"

type (
	iotestReader struct {
		err error
	}

	cancelingTokenReader struct {
		cancel context.CancelFunc
	}
)

func (r iotestReader) Read(_ []byte) (int, error) {
	return 0, r.err
}

func (r cancelingTokenReader) Read(data []byte) (int, error) {
	for index := range data {
		data[index] = 255
	}

	r.cancel()

	return len(data), nil
}
