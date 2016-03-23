package graph

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
