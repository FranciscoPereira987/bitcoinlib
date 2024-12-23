package bitcoinlib

import (
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"slices"
)

type Script struct {
	cmds []Operation
}

type ScriptPubKey struct {
	cmds []Operation
}

func P2PKHScript(hash []byte) *ScriptPubKey {
	return &ScriptPubKey{
		[]Operation{
			&OP_DUP{},
			&OP_HASH160{},
			&ScriptVal{hash},
			&OP_EQUALVERIFY{},
			&OP_CHECKSIG{},
		},
	}
}

func P2PKHSignature(der []byte, sec []byte) *Script {
	return &Script{
		[]Operation{
			&ScriptVal{der},
			&ScriptVal{sec},
		},
	}
}

func NewScriptVal(val []byte) *ScriptVal {
	return &ScriptVal{
		val,
	}
}

type CombinedScript struct {
	cmds     []Operation
	isP2SH   bool
	isP2WPKH bool
	isP2WSH  bool
}

func NewScript(cmds []Operation) *Script {
	return &Script{
		cmds,
	}
}

func NewPubkey(cmds []Operation) *ScriptPubKey {
	return &ScriptPubKey{
		cmds,
	}
}

func P2SHPubKey(hash []byte) *ScriptPubKey {

	return &ScriptPubKey{
		[]Operation{
			&OP_HASH160{},
			&ScriptVal{hash},
			&OP_EQUAL{},
		},
	}
}

func P2WPKHPubKey(hash []byte) *ScriptPubKey {
	return &ScriptPubKey{
		[]Operation{
			&OP_0{},
			&ScriptVal{hash},
		},
	}
}

func P2WSHPubKey(hash []byte) *ScriptPubKey {
	return &ScriptPubKey{
		[]Operation{
			&OP_0{},
			&ScriptVal{hash},
		},
	}
}

func (s *ScriptPubKey) isP2SH() bool {
	commands := P2SHPubKey([]byte{}).cmds
	if len(commands) != len(s.cmds) {
		return false
	}
	for index, cmd := range commands {
		if cmd.Num() != s.cmds[index].Num() {
			return false
		}
	}
	return true
}

func (s *ScriptPubKey) isP2WPKH() bool {
	commands := P2WPKHPubKey([]byte{}).cmds
	if len(commands) != len(s.cmds) {
		return false
	}
	for index, cmd := range commands {
		if cmd.Num() != s.cmds[index].Num() {
			return false
		}
	}
	isHash160, ok := s.cmds[1].(*ScriptVal)
	return ok && len(isHash160.Val) == 20
}

func (s *ScriptPubKey) isP2WSH() bool {
	commands := P2WSHPubKey([]byte{}).cmds
	if len(commands) != len(s.cmds) {
		return false
	}
	for index, cmd := range commands {
		if cmd.Num() != s.cmds[index].Num() {
			return false
		}
	}
	isHash256, ok := s.cmds[1].(*ScriptVal)
	return ok && len(isHash256.Val) == 32
}
func (t *ScriptPubKey) Combine(key Script) *CombinedScript {
	cmds := make([]Operation, len(t.cmds))
	copy(cmds, t.cmds)
	slices.Reverse(cmds)
	cmds = append(cmds, key.cmds...)
	slices.Reverse(cmds[len(t.cmds):])
	return &CombinedScript{
		cmds,
		t.isP2SH(),
		t.isP2WPKH(),
		t.isP2WSH(),
	}
}

// Evaluates the hash of the script provided
func (t *CombinedScript) EvaluateScriptHash() bool {
	//Don't need z, so just use a placeholder
	z := ""
	helperScript := &CombinedScript{
		t.cmds[:4],
		false,
		false,
		false,
	}
	return helperScript.Evaluate(z, nil)
}

// Evaluates a Redeem Script (need to parse it and then create the correct script to evaluate)
func (t *CombinedScript) EvaluateRedeemScript(z string, witness [][]byte) bool {
	if witness != nil {
		otherParse, err := parseScriptFromBytes(t.cmds[len(t.cmds)-1].(*ScriptVal).Val)
		if err != nil {
			return false
		}
		pubKey := NewPubkey(otherParse)
		privKey := NewScript([]Operation{})
		return pubKey.Combine(*privKey).Evaluate(z, witness)
	}
	script := t.cmds[len(t.cmds)-4].(*ScriptVal)
	pubKeyScript, err := parseScriptFromBytes(script.Val)
	if err != nil {
		return false
	}
	pubKey := NewPubkey(pubKeyScript)
	privKey := NewScript(t.cmds[4:])
	slices.Reverse(privKey.cmds)
	return pubKey.Combine(*privKey).Evaluate(z, witness)
}

func EvaluateP2WPSH(z string, sha string, witness [][]byte) bool {
	validation := sha256.Sum256(witness[len(witness)-1])
	if hex.EncodeToString(validation[:]) != sha {
		return false
	}
	script, err := parseScriptFromBytes(witness[len(witness)-1])
	if err != nil {
		return false
	}
	rest := []byte{}
	for i := range len(witness) - 1 {
		rest = append(rest, EncodeVarInt(uint64(len(witness[i])))...)
		rest = append(rest, witness[i]...)
	}
	pubkey, err := parseScriptFromBytes(rest)
	if err != nil {
		return false
	}
	final := append(pubkey, script...)
	slices.Reverse(final)
	return (&CombinedScript{final, false, false, false}).Evaluate(z, nil)
}

func (t *CombinedScript) Evaluate(z string, witness [][]byte) bool {
	if t.isP2WSH {
		return EvaluateP2WPSH(z, hex.EncodeToString(t.cmds[0].(*ScriptVal).Val), witness)
	}
	if t.isP2SH {
		//Evaluate P2SH
		return t.EvaluateScriptHash() && t.EvaluateRedeemScript(z, witness)
	}
	var witnesses []byte
	if witness != nil {
		for _, item := range witness {
			value := EncodeVarInt(uint64(len(item)))
			witnesses = append(witnesses, value...)
			if string(value) != string(item) {
				witnesses = append(witnesses, item...)
			}
		}
	}
	cmds := make([]Operation, len(t.cmds))
	copy(cmds, t.cmds)
	stack := make([]Operation, 0)
	altstack := make([]Operation, 0)
	for len(cmds) > 0 {
		cmd := Pop(&cmds)
		if !cmd.Operate(z, &stack, &altstack, &cmds) {
			return false
		}
		if len(stack) == 2 && t.isP2WPKH {
			witnessScript, err := parseScriptFromBytes(witnesses)
			if err != nil {
				fmt.Println(err)
			}
			h160 := Pop(&stack)

			p2wpkh := []Operation{
				&OP_DUP{},
				&OP_HASH160{},
				h160,
				&OP_EQUALVERIFY{},
				&OP_CHECKSIG{},
			}

			for len(p2wpkh) > 0 {
				Push(&cmds, Pop(&p2wpkh))
			}
			for len(witnessScript) > 0 {
				Push(&cmds, Pop(&witnessScript))
			}
			t.isP2WPKH = false
		}
	}

	if len(stack) == 0 {
		return false
	}
	op := Pop(&stack)
	return op.Num() != 0
}

func ParsePubKey(from io.Reader) (*ScriptPubKey, error) {
	length := ReadVarInt(from)
	buf := make([]byte, length)
	total, err := from.Read(buf)
	if total < int(length) {
		err = errors.Join(err, errors.New("invalid Script length decoded for pub key"))
		return nil, err
	}
	cmds, err := parseScriptFromBytes(buf)
	if err != nil {
		return nil, err
	}
	return &ScriptPubKey{
		cmds,
	}, err
}

func serializeScriptToBytes(cmds []Operation) []byte {
	result := make([]byte, 0)
	for _, op := range cmds {
		if val, ok := op.(*ScriptVal); ok {
			//If its a ScriptVal then its a value
			//That should be treated as such (length + val) or (OP_PUSHDATAx + length + VAL)
			length := len(val.Val)
			if length < 75 {
				result = append(result, byte(length))
			} else if length < 256 {
				result = append(result, 76, byte(length))
			} else if length < 520 {
				result = append(result, 77)
				result = binary.LittleEndian.AppendUint16(result, uint16(length))
			}
			result = append(result, val.Val...)
		} else {
			result = append(result, byte(op.Num()))
		}
	}
	return result
}

func parseScriptFromBytes(buf []byte) ([]Operation, error) {
	cmds := []Operation{}
	total := len(buf)
	index := 0
	for index < total {
		current := buf[index]
		index++
		if current >= 1 && current <= 75 {
			//It´s an element
			op := &ScriptVal{
				buf[index : index+int(current)],
			}
			cmds = append(cmds, op)
			index += int(current)
		} else if current == 76 {
			//OP_PUSHDATA1
			length := FromLittleEndian(buf[index : index+1])
			op := &ScriptVal{
				buf[index+1 : index+1+int(length.value.Int64())],
			}
			index += 1 + int(length.value.Int64())
			cmds = append(cmds, op)
		} else if current == 77 {
			//OP_PUSHDATA2
			length := FromLittleEndian(buf[index : index+2])
			op := &ScriptVal{
				buf[index+2 : index+2+int(length.value.Int64())],
			}
			index += 2 + int(length.value.Int64())
			cmds = append(cmds, op)
		} else {
			//Simple Operation
			op := OP_CODE_FUNCTIONS[int(current)]
			if op == nil {
				fmt.Printf("Adding undefined operation: %d\n", int(current))
				cmds = append(cmds, &UNDEFINED{int(current)})
			} else {
				cmds = append(cmds, op)
			}
		}

	}
	if index != len(buf) {
		return nil, errors.New("failed to parse script")
	}
	return cmds, nil
}

func ParseScript(from io.Reader) (*Script, error) {
	length := ReadVarInt(from)
	buf := make([]byte, length)
	total, err := from.Read(buf)
	if total < int(length) {
		err = errors.Join(err, errors.New(
			fmt.Sprintf("invalid Script length decoded: %d != %d", total, length)))
	}
	cmds, err := parseScriptFromBytes(buf)
	return &Script{
		cmds,
	}, err
}

func (t *Script) Serialize() []byte {
	val := serializeScriptToBytes(t.cmds)
	length := EncodeVarInt(uint64(len(val)))
	return append(length, val...)
}

func (t *ScriptPubKey) Serialize() []byte {
	val := serializeScriptToBytes(t.cmds)
	length := EncodeVarInt(uint64(len(val)))
	return append(length, val...)
}

func (t *Script) Height() uint64 {
	if val, ok := t.cmds[0].(*ScriptVal); ok {
		buf := make([]byte, 8)
		copy(buf, val.Val)
		return binary.LittleEndian.Uint64(buf)
	}
	return 0
}
