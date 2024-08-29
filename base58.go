package bitcoinlib

import (
	"encoding/hex"
	"math/big"
	"strings"
)

const ALPHABET = "123456789ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxyz" 
var F8 Int = FromInt(58)


func IntoBase58(num string) string {
  leading_zeroes := 0
  for  ;num[leading_zeroes] == '0' && num[leading_zeroes+1] == '0'; leading_zeroes += 2 {}
  prefix := strings.Repeat("1", leading_zeroes / 2)
  result := ""
  number, _ := hex.DecodeString(num)

  int_num := Int{
    value: big.NewInt(0).SetBytes(number),
  }
  for int_num.Ge(ZERO) {
    index := int_num.Mod(F8).value.Int64()
    int_num = int_num.Div(F8)
    result = string(ALPHABET[index]) + result 
  }

  return prefix + result 
}
