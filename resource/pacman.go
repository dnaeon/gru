// +build linux

package resource

import (
	"os/exec"
	"strings"
)

// Path to the pacman package manager
const pacmanPath = "/usr/bin/pacman"

// PacmanResource type represents the resource for package
// management on Arch Linux systems
type PacmanResource struct {
	BaseResource
	BasePackageResource
}

// NewPacmanResource creates a new resource for managing packages
// using the pacman package manager on an Arch Linux system
func NewPacmanResource(title string) Resource {
	// Create resource with defaults
	pr := &PacmanResource{
		BaseResource: BaseResource{
			Title: title,
			Type:  "pacman",
			State: StatePresent,
		},
		BasePackageResource: BasePackageResource{
			Name: title,
		},
	}

	return &pr
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
	RegisterProvider("pacman", NewPacmanResource)
}
