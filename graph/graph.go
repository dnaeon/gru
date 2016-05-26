package graph

import (
	"errors"
	"fmt"
	"io"
	"strings"

	mapset "github.com/deckarep/golang-set"
)

// ErrCircularDependency is returned when the graph cannot be
// topologically sorted because of circular dependencies
var ErrCircularDependency = errors.New("Circular dependency found in graph")

// Graph represents a DAG graph
type Graph struct {
	Nodes map[string]*Node
}

// New creates a new DAG graph
func New() *Graph {
	g := &Graph{
		Nodes: make(map[string]*Node),
	}

	return g
}

// AddNode adds nodes to the graph
func (g *Graph) AddNode(nodes ...*Node) {
	for _, node := range nodes {
		g.Nodes[node.Name] = node
	}
}

// AddEdge connects a node with other nodes in the graph
func (g *Graph) AddEdge(node *Node, edges ...*Node) {
	for _, edge := range edges {
		node.Edges = append(node.Edges, edge)
	}
}

// GetNode retrieves the node from the graph with the given name
func (g *Graph) GetNode(name string) (*Node, bool) {
	n, ok := g.Nodes[name]

	return n, ok
}

// Sort performs a topological sort of the graph
// https://en.wikipedia.org/wiki/Topological_sorting
//
// If the graph can be topologically sorted the result will
// contain the sorted nodes.
//
// If the graph cannot be sorted in case of circular dependencies,
// then the result will contain the remaining nodes from the graph,
// which are the ones causing the circular dependency.
func (g *Graph) Sort() ([]*Node, error) {
	var sorted []*Node

	// Iteratively find and remove nodes from the graph which have no edges.
	// If at some point there are still nodes in the graph and we cannot find
	// nodes without edges, that means we have a circular dependency
	for len(g.Nodes) > 0 {
		// Contains the ready nodes, which have no edges to other nodes
		//ready := make([]*Node, 0)
		ready := mapset.NewSet()

		// Find the nodes with no edges
		for _, node := range g.Nodes {
			if len(node.Edges) == 0 {
				ready.Add(node)
			}
		}

		// If there aren't any ready nodes, then we have a cicular dependency
		if ready.Cardinality() == 0 {
			// The remaining nodes in the graph are the ones causing the
			// circular dependency.
			var remaining []*Node
			for _, n := range g.Nodes {
				remaining = append(remaining, n)
			}
			return remaining, ErrCircularDependency
		}

		// Remove the ready nodes and add them to the sorted result
		for item := range ready.Iter() {
			node := item.(*Node)
			delete(g.Nodes, node.Name)
			sorted = append(sorted, node)
		}

		// Remove ready nodes from the remaining node edges as well
		for _, node := range g.Nodes {
			// Add the remaining nodes in a set
			currentEdgeSet := mapset.NewSet()
			for _, edge := range node.Edges {
				currentEdgeSet.Add(edge)
			}

			newEdgeSet := currentEdgeSet.Difference(ready)
			node.Edges = make([]*Node, 0)
			for edge := range newEdgeSet.Iter() {
				node.Edges = append(node.Edges, edge.(*Node))
			}
		}
	}

	return sorted, nil
}

// AsDot generates a DOT representation for the graph
// https://en.wikipedia.org/wiki/DOT_(graph_description_language)
func (g *Graph) AsDot(name string, w io.Writer) {
	w.Write([]byte(fmt.Sprintf("digraph %s {\n", name)))
	w.Write([]byte(fmt.Sprintf("\tlabel = %q;\n", name)))
	w.Write([]byte("\tnodesep=1.0;\n"))
	w.Write([]byte("\tnode [shape=box];\n"))
	w.Write([]byte("\tedge [style=filled];\n"))

	for _, node := range g.Nodes {
		var edges []string
		for _, edge := range node.Edges {
			edges = append(edges, fmt.Sprintf("%q", edge.Name))
		}

		if len(edges) > 0 {
			w.Write([]byte(fmt.Sprintf("\t%q -> {%s};\n", node.Name, strings.Join(edges, " "))))
		} else {
			w.Write([]byte(fmt.Sprintf("\t%q;\n", node.Name)))
		}
	}

	w.Write([]byte("}\n"))
}
