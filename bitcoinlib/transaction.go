package bitcoinlib

import (
	"encoding/binary"
	"encoding/hex"
	"errors"
	"io"
)

const VERSION_SIZE = 4

//Version, VarInt, Input\s, VarInt, Output\s
type Transaction struct{
  version Version
  inputs []*Input
  outputs []*Output
  locktime uint32
}

func (t *Transaction) GetOutputsAmount() []uint64 {
  result := make([]uint64, 0)
  for _, o := range t.outputs {
    result = append(result, o.amount)
  }
  return result
}

func (t *Transaction) GetOutputs() []string {
  result := make([]string, 0)
  for _, o := range t.outputs {
    result = append(result, o.scriptPubKey.input)
  }
  return result
}

func (t *Transaction) GetInputs() []string {
  result := make([]string, 0)
  for _, i := range t.inputs {
    result = append(result, i.scriptSig.input)
  }
  return result
}

type Version struct {
	number uint32
}

type Input struct{
	previousID string
	previousIndex uint32
	scriptSig *Script 
	sequence uint32
}

type Output struct {
  amount uint64
  scriptPubKey *ScriptPubKey
}

func parseHash(from io.Reader) (string, error) {
  buf := make([]byte, 32) //Hash256 Length in bytes
  total, err := from.Read(buf)
  if total < 32 {
    err = errors.Join(err, errors.New("Invalid bytestream for hash decoding"))
  }
  return hex.EncodeToString(buf), err

}

func parseUint32(from io.Reader) (uint32, error) {
  buf := make([]byte, 4)
  total, err := from.Read(buf)
  if total < 4 {
    err = errors.Join(err, errors.New("Could not parse uint32 from stream"))
  }
  return binary.LittleEndian.Uint32(buf), err
}

func parseUint64(from io.Reader) (uint64, error) {
  buf := make([]byte, 8)
  total, err := from.Read(buf)
  if total < 8 {
    err = errors.Join(err, errors.New("Could not parse uint64 from stream"))
  }
  return binary.LittleEndian.Uint64(buf), err
}

func NewOutputFrom(from io.Reader) (*Output, error) {
  amount, err := parseUint64(from)
  if err != nil {
    return nil, err
  }
  script, err := ParsePubKey(from)
  return &Output{
    amount,
    script,
  }, nil
}

func NewInputFrom(from io.Reader) (*Input, error) {
  previousID, err := parseHash(from)
  if err != nil {
    return nil, err
  }
  prevIndex, err := parseUint32(from)
  if err != nil {
    return nil, err
  }
  script, err := ParseScript(from)
  if err != nil {
    return nil, err
  }
  sequence, err := parseUint32(from)
  return &Input{
    previousID,
    prevIndex,
    script,
    sequence,
  }, err
}

func NewVersionFrom(from io.Reader) (*Version, error) {
	buf := make([]byte, VERSION_SIZE)
	total, err := from.Read(buf)
	if total != VERSION_SIZE || err != nil {
		err = errors.Join(err, errors.New("invalid read"))
		return nil, err
	}
	value := binary.LittleEndian.Uint32(buf)
	return &Version{
		value,
	}, nil
}

func (v Version) Eq(other Version) bool {
	return v.number == other.number
}

func (v Version) Ne(other Version) bool {
	return !v.Eq(other)
}

func NewVersion(value uint32) *Version {
	return &Version{
		value,
	}
}

func ParseTransaction(from io.Reader) (*Transaction, error) {
  version, err := NewVersionFrom(from)
  if err != nil {
    return nil, err
  }
  inputs := ReadVarInt(from)
  inputArr := make([]*Input, 0)
  for _ = range inputs {
    input, err := NewInputFrom(from)
    if err != nil {
      return nil, err
    }
    inputArr = append(inputArr, input)
  }
  outputs := ReadVarInt(from)
  outputArr := make([]*Output, 0)
  for _ = range outputs {
    output, err := NewOutputFrom(from)
    if err != nil {
      return nil, err
    }
    outputArr = append(outputArr, output)
  }
  locktime, err := parseUint32(from)
	return &Transaction{
    *version,
    inputArr,
    outputArr,
    locktime,
  }, err
}
