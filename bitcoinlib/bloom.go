package bitcoinlib

import (
	"encoding/hex"
	"fmt"
)

type BloomFilter struct {
	bitField []byte
	size     Int
}

func NewBloomFilter(size int) *BloomFilter {
	totalBytes := size / 8
	if size%8 != 0 {
		totalBytes++
	}
	return &BloomFilter{
		make([]byte, totalBytes),
		FromInt(size),
	}
}

func (bf BloomFilter) String() string {
	return hex.EncodeToString(bf.bitField)[:]
}

func (bf *BloomFilter) set(total Int) {
	bitNumber := total.Mod(bf.size).value.Int64()
	byteNumber := bitNumber / 8
	bitIndex := bitNumber % 8
	fmt.Println(bitIndex)
	bf.bitField[byteNumber] |= 0x80 >> bitIndex
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
