package bitcoinlib_test

import (
	"bitcoinlib/bitcoinlib"
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"testing"
)

func TestBlockParsingAndSerializing(t *testing.T) {
	blockString := "020000208ec39428b17323fa0ddec8e887b4a7c53b8c0a0a220cfd0000000000000000005b0750fce0a889502d40508d39576821155e9c9e3f5c3157f961db38fd8b25be1e77a759e93c0118a4ffd71d"
	asBytes, _ := hex.DecodeString(blockString)
	block := bitcoinlib.NewBlock()
	if block.Parse(bytes.NewReader(asBytes)) != nil {
		t.Fatal("Failed parsing block")
	}
	serialized := block.Serialize()
	if hex.EncodeToString(serialized) != blockString {
		t.Fatal("Failed reserializing block")
	}
}

func TestBlockHash(t *testing.T) {
	blockString := "020000208ec39428b17323fa0ddec8e887b4a7c53b8c0a0a220cfd0000000000000000005b0750fce0a889502d40508d39576821155e9c9e3f5c3157f961db38fd8b25be1e77a759e93c0118a4ffd71d"
	asBytes, _ := hex.DecodeString(blockString)
	block := bitcoinlib.NewBlock()
	if block.Parse(bytes.NewReader(asBytes)) != nil {
		t.Fatal("Failed parsing block to create hash")
	}
	expectedHash := "0000000000000000007e9e4c586439b0cdbe13b1370bdd9435d76a644d047523"
	hash := block.Hash()
	if expectedHash != hash {
		t.Fatalf("Expected block hash: %s\nBut got: %s\n", expectedHash, hash)
	}
}

func TestBlockBIP9(t *testing.T) {
	BIP9blockString, _ := hex.DecodeString("000000208ec39428b17323fa0ddec8e887b4a7c53b8c0a0a220cfd0000000000000000005b0750fce0a889502d40508d39576821155e9c9e3f5c3157f961db38fd8b25be1e77a759e93c0118a4ffd71d")
	NonBIP9BlockString, _ := hex.DecodeString("000000038ec39428b17323fa0ddec8e887b4a7c53b8c0a0a220cfd0000000000000000005b0750fce0a889502d40508d39576821155e9c9e3f5c3157f961db38fd8b25be1e77a759e93c0118a4ffd71d")
	block1 := bitcoinlib.NewBlock()
	block2 := bitcoinlib.NewBlock()

	block1.Parse(bytes.NewReader(BIP9blockString))
	block2.Parse(bytes.NewReader(NonBIP9BlockString))
	
	if !block1.BIP9() {
		t.Fatal("Block 1 was not Identified ad a BIP9 block")
	}

	if block2.BIP9() {
		t.Fatal("Block 2 was Identified as a BIP9 block")
	}
}

func TestBIP91(t *testing.T) {
	BIP91blockString, _ := hex.DecodeString("100000008ec39428b17323fa0ddec8e887b4a7c53b8c0a0a220cfd0000000000000000005b0750fce0a889502d40508d39576821155e9c9e3f5c3157f961db38fd8b25be1e77a759e93c0118a4ffd71d")
	NonBIP91BlockString, _ := hex.DecodeString("000000038ec39428b17323fa0ddec8e887b4a7c53b8c0a0a220cfd0000000000000000005b0750fce0a889502d40508d39576821155e9c9e3f5c3157f961db38fd8b25be1e77a759e93c0118a4ffd71d")
	block1 := bitcoinlib.NewBlock()
	block2 := bitcoinlib.NewBlock()

	block1.Parse(bytes.NewReader(BIP91blockString))
	block2.Parse(bytes.NewReader(NonBIP91BlockString))
	
	if !block1.BIP91() {
		t.Fatal("Block 1 was not Identified ad a BIP91 block")
	}

	if block2.BIP91() {
		t.Fatal("Block 2 was Identified as a BIP91 block")
	}	
}

func TestBIP141(t *testing.T) {
	BIP91blockString, _ := hex.DecodeString("020000008ec39428b17323fa0ddec8e887b4a7c53b8c0a0a220cfd0000000000000000005b0750fce0a889502d40508d39576821155e9c9e3f5c3157f961db38fd8b25be1e77a759e93c0118a4ffd71d")
	NonBIP91BlockString, _ := hex.DecodeString("000000038ec39428b17323fa0ddec8e887b4a7c53b8c0a0a220cfd0000000000000000005b0750fce0a889502d40508d39576821155e9c9e3f5c3157f961db38fd8b25be1e77a759e93c0118a4ffd71d")
	block1 := bitcoinlib.NewBlock()
	block2 := bitcoinlib.NewBlock()

	block1.Parse(bytes.NewReader(BIP91blockString))
	block2.Parse(bytes.NewReader(NonBIP91BlockString))
	
	if !block1.BIP141() {
		t.Fatal("Block 1 was not Identified ad a BIP141 block")
	}

	if block2.BIP141() {
		t.Fatal("Block 2 was Identified as a BIP141 block")
	}	
}

func TestBitsToTarget(t *testing.T) {
  number := binary.LittleEndian.Uint32([]byte{0xe9, 0x3c, 0x01, 0x18})
  expected := bitcoinlib.FromHexString("0x0000000000000000013ce9000000000000000000000000000000000000000000")
  
  

  got := bitcoinlib.BitsToTarget(number)
  if got.Ne(expected) {
    t.Fatalf("Expected: %s but got %s\n", expected.String(), got.String())
  }
}

func TestDifficulty(t *testing.T) {
  expected := bitcoinlib.FromInt(888171856257)
	blockString, _ := hex.DecodeString("020000008ec39428b17323fa0ddec8e887b4a7c53b8c0a0a220cfd0000000000000000005b0750fce0a889502d40508d39576821155e9c9e3f5c3157f961db38fd8b25be1e77a759e93c0118a4ffd71d")
  block := bitcoinlib.NewBlock()
  block.Parse(bytes.NewReader(blockString))
  genesis := bitcoinlib.BitsToTarget(0x1d00ffff)
  if block.Difficulty().Ne(expected) {
    t.Fatalf("Expected: %s but got: %s, diff: %s", expected.String(), block.Difficulty().String(),genesis.String()) 
  }
}

func TestPOW(t *testing.T) {
  block1String, _ := hex.DecodeString("04000000fbedbbf0cfdaf278c094f187f2eb987c86a199da22bbb20400000000000000007b7697b29129648fa08b4bcd13c9d5e60abb973a1efac9c8d573c71c807c56c3d6213557faa80518c3737ec1")
  block2String, _ := hex.DecodeString("04000000fbedbbf0cfdaf278c094f187f2eb987c86a199da22bbb20400000000000000007b7697b29129648fa08b4bcd13c9d5e60abb973a1efac9c8d573c71c807c56c3d6213557faa80518c3737ec0")

  block1 := bitcoinlib.NewBlock()
  block2 := bitcoinlib.NewBlock()

  block1.Parse(bytes.NewReader(block1String))
  block2.Parse(bytes.NewReader(block2String))

  if !block1.CheckPOW() {
    t.Fatal("Failed checking block1 Proof Of Work")
  }

  if block2.CheckPOW() {
    t.Fatal("Failed checking block2 Proof Of Work")
  }
}

func TestTargetCalculation(t *testing.T) {
  block1String, _ := hex.DecodeString("000000203471101bbda3fe307664b3283a9ef0e97d9a38a7eacd8800000000000000000010c8aba8479bbaa5e0848152fd3c2289ca50e1c3e58c9a4faaafbdf5803c5448ddb845597e8b0118e43a81d3")
  block2String, _ := hex.DecodeString("02000020f1472d9db4b563c35f97c428ac903f23b7fc055d1cfc26000000000000000000b3f449fcbe1bc4cfbcb8283a0d2c037f961a3fdf2b8bedc144973735eea707e1264258597e8b0118e5f00474")
  
  block1 := bitcoinlib.NewBlock()
  block2 := bitcoinlib.NewBlock()

  block1.Parse(bytes.NewReader(block1String))
  block2.Parse(bytes.NewReader(block2String))

  var expected uint32 = binary.LittleEndian.Uint32([]byte{0x7e, 0x8b, 0x01, 0x18}) 
  actual := block1.GetNextTarget(block2)

  if actual != expected {
    t.Fatalf("Expected new target %x but got %x", expected, actual)
  }
  
}
