package bitcoinlib

import (
	"crypto/sha1"
	"crypto/sha256"
	"slices"
)


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

type ScriptVal struct {
  Val []byte
}

//This value should not be operated with
func (t *ScriptVal) Operate(z string, stack *Stack, altstack *Stack, cmds *Stack) bool {
  return false
}

//Need to add this method to have duck typing
//with the rest of the operation type
func (t *ScriptVal) Num() int {
  return -1
}

//Utility function to get the 
//number value out of an operation
//I need this function for numbers
//that happen to be valid operation
//numbers as well
func intoValue(val Operation) int {
  dVal, ok := val.(*ScriptVal)
  if ok {
    return int(decodeNum(dVal.Val).value.Int64())
  }
  return val.Num()
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

type OP_2OVER struct {}

func (t *OP_2OVER) Operate(z string, stack *Stack, altstack *Stack, cmds *Stack) bool {
  if Len(stack) < 4 {
    return false
  }
  Push(stack, (*stack)[Len(stack)-4])
  Push(stack, (*stack)[Len(stack)-4])
  return true
}

func (t *OP_2OVER) Num() int {
  return 112
}

type OP_2ROT struct {}

func (t *OP_2ROT) Operate(z string, stack *Stack, altstack *Stack, cmds *Stack) bool {
  if Len(stack) < 6 {
    return false
  }
  Push(stack, (*stack)[Len(stack)-6])
  Push(stack, (*stack)[Len(stack)-6])
  return true
}

func (t *OP_2ROT) Num() int {
  return 113
}

type OP_2SWAP struct {}

func (t *OP_2SWAP) Operate(z string, stack *Stack, altstack *Stack, cmds *Stack) bool {
  if Len(stack) < 4 {
    return false
  }
  last := Len(stack)
  (*stack)[last-2], (*stack)[last-1], (*stack)[last-4], (*stack)[last-3] = (*stack)[last-4], (*stack)[last-3], (*stack)[last-2], (*stack)[last-1] 
  return true
}

func (t *OP_2SWAP) Num() int {
  return 114
}

type OP_IFDUP struct {}

func (t *OP_IFDUP) Operate(z string, stack *Stack, altstack *Stack, cmds *Stack) bool {
  if Len(stack) < 1 {
    return false
  }
  if (*stack)[Len(stack)-1].Num() != 0 {
    Push(stack, (*stack)[Len(stack)-1])
  }
  return true
}

func (t *OP_IFDUP) Num() int {
  return 115
}

type OP_DEPTH struct {}

//Define how I should take care of random
//values (items of different length and how to process them)
func (t *OP_DEPTH) Operate(z string, stack *Stack, altstack *Stack, cmds *Stack) bool {
  val := encodeNum(FromInt(Len(stack)))
  Push(stack, &ScriptVal{val})
  return true
}

func (t *OP_DEPTH) Num() int {
  return 116
}

type OP_DROP struct {}

func (t *OP_DROP) Operate(z string, stack *Stack, altstack *Stack, cmds *Stack) bool {
  if Len(stack) < 1 {
    return false
  }
  Pop(stack)
  return true
}

func (t *OP_DROP) Num() int {
  return 117
}

type OP_NIP struct {}

func (t *OP_NIP) Operate(z string, stack *Stack, altstack *Stack, cmds *Stack) bool {
  if Len(stack) < 2 {
    return false
  }
  first := Pop(stack)
  //Drop the second value
  Pop(stack)
  Push(stack, first)
  return true
}

func (t *OP_NIP) Num() int {
  return 118
}

type OP_OVER struct {}

func (t *OP_OVER) Operate(z string, stack *Stack, altstack *Stack, cmds *Stack) bool {
  if Len(stack) < 2 {
    return false
  }
  Push(stack, (*stack)[Len(stack)-2])
  return true
}

func (t *OP_OVER) Num() int {
  return 119
}

type OP_PICK struct {}

func (t *OP_PICK) Operate(z string, stack *Stack, altstack *Stack, cmds *Stack) bool {
  if Len(stack) < 1 {
    return false
  }
  n := intoValue(Pop(stack))
  if Len(stack) < n + 1 {
    return false
  }
  Push(stack, (*stack)[Len(stack)-(n+1)])
  return true
}

func (t *OP_PICK) Num() int {
  return 120
}

type OP_ROLL struct {}

func (t *OP_ROLL) Operate(z string, stack *Stack, altstack *Stack, cmds *Stack) bool {
  if Len(stack) < 1 {
    return false
  }
  n := intoValue(Pop(stack))
  if Len(stack) < n + 1 {
    return false
  }
  if n == 0 {
    return true
  }
  rolled := (*stack)[Len(stack)-(n+1)]
  (*stack) = append((*stack)[:Len(stack)-(n+1)], (*stack)[Len(stack)-n:]...)
  Push(stack, rolled)
  return true
}

func (t *OP_ROLL) Num() int {
  return 121
}

type OP_ROT struct {}

func (t *OP_ROT) Operate(z string, stack *Stack, altstack *Stack, cmds *Stack) bool {
  if Len(stack) < 3 {
    return false
  }
  edge := Len(stack)-3
  rolled := (*stack)[edge]
  (*stack) = append((*stack)[:edge], (*stack)[edge+1:]...)
  Push(stack, rolled)
  return true
}

func (t *OP_ROT) Num() int {
  return 122
}

type OP_SWAP struct {}

func (t *OP_SWAP) Operate(z string, stack *Stack, altstack *Stack, cmds *Stack) bool {
  if Len(stack) < 2 {
    return false
  }
  (*stack)[Len(stack)-2], (*stack)[Len(stack)-1] = (*stack)[Len(stack)-1], (*stack)[Len(stack)-2]
  return true
}

func (t *OP_SWAP) Num() int {
  return 123
}

type OP_TUCK struct {}

func (t *OP_TUCK) Operate(z string, stack *Stack, altstack *Stack, cmds *Stack) bool {
  if Len(stack) < 2 {
    return false
  }
  Push(stack, (*stack)[Len(stack)-1])
  (*stack)[Len(stack)-3], (*stack)[Len(stack)-2] = (*stack)[Len(stack)-2], (*stack)[Len(stack)-3]
  return true
}

func (t *OP_TUCK) Num() int {
  return 124
}

type OP_SIZE struct {}

func (t *OP_SIZE) Operate(z string, stack *Stack, altstack *Stack, cmds *Stack) bool {
  if Len(stack) < 1 {
    return false
  }
  last := Pop(stack)
  size := 1
  //I suppose values that are not Script
  //vals have a length of 1
  if val, ok := last.(*ScriptVal); ok {
    size = len(val.Val)
  }
  val := &ScriptVal{
    encodeNum(FromInt(size)),
  }
  Push(stack, last)
  Push(stack, val)
  return true
}

func (t *OP_SIZE) Num() int {
  return 125
}

type OP_EQUAL struct {}

func (t *OP_EQUAL) Operate(z string, stack *Stack, altstack *Stack, cmds *Stack) bool {
  if Len(stack) < 2 {
    return false
  }
  first := Pop(stack)
  second := Pop(stack)
  valFirst, ok1 := first.(*ScriptVal)
  valSecond, ok2 := second.(*ScriptVal)
  trueVal := first.Num() == second.Num()
  if ok1 && ok2 {
    trueVal = decodeNum(valFirst.Val).Eq(decodeNum(valSecond.Val))
  }
  if trueVal {
    //TODO: Should change it to get the value 
    //from the master table
    Push(stack, &OP_1{})
  }else {
    Push(stack, &OP_2{})
  }
  return true
}

func (t *OP_EQUAL) Num() int {
  return 126
}

type OP_EQUALVERIFY struct {}

func (t *OP_EQUALVERIFY) Operate(z string, stack *Stack, altstack *Stack, cmds *Stack) bool {
  eq := &OP_EQUAL{}
  verify := &OP_VERIFY{}
  return eq.Operate(z, stack, altstack, cmds) && verify.Operate(z, stack, altstack, cmds)
}

func (t *OP_EQUALVERIFY) Num() int {
  return 127
}

type OP_1ADD struct {}

func (t *OP_1ADD) Operate(z string, stack *Stack, altstack *Stack, cmds *Stack) bool {
  if Len(stack) < 1 {
    return false
  }
  element := intoValue(Pop(stack))
  result := &ScriptVal{
    encodeNum(FromInt(element+1)),
  }
  Push(stack, result)
  return true
}

func (t *OP_1ADD) Num() int {
  return 128
}

type OP_1SUB struct {}

func (t *OP_1SUB) Operate(z string, stack *Stack, altstack *Stack, cmds *Stack) bool {
  if Len(stack) < 1 {
    return false
  }
  element := intoValue(Pop(stack))
  result := &ScriptVal{
    encodeNum(FromInt(element-1)),
  }
  Push(stack, result)
  return true
}

func (t *OP_1SUB) Num() int {
  return 129
}

type OP_NEGATE struct{}

func (t *OP_NEGATE) Operate(z string, stack *Stack, altstack *Stack, cmds *Stack) bool {
  if Len(stack) < 1 {
    return false
  }
  element := intoValue(Pop(stack))
  result := &ScriptVal{
    encodeNum(FromInt(-element)),
  }
  Push(stack, result)
  return true
}

func (t *OP_NEGATE) Num() int {
  return 130
}

type OP_ABS struct {}

func (t *OP_ABS) Operate(z string, stack *Stack, altstack *Stack, cmds *Stack) bool {
  if Len(stack) < 1 {
    return false
  }
  element := intoValue(Pop(stack))
  result := &ScriptVal{
    encodeNum(FromInt(element)),
  }
  if element < 0 {
    result.Val = encodeNum(FromInt(-element))
  }
  return true
}

func (t *OP_ABS) Num() int {
  return 131
}

type OP_NOT struct {}

func (t *OP_NOT) Operate(z string, stack *Stack, altstack *Stack, cmds *Stack) bool {
  if Len(stack) < 1 {
    return false
  }
  element := Pop(stack)
  if intoValue(element) == 0{
    Push(stack, &ScriptVal{
      encodeNum(ONE),
    })
  }else {
    Push(stack, &ScriptVal{
      encodeNum(ZERO),
    })
  }
  return true
}

func (t *OP_NOT) Num() int {
  return 132
}

type OP_0NOTEQUAL struct {}

func (t *OP_0NOTEQUAL) Operate(z string, stack *Stack, altstack *Stack, cmds *Stack) bool {
  if Len(stack) < 1 {
    return false
  }
  element := Pop(stack)
  if intoValue(element) == 0 {
    Push(stack, &ScriptVal {
      encodeNum(ZERO),
    })
  }else {
    Push(stack, &ScriptVal{
      encodeNum(ONE),
    })
  }
  return true
}
func (t *OP_0NOTEQUAL) Num() int {
  return 132
}

type OP_ADD struct {}

func (t *OP_ADD) Operate(z string, stack *Stack, altstack *Stack, cmds *Stack) bool {
  if Len(stack) < 2 {
    return false
  }
  element1 := intoValue(Pop(stack))
  element2 := intoValue(Pop(stack))
  Push(stack, &ScriptVal{
    encodeNum(FromInt(element1+element2)),
  })
  return true
}

func (t *OP_ADD) Num() int {
  return 133
}

type OP_SUB struct {}

func (t *OP_SUB) Operate(z string, stack *Stack, altstack *Stack, cmds *Stack) bool {
  if Len(stack) < 2 {
    return false
  }
  element1 := intoValue(Pop(stack))
  element2 := intoValue(Pop(stack))
  Push(stack, &ScriptVal{
    encodeNum(FromInt(element2-element1)),
  })
  return true
}

func (t *OP_SUB) Num() int {
  return 134
}

type OP_MUL struct {}

func (t *OP_MUL) Operate(z string, stack *Stack, altstack *Stack, cmds *Stack) bool {
  if Len(stack) < 2 {
    return false
  }
  element1 := intoValue(Pop(stack))
  element2 := intoValue(Pop(stack))
  Push(stack, &ScriptVal{
    encodeNum(FromInt(element2 * element1)),
  })
  return true
}

func (t *OP_MUL) Num() int {
  return 135
}

type OP_BOOLAND struct {}

func (t *OP_BOOLAND) Operate(z string, stack *Stack, altstack *Stack, cmds *Stack) bool {
  if Len(stack) < 2 {
    return false
  }
  element1 := intoValue(Pop(stack))
  element2 := intoValue(Pop(stack))
  if element1 + element2 >= 2 {
    Push(stack, &ScriptVal{
      encodeNum(ONE),
    })
  }else {
    Push(stack, &ScriptVal{
      encodeNum(ZERO),
    })
  }
  return true
}

func (t *OP_BOOLAND) Num() int {
  return 136
}

type OP_BOOLOR struct{}

func (t *OP_BOOLOR) Operate(z string, stack *Stack, altstack *Stack, cmds *Stack) bool {
  if Len(stack) < 2 {
    return false
  }

  element1 := intoValue(Pop(stack))
  element2 := intoValue(Pop(stack))
  if element1 + element2 > 0 {
    Push(stack, &ScriptVal{
      encodeNum(ONE),
    })
  }else {
    Push(stack, &ScriptVal{
      encodeNum(ZERO),
    })
  }
  return true
}

func (t *OP_BOOLOR) Num() int {
  return 137
}

type OP_NUMEQUAL struct {}


func (t *OP_NUMEQUAL) Operate(z string, stack *Stack, altstack *Stack, cmds *Stack) bool {
  if Len(stack) < 2 {
    return false
  }
  element1 := intoValue(Pop(stack))
  element2 := intoValue(Pop(stack))
  if element1 == element2 {
    Push(stack, &ScriptVal{
      encodeNum(ONE),
    })
  }else {
    Push(stack, &ScriptVal{
      encodeNum(ZERO),
    })
  }
  return true
}

func (t *OP_NUMEQUAL) Num() int {
  return 138
}

type OP_NUMEQUALVERIFY struct{}


func (t *OP_NUMEQUALVERIFY) Operate(z string, stack *Stack, altstack *Stack, cmds *Stack) bool {
  equal := &OP_NUMEQUAL{}
  verify := &OP_VERIFY{}
  return equal.Operate(z, stack, altstack, cmds) && verify.Operate(z, stack, altstack, cmds)
}

func (t *OP_NUMEQUALVERIFY) Num() int {
  return 139
}

type OP_NUMNOTEQUAL struct {}

func (t *OP_NUMNOTEQUAL) Operate(z string, stack *Stack, altstack *Stack, cmds *Stack) bool {
  if Len(stack) < 2 {
    return false
  }
  element1 := intoValue(Pop(stack))
  element2 := intoValue(Pop(stack))
  if element1 == element2 {
    Push(stack, &ScriptVal{
      encodeNum(ZERO),
    })
  }else {
    Push(stack, &ScriptVal{
      encodeNum(ONE),
    })
  }
  return true
}

func (t *OP_NUMNOTEQUAL) Num() int {
  return 140
}

type OP_LESSTHAN struct {}


func (t *OP_LESSTHAN) Operate(z string, stack *Stack, altstack *Stack, cmds *Stack) bool {
  if Len(stack) < 2 {
    return false
  }
  element1 := intoValue(Pop(stack))
  element2 := intoValue(Pop(stack))
  if element2 < element1 {
    Push(stack, &ScriptVal{
      encodeNum(ONE),
    })
  }else {
    Push(stack, &ScriptVal{
      encodeNum(ZERO),
    })
  }
  return true
}

func (t *OP_LESSTHAN) Num() int {
  return 141
}

type OP_GREATERTHAN struct {}


func (t *OP_GREATERTHAN) Operate(z string, stack *Stack, altstack *Stack, cmds *Stack) bool {
  if Len(stack) < 2 {
    return false
  }
  element1 := intoValue(Pop(stack))
  element2 := intoValue(Pop(stack))
  if element2 > element1 {
    Push(stack, &ScriptVal{
      encodeNum(ONE),
    })
  }else {
    Push(stack, &ScriptVal{
      encodeNum(ZERO),
    })
  }
  return true
}

func (t *OP_GREATERTHAN) Num() int {
  return 142
}

type OP_LESSTHANOREQUAL struct {}


func (t *OP_LESSTHANOREQUAL) Operate(z string, stack *Stack, altstack *Stack, cmds *Stack) bool {
  if Len(stack) < 2 {
    return false
  }
  element1 := intoValue(Pop(stack))
  element2 := intoValue(Pop(stack))
  if element2 <= element1 {
    Push(stack, &ScriptVal{
      encodeNum(ONE),
    })
  }else {
    Push(stack, &ScriptVal{
      encodeNum(ZERO),
    })
  }
  return true
}

func (t *OP_LESSTHANOREQUAL) Num() int {
  return 143
}

type OP_GREATERTHANOREQUAL struct {}


func (t *OP_GREATERTHANOREQUAL) Operate(z string, stack *Stack, altstack *Stack, cmds *Stack) bool {
  if Len(stack) < 2 {
    return false
  }
  element1 := intoValue(Pop(stack))
  element2 := intoValue(Pop(stack))
  if element2 >= element1 {
    Push(stack, &ScriptVal{
      encodeNum(ONE),
    })
  }else {
    Push(stack, &ScriptVal{
      encodeNum(ZERO),
    })
  }
  return true
}

func (t *OP_GREATERTHANOREQUAL) Num() int {
  return 144
}

type OP_MIN struct {}


func (t *OP_MIN) Operate(z string, stack *Stack, altstack *Stack, cmds *Stack) bool {
  if Len(stack) < 2 {
    return false
  }
  element1 := intoValue(Pop(stack))
  element2 := intoValue(Pop(stack))
  if element1 < element2 {
    Push(stack, &ScriptVal{
      encodeNum(FromInt(element1)),
    })
  }else {
    Push(stack, &ScriptVal{
      encodeNum(FromInt(element2)),
    })
  }
  return true
}

func (t *OP_MIN) Num() int {
  return 145
}

type OP_MAX struct {}


func (t *OP_MAX) Operate(z string, stack *Stack, altstack *Stack, cmds *Stack) bool {
  if Len(stack) < 2 {
    return false
  }
  element1 := intoValue(Pop(stack))
  element2 := intoValue(Pop(stack))
  if element1 > element2 {
    Push(stack, &ScriptVal{
      encodeNum(FromInt(element1)),
    })
  }else {
    Push(stack, &ScriptVal{
      encodeNum(FromInt(element2)),
    })
  }
  return true
}

func (t *OP_MAX) Num() int {
  return 146
}

type OP_WITHIN struct {}


func (t *OP_WITHIN) Operate(z string, stack *Stack, altstack *Stack, cmds *Stack) bool {
  if Len(stack) < 3 {
    return false
  }
  maximum := intoValue(Pop(stack))
  minimum := intoValue(Pop(stack))
  element := intoValue(Pop(stack))
  if element < maximum && element >= minimum {
    Push(stack, &ScriptVal{
      encodeNum(ONE),
    })
  }else {
    Push(stack, &ScriptVal{
      encodeNum(ZERO),
    })
  }
  return true
}

func (t *OP_WITHIN) Num() int {
  return 147
}

type OP_RIPEMD160 struct {}


func (t *OP_RIPEMD160) Operate(z string, stack *Stack, altstack *Stack, cmds *Stack) bool {
  if Len(stack) < 1 {
    return false
  }
  element := Pop(stack)
  val, ok := element.(*ScriptVal)
  if !ok {
    val = &ScriptVal{
      encodeNum(FromInt(element.Num())),
    }
  }
  hashed := Hash160(val.Val)
  Push(stack, &ScriptVal{
    hashed,
  })
  return true
}

func (t *OP_RIPEMD160) Num() int {
  return 147
}

type OP_SHA1 struct {}


func (t *OP_SHA1) Operate(z string, stack *Stack, altstack *Stack, cmds *Stack) bool {
  if Len(stack) < 1 {
    return false
  }
  element := Pop(stack)
  val, ok := element.(*ScriptVal)
  if !ok {
    val = &ScriptVal{
      encodeNum(FromInt(element.Num())),
    }
  }
  sha := sha1.New()
  sha.Write(val.Val)
  Push(stack, &ScriptVal{
    sha.Sum(nil),
  })
  return true
}

func (t *OP_SHA1) Num() int {
  return 148
}

type OP_SHA256 struct {}


func (t *OP_SHA256) Operate(z string, stack *Stack, altstack *Stack, cmds *Stack) bool {
  if Len(stack) < 1 {
    return false
  }
  element := Pop(stack)
  val, ok := element.(*ScriptVal)
  if !ok {
    val = &ScriptVal{
      encodeNum(FromInt(element.Num())),
    }
  }
  sha := sha256.New()
  sha.Write(val.Val)
  Push(stack, &ScriptVal{
    sha.Sum(nil),
  })
  return true
}

func (t *OP_SHA256) Num() int {
  return 149
}

type OP_HASH160 struct {}


func (t *OP_HASH160) Operate(z string, stack *Stack, altstack *Stack, cmds *Stack) bool {
  if Len(stack) < 1 {
    return false
  }
  element := Pop(stack)
  val, ok := element.(*ScriptVal)
  if !ok {
    val = &ScriptVal{
      encodeNum(FromInt(element.Num())),
    }
  }
  Push(stack, &ScriptVal{
    Hash160(val.Val),
  })
  return true
}

func (t *OP_HASH160) Num() int {
  return 150
}

type OP_HASH256 struct {}


func (t *OP_HASH256) Operate(z string, stack *Stack, altstack *Stack, cmds *Stack) bool {
  if Len(stack) < 1 {
    return false
  }
  element := Pop(stack)
  val, ok := element.(*ScriptVal)
  if !ok {
    val = &ScriptVal{
      encodeNum(FromInt(element.Num())),
    }
  }
  Push(stack, &ScriptVal{
    Hash256(val.Val),
  })
  return true
}

func (t *OP_HASH256) Num() int {
  return 151
}


