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

package resource

import (
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/dnaeon/gru/utils"
	"github.com/yuin/gopher-lua"
)

// ErrInvalidType error is returned when a resource type is invalid
var ErrInvalidType = errors.New("Invalid resource type")

// ErrInvalidName error is returned when a resource name is invalid
var ErrInvalidName = errors.New("Invalid resource name")

// ErrNotImplemented error is returned when a resource does not
// implement specific functionality, e.g. the resource is not
// applicable for a refresh.
var ErrNotImplemented = errors.New("Not implemented")

// ErrInSync error is returned when a resource is in the desired state
var ErrInSync = errors.New("Resource is in sync")

// ErrResourceAbsent error is returned by properties in situations where
// setting up a property makes no sense if the resource is absent, e.g.
// setting up file permissions makes no sense if the file resource is in
// absent state.
var ErrResourceAbsent = errors.New("Resource is absent")

// TriggerMap type is a map type which keys are
// resource ids for which a resource subscribes for changes to.
// The keys of the map are Lua functions that would be executed
// when the resource state has changed.
type TriggerMap map[string]*lua.LFunction

// State type represents the current and wanted states of a resource
type State struct {
	// Current state of the resource
	Current string

	// Wanted state of the resource
	Want string

	// Outdated indicates that a property of the resource is out of date
	Outdated bool
}

// Resource is the interface type for resources.
type Resource interface {
	// ID returns the unique identifier of the resource
	ID() string

	// Initialize is used to perform any initialization prior the
	// actual resource processing, e.g. establish connection to a
	// remote API endpoint.
	Initialize() error

	// Close performs any cleanup tasks after a resource has been processed.
	Close() error

	// Validate validates the resource
	Validate() error

	// Dependencies returns the list of resource dependencies.
	// Each item in the slice is a string representing the
	// resource id for each dependency.
	Dependencies() []string

	// PresentStates returns the list of states, for which the
	// resource is considered to be present
	PresentStates() []string

	// AbsentStates returns the list of states, for which the
	// resource is considered to be absent
	AbsentStates() []string

	// IsConcurrent returns a boolean, which indicates whether
	// multiple instances of the same resource type can be
	// processed concurrently.
	IsConcurrent() bool

	// Evaluates the resource
	Evaluate() (State, error)

	// Creates the resource
	Create() error

	// Deletes the resource
	Delete() error

	// Properties returns the list of properties for the resource.
	Properties() []Property

	// SubscribedTo returns a map of the resource ids for which the
	// current resource subscribes for changes to. The keys of the
	// map are resource ids and their values are the functions to be
	// executed if the resource state changes.
	SubscribedTo() TriggerMap
}

// Config type contains various settings used by the resources
type Config struct {
	// The site repo which contains module and data files
	SiteRepo string

	// Logger used by the resources to log events
	Logger *log.Logger
}

// DefaultConfig is the default configuration used by the resources
var DefaultConfig = &Config{
	Logger: log.New(os.Stdout, "", log.LstdFlags),
}

// Logf writes an event to the default logger.
func Logf(format string, a ...interface{}) {
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

	// PresentStatesList contains the list of states, for which the
	// resource is considered to be present
	PresentStatesList []string `luar:"-"`

	// AbsentStates contains the list of states, for which the
	// resource is considered to be absent
	AbsentStatesList []string `luar:"-"`

	// Concurrent flag indicates whether multiple instances of the
	// same resource type can be processed concurrently.
	Concurrent bool `luar:"-"`

	// PropertyList contains the resource properties.
	PropertyList []Property `luar:"-"`

	// Subscribe is map whose keys are resource ids that the
	// current resource monitors for changes and the values are
	// functions that will be executed if the monitored
	// resource state has changed.
	// Subscribing to changes in other resources also automatically
	// creates an edge in the dependency graph pointing from the
	// current resource to the one that is being monitored, so that the
	// monitored resource is evaluated and processed first.
	Subscribe map[string]*lua.LFunction `luar:"subscribe"`
}

// ID returns the unique resource id
func (b *Base) ID() string {
	return fmt.Sprintf("%s[%s]", b.Type, b.Name)
}

// Initialize initializes the resource prior the actual processing, e.g.
// establishing connection to a remote API endpoint.
func (b *Base) Initialize() error {
	return nil
}

// Close is used to perform cleanup tasks after the resource has been processed.
func (b *Base) Close() error {
	return nil
}

// Validate validates the resource
func (b *Base) Validate() error {
	if b.Type == "" {
		return ErrInvalidType
	}

	if b.Name == "" {
		return ErrInvalidName
	}

	states := append(b.PresentStatesList, b.AbsentStatesList...)
	if !utils.NewList(states...).Contains(b.State) {
		return fmt.Errorf("Invalid state '%s'", b.State)
	}

	return nil
}

// Dependencies returns the list of resource dependencies.
func (b *Base) Dependencies() []string {
	return b.Require
}

// PresentStates returns the list of states, for which the
// resource is considered to be present
func (b *Base) PresentStates() []string {
	return b.PresentStatesList
}

// AbsentStates returns the list of states, for which the
// resource is considered to be absent
func (b *Base) AbsentStates() []string {
	return b.AbsentStatesList
}

// IsConcurrent returns a boolean indicating whether
// multiple instances of the same resource type can be
// processed concurrently.
func (b *Base) IsConcurrent() bool {
	return b.Concurrent
}

// SubscribedTo returns a map of resources for which the
// resource is subscribed for changes to.
func (b *Base) SubscribedTo() TriggerMap {
	return b.Subscribe
}

// Properties returns the list of properties for the resource.
func (b *Base) Properties() []Property {
	return b.PropertyList
}
