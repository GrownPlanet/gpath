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

func parsePrimitiveBlock(primitiveBlock *pbf.PrimitiveBlock, g *graph.Graph) {
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
					if okA && okB {
						d := geo.HaversineDistance(geo.NodeToPoint(a), geo.NodeToPoint(b))
						g.Edges[prev] = append(g.Edges[prev], graph.Edge{To: ref, Dist: d})
						g.Edges[ref] = append(g.Edges[ref], graph.Edge{To: prev, Dist: d})
					}
				}
				prev = ref
			}
		}
	}
}

func parseBlock(file *os.File, g *graph.Graph) error {
	blobHeader, rawData, err := getBlock(file)
	if err != nil {
		return err
	}

	headerType := blobHeader.GetType()
	if headerType == "OSMData" {
		var primitiveBlock pbf.PrimitiveBlock
		if err := proto.Unmarshal(rawData, &primitiveBlock); err != nil {
			return fmt.Errorf("failed to unmarshal PrimitiveBlockBlock: %w", err)
		}
		parsePrimitiveBlock(&primitiveBlock, g)
	} else if headerType != "OSMHeader" {
		return fmt.Errorf("unknown blobHeader type: %s", headerType)
	}

	return nil
}

func removeUnusedNodes(g *graph.Graph) {
	for id := range g.Nodes {
		if len(g.Edges[id]) == 0 {
			delete(g.Nodes, id)
		}
	}
}

func Load(filename string) (*graph.Graph, error) {
	g := graph.NewGraph()

	file, err := os.Open(filename)
	if err != nil {
		return g, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	for {
		err = parseBlock(file, g)
		if err != nil {
			if err == io.EOF {
				break
			}
			return g, err
		}
	}

	removeUnusedNodes(g)

	return g, nil
}
