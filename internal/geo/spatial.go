package geo

import (
	"math"

	"gpath/internal/graph"
)

// TODO: make this a KD-tree for final implementation
func NearestTo(g *graph.Graph, to Point) int64 {
	var nearestId int64
	dist := math.Inf(1)

	for id, node := range g.Nodes {
		newDist := HaversineDistance(Point{Lat: node.Lat, Lon: node.Lon}, to)
		if newDist < dist {
			nearestId = id
			dist = newDist
		}
	}

	return nearestId
}
