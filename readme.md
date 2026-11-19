# GPath

A simple pathfinding algorithm based A* with OSM-data written in go.

## How to run

To run this project, you need to have [go](https://go.dev/) installed, then you can run:

```sh
go run ./cmd/cli [map file].osm.pbf [start lon] [start lat] [end lon] [end lat]
```

or

```sh
go run ./cmd/server [map file].osm.pbf
```

You can call the server using curl

```sh
curl -s "http://localhost:8080/3.7251,51.0538;3.72911857634063,51.05529475737477"
```
