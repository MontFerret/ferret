package artifact

import (
	"encoding/binary"
	"fmt"
)

const (
	headerSize    = 14
	schemaVersion = uint8(1)
)

var magic = [4]byte{'F', 'B', 'C', '2'}

type header struct {
	Magic         [4]byte
	Format        FormatID
	SchemaVersion uint8
	ISAVersion    uint16
	Flags         uint16
	PayloadLength uint32
}

func decodeHeader(data []byte) (header, error) {
	if len(data) < headerSize {
		return header{}, fmt.Errorf("%w: artifact shorter than %d-byte header", ErrInvalidHeader, headerSize)
	}

	return header{
		Magic:         [4]byte{data[0], data[1], data[2], data[3]},
		Format:        FormatID(data[4]),
		SchemaVersion: data[5],
		ISAVersion:    binary.LittleEndian.Uint16(data[6:8]),
		Flags:         binary.LittleEndian.Uint16(data[8:10]),
		PayloadLength: binary.LittleEndian.Uint32(data[10:14]),
	}, nil
}

func encodeHeader(dst []byte, header header) {
	copy(dst[:4], header.Magic[:])
	dst[4] = byte(header.Format)
	dst[5] = header.SchemaVersion
	binary.LittleEndian.PutUint16(dst[6:8], header.ISAVersion)
	binary.LittleEndian.PutUint16(dst[8:10], header.Flags)
	binary.LittleEndian.PutUint32(dst[10:14], header.PayloadLength)
}
