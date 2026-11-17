package graph

import "gpath/internal/geo"

type Node struct{ Lat, Lon float64 }

type Edge struct {
	To   int
	Dist float64
}

type Graph struct {
	Nodes []Node
	Edges [][]Edge
}

func NewGraph() *Graph {
	return &Graph{Nodes: []Node{}, Edges: [][]Edge{}}
}

func (n Node) ToPoint() geo.Point {
	return geo.Point{Lat: n.Lat, Lon: n.Lon}
}
