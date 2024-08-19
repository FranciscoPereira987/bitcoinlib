package bitcoinlib

import (
	"errors"
	"fmt"
)

type Point interface {
	Eq(other Point) bool
	Ne(other Point) bool
	Add(other Point) (Point, error)
	Scale(by int) Point
	SameCurve(other Point) bool
}

type Coordinates struct {
	x   *FieldElement
	y   *FieldElement
	inf bool
}

func NewCoordinates(x, y *FieldElement) Coordinates {
	return Coordinates{
		x,
		y,
		false,
	}
}

func InfiniteCoord() Coordinates {
	return Coordinates{
		nil,
		nil,
		true,
	}
}

// Elliptic curves of the form y ** 2 = x ** 3 + a x + b
type FinitePoint struct {
	a FieldElement
	b FieldElement
	x FieldElement
	y FieldElement
}

type Infinite struct {
	a FieldElement
	b FieldElement
}

func get_ab(point Point) (a FieldElement, b FieldElement) {
	if otherInf, ok := point.(*Infinite); ok {
		a = otherInf.a
		b = otherInf.b
	}
	if point, ok := point.(*FinitePoint); ok {
		a = point.a
		b = point.b
	}
	return
}

func pointOnTheSameCurve(point_a Point, point_b Point) bool {
	a_a, b_a := get_ab(point_a)
	a_b, b_b := get_ab(point_b)

	return a_a.Eq(a_b) && b_a.Eq(b_b)
}

func (inf *Infinite) Eq(other Point) bool {
	otherInf, ok := other.(*Infinite)
	return ok && otherInf.a.Eq(inf.a) && otherInf.b.Eq(inf.b)
}

func (inf *Infinite) Ne(other Point) bool {
	return !inf.Eq(other)
}

func (inf *Infinite) Add(other Point) (Point, error) {
	return other, nil
}

func (inf *Infinite) Scale(by int) Point {
	return inf
}

func (inf *Infinite) SameCurve(other Point) bool {
	return pointOnTheSameCurve(inf, other)
}

func (p FinitePoint) evaluatex() *FieldElement {
	cubed, _ := p.x.Pow(FromInt(3))
	rest, _ := p.a.Mul(p.x)
	rest, _ = rest.Sum(p.b)
	result, _ := rest.Sum(*cubed)
	return result
}

func (p FinitePoint) evaluatey() *FieldElement {
	result, _ := p.y.Pow(FromInt(2))
	return result
}

func NewPointFromInts(order, x, y, a, b int) (Point, error) {
	a_val, err := NewFieldElement(order, a)
	if err != nil {
		return nil, err
	}
	b_val, err := NewFieldElement(order, b)
	if err != nil {
		return nil, err
	}
	x_val, err := NewFieldElement(order, x)
	if err != nil {
		return nil, err
	}
	y_val, err := NewFieldElement(order, y)
	if err != nil {
		return nil, err
	}
	return NewPoint(NewCoordinates(x_val, y_val), *a_val, *b_val)
}

func NewInfinitePoint(order, a, b int) (Point, error) {
  a_val, err := NewFieldElement(order, a)
  if err != nil {
    return nil, err
  }
  b_val, err := NewFieldElement(order, b)
  if err != nil {
    return nil, err
  }
  return &Infinite{
    a: *a_val,
    b: *b_val,
  }, nil
}

func NewInfinitePointFromInt(order, a, b Int) (Point, error) {
  a_val, err := NewFieldElementFromInt(order, a)
  if err != nil {
    return nil, err
  }
  b_val, err := NewFieldElementFromInt(order, b)
  if err != nil {
    return nil, err
  }
  return &Infinite{
    a: *a_val,
    b: *b_val,
  }, nil

}

func NewPoint(coord Coordinates, a, b FieldElement) (Point, error) {
	if coord.inf {
		return &Infinite{a, b}, nil
	}

	point := &FinitePoint{
		a,
		b,
		*coord.x,
		*coord.y,
	}

	if point.evaluatex().Ne(*point.evaluatey()) {
		return nil, errors.New(fmt.Sprintf("(%d, %d) is not on the curve", coord.x, coord.y))
	}

	return point, nil
}

func (p *FinitePoint) Eq(other Point) bool {
	finite, ok := other.(*FinitePoint)
	if !ok {

		return false
	}
	return p.a.Eq(finite.a) && p.b.Eq(finite.b) && p.x.Eq(finite.x) && p.y.Eq(finite.y)
}

func (p *FinitePoint) Ne(other Point) bool {
	return !p.Eq(other)
}

// Returns a new point when p is added to itself
func (p *FinitePoint) addOnItself() (Point, error) {
	if p.y.value.Eq(FromInt(0)) {
		return &Infinite{
      p.a,
      p.b,
    }, nil
	}
	three, _ := NewFieldElementFromInt(p.x.order, FromInt(3))
	two, _ := NewFieldElementFromInt(p.x.order, FromInt(2))
	squared, _ := p.x.Pow(FromInt(2))
	multiplied, _ := three.Mul(*squared)
	summed, _ := multiplied.Sum(p.a)
	multiplied_y, _ := two.Mul(p.y)
	s, _ := summed.Div(*multiplied_y)
	squared, _ = s.Pow(FromInt(2))
	multiplied, _ = two.Mul(p.x)
	new_x, _ := squared.Sub(*multiplied)
	sub, _ := p.x.Sub(*new_x)
	sub, _ = s.Mul(*sub)
	new_y, _ := sub.Sub(p.y)
	coords := Coordinates{
		x: new_x,
		y: new_y,
	}
	return NewPoint(coords, p.a, p.b)
}

// Returns a new Point, based on the formulae for when both
// points are Finite and different from each other
func (p *FinitePoint) addOnDifferent(other FinitePoint) (Point, error) {
	dividend, _ := other.y.Sub(p.y)
	divisor, _ := other.x.Sub(p.x)
	s, _ := dividend.Div(*divisor)
	new_x, _ := s.Pow(FromInt(2))
	new_x, _ = new_x.Sub(other.x)
	new_x, _ = new_x.Sub(p.x)
	new_y, _ := p.x.Sub(*new_x)
	new_y, _ = s.Mul(*new_y)
	new_y, _ = new_y.Sub(p.y)
	coords := Coordinates{
		x: new_x,
		y: new_y,
	}
	return NewPoint(coords, p.a, p.b)
}

func (p *FinitePoint) Add(other Point) (Point, error) {
	if otherFinite, ok := other.(*FinitePoint); ok {
		if p.Eq(other) {
			return p.addOnItself()
		}
		if otherFinite.x.Eq(p.x) {
			return &Infinite{
        p.a,
        p.b,
      }, nil
		}
		return p.addOnDifferent(*otherFinite)
	}
	return other.Add(p)
}

//Multiplies using binary expansion
func (p *FinitePoint) Scale(by int) Point {

  partial, _ := NewInfinitePointFromInt(p.a.order, p.a.value, p.b.value) 
  actual, _ := NewPoint(NewCoordinates(&p.x, &p.y), p.a, p.b)
  for by > 0 {
    //If the current value is one (add the actual number)
    if by & 1 == 1{
      partial, _ = actual.Add(partial)
    }
    //Multiplies by two before the shift
	  actual, _ = actual.Add(actual)	
    by = by >> 1
	}
	return partial
}

func (p *FinitePoint) SameCurve(other Point) bool {
	return pointOnTheSameCurve(p, other)
}
