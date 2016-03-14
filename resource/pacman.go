// +build linux

package resource

import (
	"fmt"
	"log"
	"os/exec"

	"github.com/hashicorp/hcl"
	"github.com/hashicorp/hcl/hcl/ast"
	"github.com/imdario/mergo"
)

// Path to the pacman package manager
const pacmanPath = "/usr/bin/pacman"

// Name of the resource type in HCL
const pacmanResourceTypeName = "pacman"

// PacmanResource type is the resource for managing
// packages on an Arch Linux system
type PacmanResource struct {
	// Name of the resource
	Name string `json:"name"`

	// State of the resource
	State string `json:"state"`

	// Resource dependencies
	WantResource []string `hcl:"want" json:"want,omitempty"`
}

// NewPacmanResource creates a new resource for managing packages
// using the pacman package manager on an Arch Linux system
func NewPacmanResource(obj *ast.ObjectItem) (Resource, error) {
	// Position of the resource declaration
	position := obj.Val.Pos().String()

	// Resource defaults
	defaults := &PacmanResource{
		State: ResourceStatePresent,
	}

	// Decode the object from HCL
	var p PacmanResource
	err := hcl.DecodeObject(&p, obj)
	if err != nil {
		return nil, err
	}

	// Merge in the decoded object with the resource defaults
	err = mergo.Merge(&p, defaults)

	// Sanity check the resource
	if p.Name == "" {
		return nil, fmt.Errorf("Missing resource name at %s", position)
	}

	return &p, err
}

// Type returns the resource type for this resource
func (p *PacmanResource) Type() string {
	return "pacman"
}

// ID returns the unique resource identifier
func (p *PacmanResource) ID() string {
	id := fmt.Sprintf("%s[%s]", pacmanResourceTypeName, p.Name)

	return id
}

// Evaluate evaluates the resource
func (p *PacmanResource) Evaluate() (State, error) {
	s := State{
		Current: ResourceStateUnknown,
		Want:    p.State,
	}

	_, err := exec.LookPath(pacmanPath)
	if err != nil {
		return s, err
	}

	cmd := exec.Command(pacmanPath, "--query", p.Name)
	_, err = cmd.CombinedOutput()

	if err != nil {
		s.Current = ResourceStateAbsent
	} else {
		s.Current = ResourceStatePresent
	}

	return s, nil
}

// Create creates the resource
func (p *PacmanResource) Create() error {
	cmd := exec.Command(pacmanPath, "--sync", "--noconfirm", p.Name)
	output, err := cmd.CombinedOutput()
	log.Println(string(output))

	return err
}

// Delete deletes the resource
func (p *PacmanResource) Delete() error {
	cmd := exec.Command(pacmanPath, "--remove", "--noconfirm", p.Name)
	output, err := cmd.CombinedOutput()
	log.Println(string(output))

	return err
}

// Update updates the resource
func (p *PacmanResource) Update() error {
	// Create() handles package updates as well
	return p.Create()
}

// Want returns the wanted dependencies
func (p *PacmanResource) Want() []string {
	return p.WantResource
}

func init() {
	Register(pacmanResourceTypeName, NewPacmanResource)
}
