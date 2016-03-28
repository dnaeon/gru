package resource

import (
	"fmt"
	"io"

	"github.com/hashicorp/hcl/hcl/ast"
)

// Registry contains all known resources
var Registry = make(map[string]RegistryItem)

// provider is used to create new resources from an HCL AST object item
type provider func(name string, item *ast.ObjectItem) (Resource, error)

// RegistryItem type represents an item from the registry
type RegistryItem struct {
	// Name of the resource type
	Name string

	// Short desription of the resource
	Description string

	// Resource provider
	Provider provider
}

// Register adds a resource type to the registry
func Register(item RegistryItem) error {
	_, ok := Registry[item.Name]
	if ok {
		return fmt.Errorf("Resource type '%s' is already registered", item.Name)
	}

	Registry[item.Name] = item

	return nil
}

// Resource is the interface type for resources
type Resource interface {
	// ID returns the unique identifier of a resource
	ID() string

	// Type returns the type of the resource
	Type() string

	// Returns the name of the resource
	ResourceName() string

	// Validates the resource
	Validate() error

	// Returns the wanted resources/dependencies
	Want() []string

	// Evaluates the resource and returns it's state
	Evaluate() (State, error)

	// Creates the resource
	Create(w io.Writer) error

	// Deletes the resource
	Delete(w io.Writer) error

	// Updates the resource
	Update(w io.Writer) error
}

// BaseResource is the base resource type for all resources
// The purpose of this type is to be embedded into other resources
// Partially implements the Resource interface
type BaseResource struct {
	// Name of the resource
	Name string `hcl:"name" json:"name"`

	// Desired state of the resource
	State string `hcl:"state" json:"state"`

	// Type of the resource
	ResourceType string `json:"-"`

	// Resource dependencies
	WantResource []string `hcl:"want" json:"want,omitempty"`
}

// ID returns the unique resource id
func (b *BaseResource) ID() string {
	return fmt.Sprintf("%s[%s]", b.ResourceType, b.Name)
}

// Type returns the resource type name
func (b *BaseResource) Type() string {
	return b.ResourceType
}

// ResourceName returns the resource name
func (b *BaseResource) ResourceName() string {
	return b.Name
}

// Validate checks if the resource contains valid information
func (b *BaseResource) Validate() error {
	if b.Name == "" {
		return fmt.Errorf("Missing name for resource %s", b.ID())
	}

	return nil
}

// Want returns the wanted resources/dependencies
func (b *BaseResource) Want() []string {
	return b.WantResource
}

// Printf works just like fmt.Printf except that it writes to the
// given resource writer object and prepends the
// resource id to the output
func (b *BaseResource) Printf(w io.Writer, format string, a ...interface{}) (int, error) {
	fmt.Fprintf(w, "%s ", b.ID())

	return fmt.Fprintf(w, format, a...)
}
