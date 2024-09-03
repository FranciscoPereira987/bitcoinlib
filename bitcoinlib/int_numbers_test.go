package bitcoinlib_test

import (
	"bitcoinlib/bitcoinlib"
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
  }
  expected := []string{"0000000000000000000000000000000000000000000000000000000000000001",}
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
      t.Fatalf("Failed at index %d => Expected: %s but got: %s", index, expected[index].String(), result.String())
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
    }, 
    {
      bitcoinlib.FromInt(0x34),
      bitcoinlib.FromInt(0x32),
    },
  }
  expected := []bitcoinlib.Int{bitcoinlib.FromInt(-21),
  bitcoinlib.FromArray([4]uint64{0, 0, 0, 0xfffffffffffffffe}),
  bitcoinlib.FromInt(0),
  bitcoinlib.FromArray([4]uint64{0xfffffffffffffffe, 0xffffffffffffffff, 0xffffffffffffffff, 0xffffffffffffffff}),
  bitcoinlib.FromInt(0x02),
  }
  for index, a_b := range values {
    result := a_b[0].Sub(a_b[1])
    if result.Ne(expected[index]) {
      t.Fatalf("at index %d => Expected: %s but got: %s", index, expected[index].String(), result.String())
    }
  }
}



func TestMultiplication(t *testing.T) {
  cases := []bitcoinlib.Int{
    bitcoinlib.FromInt(2),
    bitcoinlib.FromInt(10000),
    bitcoinlib.FromInt(22222),
    bitcoinlib.FromInt(2),

  }

  multipliers := []bitcoinlib.Int{
    bitcoinlib.FromInt(3),
    bitcoinlib.FromInt(1000),
    bitcoinlib.FromInt(524797),
    bitcoinlib.FromArray([4]uint64{0x4000000000000000, 0, 0, 0}),
  }

  expected := []bitcoinlib.Int{
    bitcoinlib.FromInt(6),
    bitcoinlib.FromInt(10_000_000),
    bitcoinlib.FromInt(22222 * 524797),
    bitcoinlib.FromArray([4]uint64{0x8000000000000000, 0, 0, 0}),
  }

  for index, value := range cases {
    result := value.Mul(multipliers[index])
    if result.Ne(expected[index]) {
      t.Fatalf("Failed at index %d: Expected %s but got %s", index, expected[index].String(), result.String())
    }
  }
}

func TestDivision(t *testing.T) {
  cases := [][2]bitcoinlib.Int{
    {
      bitcoinlib.FromInt(10),
      bitcoinlib.FromInt(5),
    },
    {
      bitcoinlib.FromInt(100),
      bitcoinlib.FromInt(13),
    },
    {
      bitcoinlib.FromInt(576460752303423488),
      bitcoinlib.FromInt(2),
    },
    {
      bitcoinlib.FromArray([4]uint64{0x7fffffffffffffff, 0, 0, 0}),
      bitcoinlib.FromInt(2),
    }, 
  }

  results := []bitcoinlib.Int{
    bitcoinlib.FromInt(2),
    bitcoinlib.FromInt(7),
    bitcoinlib.FromInt(288230376151711744),
    bitcoinlib.FromArray([4]uint64{0x7fffffffffffffff, 0, 0, 0}).Div(bitcoinlib.TWO),
  }

  for index, value := range cases {
    result := value[0].Div(value[1])
    if result.Ne(results[index]) {
      t.Fatalf("Failed at index %d: Expected %s but got %s", index, results[index].String(), result.String()) 
    }
  }
}

func TestMod(t *testing.T) {
  cases := []bitcoinlib.Int{
    bitcoinlib.FromInt(0),
    bitcoinlib.FromInt(10),
    bitcoinlib.FromInt(1234),
    bitcoinlib.FromArray([4]uint64{0x7fffffffffffffff, 0xffffffffffffffff, 0xffffffffffffffff, 0xffffffffffffffff}),
    bitcoinlib.FromInt(-23),
  }
  modulus := []bitcoinlib.Int{
    bitcoinlib.FromInt(10),
    bitcoinlib.FromInt(5),
    bitcoinlib.FromInt(3),
    bitcoinlib.FromInt(2),
    bitcoinlib.FromInt(5),
  }
  expected := []bitcoinlib.Int{
    bitcoinlib.FromInt(0),
    bitcoinlib.FromInt(0),
    bitcoinlib.FromInt(1),
    bitcoinlib.FromInt(1),
    bitcoinlib.FromInt(2),
  }
  for index, value := range cases {
    result := value.Mod(modulus[index])
    if result.Ne(expected[index]) {
      t.Fatalf("Failed at index %d: Expected %s but got %s", index, expected[index].String(), result.String())
    }
  }
}

