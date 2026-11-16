package geo

import (
	"gpath/internal/graph"
	"math"
)

// TODO: make this a KD-tree for final implementation
func NearestTo(g *graph.Graph, lat, lon float64) int64 {
	var nearestId int64
	dist := math.Inf(1)

	to := graph.Node{Lat: lat, Lon: lon}

	for id, node := range g.Nodes {
		newDist := HaversineDistance(node, to)
		if newDist < dist {
			nearestId = id
			dist = newDist
		}
	}

	return nearestId
}
