package bitcoinlib

import (
	"encoding/binary"
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
	cmds   []Operation
	isP2SH bool
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

func (s *ScriptPubKey) isP2SH() bool {
	fmt.Printf("Salgo por aca: %s\n", s.cmds)
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

func (t *ScriptPubKey) Combine(key Script) *CombinedScript {
	cmds := make([]Operation, len(t.cmds))
	copy(cmds, t.cmds)
	slices.Reverse(cmds)
	cmds = append(cmds, key.cmds...)
	slices.Reverse(cmds[len(t.cmds):])
	return &CombinedScript{
		cmds,
		t.isP2SH(),
	}
}

// Evaluates the hash of the script provided
func (t *CombinedScript) EvaluateScriptHash() bool {
	//Don't need z, so just use a placeholder
	z := ""
	helperScript := &CombinedScript{
		t.cmds[:4],
		false,
	}
	return helperScript.Evaluate(z)
}

// Evaluates a Redeem Script (need to parse it and then create the correct script to evaluate)
func (t *CombinedScript) EvaluateRedeemScript(z string) bool {
	script := t.cmds[len(t.cmds)-4].(*ScriptVal)
	pubKeyScript, err := parseScriptFromBytes(script.Val)
	if err != nil {
		return false
	}
	pubKey := NewPubkey(pubKeyScript)
	privKey := NewScript(t.cmds[4:])
	slices.Reverse(privKey.cmds)
	return pubKey.Combine(*privKey).Evaluate(z)
}

func (t *CombinedScript) Evaluate(z string) bool {
	if t.isP2SH {
		//Evaluate P2SH
		fmt.Println("Por aca pase")
		return t.EvaluateScriptHash() && t.EvaluateRedeemScript(z)
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
			cmds = append(cmds, op)
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
