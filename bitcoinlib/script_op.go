package bitcoinlib

import "slices"


type Operation interface {
  Operate(z string, stack *[]Operation, altstack *[]Operation, cmds *[]Operation) bool
  IsValid() bool //Returns true if its not the empty byte string
}

func encodeNum(num Int) []byte {
  if num.Eq(ZERO) {
    return []byte{}
  }
  negative := false
  if num.Le(ZERO) {
    num = num.Mul(FromInt(-1))
    negative = true
  }
  
  result := make([]byte, len(num.value.Bytes()))
  copy(result, num.value.Bytes())
  if negative && result[0] & 0x80 == 0x80{
    result = append(result, 0x80)
  }else if !negative && result[0] & 0x80 == 0x80{
    result = append(result, 0x00)
  }else if negative {
    result[0] |= 0x80
  }
  slices.Reverse(result)
  return result 
}

func decodeNum(element []byte) Int {
  if len(element) == 0 {
    return ZERO
  }
  negative := false
  if element[len(element)-1] & 0x80 == 0x80 {
    negative = true
    element[len(element)-1] = element[len(element)-1] & 0x7f 
  }
  result := make([]byte, len(element))
  copy(result, element)
  slices.Reverse(result)
  value := FromInt(0)
  value.value.FillBytes(result)
  if negative {
    value.Mul(FromInt(-1))
  }
  return value
}

type OP_0 struct{}

func (t *OP_0) Operate(z string, stack *[]Operation, altstack *[]Operation, cmds *[]Operation) bool {
  *stack = append(*stack, t)
  return true
}

func (t *OP_0) IsValid() bool {
  return false
}

type OP_1Negate struct {}

func (t *OP_1Negate) Operate(z string, stack *[]Operation, altstack *[]Operation, cmds *[]Operation) bool {
  *stack = append(*stack, t)
  return true
}

//TODO: Check if this is Valid (Or if just 1 should be valid)
func (t *OP_1Negate) IsValid() bool {
  return false
}

type OP_1 struct {}

func (t *OP_1) Operate(z string, stack *[]Operation, altstack *[]Operation, cmds *[]Operation) bool {
  *stack = append(*stack, t)
  return true
}

func (t *OP_1) IsValid() bool {
  return true 
}

type OP_2 struct {}

func (t *OP_2) Operate(z string, stack *[]Operation, altstack *[]Operation, cmds *[]Operation) bool {
  *stack = append(*stack, t)
  return true
}

func (t *OP_2) IsValid() bool {
  return false 
}

type OP_3 struct {}

func (t *OP_3) Operate(z string, stack *[]Operation, altstack *[]Operation, cmds *[]Operation) bool {
  *stack = append(*stack, t)
  return true
}

func (t *OP_3) IsValid() bool {
  return false 
}

type OP_4 struct {}

func (t *OP_4) Operate(z string, stack *[]Operation, altstack *[]Operation, cmds *[]Operation) bool {
  *stack = append(*stack, t)
  return true
}

func (t *OP_4) IsValid() bool {
  return false 
}

type OP_5 struct {}

func (t *OP_5) Operate(z string, stack *[]Operation, altstack *[]Operation, cmds *[]Operation) bool {
  *stack = append(*stack, t)
  return true
}

func (t *OP_5) IsValid() bool {
  return false 
}

type OP_6 struct {}

func (t *OP_6) Operate(z string, stack *[]Operation, altstack *[]Operation, cmds *[]Operation) bool {
  *stack = append(*stack, t)
  return true
}

func (t *OP_6) IsValid() bool {
  return false 
}

type OP_7 struct {}

func (t *OP_7) Operate(z string, stack *[]Operation, altstack *[]Operation, cmds *[]Operation) bool {
  *stack = append(*stack, t)
  return true
}

func (t *OP_7) IsValid() bool {
  return false 
}

type OP_8 struct {}

func (t *OP_8) Operate(z string, stack *[]Operation, altstack *[]Operation, cmds *[]Operation) bool {
  *stack = append(*stack, t)
  return true
}

func (t *OP_8) IsValid() bool {
  return false 
}

type OP_9 struct {}

func (t *OP_9) Operate(z string, stack *[]Operation, altstack *[]Operation, cmds *[]Operation) bool {
  *stack = append(*stack, t)
  return true
}

func (t *OP_9) IsValid() bool {
  return false 
}

type OP_10 struct {}

func (t *OP_10) Operate(z string, stack *[]Operation, altstack *[]Operation, cmds *[]Operation) bool {
  *stack = append(*stack, t)
  return true
}

func (t *OP_10) IsValid() bool {
  return false 
}

type OP_11 struct {}

func (t *OP_11) Operate(z string, stack *[]Operation, altstack *[]Operation, cmds *[]Operation) bool {
  *stack = append(*stack, t)
  return true
}

func (t *OP_11) IsValid() bool {
  return false 
}

type OP_12 struct {}

func (t *OP_12) Operate(z string, stack *[]Operation, altstack *[]Operation, cmds *[]Operation) bool {
  *stack = append(*stack, t)
  return true
}

func (t *OP_12) IsValid() bool {
  return false 
}

type OP_13 struct {}

func (t *OP_13) Operate(z string, stack *[]Operation, altstack *[]Operation, cmds *[]Operation) bool {
  *stack = append(*stack, t)
  return true
}

func (t *OP_13) IsValid() bool {
  return false 
}

type OP_14 struct {}

func (t *OP_14) Operate(z string, stack *[]Operation, altstack *[]Operation, cmds *[]Operation) bool {
  *stack = append(*stack, t)
  return true
}

func (t *OP_14) IsValid() bool {
  return false 
}

type OP_15 struct {}

func (t *OP_15) Operate(z string, stack *[]Operation, altstack *[]Operation, cmds *[]Operation) bool {
  *stack = append(*stack, t)
  return true
}

func (t *OP_15) IsValid() bool {
  return false 
}

type OP_16 struct {}

func (t *OP_16) Operate(z string, stack *[]Operation, altstack *[]Operation, cmds *[]Operation) bool {
  *stack = append(*stack, t)
  return true
}

func (t *OP_16) IsValid() bool {
  return false 
}

type OP_NOP struct {}

func (t *OP_NOP) Operate(z string, stack *[]Operation, altstack *[]Operation, cmds *[]Operation) bool {
  return true
}

func (t *OP_NOP) IsValid() bool {
  return false 
}
















