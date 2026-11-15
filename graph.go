package main

import (
	"math"
)

type Node struct{ Lat, Lon float64 }

type Edge struct {
	To   int64
	Dist float64
}

type Graph struct {
	Nodes map[int64]Node
	Edges map[int64][]Edge
}

func NewGraph() *Graph {
	return &Graph{Nodes: map[int64]Node{}, Edges: map[int64][]Edge{}}
}

func haversineDistance(a, b Node) float64 {
	const R = 6371000.0
	rad := math.Pi / 180
	p1, p2 := a.Lat*rad, b.Lat*rad
	dp, dl := (b.Lat-a.Lat)*rad, (b.Lon-a.Lon)*rad
	h := math.Sin(dp/2)*math.Sin(dp/2) +
		math.Cos(p1)*math.Cos(p2)*math.Sin(dl/2)*math.Sin(dl/2)
	return 2 * R * math.Asin(math.Sqrt(h))
}

// TODO: make this a KD-tree for final implementation
func (graph *Graph) nearestTo(to Node) int64 {
	var nearestId int64
	dist := math.Inf(1)

	for id, node := range graph.Nodes {
		newDist := haversineDistance(node, to)
		if newDist < dist {
			nearestId = id
			dist = newDist
		}
	}

	return nearestId
}
