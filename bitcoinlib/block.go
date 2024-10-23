package bitcoinlib

import (
	"encoding/binary"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"slices"
)

const BLOCK_SIZE = 80


const GENESIS_BLOCK = "0100000000000000000000000000000000000000000000000000000000000000000000003ba3edfd7a7b12b27ac72c3e67768f617fc81bc3888a51323a9fb8aa4b1e5e4a29ab5f49ffff001d1dac2b7c"
const TESTNET_GENESIS_BLOCK = "0100000000000000000000000000000000000000000000000000000000000000000000003ba3edfd7a7b12b27ac72c3e67768f617fc81bc3888a51323a9fb8aa4b1e5e4adae5494dffff001d1aa4ae18"


type Block struct {
	version uint32
	prevBlock string
	merkleRoot string
	timestamp uint32
	blockBits uint32
	nonce uint32
}

func BitsToTarget(bits uint32) Int {
  exponent := FromInt(int(bits >> 24)-3)
  coefficient := FromInt(int((bits << 8) >> 8))
  base := FromInt(256)
  return coefficient.Mul(base.Exp(exponent, MAX)) 
}

func TwoWeeks() uint32 {
  return 14*24*60*60
}

func EightWeeks() uint32 {
  return TwoWeeks() * 4
}

func ThreeDaysAndAHalf() uint32 {
  return 60*60*12*7
}

func GetUpdateCoef(from uint32, to uint32) uint32 {
  distance := to-from
  if distance < ThreeDaysAndAHalf() {
    distance = ThreeDaysAndAHalf()
  }else if distance > EightWeeks() {
    distance = EightWeeks()
  }
  return (distance) / TwoWeeks()
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
	//Reading the nonce field0xffff * 256**(0x1d-3)
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

func (b *Block) BIP9() bool {
	//A block is BIP9 if its version starts with 001
	masked := (b.version >> 29) & 1 
	return masked == 1
}

func (b *Block) BIP91() bool {
	//A BIP91 Block has bit 4 set to 1
	masked := (b.version >> 4) & 1
	
	return masked == 1
}

func (b *Block) BIP141() bool {
	//A BIP141 Block has bit 1 set to 1
	masked := (b.version >> 1) & 1
	return masked == 1
}

func (b *Block) BitsToTarget() Int {
  return BitsToTarget(b.blockBits) 
}

func (b *Block) Difficulty() Int {
  genesis := BitsToTarget(0x1d00ffff) 
   
  return genesis.Div(b.BitsToTarget()) 
}

func (b *Block) Timestamp() uint32 {
	return b.timestamp
}

func (b *Block) CheckPOW() bool {
  target := b.BitsToTarget()
  hash := b.Hash()
  hashInt := FromHexString("0x"+hash)
  return hashInt.Le(target)
}

func (b *Block) GetNextTarget(b2 *Block) uint32 {
  updateCoeff := GetUpdateCoef(b.timestamp, b2.timestamp)
  nextTarget := b2.BitsToTarget().Mul(FromInt(int(updateCoeff)))
  asBytes := nextTarget.value.Bytes()
  fmt.Println(nextTarget.String(), hex.EncodeToString(asBytes))
  exponent := len(asBytes)
  coefficient := asBytes[:3]
  if asBytes[0] > 0x7f {
    exponent++
    coefficient = append([]byte{0}, coefficient[:2]...)
  }
  slices.Reverse(coefficient)
  bits := append(coefficient, byte(exponent))
  return binary.LittleEndian.Uint32(bits)
}
