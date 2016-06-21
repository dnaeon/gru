// +build linux

package resource

import (
	"os/exec"
	"strings"
)

// Path to the pacman package manager
const pacmanPath = "/usr/bin/pacman"

// Pacman type represents the resource for package
// management on Arch Linux systems
type Pacman struct {
	BaseResource
	BasePackageResource
}

// NewPacman creates a new resource for managing packages
// using the pacman package manager on an Arch Linux system
func NewPacman(title string) Resource {
	// Create resource with defaults
	p := &Pacman{
		BaseResource: BaseResource{
			Title: title,
			Type:  "pacman",
			State: StatePresent,
		},
		BasePackageResource: BasePackageResource{
			Name: title,
		},
	}

	return &p
}

// Evaluate evaluates the state of the resource
func (p *Pacman) Evaluate() (State, error) {
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
func (p *Pacman) Create() error {
	pr.Printf("installing package\n")

	cmd := exec.Command(pacmanPath, "--sync", "--noconfirm", pr.Name)
	out, err := cmd.CombinedOutput()

	for _, line := range strings.Split(string(out), "\n") {
		pr.Printf("%s\n", line)
	}

	return err
}

// Delete deletes packages
func (p *Pacman) Delete() error {
	pr.Printf("removing package\n")

	cmd := exec.Command(pacmanPath, "--remove", "--noconfirm", pr.Name)
	out, err := cmd.CombinedOutput()

	for _, line := range strings.Split(string(out), "\n") {
		pr.Printf("%s\n", line)
	}

	return err
}

// Update updates packages
func (p *Pacman) Update() error {
	pr.Printf("updating package\n")

	return pr.Create()
}

func init() {
	RegisterProvider("pacman", NewPacman)
}
