package graph

import (
	"fmt"
	"io"
	"strings"
)

// GenerateDOT generates a DOT file for the graph
// https://en.wikipedia.org/wiki/DOT_(graph_description_language)
func (g *Graph) GenerateDOT(name string, w io.Writer) {
	dotNodes := make([]string, 0)
	dotHeader := []byte(fmt.Sprintf("digraph %s {\n", name))
	dotFooter := []byte("\n}\n")

	for _, node := range g.nodes {
		// Insert the node as the first item in the DOT object
		obj := make([]string, 0)
		obj = append(obj, fmt.Sprintf("\t%q", node.Name))

		// Now add the node edges as well
		for _, edge := range node.edges {
			obj = append(obj, fmt.Sprintf("%q", edge.Name))
		}

		dotObj := strings.Join(obj, " -> ")
		dotNodes = append(dotNodes, dotObj)
	}

	w.Write(dotHeader)
	w.Write([]byte(strings.Join(dotNodes, "\n")))
	w.Write(dotFooter)
}
