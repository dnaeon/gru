package resource

import (
	"io"
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

	Command string `hcl:"command" json:"command"`

	Creates string `hcl:"creates" json:"creates"`
}

// NewShellResource creates a new resource for executing shell commands
func NewShellResource(name string, obj *ast.ObjectItem) (Resource, error) {
	// Resource defaults
	defaults := &ShellResource{
		BaseResource: BaseResource{
			Name:  name,
			Type:  shellResourceType,
			State: StatePresent,
		},
		Command: name,
	}

	var s ShellResource
	err := hcl.DecodeObject(&s, obj)
	if err != nil {
		return nil, err
	}

	// Merge the decoded object with the resource defaults
	err = mergo.Merge(&s, defaults)

	return &s, err
}

// Evaluate evaluates the state of the resource
func (s *ShellResource) Evaluate() (State, error) {
	// Asumes that the command to be executed is idempotent
	//
	// Sets the current state to absent and wanted to be present,
	// which will cause the command to be executed.
	//
	// If the command to be executed is not idempotent on it's own,
	// in order to ensure idempotency we should specify a file,
	// that can be checked for existence.
	resourceState := State{
		Current: StateAbsent,
		Want:    s.State,
		Update:  false,
	}

	if s.Creates != "" {
		_, err := os.Stat(s.Creates)
		if os.IsNotExist(err) {
			resourceState.Current = StateAbsent
		} else {
			resourceState.Current = StatePresent
		}
	}

	return resourceState, nil
}

// Create executes the shell command
func (s *ShellResource) Create(w io.Writer) error {
	s.Printf(w, "executing command")

	args := strings.Fields(s.Command)
	cmd := exec.Command(args[0], args[1:]...)
	out, err := cmd.CombinedOutput()

	s.Printf(w, string(out))

	return err
}

// Delete is a no-op
func (s *ShellResource) Delete(w io.Writer) error {
	return nil
}

// Update is a no-op
func (s *ShellResource) Update(w io.Writer) error {
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
