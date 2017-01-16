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
		// Create edges between the nodes and the ones
		// required by it
		for _, dep := range r.Dependencies() {
			if _, ok := c[dep]; !ok {
				return g, fmt.Errorf("%s wants %s, which does not exist", id, dep)
			}
			g.AddEdge(nodes[id], nodes[dep])
		}

		// Create edges between the nodes and the resources for
		// which we subscribe for changes to
		for dep := range r.SubscribedTo() {
			if _, ok := c[dep]; !ok {
				return g, fmt.Errorf("%s subscribes to %s, which does not exist", id, dep)
			}
			g.AddEdge(nodes[id], nodes[dep])
		}
	}

	return g, nil
}
