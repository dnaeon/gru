package resource

import (
	"fmt"
	"io"
)

// Resource is the interface type for resources
type Resource interface {
	// ResourceID returns the unique identifier of a resource
	ID() string

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
	// The site repo which contains module and data files
	SiteRepo string

	// Writer used by the resources
	Writer io.Writer
}

// BaseResource is the base resource type for all resources
// The purpose of this type is to be embedded into other resources
// Partially implements the Resource interface
type BaseResource struct {
	// Type of the resource
	Type string `luar:"-"`

	// Name of the resource
	Name string `luar:"-"`

	// Resource configuration settings
	Config *Config `luar:"-"`

	// Desired state of the resource
	State string `luar:"state"`

	// Resources before which this resource should be processed
	Before []string `luar:"before"`

	// Resources after which this resource should be processed
	After []string `luar:"require"`
}

// ID returns the unique resource id
func (br *BaseResource) ID() string {
	return fmt.Sprintf("%s[%s]", br.Type, br.Name)
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
	fmt.Fprintf(br.Config.Writer, "%s ", br.ID())

	return fmt.Fprintf(br.Config.Writer, format, a...)
}
