package bitcoinlib_test

import (
	"bitcoinlib"
	"testing"
)


func TestSignature(t *testing.T) {
  z := bitcoinlib.FromHexString("0xbc62d4b80d9e36da29c16c5d4d9f11731f36052c72401a76c23c0fb5a9b74423") 
  r := bitcoinlib.FromHexString("0x37206a0610995c58074999cb9767b87af4c4978db68c06e8e6e81d282047a7c6") 
  s := bitcoinlib.FromHexString("0x8ca63759c1157ebeaec0d03cecca119fc9a75bf8e6d0fa65c841c8e2738cdaec") 

  px := bitcoinlib.FromHexString("0x04519fac3d910ca7e7138f7013706f619fa8f033e6ec6e09370ea38cee6a7574")
  py := bitcoinlib.FromHexString("0x82b51eab8c27c66e26c858a079bcdf4f1ada34cec420cafc7eac1a42216fb6c4")

  point, _ := bitcoinlib.NewS256Point(px, py)
  s_inv := s.Exp(bitcoinlib.ORDER.Sub(bitcoinlib.FromInt(2)), bitcoinlib.ORDER)
  u := z.Mul(s_inv).Mod(bitcoinlib.ORDER)
  v := r.Mul(s_inv).Mod(bitcoinlib.ORDER)
  
  result, _ := bitcoinlib.S256Mul(point, v).Add(bitcoinlib.S256Mul(bitcoinlib.G(), u))
  if  !bitcoinlib.S256Verifyr(result, r) {
    t.Fatal("Failed example verification")
  }
}
