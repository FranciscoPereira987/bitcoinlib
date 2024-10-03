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
		buf, _ = hex.DecodeString(string(buf))
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
		err = errors.Join(err, errors.New("could not parse uint32 from stream"))
	}
	return binary.LittleEndian.Uint32(buf), err
}

func parseUint64(from io.Reader) (uint64, error) {
	buf := make([]byte, 8)
	total, err := from.Read(buf)
	if total < 8 {
		err = errors.Join(err, errors.New("could not parse uint64 from stream"))
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
	segwit := inputs == 0
	if segwit {
		next := ReadVarInt(from)
		if next != 1 {
			return nil, errors.New("invalid segwit transaction")
		}
		inputs = ReadVarInt(from)
	}
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
	if segwit {
		for _ = range inputArr {
			value := []byte{0}
			from.Read(value)
			items := int(value[0])
			for _ = range items {
				length := ReadVarInt(from)
				item := make([]byte, length)
				from.Read(item)
			}
		}
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
func (tx *Transaction) Fee(testnet bool) int64 {
	var totalOutput uint64
	var totalInput uint64
	for _, val := range tx.inputs {
		valTotal, err := val.Value(testnet)
		if err != nil {
			return -1
		}
		totalInput += valTotal
	}
	for _, val := range tx.outputs {
		totalOutput += val.amount
	}
	//fmt.Printf("outputs: %d\ninputs: %d\n=====\n", totalOutput, totalInput)
	//Doing it this way to avoid overflow issues
	if totalInput < totalOutput {
		return -int64(totalOutput - totalInput)
	}
	return int64(totalInput - totalOutput)
}

// Returns the input serialization with
// the pubkey of the previous transaction
// instead of the script sig
// If empty is true, does not replace the ScripSig with the
// previous ScriptPubKey
func (in *Input) ReplaceScriptSig(empty bool, testnet bool) []byte {
	buf, _ := hex.DecodeString(in.previousID)
	slices.Reverse(buf)
	buf = binary.LittleEndian.AppendUint32(buf, in.previousIndex)
	if empty {
		buf = append(buf, 0x00)
	} else {
		pubKey, _ := FetchTransaction(in.previousID, testnet, true)
		scriptPubKey := pubKey.outputs[in.previousIndex].scriptPubKey
		buf = append(buf, scriptPubKey.Serialize()...)
	}
	buf = binary.LittleEndian.AppendUint32(buf, in.sequence)
	return buf
}

func (tx *Transaction) SigHash(input int, testnet bool) []byte {
	buf := tx.version.Serialize()
	buf = append(buf, EncodeVarInt(uint64(len(tx.inputs)))...)
	for index, val := range tx.inputs {
		buf = append(buf, val.ReplaceScriptSig(index != input, testnet)...)
	}
	buf = append(buf, EncodeVarInt(uint64(len(tx.outputs)))...)
	for _, val := range tx.outputs {
		buf = append(buf, val.Serialize()...)
	}
	buf = binary.LittleEndian.AppendUint32(buf, tx.locktime)
	buf = append(buf, 0x01, 0x00, 0x00, 0x00) //Append SIGHASH_ALL
	return Hash256(buf)
}

func (tx *Transaction) VerifyInput(input int, testnet bool) bool {
	//First of, get Z
	hash := tx.SigHash(input, testnet)
	//Get the public key that goes with this input script
	pubKey, err := tx.inputs[input].ScriptPubkey(testnet)
	if err != nil {
		return false
	}
	//Combine and evaluate the final Script
	combined := pubKey.Combine(*tx.inputs[input].scriptSig)
	return combined.Evaluate(hex.EncodeToString(hash))
}

func (tx *Transaction) Verify(testnet bool) bool {
	//Validating the fee
	if tx.Fee(testnet) < 0 {
		return false
	}
	//Need to validate the script of each input
	for i := range tx.inputs {
		if !tx.VerifyInput(i, testnet) {
			return false
		}
	}
	return true
}

func NewTransaction() *Transaction {
	return &Transaction{
		*NewVersion(1),
		[]*Input{},
		[]*Output{},
		0,
	}
}

func (tx *Transaction) AddInput(previousID string, previousIndex uint32) {
	newInput := &Input{
		previousID,
		previousIndex,
		&Script{},
		0xffffffff,
	}
	tx.inputs = append(tx.inputs, newInput)
}

func (tx *Transaction) SignInput(input int, testnet bool, key *PrivateKey) {
	z := tx.SigHash(input, testnet)
	zInt := FromHexString("0x" + hex.EncodeToString(z))
	sig := key.Sign(zInt)
	script := P2PKHSignature(append(sig.Der(), 0x01, 0x00, 0x00, 0x00), key.Sec(COMPRESSED))
	tx.inputs[input].scriptSig = script
}

func (tx *Transaction) AddOutput(amount uint64, address string) {
	pubKey := P2PKHScript([]byte(address))
	newOutput := &Output{
		amount,
		pubKey,
	}
	tx.outputs = append(tx.outputs, newOutput)
}