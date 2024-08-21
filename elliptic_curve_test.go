package bitcoinlib_test

import (
	"bitcoinlib"
	"testing"
)

func TestOnCurve(t *testing.T) {
	prime := 223
	a, _ := bitcoinlib.NewFieldElement(prime, 0)
	b, _ := bitcoinlib.NewFieldElement(prime, 7)
	valid_points := [][2]int{{192, 105}, {17, 56}, {1, 193}}
	invalid_points := [][2]int{{200, 119}, {42, 99}}
	for _, point := range valid_points {
		x, _ := bitcoinlib.NewFieldElement(prime, point[0])
		y, _ := bitcoinlib.NewFieldElement(prime, point[1])
		_, err := bitcoinlib.NewPoint(bitcoinlib.NewCoordinates(x, y), *a, *b)
		if err != nil {
			t.Fatalf("Failed at element: (%d, %d) => Expected valid but got invalid point", point[0], point[1])
		}
	}
	for _, point := range invalid_points {
		x, _ := bitcoinlib.NewFieldElement(prime, point[0])
		y, _ := bitcoinlib.NewFieldElement(prime, point[1])
		_, err := bitcoinlib.NewPoint(bitcoinlib.NewCoordinates(x, y), *a, *b)
		if err == nil {
			t.Fatalf("Failed at element: (%d, %d) => Expected invalid but got valid point", point[0], point[1])
		}
	}
}

func TestAddition(t *testing.T) {
	prime := 223
	additions := [][4]int{{170, 142, 60, 139}, {47, 71, 17, 56}, {143, 98, 76, 66}}
	results := [][2]int{{220, 181}, {215, 68}, {47, 71}}
  for index, points := range additions {
		point_a, _ := bitcoinlib.NewPointFromInts(prime, points[0], points[1], 0, 7)
		point_b, _ := bitcoinlib.NewPointFromInts(prime, points[2], points[3], 0, 7)
		expected, _ := bitcoinlib.NewPointFromInts(prime, results[index][0], results[index][1], 0, 7)
		result, err := point_a.Add(point_b)
		if err != nil {
			t.Fatalf("Failed at point number %d because of error on addition: %s", index+1, err)
		}
		if result.Ne(expected) {
			t.Fatalf("Failed at point number %d with result %s instead of %s", index+1, result, expected)
		}
	}
}

func TestScalarMultiplication(t *testing.T) {
	prime := 223
  points := [][2]int{{192, 105}, {143, 98}, {47, 71}, {47, 71}, {47, 71}}
  scalars := []int{2, 2, 2, 4, 8}
  results := [][2]int{{49, 71}, {64, 168}, {36, 111}, {194, 51}, {116, 55}}
  for index, point := range points {
	  point, _ := bitcoinlib.NewPointFromInts(prime, point[0], point[1], 0, 7)
    expected, _ := bitcoinlib.NewPointFromInts(prime, results[index][0], results[index][1], 0, 7) 
    scalar := scalars[index]
    result := point.Scale(scalar)
    if result.Ne(expected) {
      t.Fatalf("Point %d: Expected %s but got %s", index, expected, result)
    }
  } 
}

func TestScalarInfinityMultiplication(t *testing.T) {
  prime := 223

  point, _ := bitcoinlib.NewPointFromInts(prime, 47, 71, 0, 7)
  expected, _ := bitcoinlib.NewInfinitePoint(prime, 0, 7)

  result := point.Scale(21)
  if result.Ne(expected) {
    t.Fatalf("Did not get infinity: Expected %s but got %s", expected, result)
  }

}

func TestPointOnCurve(t *testing.T) {
  prime := bitcoinlib.FromArray([4]uint64{0xffffffffffffffff, 0xffffffffffffffff, 0xffffffffffffffff, 0xfffffffefffffc2f})
  gx := bitcoinlib.FromHexString("0x79be667ef9dcbbac55a06295ce870b07029bfcdb2dce28d959f2815b16f81798").Mod(prime)
  gy := bitcoinlib.FromHexString("0x483ada7726a3c4655da4fbfc0e1108a8fd17b448a68554199c47d08ffb10d4b8").Mod(prime)
  left := gy.Mul(gy).Mod(prime)
  right := gx.Mul(gx).Mod(prime).Mul(gx).Add(bitcoinlib.FromInt(7)).Mod(prime)
  if left.Ne(right) {
    t.Fatalf("Expected equals but got %s != %s", left.String(), right.String())
  } 
}


func TestPointOrder(t *testing.T) {
  prime := bitcoinlib.FromArray([4]uint64{0xffffffffffffffff, 0xffffffffffffffff, 0xffffffffffffffff, 0xfffffffefffffc2f})
  gx := bitcoinlib.FromHexString("0x79be667ef9dcbbac55a06295ce870b07029bfcdb2dce28d959f2815b16f81798").Mod(prime)
  gy := bitcoinlib.FromHexString("0x483ada7726a3c4655da4fbfc0e1108a8fd17b448a68554199c47d08ffb10d4b8").Mod(prime)
  x, _ := bitcoinlib.NewFieldElementFromInt(prime, gx)
  y, _ := bitcoinlib.NewFieldElementFromInt(prime, gy)
  seven, _ := bitcoinlib.NewFieldElementFromInt(prime, bitcoinlib.FromInt(7))
  zero, _ := bitcoinlib.NewFieldElementFromInt(prime, bitcoinlib.FromInt(0))
  infinity, _ := bitcoinlib.NewInfinitePointFromInt(prime, bitcoinlib.FromInt(0), bitcoinlib.FromInt(7))
  g, _ := bitcoinlib.NewPoint(bitcoinlib.NewCoordinates(x, y), *zero, *seven)
  n := bitcoinlib.FromHexString("0xfffffffffffffffffffffffffffffffebaaedce6af48a03bbfd25e8cd0364141")
    
  if g.ScaleInt(n).Ne(infinity) {
    t.Fatal("End result was not as expected")
  }
}

