package resource

import (
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/dnaeon/gru/utils"
)

// State type represents the current and wanted states of a resource
type State struct {
	// Current state of the resource
	Current string

	// Wanted state of the resource
	Want string

	// Outdated indicates that a property of the resource is out of date
	Outdated bool
}

// Resource is the interface type for resources
type Resource interface {
	// ID returns the unique identifier of the resource
	ID() string

	// Validate validates the resource
	Validate() error

	// Dependencies returns the list of resource dependencies.
	// Each item in the slice is a string representing the
	// resource id for each dependency.
	Dependencies() []string

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

// Base is the base resource type for all resources
// The purpose of this type is to be embedded into other resources
// Partially implements the Resource interface
type Base struct {
	// Type of the resource
	Type string `luar:"-"`

	// Name of the resource
	Name string `luar:"-"`

	// Desired state of the resource
	State string `luar:"state"`

	// Require contains the resource dependencies
	Require []string `luar:"require"`

	// PresentStates contains the list of states, for which the
	// resource is considered to be present
	PresentStates []string `luar:"-"`

	// AbsentStates contains the list of states, for which the
	// resource is considered to be absent
	AbsentStates []string `luar:"-"`
}

// ID returns the unique resource id
func (b *Base) ID() string {
	return fmt.Sprintf("%s[%s]", b.Type, b.Name)
}

// Validate validates the resource
func (b *Base) Validate() error {
	if b.Type == "" {
		return errors.New("Invalid resource type")
	}

	if b.Name == "" {
		return errors.New("Invalid resource name")
	}

	states := append(b.PresentStates, b.AbsentStates...)
	if !utils.NewList(states...).Contains(b.State) {
		return fmt.Errorf("Invalid state '%s'", b.State)
	}

	return nil
}

// Dependencies returns the list of resource dependencies.
func (b *Base) Dependencies() []string {
	return b.Require
}

// GetPresentStates returns the list of states, for which the
// resource is considered to be present
func (b *Base) GetPresentStates() []string {
	return b.PresentStates
}

// GetAbsentStates returns the list of states, for which the
// resource is considered to be absent
func (b *Base) GetAbsentStates() []string {
	return b.AbsentStates
}

// Log writes to the default config writer object and
// prepends the resource id to the output
func (b *Base) Log(format string, a ...interface{}) {
	f := fmt.Sprintf("%s %s", b.ID(), format)
	Log(f, a...)
}
