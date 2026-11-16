package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"gpath/internal/geo"
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
	loadStart := time.Now()

	graph, err := osm.Load(os.Args[1])
	if err != nil {
		log.Fatalf("error: %s\n", err)
	}

	fmt.Println("Loaded map in", time.Since(loadStart).Milliseconds(), "ms")

	start := geo.Point{Lon: startLon, Lat: startLat}
	end := geo.Point{Lon: endLon, Lat: endLat}

	fmt.Println("Finding shortest path...")
	searchStart := time.Now()

	path, err := pathfinding.ShortestPath(graph, start, end)
	if err != nil {
		log.Fatalf("error: %s\n", err)
	}

	fmt.Println(
		"Shortest path found in",
		time.Since(searchStart).Milliseconds(),
		"ms;",
		path,
	)
}
