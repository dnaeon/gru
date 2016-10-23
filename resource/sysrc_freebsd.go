// +build freebsd

package resource

import (
	"fmt"
	"os/exec"
	"regexp"
)

// SysRC is a resource which manages rc.conf variables.
//
// Example:
//   rcvar = resource.sysrc.new("keyrate")
//   rcvar.state = "present"
//   rcvar.value = "fast"
type SysRC struct {
	Base
	Value string `luar:"value"`
}

// NewSysRC creates a new resource for managing rc.conf variables
// on a FreeBSD system.
func NewSysRC(name string) (Resource, error) {
	s := &SysRC{
		Base: Base{
			Name:              name,
			Type:              "sysrc",
			State:             "present",
			Require:           make([]string, 0),
			PresentStatesList: []string{"present"},
			AbsentStatesList:  []string{"absent"},
			Concurrent:        false,
			Subscribe:         make(TriggerMap),
		},
	}

	return s, nil
}

// Evaluate evaluates the state of the resource.
func (s *SysRC) Evaluate() (State, error) {
	state := State{
		Current:  "unknown",
		Want:     s.State,
		Outdated: false,
	}

	out, err := exec.Command("sysrc", s.Name).CombinedOutput()
	if err != nil {
		state.Current = "absent"
		state.Outdated = true
		return state, nil
	}
	state.Current = "present"

	k, v, err := parseSysRCOutput(string(out))
	if err != nil {
		return state, err
	}

	if s.Name != k {
		return state, fmt.Errorf("bug: expected rcvar %v, got %v", s.Name, k)
	}

	if s.Value != v {
		state.Outdated = true
	}

	return state, nil
}

// Create adds variable to rc.conf.
func (s *SysRC) Create() error {
	return exec.Command("sysrc", s.Name, s.Value).Run()
}

// Delete removes variable from rc.conf.
func (s *SysRC) Delete() error {
	return exec.Command("sysrc", "-x", s.Value).Run()
}

// Update sets variable in rc.conf to s.Value.
func (s *SysRC) Update() error {
	return exec.Command("sysrc", s.Name, s.Value).Run()
}

var sysRCre = regexp.MustCompile("(.*): (.*)")

// ParseSysRCOutput parses output from sysrc command.
func parseSysRCOutput(out string) (k, v string, err error) {
	m := sysRCre.FindStringSubmatch(out)
	if m == nil {
		return "", "", fmt.Errorf("bug: sysrc output %q didn't match regexp", out)
	}
	return m[1], m[2], nil
}

func init() {
	sysrc := ProviderItem{
		Type:      "sysrc",
		Provider:  NewSysRC,
		Namespace: DefaultResourceNamespace,
	}

	RegisterProvider(sysrc)
}
