package bitcoinlib

import "slices"



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
