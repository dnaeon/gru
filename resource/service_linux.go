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

	// Enable specifies whether to enable or disable the
	// service during boot-time. Defaults to true.
	Enable bool `luar:"enable"`

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
			Name:              name,
			Type:              "service",
			State:             "running",
			Require:           make([]string, 0),
			PresentStatesList: []string{"present", "running"},
			AbsentStatesList:  []string{"absent", "stopped"},
			Concurrent:        true,
			Subscribe:         make(TriggerMap),
		},
		Enable: true,
		unit:   fmt.Sprintf("%s.service", name),
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
	Logf("%s starting service\n", s.ID())

	ch := make(chan string)
	jobID, err := s.conn.StartUnit(s.unit, "replace", ch)
	if err != nil {
		return err
	}

	result := <-ch
	Logf("%s systemd job id %d result: %s\n", s.ID(), jobID, result)

	return nil
}

// Delete stops the service.
func (s *Service) Delete() error {
	Logf("%s stopping service\n", s.ID())

	ch := make(chan string)
	jobID, err := s.conn.StopUnit(s.unit, "replace", ch)
	if err != nil {
		return err
	}

	result := <-ch
	Logf("%s systemd job id %d result: %s\n", s.ID(), jobID, result)

	return nil
}

// Close closes the connection to the systemd D-BUS API
func (s *Service) Close() error {
	s.conn.Close()

	return nil
}

// enableUnit enables the service unit during boot-time
func (s *Service) enableUnit() error {
	Logf("%s enabling service\n", s.ID())

	units := []string{s.unit}
	_, changes, err := s.conn.EnableUnitFiles(units, false, false)
	if err != nil {
		return err
	}

	for _, change := range changes {
		Logf("%s %s %s -> %s\n", s.ID(), change.Type, change.Filename, change.Destination)
	}

	return nil
}

// disableUnit disables the service unit during boot-time
func (s *Service) disableUnit() error {
	Logf("%s disabling service\n", s.ID())

	units := []string{s.unit}
	changes, err := s.conn.DisableUnitFiles(units, false)
	if err != nil {
		return err
	}

	for _, change := range changes {
		Logf("%s %s %s\n", s.ID(), change.Type, change.Filename)
	}

	return nil
}

// isEnableSynced determines whether the property is synced.
func (s *Service) isEnableSynced() (bool, error) {
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

	return s.Enable == enabled, nil
}

// setEnable sets the property to it's desired state.
func (s *Service) setEnable() error {
	var action func() error

	switch s.Enable {
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
