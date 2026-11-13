package main

import (
	"fmt"
	"os"
)

func main() {
	if err := Parse("gent.osm.pbf"); err != nil {
		fmt.Fprintf(os.Stderr, "error: %s\n", err)
		os.Exit(1)
	}
	fmt.Println("Done parsing map file")
}
