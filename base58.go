package bitcoinlib

import (
	"encoding/hex"
	"math/big"
	"strings"
)

const ALPHABET = "123456789ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxyz" 
var F8 Int = FromInt(58)

func firstTwoZeores(num string) bool {
  return num[:2] == "00"
}

func IntoBase58(num string) string {
  leading_zeroes := 0
  for  ;firstTwoZeores(num[2*leading_zeroes:]); leading_zeroes++ {}
  prefix := strings.Repeat("1", leading_zeroes)
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
