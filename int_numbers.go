package bitcoinlib

import (
	"encoding/binary"
)


type Int struct {
  value [32]byte
}

func FromInt(value int) Int {
  buf := make([]byte, 24)
  encoded := uint64(value)  
  buf = binary.BigEndian.AppendUint64(buf, encoded)
  buf_converted := [32]byte(buf)
  return Int{
    value: buf_converted,
  }
}

func fromArray(array [4]uint64) Int {
  value := []byte{}
  for _, number := range array {
    value = binary.BigEndian.AppendUint64(value, number)
  }
  return Int {
    value: [32]byte(value),
  }
}

func (i Int) getByteRepresentation() [4]uint64{
  reversed :=  []uint64{}
  left := 0
  right := 8
  for right <= 32 {
    reversed = append([]uint64{binary.BigEndian.Uint64(i.value[left:right])}, reversed...)
    left = right
    right += 8
  }
  return [4]uint64(reversed)
}

func addBytes(a uint64, b uint64, carry uint64) (c uint64, res uint64) {
  res = a + b + carry
  if res < a || res < b || res - a != b {
    c = 1
  }
  return
}

func (i Int) Add(other Int) Int {
  i_reversed := i.getByteRepresentation()
  other_reversed := other.getByteRepresentation()
  var carry uint64
  var partial uint64
  final := []uint64{}
  for index, i_value := range i_reversed {
    carry, partial = addBytes(i_value, other_reversed[index], carry)
    final = append([]uint64{partial}, final...)
  }
  return fromArray([4]uint64(final))
}
