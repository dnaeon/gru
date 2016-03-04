package graph

import "errors"

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

// Sort performs a topological sorting of the graph
// https://en.wikipedia.org/wiki/Topological_sorting
func (g *Graph) Sort() ([]*Node, error) {
	sorted := make([]*Node, 0)

	// Iteratively find and remove nodes from the graph which have no edges.
	// If at some point there are still nodes in the graph and we cannot find
	// nodes without edges, that means we have a circular dependency
	for len(g.nodes) > 0 {
		// Contains the ready nodes, which have no edges to other nodes
		ready := make([]*Node, 0)

		// Find the nodes with no edges
		for _, node := range g.nodes {
			if len(node.edges) == 0 {
				ready = append(ready, node)
			}
		}

		// If there aren't any ready nodes, then we have a cicular dependency
		if len(ready) == 0 {
			return nil, errors.New("Circular dependency found in graph")
		}

		// Remove the ready nodes and add them to the sorted result
		for _, node := range ready {
			delete(g.nodes, node.Name)
			sorted = append(sorted, node)
		}

		// Remove ready nodes from the remaining node edges as well
		for _, node := range g.nodes {
			// Remove the ready nodes from any remaining edges
			newEdges := make([]*Node, len(node.edges))
			copy(newEdges, node.edges)

			for i, edge := range node.edges {
				for _, r := range ready {
					if edge.Name == r.Name {
						newEdges = append(newEdges[:i], newEdges[i+1:]...)
					}
				}
			}
			node.edges = newEdges
		}
	}

	return sorted, nil
}
