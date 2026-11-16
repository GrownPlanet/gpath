package main

import (
	"fmt"
	"os"

	"gpath/internal/geo"
	"gpath/internal/osm"
)

func main() {
	graph, err := osm.Load("gent.osm.pbf")
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %s\n", err)
		os.Exit(1)
	}

	fmt.Println("Turned map into graph")
	fmt.Println(len(graph.Nodes), geo.NearestTo(graph, 50.951123900000005, 3.4592488000000003))
}
