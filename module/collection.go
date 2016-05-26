package module

import (
	"fmt"

	"github.com/dnaeon/gru/graph"
	"github.com/dnaeon/gru/resource"
)

// ResourceMap type is a map which keys are the
// resource ids and their values are the actual resources
type ResourceMap map[string]resource.Resource

// ResourceCollection creates a map of the unique resources
// contained within the provided modules.
func ResourceCollection(modules []*Module) (ResourceMap, error) {
	moduleMap := make(map[string]string)
	resourceMap := make(ResourceMap)

	for _, m := range modules {
		for _, r := range m.Resources {
			id := r.ResourceID()
			if _, ok := resourceMap[id]; ok {
				return resourceMap, fmt.Errorf("Duplicate resource %s in %s, previous declaration was in %s", id, m.Name, moduleMap[id])
			}
			moduleMap[id] = m.Name
			resourceMap[id] = r
		}
	}

	return resourceMap, nil
}

// DependencyGraph builds a dependency graph for the resource collection
func (rm ResourceMap) DependencyGraph() (*graph.Graph, error) {
	g := graph.New()

	// A map containing the resource ids and their nodes in the graph
	nodes := make(map[string]*graph.Node)
	for name := range rm {
		node := graph.NewNode(name)
		nodes[name] = node
		g.AddNode(node)
	}

	// Connect the nodes in the graph
	for name, r := range rm {
		before := r.WantBefore()
		after := r.WantAfter()

		// Connect current resource with the ones that happen after it
		for _, dep := range after {
			if _, ok := rm[dep]; !ok {
				e := fmt.Errorf("Resource %s wants %s, which does not exist", name, dep)
				return g, e
			}
			g.AddEdge(nodes[name], nodes[dep])
		}

		// Connect current resource with the ones that happen before it
		for _, dep := range before {
			if _, ok := rm[dep]; !ok {
				e := fmt.Errorf("Resource %s wants %s, which does not exist", name, dep)
				return g, e
			}
			g.AddEdge(nodes[dep], nodes[name])
		}
	}

	return g, nil
}
