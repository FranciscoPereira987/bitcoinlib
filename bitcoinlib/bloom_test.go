package bitcoinlib_test

import (
	"bitcoinlib/bitcoinlib"
	"testing"
)

func TestBloomGeneration(t *testing.T) {
	expected := "c000"
	filter := bitcoinlib.NewBloomFilter(10)
	filter.Set160([]byte("hello world"))
	filter.Set160([]byte("goodbye"))
	if filter.String() != expected {
		t.Fatalf("Expected filter to be %s but got %s", expected, filter)
	}
}
