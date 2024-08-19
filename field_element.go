package bitcoinlib

import (
	"errors"
	"fmt"
)

type FieldElement struct {
	order Int 
	value Int
}

func NewFieldElement(order, value int) (*FieldElement, error) {
	if value >= order {
		return nil, errors.New(fmt.Sprintf("value: %d if larger than the order %d", value, order))
	}
	if value < 0 {
		value = value % order
		value = order + value
	}
	return &FieldElement{
		FromInt(order),
		FromInt(value),
	}, nil
}

func NewFieldElementFromInt(order, value Int) (*FieldElement, error) {
  if value.Ge(order) {
		return nil, errors.New(fmt.Sprintf("value: %d if larger than the order %d", value, order))
  }
  return &FieldElement{
    order,
    value.Mod(order),
  }, nil
}

func (e FieldElement) PrintElement() string {
	return fmt.Sprintf("FieldElement_%d(%d)", e.value, e.order)
}

func (e FieldElement) Eq(other FieldElement) bool {
	return e.order.Eq(other.order) && e.value.Eq(other.value)
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
	result := e.value.Add(other.value).Mod(e.order)
	return NewFieldElementFromInt(e.order, result)
}

func (e FieldElement) Sub(other FieldElement) (*FieldElement, error) {
	if e.order != other.order {
		return nil, differentFieldsError()
	}
	result := e.value.Sub(other.value).Mod(e.order)
	return NewFieldElementFromInt(e.order, result)
}

func (e FieldElement) Mul(other FieldElement) (*FieldElement, error) {
	if e.order != other.order {
		return nil, differentFieldsError()
	}
	result := e.value.Mul(other.value).Mod(e.order)
	return NewFieldElementFromInt(e.order, result)
}

func (e FieldElement) Pow(by Int) (*FieldElement, error) {
	result := FromInt(1) 
	//I know that a^(p-1) == 1 (mod r)
	//This line allows me then to exponentiate by negative values
	by = by.Mod(e.order.Sub(FromInt(1)))
	for by.Ge(FromInt(0)) {
		result = result.Mul(e.value).Mod(e.order)
		by = by.Sub(FromInt(1))
	}
	return NewFieldElementFromInt(e.order, result)
}

func (e FieldElement) Inverse() *FieldElement {
	inverse, _ := e.Pow(e.order.Sub(FromInt(2)))
	return inverse
}

func (e FieldElement) Div(by FieldElement) (*FieldElement, error) {
	if e.order != by.order {
		return nil, differentFieldsError()
	}
	byInverted := by.Inverse()
	return e.Mul(*byInverted)
}
