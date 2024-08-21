package bitcoinlib

import (
	"encoding/binary"
	"encoding/hex"
	"math/bits"
	"strings"
)

var MAX_INT_64 uint64 = 0x8000000000000000

var ZERO Int = Int{
  value: [10]uint64{},
} 
var ONE Int = Int {
  value: [10]uint64{0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
} 
var TWO Int = Int {
  value: [10]uint64{0, 0, 0, 0, 0, 0, 0, 0, 0, 2},
} 

var THREE Int = Int {
  value: [10]uint64{0, 0, 0, 0, 0, 0, 0, 0, 0, 3},
} 



type Int struct {
  value [10]uint64
}

func (i Int) String() string {
  representation := []byte{}
  for _, value := range i.value[6:] {
    representation = binary.BigEndian.AppendUint64(representation, value)
  }
  return hex.EncodeToString(representation)
}

func FromInt(value int) Int {
  original := value
  if value < 0 {
    value = -1 *  value
  }
  buf := [10]uint64{}
  encoded := uint64(value)  
  buf[9] = encoded
  result := Int{
    value: buf,
  }
  if original < 0 {
    result = result.negate()
  }
  return result
}

func fromArray(array [10]uint64) Int {
  return Int {
    value: array,
  }
}

func FromArray(array [4]uint64) Int {
  value := make([]uint64, 6)
  
  return Int {
    value: [10]uint64(append(value, array[:]...)),
  }
}

//Expects a string in the format 0x<Number>
func FromHexString(str string) Int {
  str = str[2:]
  value := [80]byte{}
  if len(str) < (16 * 4) {
    str = strings.Repeat("0", (16 * 4) - len(str)) + str
  }else if len(str) > (16 * 4) {
    str = str[len(str) - (16 * 4):]
  }
  total, err := hex.Decode(value[48:], []byte(str))
  if err != nil || total != 32 {
    return ZERO 
  }
  result := []uint64{}
  for i := 0; i < 80; i += 8 {
    result = append(result, binary.BigEndian.Uint64(value[i:i+8]))
  }
  return Int{
    value:  [10]uint64(result),
  }
}

func (i Int) getComplement() Int {
  if i.value[0] & MAX_INT_64 == MAX_INT_64 {
    new_array := [10]uint64{}
    for index, value := range i.value {
      new_array[index] = ^value
    }
    new_value := Int{
      value: new_array,
    }
    return new_value.Add(ONE)
  }
  return i
}

func (i Int) GetByteRepresentation() [10]uint64{
  reversed :=  [10]uint64{}
  for index := 0; index < 10; index++ {
    reversed[index] = i.value[9-index]
  }
  return reversed
}


func addBytes(a uint64, b uint64, carry uint64) (c uint64, res uint64) {
  res, c = bits.Add64(a, b, carry) 
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
  var carry uint64
  var partial uint64
  final := [10]uint64{}
  for index := len(i.value)-1; index >= 0; index-- {
    carry, partial = addBytes(i.value[index], other.value[index], carry)
    //fmt.Printf("%d + %d + %d = %d\n", i_value, other_reversed[index], carry, partial)
    final[index] = partial 
  }
  return fromArray(final), carry
}

func (i Int) Add(other Int) Int {
  result, _ := i.add(other)
  return result
}

func (i Int) negate() Int {
  if i.value[0] & MAX_INT_64 == MAX_INT_64 {
    return i.getComplement()
  }else {
    new_array := [10]uint64{}
    for index, value := range i.value {
      new_array[index] = ^value
    }
    new_value := Int{
      value: new_array,
    }
    return new_value.Add(ONE)
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
  var carry uint64
  shifted := [10]uint64{}
  for index, value := range i.value {
    shifted[index] = value >> 1
    shifted[index] += carry
    if value & 0x01 != 0{
      carry = MAX_INT_64 
    } else {
      carry = 0
    }
  }
  return Int{
    value: shifted,
  }
}

func (i Int) Mul(other Int) Int {
  result := ZERO 
  partial := i 
  for other.Ge(ZERO) {
    if other.value[9] & 0x01 == 1 {
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
    return ZERO 
  }
  actual := i.ShiftRight()
  multiplier := TWO 
  for other.Le(actual) {
    actual = actual.ShiftRight()
    multiplier = multiplier.Mul(TWO)
  }
  if actual.Eq(other) {
    return multiplier 
  }
  right_multiplier := multiplier
  left_multiplier := multiplier.ShiftRight()
  middle_multiplier := left_multiplier.Add(right_multiplier.Sub(left_multiplier).ShiftRight())
  for right_multiplier.Ge(left_multiplier.Add(ONE)) {
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
    mul := i.getComplement().Div(other).Add(ONE)
    i = i.Add(other.Mul(mul)) 
  }
  sum := i.Div(other)
  
  return i.Sub(other.Mul(sum))
}
