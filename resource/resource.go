package resource

import (
	"fmt"

	"github.com/hashicorp/hcl/hcl/ast"
)

// Resource states
const (
	Present = "present"
	Absent  = "absent"
)

// Provider is used to create new resources from an HCL AST object item
type Provider func(item *ast.ObjectItem) (Resource, error)

// Registry contains all known resource types and their providers
var registry = make(map[string]Provider)

// Register registers a resource type and it's provider
func Register(name string, p Provider) error {
	_, ok := registry[name]
	if ok {
		return fmt.Errorf("Resource '%s' is already registered", name)
	}

	registry[name] = p

	return nil
}

// Get retrieves the provider for a given resource type
func Get(name string) (Provider, bool) {
	p, ok := registry[name]

	return p, ok
}

// State type represents the current and wanted states of a resource
type State struct {
	// Current state of the resource
	Current string

	// Wanted state of the resource
	Want string
}

// Resource is the base interface type for all resources
type Resource interface {
	// ID returns the unique identifier of a resource
	ID() string

	// Evaluates the resource and returns it's state
	Evaluate() (State, error)

	// Creates the resource
	Create() error

	// Deletes the resource
	Delete() error

	// Updates the resource
	Update() error
}
