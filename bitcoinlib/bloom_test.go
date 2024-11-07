package bitcoinlib_test

import (
	"bitcoinlib/bitcoinlib"
	"testing"
)

func TestBloomGeneration(t *testing.T) {
	expected := "0420"
	filter := bitcoinlib.NewBloomFilter(2)
	filter.Set160([]byte("hello world"))
	filter.Set160([]byte("goodbye"))
	if filter.String() != expected {
		t.Fatalf("Expected filter to be %s but got %s", expected, filter)
	}
}

func TestBloomMurmurGeneration(t *testing.T) {
	expected := "4000600a080000010940"
	filter := bitcoinlib.NewBloomFilter(10)
	params := &bitcoinlib.MurmurParams{
		FunctionCount: 5,
		Tweak:         99,
	}
	filter.Set([]byte("Hello World"), params)
	filter.Set([]byte("Goodbye!"), params)
	if filter.String() != expected {
		t.Fatalf("Expected filter to be %s but got %s instead", expected, filter)
	}
}
