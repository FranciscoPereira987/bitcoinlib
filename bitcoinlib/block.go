package bitcoinlib

import (
	"encoding/binary"
	"encoding/hex"
	"errors"
	"io"
	"slices"
)

const BLOCK_SIZE = 80

type Block struct {
	version uint32
	prevBlock string
	merkleRoot string
	timestamp uint32
	blockBits uint32
	nonce uint32
}

//Returns a clean block
func NewBlock() *Block {
	return &Block{}
}

func (b *Block) Parse(from io.Reader) error {
	buf := make([]byte, BLOCK_SIZE)
	total, err := from.Read(buf)
	if err != nil || total < BLOCK_SIZE {
		err = errors.Join(err, errors.New("invalid stream length for block"))
		return err
	}
	//Reading the version
	b.version = binary.LittleEndian.Uint32(buf)
	//Reading the previous block hash
	buf = buf[4:]
	helperBuf := make([]byte, 32)
	copy(helperBuf, buf[:32])
	slices.Reverse(helperBuf)
	b.prevBlock = hex.EncodeToString(helperBuf)
	//Reading the merkle root
	buf = buf[32:]
	copy(helperBuf, buf[:32])
	slices.Reverse(helperBuf)
	b.merkleRoot = hex.EncodeToString(helperBuf)
	//Reading the timestamp
	buf = buf[32:]
	b.timestamp = binary.LittleEndian.Uint32(buf)
	//Reading the bits field
	buf = buf[4:]
	b.blockBits = binary.LittleEndian.Uint32(buf)
	//Reading the nonce field
	buf = buf[4:]
	b.nonce = binary.LittleEndian.Uint32(buf)
	return nil
}

func (b *Block) Serialize() []byte {
	buf := make([]byte, 0)
	buf = binary.LittleEndian.AppendUint32(buf, b.version)
	//helper buf to encode prev ID and Merkle Root
	helperBuf, _ := hex.DecodeString(b.prevBlock)
	slices.Reverse(helperBuf)
	buf = append(buf, helperBuf...)
	helperBuf, _ = hex.DecodeString(b.merkleRoot)
	slices.Reverse(helperBuf)
	buf = append(buf, helperBuf...)
	//Add last fields to header
	buf = binary.LittleEndian.AppendUint32(buf, b.timestamp)
	buf = binary.LittleEndian.AppendUint32(buf, b.blockBits)
	buf = binary.LittleEndian.AppendUint32(buf, b.nonce)
	return buf
}

func (b *Block) Hash() string {
	hash := Hash256(b.Serialize())
	slices.Reverse(hash)
	return hex.EncodeToString(hash)
}