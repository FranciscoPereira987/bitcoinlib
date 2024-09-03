package bitcoinlib

import (
	"encoding/binary"
	"errors"
	"io"
)

const VERSION_SIZE = 4

type Transaction struct{}

type Version struct {
	number uint32
}

type Input struct{
	previousID string
	previousIndex uint32
	//scriptSig ScriptSig
	sequence uint32
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

	return &Transaction{}, nil
}
