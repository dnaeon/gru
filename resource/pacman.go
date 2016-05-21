// +build linux

package resource

import (
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
	Name string `hcl:"name"`
}

// NewPacmanResource creates a new resource for managing packages
// using the pacman package manager on an Arch Linux system
func NewPacmanResource(title string, obj *ast.ObjectItem, config *Config) (Resource, error) {
	// Resource defaults
	defaults := &PacmanResource{
		BaseResource: BaseResource{
			Title:  title,
			Type:   pacmanResourceType,
			State:  StatePresent,
			Config: config,
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

// Create installs packages
func (p *PacmanResource) Create() error {
	p.Printf("installing package\n")

	cmd := exec.Command(pacmanPath, "--sync", "--noconfirm", p.Name)
	out, err := cmd.CombinedOutput()
	p.Printf(string(out))

	return err
}

// Delete deletes packages
func (p *PacmanResource) Delete() error {
	p.Printf("removing package\n")

	cmd := exec.Command(pacmanPath, "--remove", "--noconfirm", p.Name)
	out, err := cmd.CombinedOutput()
	p.Printf(string(out))

	return err
}

// Update updates packages
func (p *PacmanResource) Update() error {
	p.Printf("updating package\n")

	return p.Create()
}

func init() {
	item := RegistryItem{
		Name:        pacmanResourceType,
		Description: pacmanResourceDesc,
		Provider:    NewPacmanResource,
	}

	Register(item)
}
