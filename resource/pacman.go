// +build linux

package resource

import (
	"os/exec"
	"strings"

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
	BaseResource        `hcl:",squash"`
	BasePackageResource `hcl:",squash"`
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
		BasePackageResource: BasePackageResource{
			Name: title,
		},
	}

	// Decode the object from HCL
	var pr PacmanResource
	err := hcl.DecodeObject(&pr, obj)
	if err != nil {
		return nil, err
	}

	// Merge in the decoded object with the resource defaults
	err = mergo.Merge(&pr, defaults)

	return &pr, err
}

// Evaluate evaluates the state of the resource
func (pr *PacmanResource) Evaluate() (State, error) {
	s := State{
		Current: StateUnknown,
		Want:    pr.State,
	}

	_, err := exec.LookPath(pacmanPath)
	if err != nil {
		return s, err
	}

	cmd := exec.Command(pacmanPath, "--query", pr.Name)
	err = cmd.Run()

	if err != nil {
		s.Current = StateAbsent
	} else {
		s.Current = StatePresent
	}

	return s, nil
}

// Create installs packages
func (pr *PacmanResource) Create() error {
	pr.Printf("installing package\n")

	cmd := exec.Command(pacmanPath, "--sync", "--noconfirm", pr.Name)
	out, err := cmd.CombinedOutput()

	for _, line := range strings.Split(string(out), "\n") {
		pr.Printf("%s\n", line)
	}

	return err
}

// Delete deletes packages
func (pr *PacmanResource) Delete() error {
	pr.Printf("removing package\n")

	cmd := exec.Command(pacmanPath, "--remove", "--noconfirm", pr.Name)
	out, err := cmd.CombinedOutput()

	for _, line := range strings.Split(string(out), "\n") {
		pr.Printf("%s\n", line)
	}

	return err
}

// Update updates packages
func (pr *PacmanResource) Update() error {
	pr.Printf("updating package\n")

	return pr.Create()
}

func init() {
	item := RegistryItem{
		Name:        pacmanResourceType,
		Description: pacmanResourceDesc,
		Provider:    NewPacmanResource,
	}

	Register(item)
}
