package resource

import (
	"fmt"

	"github.com/dnaeon/gru/graph"
)

// Collection type is a map which keys are the
// resource ids and their values are the actual resources
type Collection map[string]Resource

// CreateCollection creates a map from
func CreateCollection(resources []Resource) (Collection, error) {
	c := make(Collection)

	for _, r := range resources {
		id := r.ID()
		if _, ok := c[id]; ok {
			return c, fmt.Errorf("Duplicate resource declaration for %s", id)
		}
		c[id] = r
	}

	return c, nil
}

// DependencyGraph builds a dependency graph for the collection
func (c Collection) DependencyGraph() (*graph.Graph, error) {
	g := graph.New()

	// A map containing the resource ids and their nodes in the graph
	nodes := make(map[string]*graph.Node)
	for id := range c {
		node := graph.NewNode(id)
		nodes[id] = node
		g.AddNode(node)
	}

	// Connect the nodes in the graph
	for id, r := range c {
		for _, dep := range r.Dependencies() {
			if _, ok := c[dep]; !ok {
				return g, fmt.Errorf("%s wants %s, which does not exist", id, dep)
			}
			g.AddEdge(nodes[id], nodes[dep])
		}
	}

	return g, nil
}

// ReversedGraph builds a reverse dependency graph for the collection
func (c Collection) ReversedGraph() (*graph.Graph, error) {
	g := graph.New()

	// A map containing the resource ids and their nodes in the graph
	nodes := make(map[string]*graph.Node)
	for id := range c {
		node := graph.NewNode(id)
		nodes[id] = node
		g.AddNode(node)
	}

	// Connect the nodes in the graph
	for id, r := range c {
		for _, dep := range r.Dependencies() {
			if _, ok := c[dep]; !ok {
				return g, fmt.Errorf("%s wants %s, which does not exist", id, dep)
			}
			g.AddEdge(nodes[dep], nodes[id])
		}
	}

	return g, nil
}
