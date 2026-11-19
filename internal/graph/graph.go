package graph

import "gpath/internal/geo"

type Node struct{ Lat, Lon float64 }

type Edge struct {
	To   int
	Dist float64
}

type Bounds struct {
	Min, Max geo.Point
}

type Graph struct {
	Nodes  []Node
	Edges  [][]Edge
	Bounds Bounds
}

func NewGraph() *Graph {
	return &Graph{Nodes: []Node{}, Edges: [][]Edge{}}
}

func (n Node) ToPoint() geo.Point {
	return geo.Point{Lat: n.Lat, Lon: n.Lon}
}
