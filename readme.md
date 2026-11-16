# GPath

A simple pathfinding algorithm based on OSM-data written in go.

## How to run

You will need to have [go](https://go.dev/) installed, then you can run:

```sh
go run .
```

## ToDo

- [ ] k-d tree voor punt te vinden
- [ ] A* of dijkstra voor padbepaling

## File structure

pathfinder/
├── cmd/
│   └── server/
│       └── main.go
│
├── internal/
│   ├── server/
│   │   ├── server.go
│   │   └── handler.go
│   │
│   ├── osm/
│   │   ├── reader.go
│   │   └── parser.go
│   │
│   ├── graph/
│   │   ├── graph.go
│   │   └── node.go
│   │
│   ├── pathfinding/
│   │   ├── astar.go
│   │   └── dijkstra.go
│   │
│   └── geo/
│       ├── haversine.go
│       └── kdtree.go
│
├── proto/
│   └── pathfinder.proto
│
├── go.mod
└── go.sum
