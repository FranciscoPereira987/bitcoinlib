package bitcoinlib

import (
	"encoding/hex"
	"errors"
	"fmt"
	"io"
)

type Script struct{
  input string
}

type ScriptPubKey struct{
  input string
}

func ParsePubKey(from io.Reader) (*ScriptPubKey, error) {
  length := ReadVarInt(from)
  buf := make([]byte, length)
  total, err := from.Read(buf)
  if total < int(length) {
    err = errors.Join(err, errors.New("Invalid Script length decoded for pub key"))
  }
  return &ScriptPubKey{
    hex.EncodeToString(buf),
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
	return &Script{
    hex.EncodeToString(buf),
  }, err
}
