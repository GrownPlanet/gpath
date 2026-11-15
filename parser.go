package main

import (
	"bytes"
	"compress/zlib"
	"encoding/binary"
	"fmt"
	"io"
	"os"
	"slices"

	"gpath/pbf"

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

func parsePrimitiveBlock(primitiveBlock *pbf.PrimitiveBlock, graph *Graph) {
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

				graph.Nodes[id] = Node{
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
					a, okA := graph.Nodes[prev]
					b, okB := graph.Nodes[ref]
					if okA && okB {
						d := haversineDistance(a, b)
						graph.Edges[prev] = append(graph.Edges[prev], Edge{ref, d})
						graph.Edges[ref] = append(graph.Edges[ref], Edge{prev, d})
					}
				}
				prev = ref
			}
		}
	}
}

func getBlobHeader(file *os.File, headerLen uint32) (*pbf.BlobHeader, error) {
	blobHeader := &pbf.BlobHeader{}

	blobHeaderBuf := make([]byte, headerLen)
	if _, err := io.ReadFull(file, blobHeaderBuf); err != nil {
		return blobHeader, fmt.Errorf("failed to read BlobHeader bytes: %w", err)
	}

	if err := proto.Unmarshal(blobHeaderBuf, blobHeader); err != nil {
		return blobHeader, fmt.Errorf("failed to unmarshal BlobHeader: %w", err)
	}

	return blobHeader, nil
}

func getRawData(file *os.File, blobHeader *pbf.BlobHeader) ([]byte, error) {
	var rawData []byte

	blobBuf := make([]byte, blobHeader.GetDatasize())
	if _, err := io.ReadFull(file, blobBuf); err != nil {
		return rawData, fmt.Errorf("failed to read Blob bytes: %w", err)
	}

	var blob pbf.Blob
	if err := proto.Unmarshal(blobBuf, &blob); err != nil {
		return rawData, fmt.Errorf("failed to unmarshal Blob: %w", err)
	}

	if zlibData := blob.GetZlibData(); len(zlibData) > 0 {
		r, err := zlib.NewReader(bytes.NewReader(zlibData))
		if err != nil {
			return rawData, fmt.Errorf("failed to create zlib reader: %w", err)
		}
		defer r.Close()
		rawData, err = io.ReadAll(r)
		if err != nil {
			return rawData, fmt.Errorf("failed to decompress zlib data: %w", err)
		}
	} else if raw := blob.GetRaw(); len(raw) > 0 {
		rawData = raw
	} else {
		return rawData, fmt.Errorf("unsupported or empty blob data payload")
	}

	return rawData, nil
}

func parseBlock(file *os.File, graph *Graph) error {
	var headerLen uint32
	if err := binary.Read(file, binary.BigEndian, &headerLen); err != nil {
		if err == io.EOF {
			return io.EOF
		}
		return fmt.Errorf("failed to read BlobHeader length: %w", err)
	}

	blobHeader, err := getBlobHeader(file, headerLen)
	if err != nil {
		return err
	}

	rawData, err := getRawData(file, blobHeader)
	if err != nil {
		return err
	}

	headerType := blobHeader.GetType()
	if headerType == "OSMData" {
		var primitiveBlock pbf.PrimitiveBlock
		if err := proto.Unmarshal(rawData, &primitiveBlock); err != nil {
			return fmt.Errorf("failed to unmarshal PrimitiveBlockBlock: %w", err)
		}
		parsePrimitiveBlock(&primitiveBlock, graph)
	} else if headerType != "OSMHeader" {
		return fmt.Errorf("unknown blobHeader type: %s", headerType)
	}

	return nil
}

func Parse(filename string) (*Graph, error) {
	graph := NewGraph()

	file, err := os.Open(filename)
	if err != nil {
		return graph, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	for {
		err = parseBlock(file, graph)
		if err != nil {
			if err == io.EOF {
				break
			}
			return graph, err
		}
	}

	return graph, nil
}
