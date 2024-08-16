package bitcoinlib

import (
	"encoding/binary"
	"encoding/hex"
)


type Int struct {
  value [32]byte
}

func (i Int) String() string {
  complement := i.getComplement()
  return hex.EncodeToString(complement.value[:])
}

func FromInt(value int) Int {
  original := value
  if value < 0 {
    value = -1 *  value
  }
  buf := make([]byte, 24)
  encoded := uint64(value)  
  buf = binary.BigEndian.AppendUint64(buf, encoded)
  buf_converted := [32]byte(buf)
  result := Int{
    value: buf_converted,
  }
  if original < 0 {
    result = result.getComplement()
  }
  return result
}

func FromArray(array [4]uint64) Int {
  value := []byte{}
  for _, number := range array {
    value = binary.BigEndian.AppendUint64(value, number)
  }
  return Int {
    value: [32]byte(value),
  }
}

func (i Int) getComplement() Int {
  if i.value[0] & 0x80 == 0x80 {
    new_array := [32]byte{}
    for index, value := range i.value {
      new_array[index] = ^value
    }
    new_value := Int{
      value: new_array,
    }
    return new_value.Add(FromInt(1))
  }
  return i
}

func (i Int) GetByteRepresentation() [4]uint64{
  value := i
  reversed :=  []uint64{}
  left := 0
  right := 8
  for right <= 32 {
    reversed = append([]uint64{binary.BigEndian.Uint64(value.value[left:right])}, reversed...)
    left = right
    right += 8
  }
  return [4]uint64(reversed)
}


func addBytes(a uint64, b uint64, carry uint64) (c uint64, res uint64) {
  res = a + b + carry
  //The second condition covers the case where a = b = 0xffffffffffffffff and c = 1
  if (res < a || res < b) || ((res == a) && (b != 0)) {
    c = 1
  }
  return
}

func (i Int) Eq(other Int) bool {
  for index, value := range i.value {
    if value != other.value[index] {
      return false
    }
  }
  return true
}

func (i Int) Ne(other Int) bool {
  return !i.Eq(other)
}

func (i Int) Geq(other Int) bool {
  return i.Ge(other) || i.Eq(other)
}

func (i Int) Le(other Int) bool {
  return !i.Geq(other)
}

func (i Int) Leq(other Int) bool {
  return !i.Ge(other)
}

func (i Int) Ge(other Int) bool {
  for index, value := range i.value {
    if value > other.value[index] {
      return true
    }
  }
  return false
}

func (i Int) add(other Int) (Int, uint64) {
i_reversed := i.GetByteRepresentation()
  other_reversed := other.GetByteRepresentation()
  var carry uint64
  var partial uint64
  final := []uint64{}
  for index, i_value := range i_reversed {
    carry, partial = addBytes(i_value, other_reversed[index], carry)
    //fmt.Printf("%d + %d + %d = %d\n", i_value, other_reversed[index], carry, partial)
    final = append([]uint64{partial}, final...)
  }
  return FromArray([4]uint64(final)), carry
}

func (i Int) Add(other Int) Int {
  result, _ := i.add(other)
  return result
}

func (i Int) negate() Int {
  i.value[0] = i.value[0] + 0x80
  complement := i.getComplement()
  i.value[0] = i.value[0] + 0x80
  return complement
}

func (i Int) Sub(other Int) Int {
  other_complement := other.negate()
  result, carry := i.add(other_complement)
  if carry == 1 {
    return result
  }
  return result.negate()
}
