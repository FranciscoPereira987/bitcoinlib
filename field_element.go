package bitcoinlib

import (
	"errors"
	"fmt"
)

type FieldElement struct{
  order int
  value int
}

func NewFieldElement(order, value int)  (*FieldElement, error) {
  if value >= order {
    return nil, errors.New(fmt.Sprintf("value: %d if larger than the order %d", value, order))
  }

  return &FieldElement{
    order,
    value,
  }, nil
}

func (e FieldElement) PrintElement() string {
  return fmt.Sprintf("FieldElement_%d(%d)", e.value, e.order)
}

func (e FieldElement) Eq(other FieldElement) bool {
  return e.order == other.order && e.value == other.value
}

func (e FieldElement) Ne(other FieldElement) bool {
  return !e.Eq(other)
}

func differentFieldsError() error {
  return errors.New("Cannot sum elements from different fields")
}

func (e FieldElement) Sum(other FieldElement) (*FieldElement, error) {
  if e.order != other.order {
    return nil, differentFieldsError() 
  }
  result := (e.value + other.value) % e.order
  return NewFieldElement(e.order, result)
}

func (e FieldElement) Sub(other FieldElement) (*FieldElement, error) {
  if e.order != other.order {
    return nil, differentFieldsError() 
  }
  result := (e.value - other.value) % e.order
  return NewFieldElement(e.order, result)
}

func (e FieldElement) Mul(other FieldElement) (*FieldElement, error) {
  if e.order != other.order {
    return nil, differentFieldsError()
  }
  result := (e.value * other.value) % e.order
  return NewFieldElement(e.order, result)
}

func (e FieldElement) Pow(by int) (*FieldElement, error) {
  result := 1
  for by > 0 {
    result = (result * e.value) % e.order
  }
  return NewFieldElement(e.order, result)
}

func (e FieldElement) Inverse() (*FieldElement) {
  inverse, _ := e.Pow(e.order - 2)
  return inverse
}

func (e FieldElement) Div(by FieldElement) (*FieldElement, error) {
  if e.order != by.order {
    return nil, differentFieldsError()
  }
  byInverted := by.Inverse()
  return e.Mul(*byInverted)
}
