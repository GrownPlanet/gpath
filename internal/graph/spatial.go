package graph

import (
	"math"

	"gpath/internal/geo"
)

type KDNode struct {
	id          int64
	left, right *KDNode
}

func getValue(g *Graph, ids []int64, axis, i int64) float64 {
	if axis%2 == 0 {
		return g.Nodes[ids[i]].Lat
	} else {
		return g.Nodes[ids[i]].Lon
	}
}

func partition(g *Graph, ids []int64, axis, l, r, pivot int64) int64 {
	pivotValue := getValue(g, ids, axis, pivot)

	idx := l
	for i := l; i < r; i++ {
		if getValue(g, ids, axis, i) <= pivotValue {
			ids[i], ids[idx] = ids[idx], ids[i]
			idx += 1
		}
	}

	ids[r], ids[idx] = ids[idx], ids[r]

	return idx
}

func selectMedian(g *Graph, ids []int64, axis, l, r int64) (median, i int64) {
	k := l + (r-l)/2
	for {
		if l == r {
			return ids[l], l
		}
		pivot := l + (r-l)/2
		pivot = partition(g, ids, axis, l, r, pivot)
		if k == pivot {
			return ids[k], k
		} else if k < pivot {
			r = pivot - 1
		} else {
			l = pivot + 1
		}
	}
}

func build(g *Graph, ids []int64, depth int64) *KDNode {
	if len(ids) == 0 {
		return nil
	}

	axis := depth % 2
	median, i := selectMedian(g, ids, axis, 0, int64(len(ids)-1))

	return &KDNode{
		id:    median,
		left:  build(g, ids[:i], depth+1),
		right: build(g, ids[i+1:], depth+1),
	}
}

func BuildKDTree(g *Graph) *KDNode {
	ids := make([]int64, 0, len(g.Nodes))
	for k := range g.Nodes {
		ids = append(ids, k)
	}

	return build(g, ids, 0)
}

// TODO: make this a KD-tree for final implementation
func (g *Graph) NearestTo(to geo.Point) int64 {
	var nearestId int64
	dist := math.Inf(1)

	for id, node := range g.Nodes {
		newDist := geo.HaversineDistance(node.ToPoint(), to)
		if newDist < dist {
			nearestId = id
			dist = newDist
		}
	}

	return nearestId
}
