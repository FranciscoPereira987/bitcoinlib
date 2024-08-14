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
  invalid_points := [][2]int{{200,119}, {42,99}}
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


