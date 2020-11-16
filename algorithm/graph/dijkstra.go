package graph

// 单源最短路径算法
// 贪心思想

// 首先把起点到所有点的距离存下来找个最短的，然后松弛一次再找出最短的，
// 所谓的松弛操作就是，遍历一遍看通过刚刚找到的距离最短的点作为中转站会不会更近，
// 如果更近了就更新距离，这样把所有的点找遍之后就存下了起点到其他所有点的最短距离。

func (g Graph) Dijkstra(src, dst uint32) uint32 {
	return g.plan(&dijkstra{}, src, dst)
}

type dijkstra struct{}

func (d *dijkstra) pathPlan(graph Graph, src, dst uint32) uint32 {
	// for src := 0; src < graph.size(); src++ {
	graph[src] = dijkstraPlan(graph, src)
	// }
	return graph.distance(src, dst)
}

func dijkstraPlan(graph Graph, src uint32) []uint32 {
	size := graph.size()
	visit := make([]bool, size, size)
	distance := make([]uint32, size, size)

	for i := 0; i < size; i++ {
		distance[i] = graph[src][i]
	}
	visit[src] = true

	for loop := 0; loop < size-1; loop++ { // n-1次即可，保证每个点都visit过
		var min uint32 = INF
		tmpNode := 0

		// 在剩下没有遍历过的点中，找到离src最近的
		for i := 0; i < size; i++ {
			if !visit[i] && distance[i] < min {
				min = distance[i]
				tmpNode = i
			}
		}
		visit[tmpNode] = true

		// 在经过该最近点的基础之上更新dis
		for dst := 0; dst < size; dst++ {
			if graph[tmpNode][dst] == INF {
				continue
			}
			disTmp := graph[tmpNode][dst] + distance[tmpNode]
			if disTmp < distance[dst] {
				distance[dst] = disTmp
				// todo 这里可以记录下最短距离所经过的点
			}
		}
	}

	return distance
}
