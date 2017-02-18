// Copyright (c) 2015-2017 Marin Atanasov Nikolov <dnaeon@gmail.com>
// All rights reserved.
//
// Redistribution and use in source and binary forms, with or without
// modification, are permitted provided that the following conditions
// are met:
//
//  1. Redistributions of source code must retain the above copyright
//     notice, this list of conditions and the following disclaimer
//     in this position and unchanged.
//  2. Redistributions in binary form must reproduce the above copyright
//     notice, this list of conditions and the following disclaimer in the
//     documentation and/or other materials provided with the distribution.
//
// THIS SOFTWARE IS PROVIDED BY THE AUTHOR(S) ``AS IS'' AND ANY EXPRESS OR
// IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED WARRANTIES
// OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE DISCLAIMED.
// IN NO EVENT SHALL THE AUTHOR(S) BE LIABLE FOR ANY DIRECT, INDIRECT,
// INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT
// NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE,
// DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY
// THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
// (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF
// THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.

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
		Current: "unknown",
		Want:    s.State,
	}

	out, err := exec.Command("sysrc", s.Name).CombinedOutput()
	if err != nil {
		state.Current = "absent"
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
		state.Current = "absent"
	}

	return state, nil
}

// Create adds variable to rc.conf.
func (s *SysRC) Create() error {
	Logf("%s adding rcvar\n", s.ID())

	return exec.Command("sysrc", fmt.Sprintf("%s=%s", s.Name, s.Value)).Run()
}

// Delete removes variable from rc.conf.
func (s *SysRC) Delete() error {
	Logf("%s removing rcvar\n", s.ID())

	return exec.Command("sysrc", "-x", s.Name).Run()
}

// Update sets variable in rc.conf to s.Value.
func (s *SysRC) Update() error {
	Logf("%s setting rcvar to %s\n", s.ID(), s.Value)

	return exec.Command("sysrc", fmt.Sprintf("%s=%s", s.Name, s.Value)).Run()
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
