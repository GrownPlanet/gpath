package pathfinding

import (
	"fmt"
	"math"
	"slices"

	"gpath/internal/geo"
	"gpath/internal/graph"
)

func lowestScore(arr []int64, scores map[int64]float64) (i, arr_i int64) {
	var lowestIdx, arrIdx int64
	lowestScore := math.Inf(1)

	for i, arr_i := range arr {
		score, ok := scores[arr_i]
		if !ok {
			score = math.Inf(1)
		}
		if score < lowestScore {
			arrIdx = int64(i)
			lowestIdx = arr_i
			lowestScore = score
		}
	}

	return lowestIdx, arrIdx
}

func reverse(arr []int64) {
	for i, j := 0, len(arr)-1; i < j; i, j = i+1, j-1 {
		arr[i], arr[j] = arr[j], arr[i]
	}
}

func reconstructPath(came_from map[int64]int64, start, current int64) []int64 {
	path := []int64{current}
	for current != start {
		current = came_from[current]
		path = append(path, current)
	}
	reverse(path)
	return path
}

type Heuristic func(id int64) float64

func aStar(g *graph.Graph, start, end int64, h Heuristic) ([]int64, error) {
	open := []int64{start}
	came_from := map[int64]int64{}

	g_score := map[int64]float64{start: 0.0}
	f_score := map[int64]float64{start: 0.0}

	for len(open) != 0 {
		current, i := lowestScore(open, f_score)
		if current == end {
			return reconstructPath(came_from, start, current), nil
		}

		open = append(open[:i], open[i+1:]...)
		for _, neighbour := range g.Edges[current] {
			tentative_g_score := g_score[current] + neighbour.Dist
			old_g_score, ok := g_score[neighbour.To]
			if !ok {
				old_g_score = math.Inf(1)
			}

			if tentative_g_score < old_g_score {
				came_from[neighbour.To] = current
				g_score[neighbour.To] = tentative_g_score

				f_score[neighbour.To] = tentative_g_score + h(neighbour.To)
				if !slices.Contains(open, neighbour.To) {
					open = append(open, neighbour.To)
				}
			}
		}
	}

	return []int64{}, fmt.Errorf("No path found")
}

func idsToPoints(g *graph.Graph, ids []int64) []geo.Point {
	points := make([]geo.Point, len(ids))
	for i, id := range ids {
		points[i] = g.Nodes[id].ToPoint()
	}
	fmt.Println(len(ids), len(points)) // mismatch???
	return points
}

func ShortestPath(g *graph.Graph, start, end geo.Point) ([]geo.Point, error) {
	fmt.Println("Finding locations...")
	start_id := g.NearestTo(start)
	end_id := g.NearestTo(end)
	fmt.Println("Locations found:", start_id, end_id)

	goal := g.Nodes[end_id].ToPoint()
	h := func (id int64) float64 {
		return geo.HaversineDistance(g.Nodes[id].ToPoint(), goal)
	}

	path, err := aStar(g, start_id, end_id, h)
	if err != nil {
		return []geo.Point{}, err
	}

	return idsToPoints(g, path), nil
}
