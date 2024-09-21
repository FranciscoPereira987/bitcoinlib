package bitcoinlib

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"net/http"
	"slices"
)

const VERSION_SIZE = 4

// Version, VarInt, Input\s, VarInt, Output\s
type Transaction struct {
	version  Version
	inputs   []*Input
	outputs  []*Output
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
		result = append(result, hex.EncodeToString(o.scriptPubKey.Serialize()))
	}
	return result
}

func (t *Transaction) GetInputs() []string {
	result := make([]string, 0)
	for _, i := range t.inputs {
		result = append(result, hex.EncodeToString(i.Serialize()))
	}
	return result
}

type Version struct {
	number uint32
}

type Input struct {
	previousID    string
	previousIndex uint32
	scriptSig     *Script
	sequence      uint32
}

type Output struct {
	amount       uint64
	scriptPubKey *ScriptPubKey
}

var TxCache map[string]*Transaction = make(map[string]*Transaction)

func GetUrl(testnet bool) string {
	if testnet {
		return "http://testnet.programmingbitcoin.com"
	}
	return "http://mainnet.programmingbitcoin.com"
}

func FetchTransaction(tx_id string, testnet bool, fresh bool) (*Transaction, error) {
	if _, ok := TxCache[tx_id]; fresh || !ok {
		url := fmt.Sprintf("%s/tx/%s.hex", GetUrl(testnet), tx_id)
		response, err := http.Get(url)
		if err != nil {
			return nil, err
		}
		buf, err := io.ReadAll(response.Body)
		if err != nil {
			return nil, err
		}
		buf = bytes.TrimSpace(buf)
		tx, err := ParseTransaction(bytes.NewReader(buf))
		if err != nil {
			return nil, err
		}
		TxCache[tx.Id()] = tx
	}
	return TxCache[tx_id], nil
}

func (tx *Transaction) Id() string {
	hashed := Hash256(tx.Serialize())
	slices.Reverse(hashed)
	return hex.EncodeToString(hashed)
}

func parseHash(from io.Reader) (string, error) {
	buf := make([]byte, 32) //Hash256 Length in bytes
	total, err := from.Read(buf)
	if total < 32 {
		err = errors.Join(err, errors.New("Invalid bytestream for hash decoding"))
	}
	slices.Reverse(buf)
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
	if err != nil {
		return nil, err
	}
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
	for range inputs {
		input, err := NewInputFrom(from)
		if err != nil {
			return nil, err
		}
		inputArr = append(inputArr, input)
	}
	outputs := ReadVarInt(from)
	outputArr := make([]*Output, 0)
	for range outputs {
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

func (t *Input) Serialize() []byte {
	buf, _ := hex.DecodeString(t.previousID)
	slices.Reverse(buf)
	buf = binary.LittleEndian.AppendUint32(buf, t.previousIndex)
	buf = append(buf, t.scriptSig.Serialize()...)
	buf = binary.LittleEndian.AppendUint32(buf, t.sequence)
	return buf
}

func (t *Output) Serialize() []byte {
	buf := make([]byte, 0)
	buf = binary.LittleEndian.AppendUint64(buf, t.amount)
	buf = append(buf, t.scriptPubKey.Serialize()...)
	return buf
}

func (t *Version) Serialize() []byte {
	return binary.LittleEndian.AppendUint32(nil, t.number)
}

// Serializes a transaction
func (tx *Transaction) Serialize() []byte {
	buf := tx.version.Serialize()
	buf = append(buf, EncodeVarInt(uint64(len(tx.inputs)))...)
	for _, val := range tx.inputs {
		buf = append(buf, val.Serialize()...)
	}
	buf = append(buf, EncodeVarInt(uint64(len(tx.outputs)))...)
	for _, val := range tx.outputs {
		buf = append(buf, val.Serialize()...)
	}
	return binary.LittleEndian.AppendUint32(buf, tx.locktime)
}

func (t *Input) FetchTx(testnet bool) (*Transaction, error) {
	return FetchTransaction(t.previousID, testnet, false)
}

func (t *Input) Value(testnet bool) (uint64, error) {
	tx, err := t.FetchTx(testnet)
	if err != nil {
		return 0, err
	}
	return tx.GetOutputsAmount()[t.previousIndex], nil
}

func (t *Input) ScriptPubkey(testnet bool) (*ScriptPubKey, error) {
	tx, err := t.FetchTx(testnet)
	if err != nil {
		return nil, err
	}
	return tx.outputs[t.previousIndex].scriptPubKey, nil
}

// Returns the implied fee of a Transaction
func (tx *Transaction) Fee(testnet bool) uint64 {
	var totalOutput uint64
	var totalInput uint64
	for _, val := range tx.inputs {
		valTotal, err := val.Value(testnet)
		if err != nil {
			break
		}
		totalInput += valTotal
	}
	for _, val := range tx.outputs {
		totalOutput += val.amount
	}
	return totalOutput - totalInput
}

func (tx *Transaction) Verify() bool {
  return false
}
