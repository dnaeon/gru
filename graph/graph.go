package graph

import (
	"errors"
	"fmt"
	"io"
	"strings"

	mapset "github.com/deckarep/golang-set"
)

var ErrCircularDependency = errors.New("Circular dependency found in graph")

// Node represents a single node in the graph
type Node struct {
	// Name of the node
	Name string

	// Edges to other nodes in the graph
	edges []*Node
}

// NewNode creates a new node with the given name
func NewNode(name string) *Node {
	n := &Node{
		Name:  name,
		edges: make([]*Node, 0),
	}

	return n
}

// Graph represents a DAG graph
type Graph struct {
	nodes map[string]*Node
}

// NewGraph creates a new DAG graph
func NewGraph() *Graph {
	g := &Graph{
		nodes: make(map[string]*Node),
	}

	return g
}

// AddNode adds nodes to the graph
func (g *Graph) AddNode(nodes ...*Node) {
	for _, node := range nodes {
		g.nodes[node.Name] = node
	}
}

// AddEdge connects a node with other nodes in the graph
func (g *Graph) AddEdge(node *Node, edges ...*Node) {
	for _, edge := range edges {
		node.edges = append(node.edges, edge)
	}
}

func (g *Graph) NodeExists(name string) bool {
	_, ok := g.nodes[name]

	return ok
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
	sorted := make([]*Node, 0)

	// Iteratively find and remove nodes from the graph which have no edges.
	// If at some point there are still nodes in the graph and we cannot find
	// nodes without edges, that means we have a circular dependency
	for len(g.nodes) > 0 {
		// Contains the ready nodes, which have no edges to other nodes
		//ready := make([]*Node, 0)
		ready := mapset.NewSet()

		// Find the nodes with no edges
		for _, node := range g.nodes {
			if len(node.edges) == 0 {
				ready.Add(node)
			}
		}

		// If there aren't any ready nodes, then we have a cicular dependency
		if ready.Cardinality() == 0 {
			// The remaining nodes in the graph are the ones causing the
			// circular dependency.
			remaining := make([]*Node, 0)
			for _, n := range g.nodes {
				remaining = append(remaining, n)
			}
			return remaining, ErrCircularDependency
		}

		// Remove the ready nodes and add them to the sorted result
		for item := range ready.Iter() {
			node := item.(*Node)
			delete(g.nodes, node.Name)
			sorted = append(sorted, node)
		}

		// Remove ready nodes from the remaining node edges as well
		for _, node := range g.nodes {
			// Add the remaining nodes in a set
			currentEdgeSet := mapset.NewSet()
			for _, edge := range node.edges {
				currentEdgeSet.Add(edge)
			}

			newEdgeSet := currentEdgeSet.Difference(ready)
			node.edges = make([]*Node, 0)
			for edge := range newEdgeSet.Iter() {
				node.edges = append(node.edges, edge.(*Node))
			}
		}
	}

	return sorted, nil
}

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
