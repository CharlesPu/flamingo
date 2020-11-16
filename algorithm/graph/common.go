package graph

import (
	"bytes"
	"fmt"
	"math"
)

type (
	Graph [][]uint32

	pathPlanner interface {
		pathPlan(graph Graph, src, dst uint32) uint32
	}
)

const (
	INF = math.MaxUint32
)

func NewGraph(size int) Graph {
	g := make(Graph, size, size)
	for i := range g {
		g[i] = make([]uint32, size, size)
	}
	return g
}

func (g Graph) isValid() bool {
	if len(g) == 0 {
		return false
	}
	for i := range g {
		if len(g[i]) != g.size() {
			return false
		}
	}
	return true
}

func (g Graph) size() int {
	return len(g)
}

func (g Graph) plan(itf pathPlanner, src, dst uint32) uint32 {
	if !g.isValid() {
		return 0
	}

	return itf.pathPlan(g, src, dst)
}

func (g Graph) distance(src, dst uint32) uint32 {
	if !g.isValid() {
		return 0
	}
	if src >= uint32(len(g)) || dst >= uint32(len(g)) {
		return 0
	}
	return g[src][dst]
}

func (g Graph) String() string {
	var strBuf bytes.Buffer

	strBuf.WriteString("map: \n")
	for _, v := range g {
		strBuf.WriteString("[")
		for i, vv := range v {
			strTmp := fmt.Sprintf("%+v", vv)
			if vv == INF {
				strTmp = "∞"
			}
			strBuf.WriteString(strTmp)
			if i != len(v)-1 {
				strBuf.WriteString("\t")
			}
		}
		strBuf.WriteString("]\n")
	}

	return strBuf.String()
}
