package main

import (
	"fmt"
	"os"

	"gpath/internal/geo"
	"gpath/internal/osm"
	"gpath/internal/pathfinding"
)

func main() {
	fmt.Println("Loading map...")

	graph, err := osm.Load("gent.osm.pbf")
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %s\n", err)
		os.Exit(1)
	}

	fmt.Println("Done loading map")

	fmt.Println("Finding shortest path...")

	start := geo.Point{Lon: 3.7159345611930115, Lat: 51.05333770584914}
	end := geo.Point{Lon: 3.735398601423926, Lat: 51.0404909339422}

	path, err := pathfinding.ShortestPath(graph, start, end)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %s\n", err)
		os.Exit(1)
	}
	fmt.Println("Shortest path found:", path)
}
