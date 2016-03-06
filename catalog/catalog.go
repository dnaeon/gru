package catalog

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"

	"github.com/dnaeon/gru/graph"
	"github.com/dnaeon/gru/resource"
	"github.com/hashicorp/hcl"
	"github.com/hashicorp/hcl/hcl/ast"
)

// Catalog type contains resources loaded from a given HCL input
type Catalog struct {
	resources map[string]resource.Resource
}

// newCatalog creates a new empty catalog
func newCatalog() *Catalog {
	c := &Catalog{
		resources: make(map[string]resource.Resource),
	}

	return c
}

// addResource adds a new resource to the catalog
func (c *Catalog) addResource(r resource.Resource) error {
	id := r.ID()

	if c.resourceExists(id) {
		return fmt.Errorf("Resource '%s' is already declared", id)
	}

	c.resources[id] = r

	return nil
}

// resourceExists returns true if the resource id already exists in the catalog
// Otherwise it returns false
func (c *Catalog) resourceExists(id string) bool {
	_, ok := c.resources[id]

	return ok
}

// Graph returns the sorted resources DAG graph
func (c *Catalog) sortedResourceGraph() ([]*graph.Node, error) {
	// Create a DAG graph of the resources and perform
	// topological sorting of the graph to determine the
	// order of processing the resources
	g := graph.NewGraph()

	// A map containing the resource ids and their nodes in the graph
	nodes := make(map[string]*graph.Node)

	// Create the graph nodes for each resource from the catalog
	for name := range c.resources {
		node := graph.NewNode(name)
		nodes[name] = node
		g.AddNode(node)
	}

	// Connect the nodes in the graph
	for name, r := range c.resources {
		deps := r.Want()
		for _, dep := range deps {
			if !c.resourceExists(dep) {
				e := fmt.Errorf("Resource '%s' wants '%s', which is not found in catalog", name, dep)
				return nil, e
			}
			g.AddEdge(nodes[name], nodes[dep])
		}
	}

	// Perform topological sort of the graph
	sorted, err := g.Sort()
	if err != nil {
		return nil, err
	}

	return sorted, nil
}

// Run processes the resources from the catalog
func (c *Catalog) Run() error {
	// Perform topological sort of the graph
	sorted, err := c.sortedResourceGraph()
	if err != nil {
		return err
	}

	for _, node := range sorted {
		r := c.resources[node.Name]
		id := r.ID()
		log.Printf("Evaluating resource '%s'", id)
		state, err := r.Evaluate()
		if err != nil {
			log.Printf("Failed to evaluate resource '%s': %s", id, err)
			continue
		}

		if state.Want == state.Current {
			continue
		}

		log.Printf("%s is %s, should be %s", id, state.Current, state.Want)
		switch {
		case state.Want == resource.Present && state.Current == resource.Absent:
			r.Create()
		case state.Want == resource.Absent && state.Current != resource.Absent:
			r.Delete()
		case state.Want == resource.Update && state.Current == resource.Present:
			r.Update()
		default:
			// TODO: Validate resource states before evaluation them
			log.Printf("Unknown state '%s' for resource '%s'", state.Want, id)
		}
	}

	return nil
}

// Load reads a catalog from the given input and creates resources
func Load(path string) (*Catalog, error) {
	c := NewCatalog()

	input, err := ioutil.ReadFile(path)
	if err != nil {
		return c, err
	}

	// Parse configuration
	obj, err := hcl.Parse(string(input))
	if err != nil {
		return c, err
	}

	// Top-level node should be an object list
	root, ok := obj.Node.(*ast.ObjectList)
	if !ok {
		return c, errors.New("Missing root node")
	}

	// Get the resource declarations and create the actual resources
	resources := root.Filter("resource")
	for _, item := range resources.Items {
		position := item.Val.Pos().String()

		// The item is expected to have exactly one key,
		// which represents the resource type
		if len(item.Keys) != 1 {
			e := fmt.Errorf("Invalid resource declaration found at %s", position)
			return c, e
		}

		// Get the resource type and create the actual resource
		resourceType := item.Keys[0].Token.Value().(string)
		provider, ok := resource.Get(resourceType)
		if !ok {
			e := fmt.Errorf("Unknown resource type '%s' found at %s", resourceType, position)
			return c, e
		}

		// Create the actual resource
		r, err := provider(item)
		if err != nil {
			return c, err
		}

		err = c.addResource(r)
		if err != nil {
			return c, err
		}
	}

	return c, nil
}
