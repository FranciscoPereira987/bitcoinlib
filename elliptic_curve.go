package bitcoinlib

import (
	"errors"
	"fmt"
)

type Point interface {
  Eq(other Point) bool
  Ne(other Point) bool
  Add(other Point) (Point, error)
  SameCurve(other Point) bool
} 

type Coordinates struct {
  x int
  y int
  inf bool
}

//Elliptic curves of the form y ** 2 = x ** 3 + a x + b
type FinitePoint struct {
  a int
  b int
  x int
  y int
}

type Infinite struct {
  a int
  b int
}

func get_ab(point Point) (a int, b int) {
  if otherInf, ok := point.(Infinite); ok {
    a = otherInf.a
    b = otherInf.b
  }
  if point, ok := point.(FinitePoint); ok {
    a = point.a
    b = point.b
  }
  return 
}

func pointOnTheSameCurve(point_a Point, point_b Point) bool {
  a_a, b_a := get_ab(point_a)
  a_b, b_b := get_ab(point_b)

  return a_a == a_b && b_a == b_b
}

func (inf Infinite) Eq(other Point) bool {
  otherInf, ok := other.(Infinite)
  return ok && otherInf.a == inf.a && otherInf.b == inf.b
}

func (inf Infinite) Ne(other Point) bool {
  return !inf.Eq(other)
}

func (inf Infinite) Add(other Point) (Point, error) {
  return other, nil
}

func (inf Infinite) SameCurve(other Point) bool {
  return pointOnTheSameCurve(inf, other)  
}

func (p FinitePoint) evaluatex() int {
  cubed := p.x * p.x * p.x
  rest := p.a * p.x + p.b
  return cubed + rest 
}

func (p FinitePoint) evaluatey() int {
  return p.y * p.y
}

func NewPoint(coord Coordinates, a, b int) (Point, error) {
  if coord.inf {
    return &Infinite{}, nil
  }

  point := &FinitePoint{
    a,
    b,
    coord.x,
    coord.x,
  }

  if point.evaluatex() != point.evaluatey() {
    return nil, errors.New(fmt.Sprintf("(%d, %d) is not on the curve", coord.x, coord.y))
  }

  return point, nil
}

func (p FinitePoint) Eq(other Point) bool {
  finite, ok := other.(FinitePoint)
  if !ok {
    return false
  } 
  return p.a == finite.a && p.b == finite.b && p.x == finite.x && p.y == finite.y
}

func (p FinitePoint) Ne(other Point) bool {
  return !p.Eq(other)
}

func (p FinitePoint) Add(other Point) (Point, error) {
  if otherFinite, ok := other.(FinitePoint); ok {
    if otherFinite.x == p.x {
      return &Infinite{}, nil
    }

  }
  return other.Add(p)
}

func (p FinitePoint) SameCurve(other Point) bool {
  return pointOnTheSameCurve(p, other)
}
