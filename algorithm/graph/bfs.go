package graph

import (
	"github.com/CharlesPu/flamingo/container/queue"
	"github.com/CharlesPu/flamingo/plog"
)

// Breadth First Search
// 在 BFS 中，我们使用了数据结构中的一个队列（queue），
// 我们知道队列的特性是 FIFO（First In First Out），也就是先进先出。
// 正是这个 FIFO 特性，保证了我们第一个到达目标节点一定是最短路径。
// 这里注意带权值的图和不带权值的图的算法会不同，前者复杂一点

func (g Graph) Bfs(src, dst uint32) uint32 {
	return g.plan(new(bfs).init(g.size()), src, dst)
}

type bfs struct {
	q     *queue.Queue
	visit [][]bool
}

func (b *bfs) init(size int) *bfs {
	b.q = queue.New(uint32(size * size))
	b.visit = make([][]bool, size, size)
	for i := range b.visit {
		b.visit[i] = make([]bool, size, size)
	}
	return b
}

// 与dijkstra有相似之处
func (b *bfs) pathPlan(graph Graph, src, dst uint32) uint32 {
	b.q.InInt(int(src))

	distance := newDistance(graph, src)
	for !b.q.Empty() {
		curPoint := b.q.OutInt()
		plog.Debugf("----------find for src %+v", curPoint)
		for tryPoint := 0; tryPoint < graph.size(); tryPoint++ {
			plog.Debugf("try to find %+v -> %+v", curPoint, tryPoint)
			if graph[curPoint][tryPoint] != INF && // 如果能通
				distance[curPoint] != INF &&
				!b.visit[curPoint][tryPoint] && // 且没有走过
				tryPoint != curPoint { // 也不是自己
				b.visit[curPoint][tryPoint] = true
				plog.Debugf("proc %+v", tryPoint)
				dis := graph[curPoint][tryPoint] + distance[curPoint]
				if dis < distance[tryPoint] {
					plog.Debugf("find min %+v on node %+v", dis, tryPoint)
					distance[tryPoint] = dis
				}
				b.q.InInt(tryPoint)
			}
		}
	}

	return distance[dst]
}

func newDistance(graph Graph, src uint32) []uint32 {
	distance := make([]uint32, graph.size(), graph.size())
	for i := 0; i < graph.size(); i++ {
		distance[i] = graph[src][i]
	}
	return distance
}
func (g Graph) BfsNoWeight(src, dst uint32) uint32 {
	return g.plan(new(bfsNoWeight).init(g.size()), src, dst)
}

type bfsNoWeight struct {
	q      *queue.Queue
	visit  [][]bool
	parent []int
	dis    uint32
}

func (b *bfsNoWeight) init(size int) *bfsNoWeight {
	b.q = queue.New(uint32(size * size))
	b.visit = make([][]bool, size, size)
	for i := range b.visit {
		b.visit[i] = make([]bool, size, size)
	}
	b.parent = make([]int, size, size)
	for i := range b.parent {
		b.parent[i] = -1
	}
	return b
}

func (b *bfsNoWeight) pathPlan(graph Graph, src, dst uint32) uint32 {
	b.q.InInt(int(src))

	var found bool
	for !b.q.Empty() {
		curPoint := b.q.OutInt()
		plog.Debugf("----------find for src %+v", curPoint)
		for tryPoint := 0; tryPoint < graph.size(); tryPoint++ {
			plog.Debugf("try to find %+v -> %+v", curPoint, tryPoint)
			if graph[curPoint][tryPoint] != INF && // 如果能通
				!b.visit[curPoint][tryPoint] && // 且没有走过 // todo 这里只要[]bool?
				tryPoint != curPoint { // 也不是自己
				b.parent[tryPoint] = curPoint
				b.visit[curPoint][tryPoint] = true
				plog.Debugf("proc %+v", tryPoint)
				if tryPoint == int(dst) {
					found = true
					break
				}
				b.q.InInt(tryPoint)
			}
		}
		if found {
			break
		}
	}

	if !found {
		return 0
	}
	plog.Debugf("parent: %+v", b.parent)
	idx := int(dst)
	for b.parent[idx] != -1 {
		b.dis += graph[b.parent[idx]][idx]
		idx = b.parent[idx]
	}

	return b.dis
}
