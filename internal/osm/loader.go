package osm

import (
	"fmt"
	"io"
	"os"
	"slices"

	"gpath/internal/geo"
	"gpath/internal/graph"
	"gpath/internal/pbf"

	"google.golang.org/protobuf/proto"
)

// TODO: calculate dist in toGraph, so the data can be out of order and
// still parse correctly
type tempEdge struct {
	To   int64
	Dist float64
}

type tempGraph struct {
	Nodes map[int64]graph.Node
	Edges map[int64][]tempEdge
}

func (g *tempGraph) toGraph() *graph.Graph {
	nodeIDMap := make(map[int64]int)
	var nodes []graph.Node

	for id, edges := range g.Edges {
		if len(edges) == 0 {
			continue
		}
		if node, ok := g.Nodes[id]; ok {
			nodeIDMap[id] = len(nodes)
			nodes = append(nodes, node)
		}
	}

	edges := make([][]graph.Edge, len(nodes))
	for id, tempEdgeList := range g.Edges {
		fromIdx, ok := nodeIDMap[id]
		if !ok {
			continue
		}
		for _, te := range tempEdgeList {
			toIdx, okTo := nodeIDMap[te.To]
			if okTo {
				edges[fromIdx] = append(edges[fromIdx], graph.Edge{
					To:   toIdx,
					Dist: te.Dist,
				})
			}
		}
	}

	return &graph.Graph {
		Nodes: nodes,
		Edges: edges,
	}
}

var walkable = []string{
	"footway", "path", "pedestrian", "steps", "residential", "living_street",
	"service", "track", "unclassified", "tertiary", "secondary", "primary",
	"cycleway", "tertiary_link", "secondary_link", "primary_link",
}

func isWalkable(tags map[string]string) bool {
	if !slices.Contains(walkable, tags["highway"]) {
		return false
	}

	if tags["foot"] == "no" ||
		tags["access"] == "no" ||
		tags["access"] == "private" {
		return false
	}

	return true
}

func parsePrimitiveBlock(primitiveBlock *pbf.PrimitiveBlock, g *tempGraph) {
	gran := int64(primitiveBlock.GetGranularity())
	latOffset := primitiveBlock.GetLatOffset()
	lonOffset := primitiveBlock.GetLonOffset()
	stringTable := primitiveBlock.GetStringtable().GetS()

	for _, group := range primitiveBlock.GetPrimitivegroup() {
		if dense := group.GetDense(); dense != nil {
			var id, lat, lon int64
			for i := 0; i < len(dense.GetId()); i++ {
				id += dense.GetId()[i]
				lat += dense.GetLat()[i]
				lon += dense.GetLon()[i]

				g.Nodes[id] = graph.Node{
					Lat: 1e-9 * float64(latOffset+(int64(gran)*lat)),
					Lon: 1e-9 * float64(lonOffset+(int64(gran)*lon)),
				}
			}
		}

		for _, way := range group.GetWays() {
			tags := make(map[string]string, len(way.GetKeys()))
			for i, k := range way.GetKeys() {
				tags[string(stringTable[k])] = string(stringTable[way.GetVals()[i]])
			}
			if !isWalkable(tags) {
				continue
			}

			var ref, prev int64
			for i, delta := range way.GetRefs() {
				ref += delta
				if i > 0 {
					a, okA := g.Nodes[prev]
					b, okB := g.Nodes[ref]
					if !(okA && okB) {
						fmt.Fprintf(os.Stderr, "error: nodes %d and %d do not exist (yet), skipping", prev, ref)
						continue
					}
					d := geo.HaversineDistance(a.ToPoint(), b.ToPoint())
					g.Edges[prev] = append(g.Edges[prev], tempEdge{
						To:   ref,
						Dist: d,
					})
					g.Edges[ref] = append(g.Edges[ref], tempEdge{
						To:   prev,
						Dist: d,
					})
				}
				prev = ref
			}
		}
	}
}

func parseBlock(file *os.File, g *tempGraph) error {
	blobHeader, rawData, err := getBlock(file)
	if err != nil {
		return err
	}

	headerType := blobHeader.GetType()
	if headerType == "OSMData" {
		var primitiveBlock pbf.PrimitiveBlock
		if err := proto.Unmarshal(rawData, &primitiveBlock); err != nil {
			return fmt.Errorf("failed to unmarshal PrimitiveBlock: %w", err)
		}
		parsePrimitiveBlock(&primitiveBlock, g)
	} else if headerType != "OSMHeader" {
		return fmt.Errorf("unknown blobHeader type: %s", headerType)
	}

	return nil
}

func Load(filename string) (*graph.Graph, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	g := tempGraph{
		Nodes: make(map[int64]graph.Node),
		Edges: make(map[int64][]tempEdge),
	}

	for {
		err = parseBlock(file, &g)
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
	}

	return g.toGraph(), nil
}
