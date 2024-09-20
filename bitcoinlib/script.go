package bitcoinlib

import (
	"errors"
	"fmt"
	"io"
)

type Script struct {
	input []byte
	cmds  []Operation
}

type ScriptPubKey struct {
	input []byte
}

func ParsePubKey(from io.Reader) (*ScriptPubKey, error) {
	length := ReadVarInt(from)
	buf := make([]byte, length)
	total, err := from.Read(buf)
	if total < int(length) {
		err = errors.Join(err, errors.New("Invalid Script length decoded for pub key"))
	}
	return &ScriptPubKey{
		buf,
	}, err
}

func ParseScript(from io.Reader) (*Script, error) {
	length := ReadVarInt(from)
	buf := make([]byte, length)
	total, err := from.Read(buf)
	if total < int(length) {
		err = errors.Join(err, errors.New(
			fmt.Sprintf("Invalid Script length decoded: %d != %d", total, length)))
	}
	cmds := []Operation{}
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
	return &Script{
		buf,
		cmds,
	}, err
}

func (t *Script) Serialize() []byte {
	val := t.input
	val = append(EncodeVarInt(uint64(len(val))), val...)
	return val
}

func (t *ScriptPubKey) Serialize() []byte {
	val := t.input
	val = append(EncodeVarInt(uint64(len(val))), val...)
	return val
}
