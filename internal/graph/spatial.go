package graph

import (
	"math"

	"gpath/internal/geo"
)

// TODO: make this a KD-tree for final implementation
func (g *Graph) NearestTo(to geo.Point) int64 {
	var nearestId int64
	dist := math.Inf(1)

	for id, node := range g.Nodes {
		newDist := geo.HaversineDistance(geo.Point{Lat: node.Lat, Lon: node.Lon}, to)
		if newDist < dist {
			nearestId = id
			dist = newDist
		}
	}

	return nearestId
}
