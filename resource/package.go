// +build !windows

package resource

import (
	"errors"
	"os/exec"
	"strings"

	"github.com/dnaeon/gru/utils"
)

// ErrNoPackageProviderFound is returned when no suitable provider is found
var ErrNoPackageProviderFound = errors.New("No suitable package provider found")

// BasePackage is the base resource type for package management
// It's purpose is to be embedded into other package resource providers.
type BasePackage struct {
	Base

	// Name of the package to manage. Defaults to the resource name.
	Package string `luar:"-"`

	// Version of the package.
	Version string `luar:"version"`

	// Package manager to use
	manager string `luar:"-"`

	// Arguments to use when quering a package
	queryArgs []string `luar:"-"`

	// Arguments to use when installing a package
	installArgs []string `luar:"-"`

	// Arguments to use when deinstalling a package
	deinstallArgs []string `luar:"-"`

	// Arguments to use when updating a package
	updateArgs []string `luar:"-"`
}

// Evaluate evaluates the state of the package
func (bp *BasePackage) Evaluate() (State, error) {
	s := State{
		Current: StateUnknown,
		Want:    bp.State,
	}

	_, err := exec.LookPath(bp.manager)
	if err != nil {
		return s, err
	}

	bp.queryArgs = append(bp.queryArgs, bp.Package)
	cmd := exec.Command(bp.manager, bp.queryArgs...)
	err = cmd.Run()

	if err != nil {
		s.Current = StateAbsent
	} else {
		s.Current = StatePresent
	}

	return s, nil
}

// Create installs the package
func (bp *BasePackage) Create() error {
	bp.Log("installing package\n")

	bp.installArgs = append(bp.installArgs, bp.Package)
	cmd := exec.Command(bp.manager, bp.installArgs...)
	out, err := cmd.CombinedOutput()

	for _, line := range strings.Split(string(out), "\n") {
		bp.Log("%s\n", line)
	}

	return err
}

// Delete deletes the package
func (bp *BasePackage) Delete() error {
	bp.Log("removing package\n")

	bp.deinstallArgs = append(bp.deinstallArgs, bp.Package)
	cmd := exec.Command(bp.manager, bp.deinstallArgs...)
	out, err := cmd.CombinedOutput()

	for _, line := range strings.Split(string(out), "\n") {
		bp.Log("%s\n", line)
	}

	return err
}

// Update updates the package
func (bp *BasePackage) Update() error {
	bp.Log("updating package\n")

	bp.updateArgs = append(bp.updateArgs, bp.Package)
	cmd := exec.Command(bp.manager, bp.updateArgs...)
	out, err := cmd.CombinedOutput()

	for _, line := range strings.Split(string(out), "\n") {
		bp.Log("%s\n", line)
	}

	return err
}

// NewPackage creates a new resource for managing packages.
// This provider tries to determine the most appropriate
// package provider for you, so it is more like a meta-provider.
//
// Example:
//   pkg = package.new("tmux")
//   pkg.state = "present"
func NewPackage(name string) (Resource, error) {
	// Releases files used by the various GNU/Linux distros
	releases := map[string]Provider{
		"/etc/arch-release":   NewPacman,
		"/etc/centos-release": NewYum,
		"/etc/redhat-release": NewYum,
	}

	// Do our best to determine the proper provider
	for release, provider := range releases {
		dst := utils.NewFileUtil(release)
		if dst.Exists() {
			return provider(name)
		}
	}

	return nil, ErrNoPackageProviderFound
}

// Pacman type represents the resource for package management on
// Arch Linux systems.
//
// Example:
//   pkg = pacman.new("tmux")
//   pkg.state = "present"
type Pacman struct {
	BasePackage
}

// NewPacman creates a new resource for managing packages
// using the pacman package manager on an Arch Linux system
func NewPacman(name string) (Resource, error) {
	p := &Pacman{
		BasePackage: BasePackage{
			Base: Base{
				Name:   name,
				Type:   "pkg",
				State:  StatePresent,
				After:  make([]string, 0),
				Before: make([]string, 0),
			},
			Package:       name,
			Version:       "",
			manager:       "/usr/bin/pacman",
			queryArgs:     []string{"--query"},
			installArgs:   []string{"--sync", "--noconfirm"},
			deinstallArgs: []string{"--remove", "--noconfirm"},
			updateArgs:    []string{"--sync", "--noconfirm"},
		},
	}

	return p, nil
}

// Yum type represents the resource for package management on
// RHEL and CentOS systems.
//
// Example:
//   pkg = yum.new("emacs")
//   pkg.state = "present"
type Yum struct {
	BasePackage
}

// NewYum creates a new resource for managing packages
// using the yum package manager on RHEL and CentOS systems
func NewYum(name string) (Resource, error) {
	y := &Yum{
		BasePackage: BasePackage{
			Base: Base{
				Name:   name,
				Type:   "pkg",
				State:  StatePresent,
				After:  make([]string, 0),
				Before: make([]string, 0),
			},
			Package:       name,
			manager:       "/usr/bin/yum",
			queryArgs:     []string{"-q", "--noplugins", "list", "installed"},
			installArgs:   []string{"--assumeyes", "install"},
			deinstallArgs: []string{"--assumeyes", "remove"},
			updateArgs:    []string{"--assumeyes", "install"},
		},
	}

	return y, nil
}

func init() {
	RegisterProvider("pkg", NewPackage)
	RegisterProvider("yum", NewYum)
	RegisterProvider("pacman", NewPacman)
}
