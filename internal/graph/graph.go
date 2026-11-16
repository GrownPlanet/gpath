package graph

import "gpath/internal/geo"

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

func (n Node) ToPoint() geo.Point {
	return geo.Point{Lat: n.Lat, Lon: n.Lon}
}
