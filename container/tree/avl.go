package tree

import (
	"fmt"

	"github.com/CharlesPu/flamingo/operate"
)

type (
	AVLTree struct {
		root *avlNode
	}
	avlNode struct {
		value  int
		height int
		left   *avlNode
		right  *avlNode
	}
)

func NewAVLTree(values ...int) *AVLTree {
	var r int
	if len(values) > 0 {
		r = values[0]
	}
	tree := &AVLTree{root: &avlNode{value: r, height: 1}}
	for i := 1; i < len(values); i++ {
		tree.InsertNode(values[i])
	}
	return tree
}

func (at *AVLTree) InsertNode(value int) {
	at.root = at.root.insertNode(value)
}

func (an *avlNode) insertNode(value int) *avlNode {
	if an == nil {
		return &avlNode{value: value, height: 1}
	}
	if value < an.value {
		an.left = an.left.insertNode(value)
		an.height = operate.MaxInt(an.left.getHeight(), an.right.getHeight()) + 1
		if an.left.getHeight()-an.right.getHeight() >= 2 { // need rotate
			if value < an.left.value {
				return an.llRotateR()
			} else {
				return an.lrRotateLr()
			}
		}
	} else {
		an.right = an.right.insertNode(value)
		an.height = operate.MaxInt(an.left.getHeight(), an.right.getHeight()) + 1
		if an.right.getHeight()-an.left.getHeight() >= 2 {
			if value > an.right.value {
				return an.rrRotateL()
			} else {
				return an.rlRotateRl()
			}
		}
	}
	return an
}

func (an *avlNode) llRotateR() *avlNode {
	fmt.Printf("llRotateR on node %+v\n", an.value)
	l := an.left
	an.left = l.right
	l.right = an

	// update height after rotate
	an.height = operate.MaxInt(an.left.getHeight(), an.right.getHeight()) + 1
	l.height = operate.MaxInt(l.left.getHeight(), l.right.getHeight()) + 1

	return l
}

func (an *avlNode) lrRotateLr() *avlNode {
	fmt.Printf("lrRotateLr on node %+v\n", an.value)
	an.left = an.left.rrRotateL()
	return an.llRotateR()
}

func (an *avlNode) rrRotateL() *avlNode {
	fmt.Printf("rrRotateL on node %+v\n", an.value)
	l := an.right
	an.right = l.left
	l.left = an

	an.height = operate.MaxInt(an.left.getHeight(), an.right.getHeight()) + 1
	l.height = operate.MaxInt(l.left.getHeight(), l.right.getHeight()) + 1

	return l
}

func (an *avlNode) rlRotateRl() *avlNode {
	fmt.Printf("rlRotateRl on node %+v\n", an.value)
	an.right.llRotateR()
	return an.rrRotateL()
}

func (an *avlNode) getHeight() int {
	if an == nil {
		return 0
	}
	return an.height
}

func (at *AVLTree) String() string {
	var res []string

	var queue []*avlNode
	queue = append(queue, at.root)
	res = append(res, fmt.Sprintf("{%+v|%+v}", at.root.value, at.root.height))
	for len(queue) != 0 {
		p := queue[0]
		queue = queue[1:]
		if p.left == nil && p.right == nil {
			continue
		}
		if p.left != nil {
			res = append(res, fmt.Sprintf("{%+v|%+v}", p.left.value, p.left.height))
			queue = append(queue, p.left)
		} else {
			res = append(res, "nil")
		}
		if p.right != nil {
			res = append(res, fmt.Sprintf("{%+v|%+v}", p.right.value, p.right.height))
			queue = append(queue, p.right)
		} else {
			res = append(res, "nil")
		}
	}

	return fmt.Sprintf("%+v", res)
}
