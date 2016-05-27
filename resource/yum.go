// +build linux

package resource

import (
	"os/exec"
	"strings"

	"github.com/hashicorp/hcl"
	"github.com/hashicorp/hcl/hcl/ast"
	"github.com/imdario/mergo"
)

// Path to the Yum package manager
const yumPath = "/usr/bin/yum"

// Name and description of the resource
const yumResourceType = "yum"
const yumResourceDesc = "manages packages using the yum package manager"

// YumResource type represents the resource for
// package management on RHEL/CentOS systems
type YumResource struct {
	BaseResource        `hcl:",squash"`
	BasePackageResource `hcl:",squash"`
}

// NewYumResource creates a new resource for managing packages
// using the yum package manager on RHEL/CentOS systems
func NewYumResource(title string, obj *ast.ObjectItem, config *Config) (Resource, error) {
	// Resource defaults
	defaults := &YumResource{
		BaseResource: BaseResource{
			Title:  title,
			Type:   yumResourceType,
			State:  StatePresent,
			Config: config,
		},
		BasePackageResource: BasePackageResource{
			Name: title,
		},
	}

	// Decode the object from HCL
	var yr YumResource
	err := hcl.DecodeObject(&yr, obj)
	if err != nil {
		return nil, err
	}

	// Merge in the decoded object with the resource defaults
	err = mergo.Merge(&yr, defaults)

	return &yr, err
}

// Evaluate evaluates the state of the package resource
func (yr *YumResource) Evaluate() (State, error) {
	s := State{
		Current: StateUnknown,
		Want:    yr.State,
	}

	_, err := exec.LookPath(yumPath)
	if err != nil {
		return s, err
	}

	cmd := exec.Command(yumPath, "-q", "--noplugins", "list", "installed", yr.Name)
	err = cmd.Run()

	if err != nil {
		s.Current = StateAbsent
	} else {
		s.Current = StatePresent
	}

	return s, nil
}

// Create installs the package managed by the resource
func (yr *YumResource) Create() error {
	yr.Printf("installing package\n")

	cmd := exec.Command(yumPath, "--assumeyes", "install", yr.Name)
	out, err := cmd.CombinedOutput()

	for _, line := range strings.Split(string(out), "\n") {
		yr.Printf("%s\n", line)
	}

	return err
}

// Delete deletes the package managed by the resource
func (yr *YumResource) Delete() error {
	yr.Printf("removing package\n")

	cmd := exec.Command(yumPath, "--assumeyes", "remove", yr.Name)
	out, err := cmd.CombinedOutput()

	for _, line := range strings.Split(string(out), "\n") {
		yr.Printf("%s\n", line)
	}

	return err
}

// Update updates packages
func (yr *YumResource) Update() error {
	// This method is a no-op for now
	//
	// TODO: Be able to handle upgrades/downgrades
}

func init() {
	item := RegistryItem{
		Name:        yumResourceType,
		Description: yumResourceDesc,
		Provider:    NewYumResource,
	}

	Register(item)
}
