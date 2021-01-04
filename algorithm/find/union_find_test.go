package find

import "testing"

func TestUnionFind(t *testing.T) {
	uf := NewUnionFind(4)

	uf.Union(1, 3)

	t.Log(uf)
	uf.Union(1, 2)
	t.Log(uf)
	t.Log(uf.Connected(1, 3), uf.Connected(1, 0))
}
