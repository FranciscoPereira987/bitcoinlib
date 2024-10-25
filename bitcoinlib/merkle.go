package bitcoinlib

import (
	"slices"
)



type MerkleTree struct {
  levels [][][]byte
  currentDepth int
  currentIndex int
  total int
}

func NewMerkleTree(total int) *MerkleTree {
  levels := make([][][]byte, 0)
  actual := total
  for actual > 1 {
    levels = append(levels, make([][]byte, actual))
    if actual % 2 == 1 {
      actual++
    }
    actual /= 2
  }
  levels = append(levels, make([][]byte, 1))
  slices.Reverse(levels)
  return &MerkleTree{
    levels,
    0,
    0,
    total,
  }
}

func MerkleParent(a, b []byte) []byte {
  return Hash256(append(a, b...)) 
}

func MerkleParentLevel(children [][]byte) [][]byte {
  if len(children) % 2 == 1 {
    children = append(children, children[len(children)-1])
  }
  parentLevel := make([][]byte, len(children)/2)
  for i := range (len(children)/2) {
    j := i * 2
    parentLevel[i] = MerkleParent(children[j], children[j+1])
  }
  return parentLevel
}

func copyAndReverseLeaves(leaves [][]byte) [][]byte {
  actual := make([][]byte, len(leaves))
  copy(actual, leaves)
  for _, leaf := range leaves {
    slices.Reverse(leaf)
  }
  return actual
}

func MerkleRoot(leaves [][]byte) []byte {
  actual := leaves
  for len(actual) > 1 {
    actual = MerkleParentLevel(actual)
  }
  return actual[0]
}
