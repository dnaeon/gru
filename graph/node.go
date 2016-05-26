package graph

// Node represents a single node in the graph
type Node struct {
	// Name of the node
	Name string

	// Edges to other nodes in the graph
	Edges []*Node
}

// NewNode creates a new node with the given name
func NewNode(name string) *Node {
	n := &Node{
		Name:  name,
		Edges: make([]*Node, 0),
	}

	return n
}
