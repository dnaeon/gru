// +build linux

package resource

import (
	"io"
	"os/exec"

	"github.com/hashicorp/hcl"
	"github.com/hashicorp/hcl/hcl/ast"
	"github.com/imdario/mergo"
)

// Path to the pacman package manager
const pacmanPath = "/usr/bin/pacman"

// Name and description of the resource
const pacmanResourceType = "pacman"
const pacmanResourceDesc = "manages packages using the pacman package manager"

// PacmanResource type represents the resource for
// package management on Arch Linux systems
type PacmanResource struct {
	BaseResource `hcl:",squash"`

	// Name of the package
	Name string `hcl:"name" json:"name"`
}

// NewPacmanResource creates a new resource for managing packages
// using the pacman package manager on an Arch Linux system
func NewPacmanResource(title string, obj *ast.ObjectItem) (Resource, error) {
	// Resource defaults
	defaults := &PacmanResource{
		BaseResource: BaseResource{
			Title: title,
			Type:  pacmanResourceType,
			State: StatePresent,
		},
		Name: title,
	}

	// Decode the object from HCL
	var p PacmanResource
	err := hcl.DecodeObject(&p, obj)
	if err != nil {
		return nil, err
	}

	// Merge in the decoded object with the resource defaults
	err = mergo.Merge(&p, defaults)

	return &p, err
}

// Evaluate evaluates the state of the resource
func (p *PacmanResource) Evaluate() (State, error) {
	s := State{
		Current: StateUnknown,
		Want:    p.State,
	}

	_, err := exec.LookPath(pacmanPath)
	if err != nil {
		return s, err
	}

	cmd := exec.Command(pacmanPath, "--query", p.Name)
	err = cmd.Run()

	if err != nil {
		s.Current = StateAbsent
	} else {
		s.Current = StatePresent
	}

	return s, nil
}

// Create creates the resource
func (p *PacmanResource) Create(w io.Writer) error {
	p.Printf(w, "installing package\n")

	cmd := exec.Command(pacmanPath, "--sync", "--noconfirm", p.Name)
	out, err := cmd.CombinedOutput()
	p.Printf(w, string(out))

	return err
}

// Delete deletes the resource
func (p *PacmanResource) Delete(w io.Writer) error {
	p.Printf(w, "removing package\n")

	cmd := exec.Command(pacmanPath, "--remove", "--noconfirm", p.Name)
	out, err := cmd.CombinedOutput()
	p.Printf(w, string(out))

	return err
}

// Update updates the resource
func (p *PacmanResource) Update(w io.Writer) error {
	p.Printf(w, "updating package\n")

	return p.Create(w)
}

func init() {
	item := RegistryItem{
		Name:        pacmanResourceType,
		Description: pacmanResourceDesc,
		Provider:    NewPacmanResource,
	}

	Register(item)
}
