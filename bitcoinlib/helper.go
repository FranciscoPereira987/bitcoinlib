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

func Murmur3Seed(hashNumber int, tweak int) int {
	return hashNumber*0xfba4c795 + tweak
}

func Murmur3(data []byte, seed int) int {
	c1 := 0xcc9e2d51
	c2 := 0x1b873593
	length := len(data)
	h1 := seed
	roundedEnd := (length & 0xfffffffc) // round down to 4 byte block
	for i := 0; i < roundedEnd; i += 4 {
		// little endian load order
		k1 := int((data[i] & 0xff)) | (int((data[i+1] & 0xff)) << 8) | (int((data[i+2] & 0xff)) << 16) | (int(data[i+3]) << 24)
		k1 *= c1
		k1 = (k1 << 15) | ((k1 & 0xffffffff) >> 17) // ROTL32(k1,15)
		k1 *= c2
		h1 ^= k1
		h1 = (h1 << 13) | ((h1 & 0xffffffff) >> 19) // ROTL32(h1,13)
		h1 = h1*5 + 0xe6546b64
	}
	k1 := 0
	val := length & 0x03
	if val == 3 {
		k1 = int(data[roundedEnd+2]&0xff) << 16
	}
	//fallthrough
	if val == 2 || val == 3 {
		k1 |= int(data[roundedEnd+1]&0xff) << 8
	}
	//fallthrough
	if val == 1 || val == 2 || val == 3 {
		k1 |= int(data[roundedEnd] & 0xff)
		k1 *= c1
		k1 = (k1 << 15) | ((k1 & 0xffffffff) >> 17) // ROTL32(k1,15)
		k1 *= c2
		h1 ^= k1
	}
	// finalization
	h1 ^= length
	// fmix(h1)
	h1 ^= ((h1 & 0xffffffff) >> 16)
	h1 *= 0x85ebca6b
	h1 ^= ((h1 & 0xffffffff) >> 13)
	h1 *= 0xc2b2ae35
	h1 ^= ((h1 & 0xffffffff) >> 16)
	return h1 & 0xffffffff
}
