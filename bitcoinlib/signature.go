package bitcoinlib

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"strings"
)

type SecStart uint8

const (
	COMPRESSED SecStart = 1 + iota
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
	s = s.ExpNeg(ORDER)
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

func Address(point Point, secType SecStart, testnet bool) string {
	secVal := sec(point, secType)
	hashed := Hash160(secVal)
	prefix := []byte{0x00}
	if testnet {
		prefix[0] = 0x6f
	}
	hashed = append(prefix, hashed...)
	checksum := Hash256(hashed)
	hashed = append(hashed, checksum[:4]...)
	return IntoBase58(hex.EncodeToString(hashed))
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
		return nil, errors.New("stream too short")
	}
	first_byte := SecStart(stream[0])

	if first_byte == UNCOMPRESSED {
		return parseUncompressed(stream)
	}
	return parseCompressed(stream)
}

func parseCompressed(stream []byte) (Point, error) {
	if len(stream) != 33 {
		return nil, errors.New("invalid stream length for compressed sec format")
	}
	even := SecStart(stream[0]) == EVEN_Y
	x := FromHexString("0x"+hex.EncodeToString(stream[1:]))
	return solveY(x, even), nil
}

func parseUncompressed(stream []byte) (Point, error) {
	if len(stream) != 65 {
		return nil, errors.New("invalid stream length for uncompressed sec format")
	}
	x_stream := stream[1:33]
	y_stream := stream[33:]
	x := FromHexString("0x"+hex.EncodeToString(x_stream))
	y := FromHexString("0x"+hex.EncodeToString(y_stream))
	return NewS256Point(x, y)
}

func encodeIntToDer(value Int) []byte {
	buf := make([]byte, 32)
	value.value.FillBytes(buf)
	if len(buf) > 0 && buf[0]&0x80 != 0 {
		buf = append([]byte{0}, buf...)
	}
	start := []byte{0x02, byte(len(buf))}
	return append(start, buf...)
}

func (sg *Signature) Der() []byte {
	der := encodeIntToDer(sg.r)
	der = append(der, encodeIntToDer(sg.s.ExpNeg(ORDER))...)
	return append([]byte{0x30, byte(len(der))}, der...)
}

func ParseFromDer(pubkey Point, sign []byte) (*Signature, error){
  if sign[0] != 0x30 {
    return nil, errors.New("invalid der signature")
  }
  sign = sign[1:]
  if len(sign[1:]) != int(sign[0]) {
    return nil, fmt.Errorf("invalid der signature length: %d vs %d", sign[0], len(sign[1:]))
  }
  sign = sign[1:]
  marker := sign[0]
  if marker != 0x02 {
    return nil, errors.New("invalid marker")
  }
  rLength := sign[1]
  sign = sign[2:]
  r := sign[:rLength]
  sign = sign[rLength:]
  marker = sign[0]
  if marker != 0x02 {
    return nil, errors.New("invalid marker")
  }
  sLength := sign[1]
  sign = sign[2:]
  s := sign[:sLength]
  return &Signature{
    s: FromLittleEndian(s),
    r: FromLittleEndian(r),
    p: pubkey,
  }, nil

}

func sec(p Point, secType SecStart) []byte {
	if point, ok := p.(*FinitePoint); ok {
		switch secType {
		case COMPRESSED:
			return compressedSec(point)
		case UNCOMPRESSED:
			return uncompressedSec(point)
		default:
			return compressedSec(point)
		}
	} else {
		buf := make([]byte, 65)
		buf[0] = byte(UNCOMPRESSED)
		return buf
	}
}

func (p *PrivateKey) Sec(secType SecStart) []byte {
	return sec(p.p, secType)
}

func (p *PrivateKey) Address(secType SecStart, testnet bool) string {
	return Address(p.p, secType, testnet)
}

func (p *PrivateKey) WIF(secType SecStart, testnet bool) string {
	wif := []byte{0x80}
	num := p.e.IntoBytes()
	suffix := []byte{}
	if testnet {
		wif[0] = 0xef
	}
	if secType != UNCOMPRESSED {
		suffix = append(suffix, 0x01)
	}
	wif = append(wif, num[:]...)
	wif = append(wif, suffix...)
	
	checksum := Hash256(wif)
	wif = append(wif, checksum[:4]...)
	return IntoBase58(hex.EncodeToString(wif))
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
	for {
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
