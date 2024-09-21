package bitcoinlib

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"
)

type Script struct {
	cmds  []Operation
}

type ScriptPubKey struct {
  cmds []Operation
}

type CombinedScript struct {
  cmds []Operation
}

func (t *ScriptPubKey) Combine(key Script) *CombinedScript {
  cmds := make([]Operation, 0)
  i := len(t.cmds)-1
  for ;i >= 0;i-- {
    Push(&cmds, t.cmds[i])
  }
  i = len(key.cmds)-1
  for ;i>=0;i-- {
    Push(&cmds, key.cmds[i])
  }
  return &CombinedScript{
    cmds,
  }
}

func (t *CombinedScript) Evaluate(z string) bool {
  cmds := make([]Operation, len(t.cmds))
  copy(cmds, t.cmds)
  stack := make([]Operation, 0)
  altstack := make([]Operation, 0)
  fmt.Printf("%s\n", cmds)
  for len(cmds) > 0 {
    cmd := Pop(&cmds)
    if !cmd.Operate(z, &stack, &altstack, &cmds){
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
      }else if length < 256 {
        result = append(result, 76, byte(length))
      }else if length < 520 {
        result = append(result, 77)
        result = binary.LittleEndian.AppendUint16(result, uint16(length))
      }
      result = append(result, val.Val...)
    }else {
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
