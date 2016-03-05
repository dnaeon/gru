package graph

import (
	"reflect"
	"testing"
)

func TestWorkingGraph(t *testing.T) {
	g := NewGraph()

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

	gotSorted := make([]string, 0)
	for _, node := range gotNodes {
		gotSorted = append(gotSorted, node.Name)
	}

	if !reflect.DeepEqual(wantSorted, gotSorted) {
		t.Errorf("Want %q, got %q graph", wantSorted, gotSorted)
	}
}

func TestCircularGraph(t *testing.T) {
	g := NewGraph()

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
