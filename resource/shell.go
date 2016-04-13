package resource

import (
	"fmt"
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

// Validate validates the shell resource
func (s *ShellResource) Validate() error {
	err := s.BaseResource.Validate()
	if err != nil {
		return err
	}

	if s.Creates == "" {
		return fmt.Errorf("Missing required field 'creates' in %s", s.ResourceID())
	}

	return nil
}

// Evaluate evaluates the state of the resource
func (s *ShellResource) Evaluate() (State, error) {
	resourceState := State{
		Current: StateUnknown,
		Want:    s.State,
		Update:  false,
	}

	_, err := os.Stat(s.Creates)
	if os.IsNotExist(err) {
		resourceState.Current = StateAbsent
	} else {
		resourceState.Current = StatePresent
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
