package pathfinding

import (
	"container/heap"
	"fmt"
	"math"

	"gpath/internal/geo"
	"gpath/internal/graph"
)

func reverse(arr []int64) {
	for i, j := 0, len(arr)-1; i < j; i, j = i+1, j-1 {
		arr[i], arr[j] = arr[j], arr[i]
	}
}

func reconstructPath(cameFrom map[int64]int64, start, current int64) []int64 {
	path := []int64{current}
	for current != start {
		current = cameFrom[current]
		path = append(path, current)
	}
	reverse(path)
	return path
}

type Heuristic func(id int64) float64

func aStar(g *graph.Graph, start, end int64, h Heuristic) ([]int64, error) {
	open := &pq{{id: start, g: 0, f: h(start)}}
	cameFrom := map[int64]int64{}
	gScore := map[int64]float64{start: 0.0}

	for open.Len() != 0 {
		cur := heap.Pop(open).(item)
		if cur.g > gScore[cur.id] {
			continue // stale entry
		}
		if cur.id == end {
			return reconstructPath(cameFrom, start, end), nil
		}

		for _, neighbour := range g.Edges[cur.id] {
			tentative := cur.g + neighbour.Dist
			old, ok := gScore[neighbour.To]
			if !ok {
				old = math.Inf(1)
			}

			if tentative >= old {
				continue
			}

			cameFrom[neighbour.To] = cur.id
			gScore[neighbour.To] = tentative

			heap.Push(open, item{
				id: neighbour.To,
				g:  tentative,
				f:  tentative + h(neighbour.To)},
			)
		}
	}

	return []int64{}, fmt.Errorf("no path found")
}

func idsToPoints(g *graph.Graph, ids []int64) []geo.Point {
	points := make([]geo.Point, len(ids))
	for i, id := range ids {
		points[i] = g.Nodes[id].ToPoint()
	}
	return points
}

func ShortestPath(g *graph.Graph, kdTree *graph.KDNode, start, end geo.Point) ([]geo.Point, error) {
	startId := kdTree.NearestTo(g, start)
	endId := kdTree.NearestTo(g, end)

	goal := g.Nodes[endId].ToPoint()
	h := func(id int64) float64 {
		return geo.HaversineDistance(g.Nodes[id].ToPoint(), goal)
	}

	path, err := aStar(g, startId, endId, h)
	if err != nil {
		return []geo.Point{}, err
	}

	return idsToPoints(g, path), nil
}
