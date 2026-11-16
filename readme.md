# GPath

A simple pathfinding algorithm based on OSM-data written in go.

## How to run

You will need to have [go](https://go.dev/) installed, then you can run:

```sh
go run ./cmd/cli [map file].osm.pbf [start lon] [start lat] [end lon] [end lat]
```

or

```sh
go run ./cmd/server [map file].osm.pbf
```

## TODO

- [ ] k-d tree voor punt te vinden
- [x] A* of dijkstra voor padbepaling
- [x] een server van maken
- [x] server code cleanup
- [x] meer debug voor server
