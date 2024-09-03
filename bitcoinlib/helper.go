package bitcoinlib

import (
	"encoding/binary"
	"io"
)

const TWO_BYTES = 0xfd
const FOUR_BYTES = 0xfe
const EIGHT_BYTES = 0xff

const ONE_BYTE_LIMIT = 0xfc
const TWO_BYTE_LIMIT = 0xffff
const FOUR_BYTES_LIMIT = 0xffffffff

func ReadVarInt(from io.Reader) uint64 {
	buf := make([]byte, 1)
	_, _ = from.Read(buf)
	var result uint64
	switch buf[0] {
	case TWO_BYTES:
		{
			buf = make([]byte, 2)
			from.Read(buf)
			result = uint64(binary.LittleEndian.Uint16(buf))
		}
	case FOUR_BYTES:
		{
			buf = make([]byte, 4)
			from.Read(buf)
			result = uint64(binary.LittleEndian.Uint32(buf))
		}
	case EIGHT_BYTES:
		{
			buf = make([]byte, 8)
			from.Read(buf)
			result = uint64(binary.LittleEndian.Uint64(buf))
		}
	default:
		result = uint64(buf[0])
	}
	return result
}

func EncodeVarInt(value uint64) []byte {
	encoded := binary.LittleEndian.AppendUint64(nil, value)
	var number []byte
	var prefix []byte
	if value <= ONE_BYTE_LIMIT {
		return encoded[:1]
	}
	if value <= TWO_BYTE_LIMIT {
		number = encoded[:2]
		prefix = []byte{TWO_BYTES}
	} else if value <= FOUR_BYTES_LIMIT {
		number = encoded[:4]
		prefix = []byte{FOUR_BYTES}
	} else {
		number = encoded
		prefix = []byte{EIGHT_BYTES}
	}
	return append(prefix, number...)

}
