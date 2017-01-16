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
	"reflect"
	"testing"
)

func TestWorkingGraph(t *testing.T) {
	g := New()

	// Graph node names
	nodeNames := []string{
		"A",
		"B",
		"C",
		"D",
		"E",
	}

	// Map containing the node names and the actual node instance
	nodes := make(map[string]*Node)

	// Create nodes and add them to the graph
	for _, name := range nodeNames {
		n := NewNode(name)
		nodes[name] = n
		g.AddNode(n)
	}

	// Connect the nodes in the graph
	//
	// A
	// B -> C
	// C -> D
	// D -> E
	// E -> A
	//
	g.AddEdge(nodes["B"], nodes["C"])
	g.AddEdge(nodes["C"], nodes["D"])
	g.AddEdge(nodes["D"], nodes["E"])
	g.AddEdge(nodes["E"], nodes["A"])

	// Excepted result after topological sort
	wantSorted := []string{
		"A",
		"E",
		"D",
		"C",
		"B",
	}

	gotNodes, err := g.Sort()
	if err != nil {
		t.Error(err)
	}

	var gotSorted []string
	for _, node := range gotNodes {
		gotSorted = append(gotSorted, node.Name)
	}

	if !reflect.DeepEqual(wantSorted, gotSorted) {
		t.Errorf("Want %q, got %q graph", wantSorted, gotSorted)
	}
}

func TestCircularGraph(t *testing.T) {
	g := New()

	// Node names
	nodeNames := []string{
		"A",
		"B",
		"C",
		"D",
		"E",
	}

	// Map containing the node names and the actual node instance
	nodes := make(map[string]*Node)

	// Create nodes and add them to the graph
	for _, name := range nodeNames {
		n := NewNode(name)
		nodes[name] = n
		g.AddNode(n)
	}

	// Connect the nodes in the graph
	//
	// A
	// B -> C
	// C -> D
	// D -> E
	// E -> D  <- Circular dependency here
	//
	g.AddEdge(nodes["B"], nodes["C"])
	g.AddEdge(nodes["C"], nodes["D"])
	g.AddEdge(nodes["D"], nodes["E"])
	g.AddEdge(nodes["E"], nodes["D"])

	_, err := g.Sort()
	if err != ErrCircularDependency {
		t.Errorf("want a circular dependency error, got %s", err)
	}
}
