package graph

import "math"

// 多源最短路径算法

// 最开始只允许经过1号顶点进行中转，
// 接下来只允许经过1和2号顶点进行中转……允许经过1~n号所有顶点进行中转，
// 求任意两点之间的最短路程。
// 用一句话概括就是：从i号顶点到j号顶点只经过k号点的最短路程。

func (g Graph) Floyd(src, dst uint32) uint32 {
	return g.plan(&floyd{}, src, dst)
}

type floyd struct{}

func (f *floyd) pathPlan(graph Graph, src, dst uint32) uint32 {
	size := graph.size()
	for pass := 0; pass < size; pass++ {
		for src := 0; src < size; src++ {
			for dst := 0; dst < size; dst++ {
				if graph[src][pass] == math.MaxUint32 ||
					graph[pass][dst] == math.MaxUint32 {
					continue
				}
				tmp := graph[src][pass] + graph[pass][dst]
				if tmp < graph[src][dst] {
					graph[src][dst] = tmp
				}
			}
		}
	}
	return graph.distance(src, dst)
}
