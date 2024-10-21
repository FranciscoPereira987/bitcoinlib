package bitcoinlib_test

import (
	"bitcoinlib/bitcoinlib"
	"bytes"
	"encoding/hex"
	"testing"
)

func TestParseExcercise1(t *testing.T) {
	message, _ := hex.DecodeString("f9beb4d976657261636b000000000000000000005df6e0e2")
	blockMessage := bitcoinlib.NewNetworkMessage(false)
	err := blockMessage.Parse(bytes.NewReader(message))
	if err != nil {
		t.Fatalf("Error parsing messag: %s", err)
	}
	expectedCommand := "verack"
	parsedCommand := blockMessage.GetCommand()
	if expectedCommand != parsedCommand {
		t.Fatalf("Expected command: %s; But got: %s", expectedCommand, parsedCommand)
	}
	expectedPayload := hex.EncodeToString([]byte{})
	parsedPayload := hex.EncodeToString(blockMessage.GetPayload())
	if expectedPayload != parsedPayload {
		t.Fatalf("Expected payload: %s; But got: %s", expectedCommand, parsedCommand)
	}
}	

func TestSerializing(t *testing.T) {
	msg, _ := hex.DecodeString("f9beb4d976657261636b000000000000000000005df6e0e2")
	nmsg := bitcoinlib.NewNetworkMessage(false)
	if nmsg.Parse(bytes.NewReader(msg)) != nil {
		t.Fatal("Failed parsing first message")
	}
	if hex.EncodeToString(nmsg.Serialize()) != hex.EncodeToString(msg) {
		t.Fatal("First message different from original one")
	}
	msg, _ = hex.DecodeString("f9beb4d976657273696f6e0000000000650000005f1a69d2721101000100000000000000bc8f5e5400000000010000000000000000000000000000000000ffffc61b6409208d010000000000000000000000000000000000ffffcb0071c0208d128035cbc97953f80f2f5361746f7368693a302e392e332fcf05050001")
	if nmsg.Parse(bytes.NewReader(msg)) != nil {
		t.Fatal("Failed parsing second message")
	}
	if hex.EncodeToString(nmsg.Serialize()) != hex.EncodeToString(msg) {
		t.Fatal("Second serialization was different from actual message")
	}
}

func TestVersionMessage(t *testing.T) {
  m := bitcoinlib.NewVersionMessage()
	m.Nonce = 0
  m.Timestamp = 0 
  expected := "7f11010000000000000000000000000000000000000000000000000000000000000000000000ffff00000000208d000000000000000000000000000000000000ffff00000000208d0000000000000000182f70726f6772616d6d696e67626974636f696e3a302e312f0000000000"
	serialized := hex.EncodeToString(m.Serialize())
	if expected != serialized {
		t.Fatalf("Expected serialization\n%s\nBut got instead\n%s\n", expected, serialized)
	}
}

func TestGetHeadersMessage(t *testing.T) {
  m := bitcoinlib.NewGetHeadersMessage("0000000000000000001237f46acddf58578a37e213d2a6edc4884a2fcad05ba3", "")
  toParse := "7f11010001a35bd0ca2f4a88c4eda6d213e2378a5758dfcd6af437120000000000000000000000000000000000000000000000000000000000000000000000000000000000" 
  serialized := hex.EncodeToString(m.Serialize())
  if serialized != toParse {
    t.Fatalf("Serialization fail \n %s\n!=\n%s", toParse, serialized)
  }
  m = bitcoinlib.NewGetHeadersMessage("", "")
  asBytes, _ := hex.DecodeString(toParse)
  r,_ := m.Parse(asBytes)
  if hex.EncodeToString(r.Serialize()) != toParse {
    t.Fatalf("Failed parsing")
  }
}
