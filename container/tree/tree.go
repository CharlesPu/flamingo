package tree

type (
	Tree struct {
		root *Node
	}
	Node struct {
		value interface{}
		left  *Node
		right *Node
	}
)

func NewTree() *Tree {
	return &Tree{}
}
