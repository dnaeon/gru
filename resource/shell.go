package resource

import (
	"os"
	"os/exec"
	"strings"

	"github.com/hashicorp/hcl"
	"github.com/hashicorp/hcl/hcl/ast"
	"github.com/imdario/mergo"
)

// Name and description of the resource
const shellResourceType = "shell"
const shellResourceDesc = "executes shell commands"

// ShellResource type is a resource which executes shell commands
type ShellResource struct {
	BaseResource `hcl:",squash"`

	Command string `hcl:"command"`

	Creates string `hcl:"creates"`
}

// NewShellResource creates a new resource for executing shell commands
func NewShellResource(title string, obj *ast.ObjectItem, config *Config) (Resource, error) {
	// Resource defaults
	defaults := &ShellResource{
		BaseResource: BaseResource{
			Title:  title,
			Type:   shellResourceType,
			State:  StatePresent,
			Config: config,
		},
		Command: title,
	}

	var sr ShellResource
	err := hcl.DecodeObject(&sr, obj)
	if err != nil {
		return nil, err
	}

	// Merge the decoded object with the resource defaults
	err = mergo.Merge(&sr, defaults)

	return &sr, err
}

// Evaluate evaluates the state of the resource
func (sr *ShellResource) Evaluate() (State, error) {
	// Assumes that the command to be executed is idempotent
	//
	// Sets the current state to absent and wanted to be present,
	// which will cause the command to be executed.
	//
	// If the command to be executed is not idempotent on it's own,
	// in order to ensure idempotency we should specify a file,
	// that can be checked for existence.
	rs := State{
		Current: StateAbsent,
		Want:    sr.State,
		Update:  false,
	}

	if sr.Creates != "" {
		_, err := os.Stat(sr.Creates)
		if os.IsNotExist(err) {
			rs.Current = StateAbsent
		} else {
			rs.Current = StatePresent
		}
	}

	return rs, nil
}

// Create executes the shell command
func (sr *ShellResource) Create() error {
	sr.Printf("executing command\n")

	args := strings.Fields(sr.Command)
	cmd := exec.Command(args[0], args[1:]...)
	out, err := cmd.CombinedOutput()
	sr.Printf(string(out))

	return err
}

// Delete is a no-op
func (sr *ShellResource) Delete() error {
	return nil
}

// Update is a no-op
func (sr *ShellResource) Update() error {
	return nil
}

func init() {
	item := RegistryItem{
		Name:        shellResourceType,
		Description: shellResourceDesc,
		Provider:    NewShellResource,
	}

	Register(item)
}
