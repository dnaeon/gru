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
	"os"
	"os/exec"
	"strings"
)

// Shell type is a resource which executes shell commands.
//
// The command that is to be executed should be idempotent.
// If the command that is to be executed is not idempotent on it's own,
// in order to achieve idempotency of the resource you should set the
// "creates" field to a filename that can be checked for existence.
//
// Example:
//   sh = resource.shell.new("touch /tmp/foo")
//   sh.creates = "/tmp/foo"
//
// Same example as the above one, but written in a different way.
//
// Example:
//   sh = resource.shell.new("creates the /tmp/foo file")
//   sh.command = "/usr/bin/touch /tmp/foo"
//   sh.creates = "/tmp/foo"
type Shell struct {
	Base

	// Command to be executed. Defaults to the resource name.
	Command string `luar:"command"`

	// File to be checked for existence before executing the command.
	Creates string `luar:"creates"`

	// Mute flag indicates whether output from the command should be
	// dislayed or suppressed
	Mute bool `luar:"mute"`
}

// NewShell creates a new resource for executing shell commands
func NewShell(name string) (Resource, error) {
	s := &Shell{
		Base: Base{
			Name:              name,
			Type:              "shell",
			State:             "present",
			Require:           make([]string, 0),
			PresentStatesList: []string{"present"},
			AbsentStatesList:  []string{"absent"},
			Concurrent:        true,
			Subscribe:         make(TriggerMap),
		},
		Command: name,
		Creates: "",
		Mute:    false,
	}

	return s, nil
}

// Evaluate evaluates the state of the resource
func (s *Shell) Evaluate() (State, error) {
	// Assumes that the command to be executed is idempotent
	//
	// Sets the current state to absent and wanted to be present,
	// which will cause the command to be executed.
	//
	// If the command to be executed is not idempotent on it's own,
	// in order to ensure idempotency we should specify a file,
	// that can be checked for existence.
	state := State{
		Current:  "absent",
		Want:     s.State,
		Outdated: false,
	}

	if s.Creates != "" {
		_, err := os.Stat(s.Creates)
		if os.IsNotExist(err) {
			state.Current = "absent"
		} else {
			state.Current = "present"
		}
	}

	return state, nil
}

// Create executes the shell command
func (s *Shell) Create() error {
	Logf("%s executing command\n", s.ID())

	args := strings.Fields(s.Command)
	cmd := exec.Command(args[0], args[1:]...)
	out, err := cmd.CombinedOutput()

	if !s.Mute {
		for _, line := range strings.Split(string(out), "\n") {
			Logf("%s %s\n", s.ID(), line)
		}
	}

	return err
}

// Delete is a no-op
func (s *Shell) Delete() error {
	return nil
}

// Update is a no-op
func (s *Shell) Update() error {
	return nil
}

func init() {
	item := ProviderItem{
		Type:      "shell",
		Provider:  NewShell,
		Namespace: DefaultResourceNamespace,
	}

	RegisterProvider(item)
}
