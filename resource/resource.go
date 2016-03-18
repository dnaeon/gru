package resource

import (
	"fmt"
	"io"

	"github.com/hashicorp/hcl/hcl/ast"
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

// Resource is the base interface type for all resources
type Resource interface {
	// Type of the resource
	Type() string

	// ID returns the unique identifier of a resource
	ID() string

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

// Type returns the resource type name
func (b *BaseResource) Type() string {
	return b.ResourceType
}

// ID returns the unique resource id
func (b *BaseResource) ID() string {
	return fmt.Sprintf("%s[%s]", b.ResourceType, b.Name)
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
