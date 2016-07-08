package resource

import (
	"fmt"
	"log"
	"os"
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

	// Log logs events
	Log(format string, a ...interface{})
}

// Config type contains various settings used by the resources
type Config struct {
	// The site repo which contains module and data files
	SiteRepo string

	// Logger used by the resources to log events
	Logger *log.Logger
}

// DefaultLogger is the default logger instance used for
// logging events from the resources
var DefaultLogger = log.New(os.Stdout, "", log.LstdFlags)

// DefaultConfig is the default configuration used by the resources
var DefaultConfig = &Config{
	Logger: DefaultLogger,
}

// Log logs an event using the default resource logger
func Log(format string, a ...interface{}) {
	DefaultConfig.Logger.Printf(format, a...)
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
	After []string `luar:"after"`
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

// Log writes to the default config writer object and
// prepends the resource id to the output
func (br *BaseResource) Log(format string, a ...interface{}) {
	f := fmt.Sprintf("%s %s", br.ID(), format)
	Log(f, a...)
}
