package bitcoinlib_test

import (
	"bitcoinlib/bitcoinlib"
	"bytes"
	"testing"
)

func TestVersionParsing(t *testing.T) {
	version_buf := bytes.NewReader([]byte{0x01, 0x00, 0x00, 0x00, 0xff})
	result, err := bitcoinlib.NewVersionFrom(version_buf)
	if err != nil {
		t.Fatalf("Faile with error %s", err)
	}
	expected := bitcoinlib.NewVersion(1)
	if expected.Ne(*result) {
		t.Fatalf("%v is different from %v", result, expected)
	}
}
