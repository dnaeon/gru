package resource

import (
	"fmt"

	"github.com/hashicorp/hcl/hcl/ast"
)

// Resource states
const (
	ResourceStateUnknown = "unknown"
	ResourceStatePresent = "present"
	ResourceStateAbsent  = "absent"
	ResourceStateUpdate  = "update"
)

// Provider is used to create new resources from an HCL AST object item
type Provider func(item *ast.ObjectItem) (Resource, error)

// Registry contains all known resource types and their providers
var registry = make(map[string]Provider)

// Register registers a resource type and it's provider
func Register(name string, p Provider) error {
	_, ok := registry[name]
	if ok {
		return fmt.Errorf("Resource provider for '%s' is already registered", name)
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
	// Type of the resource
	Type() string

	// ID returns the unique identifier of a resource
	ID() string

	// Returns the wanted resources/dependencies
	Want() []string

	// Evaluates the resource and returns it's state
	Evaluate() (State, error)

	// Creates the resource
	Create() error

	// Deletes the resource
	Delete() error

	// Updates the resource
	Update() error
}

// BaseResource is the base resource type for all resources
// The purpose of this type is to be embedded into other resources
// Partially implements the Resource interface
type BaseResource struct {
	// Name of the resource
	Name string `json:"name"`

	// Desired state of the resource
	State string `json:"state"`

	// Type of the resource
	ResourceType string `json:"-"`

	// Resource dependencies
	WantResource []string `json:"want,omitempty" hcl:"want"`
}

// Type returns the resource type name
func (b *BaseResource) Type() string {
	return b.ResourceType
}

// ID returns the unique resource id
func (b *BaseResource) ID() string {
	return fmt.Sprintf("%s[%s]", b.ResourceType, b.Name)
}

// Want returns the wanted resources/dependencies
func (b *BaseResource) Want() []string {
	return b.WantResource
}
