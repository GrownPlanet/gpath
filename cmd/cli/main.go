package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"gpath/internal/geo"
	"gpath/internal/graph"
	"gpath/internal/osm"
	"gpath/internal/pathfinding"
)

func main() {
	if len(os.Args) != 6 {
		log.Fatalf(
			"Usage: %s [map file] [start lon] [start lat] [end lon] [end lat]",
			os.Args[0],
		)
	}

	startLon, err2 := strconv.ParseFloat(os.Args[2], 64)
	startLat, err1 := strconv.ParseFloat(os.Args[3], 64)
	endLon, err4 := strconv.ParseFloat(os.Args[4], 64)
	endLat, err3 := strconv.ParseFloat(os.Args[5], 64)

	if err1 != nil || err2 != nil || err3 != nil || err4 != nil {
		log.Fatalf("Invalid coordinate values. Expected floating point numbers")
	}

	fmt.Println("Loading map...")
	timer := time.Now()

	g, err := osm.Load(os.Args[1])
	if err != nil {
		log.Fatalf("error: %s\n", err)
	}

	fmt.Println("Loaded map in", time.Since(timer).Milliseconds(), "ms")

	fmt.Println("Building k-d tree...")
	timer = time.Now()

	kdTree := graph.BuildKDTree(g)

	fmt.Println("Built tree in", time.Since(timer).Milliseconds(), "ms")

	start := geo.Point{Lon: startLon, Lat: startLat}
	end := geo.Point{Lon: endLon, Lat: endLat}

	fmt.Println("Finding shortest path...")
	timer = time.Now()

	path, err := pathfinding.ShortestPath(g, start, end)
	if err != nil {
		log.Fatalf("error: %s\n", err)
	}

	fmt.Println(
		"Shortest path found in",
		time.Since(timer).Milliseconds(),
		"ms;",
		path,
	)
}
