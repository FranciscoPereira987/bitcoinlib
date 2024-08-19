package bitcoinlib

import (
	"encoding/binary"
	"encoding/hex"
)


type Int struct {
  value [32]byte
}

func (i Int) String() string {
  return hex.EncodeToString(i.value[:])
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
    result = result.negate()
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
    }else if value < other.value[index] {
      return false
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
  if i.value[0] & 0x80 == 0x80 {
    return i.getComplement()
  }else {
    new_array := [32]byte{}
    for index, value := range i.value {
      new_array[index] = ^value
    }
    new_value := Int{
      value: new_array,
    }
    return new_value.Add(FromInt(1))
  }
}

func (i Int) Sub(other Int) Int {
  other_complement := other.negate()
  result, carry := i.add(other_complement)
  if carry == 1 {
    return result
  }
  return result
}

//Shifts the number to the right by one
func (i Int) ShiftRight() Int {
  var carry uint8
  shifted := [32]byte{}
  for index, value := range i.value {
    shifted[index] = value >> 1
    shifted[index] += carry
    if value & 0x01 != 0{
      carry = 0x80 
    } else {
      carry = 0
    }
  }
  return Int{
    value: shifted,
  }
}

func (i Int) Mul(other Int) Int {
  result := FromInt(0)
  partial := i 
  for other.Ge(FromInt(0)) {
    if other.value[31] & 0x01 == 1 {
      result = result.Add(partial)
    }
    partial = partial.Add(partial) 
    other = other.ShiftRight()
  }
  return result
}

//Performs integer division for positive numbers
func (i Int) Div(other Int) Int {
  if other.Ge(i) {
    return FromInt(0)
  }
  actual := i.ShiftRight()
  multiplier := FromInt(2)
  for other.Le(actual) {
    actual = actual.ShiftRight()
    multiplier = multiplier.Mul(FromInt(2))
  }
  if actual.Eq(other) {
    return multiplier 
  }
  right_multiplier := multiplier
  left_multiplier := multiplier.ShiftRight()
  middle_multiplier := left_multiplier.Add(right_multiplier.Sub(left_multiplier).ShiftRight())
  for right_multiplier.Ge(left_multiplier.Add(FromInt(1))) {
    value := other.Mul(middle_multiplier)
    if value.Le(i) {
      left_multiplier = middle_multiplier
    }else if value.Ge(i) {
      right_multiplier = middle_multiplier 
    } else {
      return middle_multiplier 
    }
    middle_multiplier = left_multiplier.Add(right_multiplier.Sub(left_multiplier).ShiftRight())
  }
  return middle_multiplier
}

//Performs modulus other Integer
func (i Int) Mod(other Int) Int {
  if i.Ne(i.getComplement()) {
    mul := i.getComplement().Div(other).Add(FromInt(1))
    i = i.Add(other.Mul(mul)) 
  }
  sum := i.Div(other)
  
  return i.Sub(other.Mul(sum))
}
