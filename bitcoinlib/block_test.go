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
