package main

import "fmt"

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

func (graph *Graph) Print() {
	fmt.Println("Nodes:")
	for k, node := range graph.Nodes {
		fmt.Println(k, " -> ", node.Lat, ",", node.Lon)
	}
	fmt.Println()
	fmt.Println("Edges:")
	for k, edges := range graph.Edges {
		fmt.Print(k, " -> ")
		for _, edge := range edges {
			fmt.Print(edge.To, "(", edge.Dist, "); ")
		}
		fmt.Println()
	}
}
