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
			Name:              name,
			Type:              "service",
			State:             "running",
			Require:           make([]string, 0),
			PresentStatesList: []string{"present", "running"},
			AbsentStatesList:  []string{"absent", "stopped"},
			Concurrent:        false,
			Subscribe:         make(TriggerMap),
		},
		Enable: true,
		RCVar:  fmt.Sprintf("%v_enable", name),
	}

	// Set resource properties
	s.PropertyList = []Property{
		&ResourceProperty{
			PropertyName:         "enable",
			PropertySetFunc:      s.setEnable,
			PropertyIsSyncedFunc: s.isEnableSynced,
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
	Logf("%s starting service\n", s.ID())

	return exec.Command("service", s.Name, "onestart").Run()
}

// Delete stops the service.
func (s *Service) Delete() error {
	Logf("%s stopping service\n", s.ID())

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

	return enabled == s.Enable, nil
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
	err := exec.Command("sysrc", fmt.Sprintf(`%s="%s"`, s.RCVar, rcValue)).Run()

	return err
}

func init() {
	service := ProviderItem{
		Type:      "service",
		Provider:  NewService,
		Namespace: DefaultResourceNamespace,
	}

	RegisterProvider(service)
}
