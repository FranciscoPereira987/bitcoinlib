package bitcoinlib

import "slices"


type Operation interface {
  Operate(z string, stack *Stack, altstack *Stack, cmds *Stack) bool
  Num() int
}

type Stack = []Operation

func Push(st *Stack, el Operation) {
  *st = append(*st, el)
}

func Pop(st *Stack)  Operation{
  if Len(st) == 0 {
    return nil
  }
  e := (*st)[Len(st)-1]
  *st = (*st)[:Len(st)-1]
  return e
}

func Len(st *Stack) int {
  return len(*st)
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

func (t *OP_0) Operate(z string, stack *Stack, altstack *Stack, cmds *Stack) bool {
  *stack = append(*stack, t)
  return true
}


func (t *OP_0) Num() int {
  return 0
}

type OP_1Negate struct {}

func (t *OP_1Negate) Operate(z string, stack *Stack, altstack *Stack, cmds *Stack) bool {
  *stack = append(*stack, t)
  return true
}

func (t *OP_1Negate) Num() int {
  return 79
}


type OP_1 struct {}

func (t *OP_1) Operate(z string, stack *Stack, altstack *Stack, cmds *Stack) bool {
  *stack = append(*stack, t)
  return true
}


func (t *OP_1) Num() int {
  return 81 
}


type OP_2 struct {}

func (t *OP_2) Operate(z string, stack *Stack, altstack *Stack, cmds *Stack) bool {
  *stack = append(*stack, t)
  return true
}

func (t *OP_2) Num() int {
  return 82 
}


type OP_3 struct {}

func (t *OP_3) Operate(z string, stack *Stack, altstack *Stack, cmds *Stack) bool {
  *stack = append(*stack, t)
  return true
}

func (t *OP_3) Num() int {
  return 83 
}

type OP_4 struct {}

func (t *OP_4) Operate(z string, stack *Stack, altstack *Stack, cmds *Stack) bool {
  *stack = append(*stack, t)
  return true
}

func (t *OP_4) Num() int {
  return 84 
}


type OP_5 struct {}

func (t *OP_5) Operate(z string, stack *Stack, altstack *Stack, cmds *Stack) bool {
  *stack = append(*stack, t)
  return true
}

func (t *OP_5) Num() int {
  return 85 
}

type OP_6 struct {}

func (t *OP_6) Operate(z string, stack *Stack, altstack *Stack, cmds *Stack) bool {
  *stack = append(*stack, t)
  return true
}


func (t *OP_6) Num() int {
  return 86 
}

type OP_7 struct {}

func (t *OP_7) Operate(z string, stack *Stack, altstack *Stack, cmds *Stack) bool {
  *stack = append(*stack, t)
  return true
}


func (t *OP_7) Num() int {
  return 87 
}

type OP_8 struct {}

func (t *OP_8) Operate(z string, stack *Stack, altstack *Stack, cmds *Stack) bool {
  *stack = append(*stack, t)
  return true
}


func (t *OP_8) Num() int {
  return 88 
}

type OP_9 struct {}

func (t *OP_9) Operate(z string, stack *Stack, altstack *Stack, cmds *Stack) bool {
  *stack = append(*stack, t)
  return true
}


func (t *OP_9) Num() int {
  return 89 
}


type OP_10 struct {}

func (t *OP_10) Operate(z string, stack *Stack, altstack *Stack, cmds *Stack) bool {
  *stack = append(*stack, t)
  return true
}


func (t *OP_10) Num() int {
  return 90 
}


type OP_11 struct {}

func (t *OP_11) Operate(z string, stack *Stack, altstack *Stack, cmds *Stack) bool {
  *stack = append(*stack, t)
  return true
}

func (t *OP_11) Num() int {
  return 91 
}


type OP_12 struct {}

func (t *OP_12) Operate(z string, stack *Stack, altstack *Stack, cmds *Stack) bool {
  *stack = append(*stack, t)
  return true
}

func (t *OP_12) Num() int {
  return 92 
}

type OP_13 struct {}

func (t *OP_13) Operate(z string, stack *Stack, altstack *Stack, cmds *Stack) bool {
  *stack = append(*stack, t)
  return true
}

func (t *OP_13) Num() int {
  return 93 
}

type OP_14 struct {}

func (t *OP_14) Operate(z string, stack *Stack, altstack *Stack, cmds *Stack) bool {
  *stack = append(*stack, t)
  return true
}

func (t *OP_14) Num() int {
  return 94 
}


type OP_15 struct {}

func (t *OP_15) Operate(z string, stack *Stack, altstack *Stack, cmds *Stack) bool {
  *stack = append(*stack, t)
  return true
}

func (t *OP_15) Num() int {
  return 95 
}


type OP_16 struct {}

func (t *OP_16) Operate(z string, stack *Stack, altstack *Stack, cmds *Stack) bool {
  *stack = append(*stack, t)
  return true
}

func (t *OP_16) Num() int {
  return 96 
}


type OP_NOP struct {}

func (t *OP_NOP) Operate(z string, stack *Stack, altstack *Stack, cmds *Stack) bool {
  return true
}

func (t *OP_NOP) Num() int {
  return 97 
}


type OP_IF struct {}

//This function manipulatesc cmds to eliminate or "Prune" the branched values 
//that should not be executed based on the condition in the stack.
func (t *OP_IF) Operate(z string, stack *Stack, altstack *Stack, cmds *Stack) bool {
  if len(*stack) < 1 {
    return false
  }
  trueItems := make(Stack, 0)
  falseItems := make(Stack, 0)
  found := false
  numEndifsNeeded := 1
  currentArray := &trueItems
  for Len(cmds) > 0 {
    item := Pop(cmds) 
    if item.Num() == 99 || item.Num() == 100 {
      // Nested if
      numEndifsNeeded++
      Push(currentArray, item)
    }else if numEndifsNeeded == 1 && item.Num() == 103 {
      currentArray = &falseItems
    }else if item.Num() == 104 {
      found = numEndifsNeeded == 1
      if found {
        break
      }else {
        numEndifsNeeded--
        Push(currentArray, item)
      }
    }else {
      Push(currentArray, item)
    }
  }
  if !found {
    return false
  }
  element := Pop(stack)
  if element.Num() == 0 {
    *cmds = append(*cmds, falseItems...)
  }else {
    *cmds = append(*cmds, trueItems...)
  }
  return true
}

func (t *OP_IF) Num() int {
  return 99
}

type OP_NOTIF struct{}

//Same as OP_IF, but switches the branches that are reinserted into cmds 
func (t *OP_NOTIF) Operate(z string, stack *Stack, altstack *Stack, cmds *Stack) bool {
  if len(*stack) < 1 {
    return false
  }
  trueItems := make(Stack, 0)
  falseItems := make(Stack, 0)
  found := false
  numEndifsNeeded := 1
  currentArray := &trueItems
  for Len(cmds) > 0 {
    item := Pop(cmds) 
    if item.Num() == 99 || item.Num() == 100 {
      // Nested if
      numEndifsNeeded++
      Push(currentArray, item)
    }else if numEndifsNeeded == 1 && item.Num() == 103 {
      currentArray = &falseItems
    }else if item.Num() == 104 {
      found = numEndifsNeeded == 1
      if found {
        break
      }else {
        numEndifsNeeded--
        Push(currentArray, item)
      }
    }else {
      Push(currentArray, item)
    }
  }
  if !found {
    return false
  }
  element := Pop(stack)
  if element.Num() == 0 {
    *cmds = append(*cmds, trueItems...)
  }else {
    *cmds = append(*cmds, falseItems...)
  }
  return true
}

func (t *OP_NOTIF) Num() int {
  return 100
}

type OP_VERIFY struct {}

func (t *OP_VERIFY) Operate(z string, stack *Stack, altstack *Stack, cmds *Stack) bool {
  if Len(stack) < 1 {
    return false
  }
  element := Pop(stack)
  return element.Num() != 0 
}

func (t *OP_VERIFY) Num() int {
  return 105
}

type OP_RETURN struct {}


func (t *OP_RETURN) Operate(z string, stack *Stack, altstack *Stack, cmds *Stack) bool {
  return false
}

func (t *OP_RETURN) Num() int {
  return 106
}

type OP_TOALTSTACK struct {}


func (t *OP_TOALTSTACK) Operate(z string, stack *Stack, altstack *Stack, cmds *Stack) bool {
  if Len(altstack) < 1 {
    return false
  }
  Push(altstack, Pop(stack))
  return true
}

func (t *OP_TOALTSTACK) Num() int {
  return 107
}

type OP_FROMALTSTACK struct {}


func (t *OP_FROMALTSTACK) Operate(z string, stack *Stack, altstack *Stack, cmds *Stack) bool {
  if Len(altstack) < 1 {
    return false
  }
  Push(stack, Pop(altstack))
  return true
}

func (t *OP_FROMALTSTACK) Num() int {
  return 108
}

type OP_2DROP struct {}


func (t *OP_2DROP) Operate(z string, stack *Stack, altstack *Stack, cmds *Stack) bool {
  if Len(stack) < 2 {
    return false
  }
  Pop(stack)
  Pop(stack)
  return true
}

func (t *OP_2DROP) Num() int {
  return 109
}

type OP_2DUP struct {}


func (t *OP_2DUP) Operate(z string, stack *Stack, altstack *Stack, cmds *Stack) bool {
  if Len(stack) < 2 {
    return false
  }
  Push(stack, (*stack)[Len(stack)-2])
  Push(stack, (*stack)[Len(stack)-2])
  return true 
}

func (t *OP_2DUP) Num() int {
  return 110
}

type OP_3DUP struct {}

func (t *OP_3DUP) Operate(z string, stack *Stack, altstack *Stack, cmds *Stack) bool {
  if Len(stack) < 3 {
    return false
  }
  Push(stack, (*stack)[Len(stack)-3])
  Push(stack, (*stack)[Len(stack)-3])
  Push(stack, (*stack)[Len(stack)-3])
  return true 
}

func (t *OP_3DUP) Num() int {
  return 111
}
