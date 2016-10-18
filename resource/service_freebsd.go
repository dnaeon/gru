// +build freebsd

package resource

import (
	"fmt"
	"os/exec"
)

// Service type is a resource which manages services on a
// FreeBSD system.
//
// Example:
//   svc = resource.service.new("nginx")
//   svc.state = "running"
//   svc.enable = true
//   svc.rcvar = "nginx_enable"
type Service struct {
	Base

	// If true then enable the service during boot-time
	Enable bool `luar:"enable"`

	// RCVar (see rc.subr(8)), set to {svcname}_enable by default.
	// If service doesn't define rcvar, you should set svc.rcvar = "".
	RCVar string `luar:"rcvar"`
}

// NewService creates a new resource for managing services
// on a FreeBSD system.
func NewService(name string) (Resource, error) {
	s := &Service{
		Base: Base{
			Name:          name,
			Type:          "service",
			State:         "running",
			Require:       make([]string, 0),
			PresentStates: []string{"present", "running"},
			AbsentStates:  []string{"absent", "stopped"},
			Concurrent:    false,
			Subscribe:     make(TriggerMap),
		},
		Enable: true,
		RCVar:  fmt.Sprintf("%v_enable", name),
	}

	// Set resource properties
	s.Properties = []Property{
		Property{
			Name:     "enable",
			Set:      s.setEnable,
			IsSynced: s.isEnableSynced,
		},
	}

	return s, nil
}

// Evaluate evaluates the state of the resource.
func (s *Service) Evaluate() (State, error) {
	state := State{
		Current: "unknown",
		Want:    s.State,
	}

	// TODO: handle non existent service
	err := exec.Command("service", s.Name, "onestatus").Run()
	if err != nil {
		state.Current = "stopped"
	} else {
		state.Current = "running"
	}

	return state, nil
}

// Create starts the service.
func (s *Service) Create() error {
	Log(s, "starting service\n")

	return exec.Command("service", s.Name, "onestart").Run()
}

// Delete stops the service.
func (s *Service) Delete() error {
	Log(s, "stopping service\n")

	return exec.Command("service", s.Name, "onestop").Run()
}

// isEnableSynced checks whether the service is in the desired state.
func (s *Service) isEnableSynced() (bool, error) {
	var enabled bool

	err := exec.Command("service", s.Name, "enabled").Run()
	switch err {
	case nil:
		enabled = true
	default:
		enabled = false
	}

	return enabled != s.Enable
}

// setEnable enables or disables the service during boot-time.
func (s *Service) setEnable() error {
	if s.RCVar == "" {
		return nil
	}

	var rcValue string
	switch s.Enable {
	case true:
		rcValue = "YES"
	case false:
		rcValue = "NO"
	}

	// TODO: rcvar should probably be deleted from rc.conf, when disabling service.
	// Compare default value (sysrc -D) with requested (rcValue) and if they match, delete rcvar.
	// Currently we just set it to NO.
	return exec.Command("sysrc", fmt.Sprintf(`%s="%s"`, s.RCVar, rcValue)).Run()
}

func init() {
	service := ProviderItem{
		Type:      "service",
		Provider:  NewService,
		Namespace: DefaultResourceNamespace,
	}

	RegisterProvider(service)
}
