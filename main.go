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
	graph.Print()
	fmt.Println("Turned map into graph")
}
