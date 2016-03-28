package module

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"

	"github.com/dnaeon/gru/resource"
	"github.com/hashicorp/hcl"
	"github.com/hashicorp/hcl/hcl/ast"
)

// ErrMultipleImport error is returned if there are multiple import
// declarations in the same module
var ErrMultipleImport = errors.New("Multiple import declarations found")

// Import type represents a single import declaration from HCL/JSON
type ImportType struct {
	// Module names being imported
	Module []string
}

// Module type represents a collection of resources and module imports
type Module struct {
	// Name of the module
	Name string

	// Resources loaded from the module
	Resources []resource.Resource

	// Module imports
	ModuleImport ImportType
}

// New creates a new empty module
func New(name string) *Module {
	m := &Module{
		Name:      name,
		Resources: make([]resource.Resource, 0),
		ModuleImport: ImportType{
			Module: make([]string, 0),
		},
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

	// We expect to have exactly one import declaration per module file
	if len(hclImport.Items) > 1 {
		return fmt.Errorf("Multiple import declarations found in %s", m.Name)
	}

	if len(hclImport.Items) == 0 {
		return nil
	}

	err := hcl.DecodeObject(&m.ModuleImport, hclImport.Items[0])
	if err != nil {
		return err
	}

	return nil
}
