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
		},
		Enable: true,
		RCVar:  fmt.Sprintf("%v_enable", name),
	}

	return s, nil
}

// isEnabled returns true if service is set to start at boot.
func (s *Service) isEnabled() bool {
	if err := exec.Command("service", s.Name, "enabled").Run(); err != nil {
		return false
	}

	return true
}

// Evaluate evaluates the state of the resource.
func (s *Service) Evaluate() (State, error) {
	state := State{
		Current:  "unknown",
		Want:     s.State,
		Outdated: false,
	}

	// TODO: handle non existent service
	err := exec.Command("service", s.Name, "onestatus").Run()
	if err != nil {
		state.Current = "stopped"
	} else {
		state.Current = "running"
	}

	if s.Enable != s.isEnabled() {
		state.Outdated = true
	}

	return state, nil
}

// Create starts the service.
func (s *Service) Create() error {
	return exec.Command("service", s.Name, "onestart").Run()
}

// Delete stops the service.
func (s *Service) Delete() error {
	return exec.Command("service", s.Name, "onestop").Run()
}

// Update updates the service's rcvar.
func (s *Service) Update() error {
	if s.RCVar == "" {
		return nil
	}

	rcValue := "YES"
	if !s.Enable {
		rcValue = "NO"
		if s.RCVar == "sendmail_enable" {
			// I think sendmail is the only service, that requires NONE to be disabled.
			rcValue = "NONE"
		}
	}

	// TODO: rcvar should probably be deleted from rc.conf, when disabling service (except for sendmail).
	// Currently we set it to NO.
	return exec.Command("sysrc", fmt.Sprintf(`%s="%s"`, s.RCVar, rcValue)).Run()
}

func init() {
	service := RegistryItem{
		Type:      "service",
		Provider:  NewService,
		Namespace: DefaultNamespace,
	}

	Register(service)
}
