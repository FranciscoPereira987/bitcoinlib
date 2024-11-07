package bitcoinlib

import (
	"encoding/binary"
	"encoding/hex"
)

type BloomFilter struct {
	bitField []byte
	size     Int
	params   *MurmurParams
}

type MurmurParams struct {
	FunctionCount int
	Tweak         int
}

func NewBloomFilter(size int) *BloomFilter {
	totalBytes := size
	return &BloomFilter{
		make([]byte, totalBytes),
		FromInt(size * 8),
		nil,
	}
}

func (bf BloomFilter) String() string {
	return hex.EncodeToString(bf.bitField)[:]
}

func (bf *BloomFilter) set(total Int) {
	bitNumber := total.Mod(bf.size).value.Int64()
	byteNumber := bitNumber / 8
	bitIndex := bitNumber % 8
	bf.bitField[byteNumber] |= 0x01 << bitIndex
}

/*
Sets the correspoding bit based on the Hash160 of the value
passed to this function
*/
func (bf *BloomFilter) Set160(stream []byte) {
	hashed := "0x" + hex.EncodeToString(Hash160(stream))
	total := FromHexString(hashed)
	bf.set(total)
}

/*
Sets the bloom filter using Murmur3 based on function Count and tweak
*/
func (bf *BloomFilter) Set(value []byte, params *MurmurParams) {
	if bf.params != nil {
		params = bf.params
	} else {
		bf.params = params
	}
	for i := range params.FunctionCount {
		seed := Murmur3Seed(i, params.Tweak)
		total := Murmur3(value, seed)
		bf.set(FromInt(total))
	}
}

func (bf *BloomFilter) FilterLoad() []byte {
	buf := EncodeVarInt(bf.size.value.Uint64() / 8)
	buf = append(buf, bf.bitField...)
	buf = binary.LittleEndian.AppendUint32(buf, uint32(bf.params.FunctionCount))
	buf = binary.LittleEndian.AppendUint32(buf, uint32(bf.params.Tweak))
	buf = append(buf, 1)
	return buf
}
