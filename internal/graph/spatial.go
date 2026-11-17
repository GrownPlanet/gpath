package graph

import (
	"math"

	"gpath/internal/geo"
)

type KDNode struct {
	id          int
	left, right *KDNode
}

func getValue(g *Graph, ids []int, axis, i int) float64 {
	if axis == 0 {
		return g.Nodes[ids[i]].Lat
	} else {
		return g.Nodes[ids[i]].Lon
	}
}

func partition(g *Graph, ids []int, axis, l, r, pivot int) int {
	pivotValue := getValue(g, ids, axis, pivot)
	ids[pivot], ids[r] = ids[r], ids[pivot]

	idx := l
	for i := l; i < r; i++ {
		if getValue(g, ids, axis, i) < pivotValue {
			ids[i], ids[idx] = ids[idx], ids[i]
			idx++
		}
	}

	ids[r], ids[idx] = ids[idx], ids[r]

	return idx
}

func selectMedian(g *Graph, ids []int, axis, l, r int) int {
	k := l + (r-l)/2
	for {
		if l == r {
			return l
		}

		pivot := l + (r-l)/2
		pivot = partition(g, ids, axis, l, r, pivot)

		if k == pivot {
			return k
		} else if k < pivot {
			r = pivot - 1
		} else {
			l = pivot + 1
		}
	}
}

func build(g *Graph, ids []int, depth int) *KDNode {
	if len(ids) == 0 {
		return nil
	}

	axis := depth % 2
	i := selectMedian(g, ids, axis, 0, len(ids)-1)

	return &KDNode{
		id:    ids[i],
		left:  build(g, ids[:i], depth+1),
		right: build(g, ids[i+1:], depth+1),
	}
}

func BuildKDTree(g *Graph) *KDNode {
	ids := make([]int, 0, len(g.Nodes))
	for k := range g.Nodes {
		ids = append(ids, k)
	}

	return build(g, ids, 0)
}

func (root *KDNode) NearestTo(g *Graph, target geo.Point) int {
	var bestID int
	bestDist := math.Inf(1)

	var search func(node *KDNode, depth int)
	search = func(node *KDNode, depth int) {
		if node == nil {
			return
		}

		cur := g.Nodes[node.id]
		dist := geo.HaversineDistance(cur.ToPoint(), target)

		if dist < bestDist {
			bestDist = dist
			bestID = node.id
		}

		axis := depth % 2
		var diff float64
		var splitPoint geo.Point

		if axis == 0 {
			diff = target.Lat - cur.Lat
			splitPoint = geo.Point{Lat: cur.Lat, Lon: target.Lon}
		} else {
			diff = target.Lon - cur.Lon
			splitPoint = geo.Point{Lat: target.Lat, Lon: cur.Lon}
		}

		far, near := node.left, node.right
		if diff < 0 {
			far, near = node.right, node.left
		}

		search(near, depth+1)

		if geo.HaversineDistance(target, splitPoint) < bestDist {
			search(far, depth+1)
		}
	}

	search(root, 0)

	return bestID
}

// Stupid version
// func (g *Graph) NearestTo(point geo.Point) int64 {
// 	var nearestId int64
// 	dist := math.Inf(1)
//
// 	for id, node := range g.Nodes {
// 		newDist := geo.HaversineDistance(node.ToPoint(), point)
// 		if newDist < dist {
// 			nearestId = id
// 			dist = newDist
// 		}
// 	}
//
// 	return nearestId
// }
