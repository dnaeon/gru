package resource

import (
	"fmt"
	"io"

	"github.com/hashicorp/hcl/hcl/ast"
)

// Registry contains all known resources
var Registry = make(map[string]RegistryItem)

// Provider type creates new resources from an HCL AST object item
type Provider func(title string, item *ast.ObjectItem, config *Config) (Resource, error)

// RegistryItem type represents an item from the registry
type RegistryItem struct {
	// Name of the resource type
	Name string

	// Short desription of the resource
	Description string

	// Resource provider
	Provider Provider
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
	// ResourceID returns the unique identifier of a resource
	ResourceID() string

	// Returns the resources before which this resource shoud be processed
	WantBefore() []string

	// Returns the resources after which this resource should be processed
	WantAfter() []string

	// Evaluates the resource
	Evaluate() (State, error)

	// Creates the resource
	Create() error

	// Deletes the resource
	Delete() error

	// Updates the resource
	Update() error
}

// Config type contains various settings used by the resources
type Config struct {
	// The site directory which contains module and data files
	SiteDir string

	// Writer used by the resources
	Writer io.Writer
}

// BaseResource is the base resource type for all resources
// The purpose of this type is to be embedded into other resources
// Partially implements the Resource interface
type BaseResource struct {
	// Type of the resource
	Type string `hcl:"-"`

	// Title of the resource
	Title string `hcl:"-"`

	// Resource configuration settings
	Config *Config `hcl:"-"`

	// Desired state of the resource
	State string `hcl:"state"`

	// Resources before which this resource should be processed
	Before []string `hcl:"before"`

	// Resources after which this resource should be processed
	After []string `hcl:"after"`
}

// ResourceID returns the unique resource id
func (br *BaseResource) ResourceID() string {
	return fmt.Sprintf("%s[%s]", br.Type, br.Title)
}

// WantBefore returns the resources before which this resource
// should be processed
func (br *BaseResource) WantBefore() []string {
	return br.Before
}

// WantAfter returns the resources after which this resource
// should be processed
func (br *BaseResource) WantAfter() []string {
	return br.After
}

// Printf works just like fmt.Printf except that it writes to the
// given resource writer object and prepends the
// resource id to the output
func (br *BaseResource) Printf(format string, a ...interface{}) (int, error) {
	fmt.Fprintf(br.Config.Writer, "%s ", br.ResourceID())

	return fmt.Fprintf(br.Config.Writer, format, a...)
}
