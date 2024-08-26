package bitcoinlib

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"strings"
)


type SecStart uint8

const (
  COMPRESSED SecStart =  1 + iota 
  EVEN_Y 
  ODD_Y 
  UNCOMPRESSED 
)


type Signature struct {
  r Int
  s Int
  p Point
}

func (s *Signature) String() string {
  return s.r.String() + ":" + s.s.String()
}

type PrivateKey struct {
  e Int
  p Point
}

func NewSignature(r, s Int, p Point) *Signature {
  s = s.Exp(ORDER.Sub(TWO), ORDER)
  return &Signature{
    r,
    s,
    p,
  }
}

func NewPrivateKey(e Int) *PrivateKey {
  return &PrivateKey{
    e,
    G().ScaleInt(e),
  }
}

func uncompressedSec(p *FinitePoint) []byte {
  sec := []byte{byte(UNCOMPRESSED)}
  buf := make([]byte, 32)
  p.x.value.value.FillBytes(buf)
  sec = append(sec, buf...)
  p.y.value.value.FillBytes(buf)
  return append(sec, buf...)
}

func compressedSec(p *FinitePoint) []byte {
  oddity := p.y.value.Mod(TWO)
  sec := []byte{byte(EVEN_Y)}
  if oddity.Eq(ONE) {
    sec[0] = byte(ODD_Y)
  }
  buf := make([]byte, 32)
  p.x.value.value.FillBytes(buf)
  return append(sec, buf...)
}

func ParseFromSec(stream []byte) (Point, error) {
  if len(stream) < 33 {
    return nil, errors.New("Stream too short")
  }
  first_byte := SecStart(stream[0])
  
  if first_byte == UNCOMPRESSED{
    return parseUncompressed(stream)
  } 
  return parseCompressed(stream)
}

func parseCompressed(stream []byte) (Point, error) {
  if len(stream) != 33 {
    return nil, errors.New("Invalid stream length for compressed sec format")
  }
  even := SecStart(stream[0]) == EVEN_Y
  x := FromHexString(hex.EncodeToString(stream[1:]))
  return solveY(x, even), nil
}

func parseUncompressed(stream []byte) (Point, error) {
  if len(stream) != 65 {
    return nil, errors.New("Invalid stream length for uncompressed sec format")
  }
  x_stream := stream[1:33]
  y_stream := stream[33:]
  x := FromHexString(hex.EncodeToString(x_stream))
  y := FromHexString(hex.EncodeToString(y_stream))
  return NewS256Point(x, y)
}

func (p *PrivateKey) Sec(secType SecStart) []byte {
  if point, ok := p.p.(*FinitePoint); ok {
    switch secType {
      case COMPRESSED: return compressedSec(point)
      case UNCOMPRESSED: return uncompressedSec(point)
      default: return compressedSec(point) 
    } 
  }else {
    buf := make([]byte, 65)
    buf[0] = byte(UNCOMPRESSED)
    return buf
  } 
}

func GetR(k Int) Int {
  k_g := G().ScaleInt(k)
  p, _ := k_g.(*FinitePoint)
  return p.x.value
}

func (s *Signature) Verify(z Int) bool {
  u := z.Mul(s.s).Mod(ORDER)
  v := s.r.Mul(s.s).Mod(ORDER)
  total, _ := G().ScaleInt(u).Add(s.p.ScaleInt(v))
  result, ok := total.(*FinitePoint)
  return ok && result.x.value.Eq(s.r)
}

func (pk *PrivateKey) generateRandom(z Int) Int {
  k := []byte(strings.Repeat("\x00", 32)) 
  v := []byte(strings.Repeat("\x01", 32))
  if z.Ge(ORDER) {
    z = z.Sub(ORDER)
  }
  z_bytes := z.value.Bytes()
  e_bytes := pk.e.value.Bytes()
  k = hmac.New(sha256.New, k).Sum(append(append([]byte{0}, e_bytes...), z_bytes...))
  v = hmac.New(sha256.New, k).Sum(v)

  k = hmac.New(sha256.New, k).Sum(append(append([]byte{0}, e_bytes...), z_bytes...))
  v = hmac.New(sha256.New, k).Sum(v)
  for  {
    v = hmac.New(sha256.New, k).Sum(v)
    candidate := FromHexString("0x" + hex.EncodeToString(v))
    if candidate.Ge(ONE) && candidate.Le(ORDER) {
      return candidate
    }
    k = hmac.New(sha256.New, k).Sum(append(v, 0))
    v = hmac.New(sha256.New, k).Sum(v)
  }
} 

func (pk *PrivateKey) Sign(z Int) *Signature {
  k := pk.generateRandom(z)
  r := GetR(k)
  k_inv := k.Exp(ORDER.Sub(TWO), ORDER)
  s := z.Add(r.Mul(pk.e)).Mul(k_inv).Mod(ORDER)
  if s.Ge(ORDER.Div(TWO)) {
    s = ORDER.Sub(s)
  }
  return NewSignature(r, s, pk.p)
}
