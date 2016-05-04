package module

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"

	"github.com/dnaeon/gru/graph"
	"github.com/dnaeon/gru/resource"
	"github.com/hashicorp/hcl"
	"github.com/hashicorp/hcl/hcl/ast"
)

// ErrMultipleImport error is returned if there are multiple import
// declarations in the same module
var ErrMultipleImport = errors.New("Multiple import declarations found")

// ValidKeys contains a map of valid keys that can be used in modules
var ValidKeys = validKeys()

// Module type represents a collection of resources and module imports
type Module struct {
	// Name of the module
	Name string

	// Resources loaded from the module
	Resources []resource.Resource

	// Module imports
	Imports []Import

	// Unknown keys found in the module
	UnknownKeys []string
}

// Import type represents an import declaration
type Import struct {
	// Name of the module that is imported
	Name string `hcl:"name"`

	// Path to the module file
	Path string `hcl:"path"`
}

// validKeys returns a map of valid keys which can be used in modules
func validKeys() map[string]struct{} {
	// All resource types found in resource.Registry are considered
	// valid keys to be used in modules.
	keys := make(map[string]struct{})

	for name := range resource.Registry {
		keys[name] = struct{}{}
	}

	// Others keys considered as valid
	keys["import"] = struct{}{}

	return keys
}

// New creates a new empty module
func New(name string) *Module {
	m := &Module{
		Name:      name,
		Resources: make([]resource.Resource, 0),
		Imports:   make([]Import, 0),
	}

	return m
}

// Load loads a module from the given HCL or JSON input
func Load(name string, r io.Reader) (*Module, error) {
	m := New(name)

	input, err := ioutil.ReadAll(r)
	if err != nil {
		return m, err
	}

	// Parse configuration
	obj, err := hcl.Parse(string(input))
	if err != nil {
		return m, err
	}

	// Top-level node should be an object list
	root, ok := obj.Node.(*ast.ObjectList)
	if !ok {
		return m, fmt.Errorf("Missing root node in %s", name)
	}

	err = m.hclLoadImport(root)
	if err != nil {
		return m, err
	}

	// Load all known resource types from the given input
	for name := range resource.Registry {
		err = m.hclLoadResources(name, root)
		if err != nil {
			return m, err
		}
	}

	// Check for unknown keys in the provided input
	//
	// For now the only valid keys are the resource types,
	// which can be found in resource.Registry.
	for _, item := range root.Items {
		key := item.Keys[0].Token.Value().(string)
		if _, ok := ValidKeys[key]; !ok {
			m.UnknownKeys = append(m.UnknownKeys, key)
		}
	}

	return m, nil
}

// hclLoadResources loads all declarations with the
// given resource type from the provided HCL input
func (m *Module) hclLoadResources(resourceType string, root *ast.ObjectList) error {
	hclResources := root.Filter(resourceType)
	for _, item := range hclResources.Items {
		position := item.Val.Pos().String()

		// The item is expected to exactly one key which
		// represent the resource name
		if len(item.Keys) != 1 {
			e := fmt.Errorf("Invalid resource declaration found in %s:%s", m.Name, position)
			return e
		}

		// Get the resource from registry and create the actual resource
		resourceName := item.Keys[0].Token.Value().(string)
		registryItem, ok := resource.Registry[resourceType]
		if !ok {
			e := fmt.Errorf("Unknown resource type '%s' found in %s:%s", resourceType, m.Name, position)
			return e
		}

		// Create the actual resource by calling it's provider
		r, err := registryItem.Provider(resourceName, item)
		if err != nil {
			return err
		}

		m.Resources = append(m.Resources, r)
	}

	return nil
}

// hclLoadImport loads all import declarations from the given HCL input
func (m *Module) hclLoadImport(root *ast.ObjectList) error {
	hclImport := root.Filter("import")

	for _, item := range hclImport.Items {
		position := item.Val.Pos().String()

		if len(item.Keys) != 0 {
			e := fmt.Errorf("Invalid module import found in %s:%s", m.Name, position)
			return e
		}

		var i Import
		err := hcl.DecodeObject(&i, item)
		if err != nil {
			return err
		}

		m.Imports = append(m.Imports, i)
	}

	return nil
}

// ImportGraph creates a DAG graph of the
// module imports for a given module.
// The resulting DAG graph can be used to determine the
// proper ordering of modules and also to detect whether
// we have circular imports in our modules.
func ImportGraph(main, path string) (*graph.Graph, error) {
	g := graph.NewGraph()

	modules, err := DiscoverAndLoad(path)
	if err != nil {
		return g, err
	}

	if _, ok := modules[main]; !ok {
		return g, fmt.Errorf("Module %s not found in module path", main)
	}

	// A map containing the modules as graph nodes
	// The graph can be used to determine if we have
	// circular module imports and also to provide the
	// proper ordering of loading modules after a
	// topological sort of the graph nodes
	nodes := make(map[string]*graph.Node)
	for n := range modules {
		node := graph.NewNode(n)
		nodes[n] = node
	}

	// Recursively find all imports that the main module has
	var buildImportGraphFunc func(m *Module) error
	buildImportGraphFunc = func(m *Module) error {
		// Add the node to the graph if it is not present already
		if _, ok := g.GetNode(m.Name); !ok {
			g.AddNode(nodes[m.Name])
		} else {
			return nil
		}

		// Build the import graph for each imported module
		for _, mi := range m.Imports {
			if _, ok := modules[mi.Name]; !ok {
				return fmt.Errorf("Module %s imports %s, which is not in the module path", m.Name, mi.Name)
			}

			// Build the dependencies of imported modules as well
			buildImportGraphFunc(modules[mi.Name])

			// Finally connect the nodes in the graph
			g.AddEdge(nodes[m.Name], nodes[mi.Name])
		}

		return nil
	}

	if err := buildImportGraphFunc(modules[main]); err != nil {
		return g, err
	}

	return g, nil
}
