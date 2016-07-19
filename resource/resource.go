package resource

import (
	"errors"
	"fmt"
	"log"
	"os"
)

// Resource is the interface type for resources
type Resource interface {
	// ID returns the unique identifier of the resource
	ID() string

	// Validate validates the resource
	Validate() error

	// GetBefore returns the list of resources before which this
	// resource shoud be processed
	GetBefore() []string

	// GetAfter returns the list of resources after which this
	// resource should be processed
	GetAfter() []string

	// GetPresentStates returns the list of states, for which the
	// resource is considered to be present
	GetPresentStates() []string

	// GetAbsentStates returns the list of states, for which the
	// resource is considered to be absent
	GetAbsentStates() []string

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

// Validate validates the resource
func (br *BaseResource) Validate() error {
	if br.Type == "" {
		return errors.New("Invalid resource type")
	}

	if br.Name == "" {
		return errors.New("Invalid resource name")
	}
}

// GetBefore returns the list of resources before which this resource
// should be processed
func (br *BaseResource) GetBefore() []string {
	return br.Before
}

// GetAfter returns the list of resources after which this resource
// should be processed
func (br *BaseResource) GetAfter() []string {
	return br.After
}

// Log writes to the default config writer object and
// prepends the resource id to the output
func (br *BaseResource) Log(format string, a ...interface{}) {
	f := fmt.Sprintf("%s %s", br.ID(), format)
	Log(f, a...)
}
