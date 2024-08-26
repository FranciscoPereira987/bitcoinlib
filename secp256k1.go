package bitcoinlib



var PRIME Int = FromHexString("0xfffffffffffffffffffffffffffffffffffffffffffffffffffffffefffffc2f")
var ORDER Int = FromHexString("0xfffffffffffffffffffffffffffffffebaaedce6af48a03bbfd25e8cd0364141")

var SQRT_EXP Int = PRIME.Add(ONE).Div(FOUR)

func A() *FieldElement {
  a, _ := NewS256Field(FromInt(0))
  return a
}

func B() *FieldElement {
  b, _ := NewS256Field(FromInt(7))
  return b
}

func G() Point {
  gx := FromHexString("0x79be667ef9dcbbac55a06295ce870b07029bfcdb2dce28d959f2815b16f81798")
  gy := FromHexString("0x483ada7726a3c4655da4fbfc0e1108a8fd17b448a68554199c47d08ffb10d4b8")
  val, _ := NewS256Point(gx, gy)
  return val
}

func NewS256Field(value Int) (*FieldElement, error) {
  return NewFieldElementFromInt(PRIME, value)
}


func NewS256Point(x Int, y Int) (Point, error) {
  x_field, err := NewS256Field(x)
  if err != nil {
    return nil, err
  }
  y_field, err := NewS256Field(y)
  if err != nil {
    return nil, err
  }
  coords := NewCoordinates(x_field, y_field)

  return NewPoint(coords, *A(), *B()) 
}

func NewS256Infinite() Point {
  val, _ := NewInfinitePointFromInt(PRIME, FromInt(0), FromInt(7))
  return val
}

func solveY(x Int, even bool) Point {
  alpha := x.Exp(THREE, PRIME).Add(B().value)
  beta := alpha.Exp(SQRT_EXP, PRIME)
  even_beta := beta
  odd_beta := PRIME.Sub(beta).Mod(PRIME)
  if beta.Mod(TWO).Eq(ONE) {
    even_beta, odd_beta = odd_beta, even_beta 
  }
  if even {
    result, _ :=  NewS256Point(x, even_beta)
    return result
  }
  result, _ := NewS256Point(x, odd_beta)
  return result
}

//Wraps Scalar of Point
func S256Mul(a Point, by Int) Point {
  return a.ScaleInt(by)
}

func S256Verifyr(a Point, r Int) bool {
  if a_trans, ok := a.(*FinitePoint); ok {
    return a_trans.x.value.Eq(r)
  }
  return false
}


