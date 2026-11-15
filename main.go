package main

import (
	"fmt"
	"os"
)

func main() {
	graph, err := Parse("gent.osm.pbf")
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %s\n", err)
		os.Exit(1)
	}
	fmt.Println("Turned map into graph")
	fmt.Println(len(graph.Nodes), graph.nearestTo(Node{Lat: 50.951123900000005, Lon: 3.4592488000000003}))
}
