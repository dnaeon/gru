package resource

import (
	"fmt"
	"io"

	"github.com/hashicorp/hcl/hcl/ast"
)

// Registry contains all known resources
var Registry = make(map[string]RegistryItem)

// provider is used to create new resources from an HCL AST object item
type provider func(title string, item *ast.ObjectItem) (Resource, error)

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
	ResourceID() string

	// Type returns the type of the resource
	ResourceType() string

	// Returns the title of the resource
	ResourceTitle() string

	// Validates the resource
	Validate() error

	// Returns the resources before which this resource shoud be processed
	WantBefore() []string

	// Returns the resources after which this resource should be processed
	WantAfter() []string

	// Evaluates the resource
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
	// Type of the resource
	Type string `json:"-"`

	// Title of the resource
	Title string `json:"-"`

	// Name of the resource
	Name string `hcl:"name" json:"name"`

	// Desired state of the resource
	State string `hcl:"state" json:"state"`

	// Resources before which this resource should be processed
	Before []string `hcl:"before" json:"before,omitempty"`

	// Resources after which this resource should be processed
	After []string `hcl:"after" json:"after,omitempty"`
}

// ResourceID returns the unique resource id
func (b *BaseResource) ResourceID() string {
	return fmt.Sprintf("%s[%s]", b.Type, b.Title)
}

// ResourceType returns the resource type
func (b *BaseResource) ResourceType() string {
	return b.Type
}

// ResourceTitle returns the resource title
func (b *BaseResource) ResourceTitle() string {
	return b.Title
}

// Validate checks if the resource contains valid information
func (b *BaseResource) Validate() error {
	if b.Title == "" {
		return fmt.Errorf("Missing title for resource %s", b.ResourceID())
	}

	return nil
}

// WantBefore returns the resources before which this resource
// should be processed
func (b *BaseResource) WantBefore() []string {
	return b.Before
}

// WantAfter returns the resources after which this resource
// should be processed
func (b *BaseResource) WantAfter() []string {
	return b.After
}

// Printf works just like fmt.Printf except that it writes to the
// given resource writer object and prepends the
// resource id to the output
func (b *BaseResource) Printf(w io.Writer, format string, a ...interface{}) (int, error) {
	fmt.Fprintf(w, "%s ", b.ResourceID())

	return fmt.Fprintf(w, format, a...)
}
