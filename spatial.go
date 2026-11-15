package main

import "math"

// TODO: make this a KD-tree for final implementation
func (graph *Graph) nearestTo(to Node) int64 {
	var nearestId int64
	dist := math.Inf(1)

	for id, node := range graph.Nodes {
		newDist := HaversineDistance(node, to)
		if newDist < dist {
			nearestId = id
			dist = newDist
		}
	}

	return nearestId
}
