package tree

import "testing"

func TestAVL(t *testing.T) {
	at := NewAVLTree(1, 2, 3, 4, 5, 6, 7, 8)

	t.Log(at.String())
}
