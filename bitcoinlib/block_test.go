package bitcoinlib_test

import (
	"bitcoinlib/bitcoinlib"
	"bytes"
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