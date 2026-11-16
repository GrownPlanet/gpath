package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"gpath/internal/geo"
)

type RouteResponse struct {
	Code   string  `json:"code"`
	Routes []Route `json:"routes"`
}

type Route struct {
	Geometry Geometry `json:"geometry"`
	Distance float64  `json:"distance"`
	Duration float64  `json:"duration"`
}

type Geometry struct {
	Coordinates [][2]float64 `json:"coordinates"`
}

func ParseCoordinates(path string) (start, end geo.Point, err error) {
	path = strings.TrimPrefix(path, "/")
	if path == "" {
		return start, end, fmt.Errorf("Please provide coordinates in the format: /lat,lon;lat,lon")
	}

	coordinates := strings.Split(path, ";")
	if len(coordinates) != 2 {
		return start, end, fmt.Errorf("Invalid format. Expected two coordinate pairs separated by ';'")
	}

	startStrings := strings.Split(coordinates[0], ",")
	endStrings := strings.Split(coordinates[1], ",")

	if len(startStrings) != 2 || len(endStrings) != 2 {
		return start, end, fmt.Errorf("Invalid coordinate values. Expected lat,lon")
	}

	startLon, err2 := strconv.ParseFloat(startStrings[0], 64)
	startLat, err1 := strconv.ParseFloat(startStrings[1], 64)
	endLon, err4 := strconv.ParseFloat(endStrings[0], 64)
	endLat, err3 := strconv.ParseFloat(endStrings[1], 64)
	if err1 != nil || err2 != nil || err3 != nil || err4 != nil {
		return start, end, fmt.Errorf("Invalid coordinate values. Expected floating point numbers")
	}

	start = geo.Point{Lon: startLon, Lat: startLat}
	end = geo.Point{Lon: endLon, Lat: endLat}

	return start, end, nil
}

func FlattenPoints(points []geo.Point) [][2]float64 {
	flatPoints := make([][2]float64, len(points))
	for i, p := range points {
		flatPoints[i] = [2]float64{p.Lat, p.Lon}
	}
	return flatPoints
}

func SendResponse(
	w http.ResponseWriter,
	path [][2]float64,
	distance, duration float64,
) {
	response := RouteResponse{
		Code: "Ok",
		Routes: []Route{
			{
				Geometry: Geometry{Coordinates: path},
				Distance: distance,
				Duration: duration,
			},
		},
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
