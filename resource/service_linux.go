// +build linux

package resource

import (
	"errors"
	"fmt"

	"github.com/coreos/go-systemd/dbus"
	"github.com/coreos/go-systemd/util"
)

// ErrNoSystemd error is returned when the system is detected to
// have no support for systemd.
var ErrNoSystemd = errors.New("No systemd support found")

// Service type is a resource which manages services on a
// GNU/Linux system running with systemd.
//
// Example:
//   svc = resource.service.new("nginx")
//   svc.state = "running"
//   svc.enable = true
type Service struct {
	Base

	// EnableProperty specifies whether to enable or disable the
	// service during boot-time. Defaults to true.
	EnableProperty bool `luar:"enable"`

	// Systemd unit name
	unit string `luar:"-"`

	conn *dbus.Conn `luar:"-"`
}

// NewService creates a new resource for managing services
// using systemd on a GNU/Linux system
func NewService(name string) (Resource, error) {
	if !util.IsRunningSystemd() {
		return nil, ErrNoSystemd
	}

	s := &Service{
		Base: Base{
			Name:          name,
			Type:          "service",
			State:         "running",
			Require:       make([]string, 0),
			PresentStates: []string{"present", "running"},
			AbsentStates:  []string{"absent", "stopped"},
			Concurrent:    true,
			Subscribe:     make(TriggerMap),
		},
		EnableProperty: true,
		unit:           fmt.Sprintf("%s.service", name),
	}

	// Set resource properties
	s.Properties = []Property{
		Property{
			Name:     "enable",
			Set:      s.setEnableProperty,
			IsSynced: s.isEnablePropertySynced,
		},
	}

	return s, nil
}

// Initialize initializes the service resource by establishing a
// connection the systemd D-BUS API
func (s *Service) Initialize() error {
	conn, err := dbus.New()
	s.conn = conn

	return err
}

// Evaluate evaluates the state of the resource
func (s *Service) Evaluate() (State, error) {
	state := State{
		Current: "unknown",
		Want:    s.State,
	}

	// Check if the unit is started/stopped
	activeState, err := s.conn.GetUnitProperty(s.unit, "ActiveState")
	if err != nil {
		return state, err
	}

	// TODO: Handle cases where the unit is not found

	value := activeState.Value.Value().(string)
	switch value {
	case "active", "reloading", "activating":
		state.Current = "running"
	case "inactive", "failed", "deactivating":
		state.Current = "stopped"
	}

	return state, nil
}

// Create starts the service.
func (s *Service) Create() error {
	Log(s, "starting service\n")

	ch := make(chan string)
	jobID, err := s.conn.StartUnit(s.unit, "replace", ch)
	if err != nil {
		return err
	}

	result := <-ch
	Log(s, "systemd job id %d result: %s\n", jobID, result)

	return nil
}

// Delete stops the service.
func (s *Service) Delete() error {
	Log(s, "stopping service\n")

	ch := make(chan string)
	jobID, err := s.conn.StopUnit(s.unit, "replace", ch)
	if err != nil {
		return err
	}

	result := <-ch
	Log(s, "systemd job id %d result: %s\n", jobID, result)

	return nil
}

// Close closes the connection to the systemd D-BUS API
func (s *Service) Close() error {
	s.conn.Close()

	return nil
}

// enableUnit enables the service unit during boot-time
func (s *Service) enableUnit() error {
	Log(s, "enabling service\n")

	units := []string{s.unit}
	_, changes, err := s.conn.EnableUnitFiles(units, false, false)
	if err != nil {
		return err
	}

	for _, change := range changes {
		Log(s, "%s %s -> %s\n", change.Type, change.Filename, change.Destination)
	}

	return nil
}

// disableUnit disables the service unit during boot-time
func (s *Service) disableUnit() error {
	Log(s, "disabling service\n")

	units := []string{s.unit}
	changes, err := s.conn.DisableUnitFiles(units, false)
	if err != nil {
		return err
	}

	for _, change := range changes {
		Log(s, "%s %s\n", change.Type, change.Filename)
	}

	return nil
}

// isEnablePropertySynced determines whether the property is synced.
func (s *Service) isEnablePropertySynced() (bool, error) {
	unitState, err := s.conn.GetUnitProperty(s.unit, "UnitFileState")
	if err != nil {
		return false, err
	}

	var enabled bool
	value := unitState.Value.Value().(string)
	switch value {
	case "enabled", "static", "enabled-runtime", "linked", "linked-runtime":
		enabled = true
	case "disabled", "masked", "masked-runtime":
		enabled = false
	case "invalid":
		fallthrough
	default:
		return false, errors.New("Invalid unit state")
	}

	return s.EnableProperty == enabled, nil
}

// setEnableProperty sets the property to it's desired state.
func (s *Service) setEnableProperty() error {
	var action func() error

	switch s.EnableProperty {
	case true:
		action = s.enableUnit
	case false:
		action = s.disableUnit
	}

	if err := action(); err != nil {
		return err
	}

	return s.conn.Reload()
}

func init() {
	item := ProviderItem{
		Type:      "service",
		Provider:  NewService,
		Namespace: DefaultResourceNamespace,
	}

	RegisterProvider(item)
}
