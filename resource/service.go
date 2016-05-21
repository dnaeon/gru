// +build linux

package resource

import (
	"errors"
	"fmt"

	"github.com/coreos/go-systemd/dbus"
	"github.com/coreos/go-systemd/util"
	"github.com/hashicorp/hcl"
	"github.com/hashicorp/hcl/hcl/ast"
	"github.com/imdario/mergo"
)

// Name and description of the resource
const serviceResourceType = "service"
const serviceResourceDesc = "manages services using systemd"

// ServiceResource type is a resource which manages
// services on a GNU/Linux system running systemd
type ServiceResource struct {
	BaseResource `hcl:",squash"`

	// Name of the service
	Name string `hcl:"name"`

	// If true then enable service during boot-time
	Enable bool `hcl:"enable"`

	// Systemd unit name
	UnitName string `hcl:"-"`
}

// NewServiceResource creates a new resource for managing services
// using systemd on a GNU/Linux system
func NewServiceResource(title string, obj *ast.ObjectItem, config *Config) (Resource, error) {
	// Resource defaults
	defaults := &ServiceResource{
		BaseResource: BaseResource{
			Title:  title,
			Type:   serviceResourceType,
			State:  StateRunning,
			Config: config,
		},
		Name:   title,
		Enable: false,
	}

	var sr ServiceResource
	err := hcl.DecodeObject(&sr, obj)
	if err != nil {
		return nil, err
	}

	// Merge the decoded object with the resource defaults
	err = mergo.Merge(&sr, defaults)

	// Set the unit name for the service we manage
	sr.UnitName = fmt.Sprintf("%s.service", sr.Name)

	return &sr, err
}

// unitProperty retrieves the requested property for the service unit
func (sr *ServiceResource) unitProperty(propertyName string) (*dbus.Property, error) {
	conn, err := dbus.New()
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	property, err := conn.GetUnitProperty(sr.UnitName, propertyName)

	return property, err
}

// unitIsEnabled checks if the unit is enabled or disabled
func (sr *ServiceResource) unitIsEnabled() (bool, error) {
	unitState, err := sr.unitProperty("UnitFileState")
	if err != nil {
		return false, err
	}

	value := unitState.Value.Value().(string)
	switch value {
	case "enabled", "static", "enabled-runtime", "linked", "linked-runtime":
		return true, nil
	case "disabled", "masked", "masked-runtime":
		return false, nil
	case "invalid":
		fallthrough
	default:
		return false, errors.New("Invalid unit state")
	}
}

// enableUnit enables the service unit during boot-time
func (sr *ServiceResource) enableUnit() error {
	conn, err := dbus.New()
	if err != nil {
		return err
	}
	defer conn.Close()

	sr.Printf("enabling service\n")

	units := []string{sr.UnitName}
	_, changes, err := conn.EnableUnitFiles(units, false, false)
	if err != nil {
		return err
	}

	for _, change := range changes {
		sr.Printf("%s %s -> %s\n", change.Type, change.Filename, change.Destination)
	}

	return nil
}

// disableUnit disables the service unit during boot-time
func (sr *ServiceResource) disableUnit() error {
	conn, err := dbus.New()
	if err != nil {
		return err
	}
	defer conn.Close()

	sr.Printf("disabling service\n")

	units := []string{sr.UnitName}
	changes, err := conn.DisableUnitFiles(units, false)
	if err != nil {
		return err
	}

	for _, change := range changes {
		sr.Printf("%s %s\n", change.Type, change.Filename)
	}

	return nil
}

// daemonReload instructs systemd to scan for and reload unit files
func (sr *ServiceResource) daemonReload() error {
	conn, err := dbus.New()
	if err != nil {
		return err
	}
	defer conn.Close()

	return conn.Reload()
}

// Evaluate evaluates the state of the resource
func (sr *ServiceResource) Evaluate() (State, error) {
	rs := State{
		Current: StateUnknown,
		Want:    sr.State,
		Update:  false,
	}

	// Check if the unit is started/stopped
	activeState, err := sr.unitProperty("ActiveState")
	if err != nil {
		return rs, err
	}

	// TODO: Handle cases where the unit is not found

	value := activeState.Value.Value().(string)
	switch value {
	case "active", "reloading", "activating":
		rs.Current = StateRunning
	case "inactive", "failed", "deactivating":
		rs.Current = StateStopped
	}

	// Check if the unit is enabled/disabled
	enabled, err := sr.unitIsEnabled()
	if err != nil {
		return rs, err
	}

	// Check if the resource needs to be updated
	if sr.Enable != enabled {
		rs.Update = true
	}

	return rs, nil
}

// Create starts the service unit
func (sr *ServiceResource) Create() error {
	conn, err := dbus.New()
	if err != nil {
		return err
	}
	defer conn.Close()

	sr.Printf("starting service\n")

	ch := make(chan string)
	jobID, err := conn.StartUnit(sr.UnitName, "replace", ch)
	if err != nil {
		return err
	}

	result := <-ch
	sr.Printf("systemd job id %d result: %s\n", jobID, result)

	return nil
}

// Delete stops the service unit
func (sr *ServiceResource) Delete() error {
	conn, err := dbus.New()
	if err != nil {
		return err
	}
	defer conn.Close()

	sr.Printf("stopping service\n")

	ch := make(chan string)
	jobID, err := conn.StopUnit(sr.UnitName, "replace", ch)
	if err != nil {
		return err
	}

	result := <-ch
	sr.Printf("systemd job id %d result: %s\n", jobID, result)

	return nil
}

// Update updates the service unit state
func (sr *ServiceResource) Update() error {
	enabled, err := sr.unitIsEnabled()
	if err != nil {
		return err
	}

	if sr.Enable && !enabled {
		sr.enableUnit()
	} else {
		sr.disableUnit()
	}

	return sr.daemonReload()
}

func init() {
	if util.IsRunningSystemd() {
		item := RegistryItem{
			Name:        serviceResourceType,
			Description: serviceResourceDesc,
			Provider:    NewServiceResource,
		}

		Register(item)
	}
}
