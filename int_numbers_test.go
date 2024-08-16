package bitcoinlib_test

import (
	"bitcoinlib"
	"testing"
)

func compare(a [4]uint64, b [4]uint64) bool {
  for index, value := range a {
    if value != b[index] {
      return false
    }
  }
  return true
}

func TestRepresentation(t *testing.T) {
  values := []bitcoinlib.Int{
    bitcoinlib.FromInt(1),
    bitcoinlib.FromArray([4]uint64{0x8000000000000000, 0, 0, 1}),
  }
  expected := []string{"0000000000000000000000000000000000000000000000000000000000000001","7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff"}
  for index, value := range values {
    if value.String() != expected[index] {
      t.Fatalf("compared Failed: %s vs %s", expected[index], value.String())
    }
  }
}

func TestAddtion(t *testing.T) {
  values := [][2]bitcoinlib.Int{{
    bitcoinlib.FromInt(234),
    bitcoinlib.FromInt(255),
  },
    {
      bitcoinlib.FromArray([4]uint64{0, 0, 0, 0xffffffffffffffff}),
        bitcoinlib.FromInt(1),
    },
    {
      bitcoinlib.FromArray([4]uint64{0, 0xffffffffffffffff, 0xffffffffffffffff, 0xffffffffffffffff}),
      bitcoinlib.FromArray([4]uint64{0, 0xffffffffffffffff, 0xffffffffffffffff, 0xffffffffffffffff}),
    },
      }
  expected := []bitcoinlib.Int{bitcoinlib.FromInt(489),
  bitcoinlib.FromArray([4]uint64{0, 0, 1, 0}),
  bitcoinlib.FromArray([4]uint64{1, 0xffffffffffffffff, 0xffffffffffffffff, 0xfffffffffffffffe}),

  }
  for index, a_b := range values {
    result := a_b[0].Add(a_b[1])
    if result.Ne(expected[index]) {
      t.Fatalf("Expected: %s but got: %s", expected[index].String(), result.String())
    }
  }
}

func TestSub(t *testing.T) {
values := [][2]bitcoinlib.Int{{
    bitcoinlib.FromInt(234),
    bitcoinlib.FromInt(255),
  },
    {
      bitcoinlib.FromArray([4]uint64{0, 0, 0, 0xffffffffffffffff}),
        bitcoinlib.FromInt(1),
    },
    {
      bitcoinlib.FromArray([4]uint64{0, 0xffffffffffffffff, 0xffffffffffffffff, 0xffffffffffffffff}),
      bitcoinlib.FromArray([4]uint64{0, 0xffffffffffffffff, 0xffffffffffffffff, 0xffffffffffffffff}),
    },
    {
      bitcoinlib.FromArray([4]uint64{0xffffffffffffffff, 0, 0, 0}),
      bitcoinlib.FromArray([4]uint64{0, 0, 0, 1}),
    }, }
  expected := []bitcoinlib.Int{bitcoinlib.FromInt(-21),
  bitcoinlib.FromArray([4]uint64{0, 0, 0, 0xfffffffffffffffe}),
  bitcoinlib.FromInt(0),
  bitcoinlib.FromArray([4]uint64{0x7ffffffffffffffe, 0xffffffffffffffff, 0xffffffffffffffff, 0xffffffffffffffff}),
  }
  for index, a_b := range values {
    result := a_b[0].Sub(a_b[1])
    if result.Ne(expected[index]) {
      t.Fatalf("Expected: %s but got: %s", expected[index].String(), result.String())
    }
  }
}
