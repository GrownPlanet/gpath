package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"gpath/internal/graph"
	"gpath/internal/osm"
	"gpath/internal/pathfinding"
	"gpath/internal/server"
)

type Server struct {
	g      *graph.Graph
	kdTree *graph.KDNode
}

func (s Server) coordinateHandler(w http.ResponseWriter, r *http.Request) {
	timer := time.Now()

	start, end, err := server.ParseCoordinates(r.URL.Path)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	shortestPath, err := pathfinding.ShortestPath(s.g, s.kdTree, start, end)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %s\n", err)
		http.Error(w, "Failed to find path", http.StatusBadRequest)
		return
	}

	fmt.Println()
	fmt.Println("calculated path")
	fmt.Println(" > start:", start)
	fmt.Println(" > end:", end)
	fmt.Println(" > duration:", time.Since(timer).Milliseconds(), "ms")

	server.SendResponse(w, shortestPath, 0, 0)
}

func main() {
	if len(os.Args) < 2 {
		log.Fatalf("Usage: %s [map file]", os.Args[0])
	}

	fmt.Println("Loading map...")
	timer := time.Now()

	g, err := osm.Load(os.Args[1])
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %s\n", err)
		os.Exit(1)
	}

	fmt.Println("Loaded map in", time.Since(timer).Milliseconds(), "ms")
	fmt.Println("Building k-d tree...")
	timer = time.Now()

	kdTree := graph.BuildKDTree(g)

	fmt.Println("Built tree in", time.Since(timer).Milliseconds(), "ms")

	server := Server{g: g, kdTree: kdTree}

	http.HandleFunc("/", server.coordinateHandler)

	fmt.Println("Server is running on http://localhost:8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
