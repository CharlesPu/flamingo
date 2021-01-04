package find

type (
	// 并查集算法
	UnionFind interface {
		Union(p, q int)
		Connected(p, q int) bool
		UnionCount() int

		findRoot(x int) int
	}

	unionFind struct {
		count int

		parent []int
		size   []int // i为根时的树的大小
	}
)

func NewUnionFind(n int) UnionFind {
	uf := &unionFind{}
	uf.count = n
	uf.parent = make([]int, n)
	uf.size = make([]int, n)

	for i := 0; i < n; i++ {
		uf.parent[i] = i
		uf.size[i] = 1
	}

	return uf
}

func (uf *unionFind) Union(p, q int) {
	pRoot := uf.findRoot(p)
	qRoot := uf.findRoot(q)
	if pRoot == qRoot {
		return
	}

	if uf.size[qRoot] > uf.size[pRoot] { // merge p to q
		uf.parent[pRoot] = qRoot
		uf.size[qRoot] += uf.size[pRoot]
	} else {
		uf.parent[qRoot] = pRoot
		uf.size[pRoot] += uf.size[qRoot]
	}
	uf.count--
}

func (uf *unionFind) Connected(p, q int) bool {
	return uf.findRoot(p) == uf.findRoot(q)
}

func (uf *unionFind) UnionCount() int {
	return uf.count
}

func (uf *unionFind) findRoot(x int) int {
	for uf.parent[x] != x {
		uf.parent[x] = uf.parent[uf.parent[x]]
		x = uf.parent[x]
	}
	return x
}
