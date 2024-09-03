package bitcoinlib

import (
	"encoding/binary"
	"encoding/hex"
	"math/big"
	"strings"
)

var MAX_INT_64 uint64 = 0x8000000000000000

var MAX_INT_32 string = "0xffffffffffffffffffffffffffffffffffffffffffffffff"

var ZERO Int = Int{
	value: big.NewInt(0),
}
var ONE Int = Int{
	value: big.NewInt(1),
}
var TWO Int = Int{
	value: big.NewInt(2),
}

var THREE Int = Int{
	value: big.NewInt(3),
}

var FOUR Int = Int{
	value: big.NewInt(4),
}

type Int struct {
	value *big.Int
}

func (i Int) String() string {
	representation := make([]byte, 32)
	return hex.EncodeToString(i.value.FillBytes(representation))

}

func FromInt(value int) Int {
	return Int{big.NewInt(int64(value))}
}

func FromArray(array [4]uint64) Int {
	value := [4]big.Word{}
	for index, number := range array {
		value[3-index] = big.Word(number)
	}

	return Int{
		value: big.NewInt(0).SetBits(value[:]),
	}
}

// Expects a string in the format 0x<Number>
func FromHexString(str string) Int {
	str = str[2:]
	value := [80]byte{}
	if len(str) < (16 * 4) {
		str = strings.Repeat("0", (16*4)-len(str)) + str
	} else if len(str) > (16 * 4) {
		str = str[len(str)-(16*4):]
	}
	total, err := hex.Decode(value[48:], []byte(str))
	if err != nil || total != 32 {
		return ZERO
	}
	result := []big.Word{}
	for i := 80; i > 0; i -= 8 {
		result = append(result, big.Word(binary.BigEndian.Uint64(value[i-8:i])))
	}

	number := Int{
		value: big.NewInt(0).SetBits(result),
	}
	return number
}

func (i Int) IntoBytes() [32]byte {
	buf := make([]byte, 32)
	i.value.FillBytes(buf)
	return [32]byte(buf)
}

func (i Int) Eq(other Int) bool {
	return i.value.Cmp(other.value) == 0
}

func (i Int) Ne(other Int) bool {
	return !i.Eq(other)
}

func (i Int) Geq(other Int) bool {
	return i.value.Cmp(other.value) >= 0
}

func (i Int) Le(other Int) bool {
	return !i.Geq(other)
}

func (i Int) Leq(other Int) bool {
	return !i.Ge(other)
}

func (i Int) Ge(other Int) bool {
	return i.value.Cmp(other.value) > 0
}

func (i Int) Add(other Int) Int {
	return Int{
		value: big.NewInt(0).Add(i.value, other.value),
	}
}
func (i Int) Sub(other Int) Int {
	return Int{
		value: big.NewInt(0).Sub(i.value, other.value),
	}
}

func (i Int) Mul(other Int) Int {
	return Int{
		value: big.NewInt(0).Mul(i.value, other.value),
	}
}

// Performs integer division for positive numbers
func (i Int) Div(other Int) Int {
	return Int{
		value: big.NewInt(0).Div(i.value, other.value),
	}
}

// Performs modulus other Integer
func (i Int) Mod(other Int) Int {
	return Int{
		value: big.NewInt(0).Mod(i.value, other.value),
	}
}

func (i Int) Exp(other Int, mod Int) Int {
	return Int{
		value: big.NewInt(0).Exp(i.value, other.value, mod.value),
	}
}

// Raises a number by -1
func (i Int) ExpNeg(base Int) Int {
	return i.Exp(base.Sub(TWO), base)
}
