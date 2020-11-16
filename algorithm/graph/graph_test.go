package graph

import (
	"testing"
)

var g = Graph{
	{0, 2, 6, 4},
	{INF, 0, 3, INF},
	{7, INF, 0, 1},
	{5, INF, 12, 0},
}

// result
// [0	2	5	4]
// [9	0	3	4]
// [6	8	0	1]
// [5	7	10	0]

func TestFloyd(t *testing.T) {
	t.Log(g.Floyd(2, 3))
}

func TestDijkstra(t *testing.T) {
	t.Log(g.Dijkstra(3, 2))
}

func TestBfs(t *testing.T) {
	t.Log(g.Bfs(3, 2))
}

func TestBfsNoWeight(t *testing.T) {
	t.Log(g.BfsNoWeight(3, 2))
}
