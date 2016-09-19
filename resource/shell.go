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
			Name:          name,
			Type:          "shell",
			State:         "present",
			Require:       make([]string, 0),
			PresentStates: []string{"present"},
			AbsentStates:  []string{"absent"},
			Concurrent:    true,
			Subscribe:     make(TriggerMap),
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
	Log(s, "executing command\n")

	args := strings.Fields(s.Command)
	cmd := exec.Command(args[0], args[1:]...)
	out, err := cmd.CombinedOutput()

	if !s.Mute {
		for _, line := range strings.Split(string(out), "\n") {
			Log(s, "%s\n", line)
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
