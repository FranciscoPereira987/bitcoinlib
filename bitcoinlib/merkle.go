package bitcoinlib

import (
	"fmt"
	"slices"
)

type MerkleTree struct {
	levels       [][][]byte
	currentDepth int
	currentIndex int
	total        int
}

func NewMerkleTree(total int) *MerkleTree {
	levels := make([][][]byte, 0)
	actual := total
	for actual > 1 {
		levels = append(levels, make([][]byte, actual))
		if actual%2 == 1 {
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

func (m *MerkleTree) Fill(hashes [][]byte) error {
	if len(hashes) != m.total {
		return fmt.Errorf("invalid hashes, expected array of length %d instead of %d", m.total, len(hashes))
	}
	m.levels[len(m.levels)-1] = hashes
	for i := len(m.levels) - 2; i >= 0; i-- {
		m.levels[i] = MerkleParentLevel(m.levels[i+1])
	}
	return nil
}

func (m *MerkleTree) Up() error {
	if m.currentDepth == 0 {
		return fmt.Errorf("Already at level 0")
	}
	m.currentDepth--
	m.currentIndex /= 2
	return nil
}

func (m *MerkleTree) Left() error {
	if m.currentDepth == len(m.levels) {
		return fmt.Errorf("Already at the lowest level")
	}
	m.currentDepth++
	m.currentIndex *= 2
	return nil
}

func (m *MerkleTree) Right() error {
	if m.currentDepth == len(m.levels) {
		return fmt.Errorf("Already at the lowest level")
	}
	m.currentDepth++
	m.currentIndex *= 2
	m.currentIndex++
	if m.currentIndex >= len(m.levels[m.currentDepth]) {
		m.currentDepth--
		m.currentIndex /= 2
		return fmt.Errorf("No right node")
	}
	return nil
}

func (m *MerkleTree) Root() []byte {
	return m.levels[0][0]
}

func (m *MerkleTree) SetCurrentNode(value []byte) {
	m.levels[m.currentDepth][m.currentIndex] = value
}

func (m *MerkleTree) GetCurrentNode() []byte {
	return m.levels[m.currentDepth][m.currentIndex]
}

// Returns the left and right children of the current Node
// If the current Node is a left, returns nil, nil
func (m *MerkleTree) GetChildren() (left []byte, right []byte) {
	if m.IsLeaf() {
		return
	}
	m.Left()
	left = m.GetCurrentNode()
	m.Up()
	if m.RightExists() {
		m.Right()
		right = m.GetCurrentNode()
		m.Up()
	} else {
		right = left
	}
	return
}

func (m *MerkleTree) IsLeaf() bool {
	return m.currentDepth == len(m.levels)-1
}

func (m *MerkleTree) RightExists() bool {
	return !m.IsLeaf() && len(m.levels[m.currentDepth+1])-1 > 2*m.currentIndex
}

func (m *MerkleTree) populateTree(flags []bool, hashes [][]byte) {
	//Start by positioning myself on the root
	m.currentDepth = 0
	m.currentIndex = 0
	for _, hash := range hashes {
		for i, flag := range flags {
			if !flag || m.IsLeaf() {
				m.SetCurrentNode(hash)
				flags = flags[i+1:]
				m.Up()
				for m.GetCurrentNode() != nil && m.currentDepth != 0 {
					m.Up()
					if left, right := m.GetChildren(); len(left) == 32 && len(right) == 32 {
						m.SetCurrentNode(MerkleParent(left, right))
					}
				}
				m.SetCurrentNode([]byte{})
				m.Right()
				break
			}
			m.Left()
		}
	}
}

func MerkleParent(a, b []byte) []byte {
	return Hash256(append(a, b...))
}

func MerkleParentLevel(children [][]byte) [][]byte {
	if len(children)%2 == 1 {
		children = append(children, children[len(children)-1])
	}
	parentLevel := make([][]byte, len(children)/2)
	for i := range len(children) / 2 {
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
