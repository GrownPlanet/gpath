package main

import (
	"fmt"
	"log"
	"net/http"
	"encoding/json"
	"os"
	"strconv"
	"strings"

	"gpath/internal/geo"
	"gpath/internal/pathfinding"
	"gpath/internal/graph"
	"gpath/internal/osm"
)

type RouteResponse struct {
	Code   string  `json:"code"`
	Routes []Route `json:"routes"`
}

type Route struct {
	Geometry   Geometry `json:"geometry"`
	Distance   float64  `json:"distance"`
	Duration   float64  `json:"duration"`
}

type Geometry struct {
	Coordinates [][2]float64 `json:"coordinates"`
}

func coordinateHandler(w http.ResponseWriter, r *http.Request, g *graph.Graph) {
	path := strings.TrimPrefix(r.URL.Path, "/")
	if path == "" {
		http.Error(w, "Please provide coordinates in the format: /lat,lon;lat,lon", http.StatusBadRequest)
		return
	}

	coordinates := strings.Split(path, ";")
	if len(coordinates) != 2 {
		http.Error(w, "Invalid format. Expected two coordinate pairs separated by ';'", http.StatusBadRequest)
		return
	}

	start_strings := strings.Split(coordinates[0], ",")
	end_strings := strings.Split(coordinates[1], ",")

	if len(start_strings) != 2 || len(end_strings) != 2 {
		http.Error(w, "Invalid coordinate values. Expected lat,lon", http.StatusBadRequest)
		return
	}

	startLon, err2 := strconv.ParseFloat(start_strings[0], 64)
	startLat, err1 := strconv.ParseFloat(start_strings[1], 64)
	endLon, err4 := strconv.ParseFloat(end_strings[0], 64)
	endLat, err3 := strconv.ParseFloat(end_strings[1], 64)
	if err1 != nil || err2 != nil || err3 != nil || err4 != nil {
		http.Error(w, "Invalid coordinate values. Expected floating point numbers", http.StatusBadRequest)
		return
	}

	start := geo.Point{Lat: startLat, Lon: startLon}
	end := geo.Point{Lat: endLat, Lon: endLon}

	shortestPath, err := pathfinding.ShortestPath(g, start, end)
	if err != nil {
        fmt.Fprintf(os.Stderr, "error: %s\n", err)
		http.Error(w, "Failed to find path", http.StatusBadRequest)
		return
	}

	response := RouteResponse{
		Code: "Ok",
		Routes: []Route{
			{
				Geometry: Geometry{ Coordinates: shortestPath },
				Distance: 0,
				Duration: 0,
			},
		},
	}

	// Tell the client we are returning JSON
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	// Encode the Go struct directly into the HTTP response stream
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func main() {
	fmt.Println("Loading map...")

	g, err := osm.Load("gent.osm.pbf")
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %s\n", err)
		os.Exit(1)
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		coordinateHandler(w, r, g)
	})

	fmt.Println("Server is running on http://localhost:8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
