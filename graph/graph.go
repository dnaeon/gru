// Copyright (c) 2015-2017 Marin Atanasov Nikolov <dnaeon@gmail.com>
// All rights reserved.
//
// Redistribution and use in source and binary forms, with or without
// modification, are permitted provided that the following conditions
// are met:
//
//  1. Redistributions of source code must retain the above copyright
//     notice, this list of conditions and the following disclaimer
//     in this position and unchanged.
//  2. Redistributions in binary form must reproduce the above copyright
//     notice, this list of conditions and the following disclaimer in the
//     documentation and/or other materials provided with the distribution.
//
// THIS SOFTWARE IS PROVIDED BY THE AUTHOR(S) ``AS IS'' AND ANY EXPRESS OR
// IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED WARRANTIES
// OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE DISCLAIMED.
// IN NO EVENT SHALL THE AUTHOR(S) BE LIABLE FOR ANY DIRECT, INDIRECT,
// INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT
// NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE,
// DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY
// THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
// (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF
// THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.

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
	fmt.Fprintf(w, "digraph %s {\n", name)
	fmt.Fprintf(w, "\tlabel = %q;\n", name)
	fmt.Fprintf(w, "\tnodesep=1.0;\n")
	fmt.Fprintf(w, "\tnode [shape=box];\n")
	fmt.Fprintf(w, "\tedge [style=filled];\n")

	for _, node := range g.Nodes {
		var edges []string
		for _, edge := range node.Edges {
			edges = append(edges, fmt.Sprintf("%q", edge.Name))
		}

		if len(edges) > 0 {
			fmt.Fprintf(w, "\t%q -> {%s};\n", node.Name, strings.Join(edges, " "))
		} else {
			fmt.Fprintf(w, "\t%q;\n", node.Name)
		}
	}

	fmt.Fprintf(w, "}\n")
}

// Reversed creates the reversed representation of the graph
func (g *Graph) Reversed() *Graph {
	reversed := New()

	// Create a map of the graph nodes
	nodes := make(map[string]*Node)
	for _, n := range g.Nodes {
		node := NewNode(n.Name)
		nodes[n.Name] = node
		reversed.AddNode(node)
	}

	// Connect the nodes in the graph
	for _, node := range g.Nodes {
		for _, edge := range node.Edges {
			reversed.AddEdge(nodes[edge.Name], nodes[node.Name])
		}
	}

	return reversed
}
