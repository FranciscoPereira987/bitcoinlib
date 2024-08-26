package bitcoinlib_test

import (
	"bitcoinlib"
	"crypto/sha256"
	"encoding/hex"
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

func TestVerifyingSignatures(t *testing.T) {
  e := bitcoinlib.FromInt(12345)
  z_sign := sha256.Sum256([]byte("Programming Bitcoin!"))
  z := bitcoinlib.FromHexString("0x" + hex.EncodeToString(z_sign[:]))
  pk := bitcoinlib.NewPrivateKey(e)
  signature := pk.Sign(z)
  if !signature.Verify(z) {
    t.Fatalf("Failed to verify something I've just signed\n%s", signature)
  }
}

func TestSecValuesUncompressed(t *testing.T) {
  e_vals := []bitcoinlib.Int{
    bitcoinlib.FromInt(5000),
    bitcoinlib.FromInt(2018).Exp(bitcoinlib.FromInt(5), bitcoinlib.ORDER),
    bitcoinlib.FromHexString("0xdeadbeef12345"),
  }
  results := []string{
    "04ffe558e388852f0120e46af2d1b370f85854a8eb0841811ece0e3e03d282d57c315dc72890a4f10a1481c031b03b351b0dc79901ca18a00cf009dbdb157a1d10",
    "04027f3da1918455e03c46f659266a1bb5204e959db7364d2f473bdf8f0a13cc9dff87647fd023c13b4a4994f17691895806e1b40b57f4fd22581a4f46851f3b06",
    "04d90cd625ee87dd38656dd95cf79f65f60f7273b67d3096e68bd81e4f5342691f842efa762fd59961d0e99803c61edba8b3e3f7dc3a341836f97733aebf987121",
  }
  for index, value := range e_vals {
    key := bitcoinlib.NewPrivateKey(value)
    uncompressed_sec := key.Sec(bitcoinlib.UNCOMPRESSED)
    result_string := hex.EncodeToString(uncompressed_sec)
    if results[index] !=  result_string{
      t.Fatalf("Failed at index %d => %s \n != \n %s", index, results[index], result_string)
    }
  }
}

func TestSecValuesCompressed(t *testing.T) {
  e_vals := []bitcoinlib.Int{
    bitcoinlib.FromInt(5001),
    bitcoinlib.FromInt(2019).Exp(bitcoinlib.FromInt(5), bitcoinlib.ORDER),
    bitcoinlib.FromHexString("0xdeadbeef54321"),
  }
  results := []string {
    "0357a4f368868a8a6d572991e484e664810ff14c05c0fa023275251151fe0e53d1",
    "02933ec2d2b111b92737ec12f1c5d20f3233a0ad21cd8b36d0bca7a0cfa5cb8701",
    "0296be5b1292f6c856b3c5654e886fc13511462059089cdf9c479623bfcbe77690",
  }
  for index, value := range e_vals {
    key := bitcoinlib.NewPrivateKey(value)
    uncompressed_sec := key.Sec(bitcoinlib.COMPRESSED)
    result_string := hex.EncodeToString(uncompressed_sec)
    if results[index] !=  result_string{
      t.Fatalf("Failed at index %d => %s \n != \n %s", index, results[index], result_string)
    }
  }

}
