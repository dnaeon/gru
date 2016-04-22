package resource

import (
	"fmt"
	"io"
	"os"

	"github.com/hashicorp/hcl"
	"github.com/hashicorp/hcl/hcl/ast"
	"github.com/imdario/mergo"
)

// Name and description of the resource
const dirResourceType = "directory"
const dirResourceDesc = "manages directories"

// DirectoryResource is a resource which manages directories
type DirectoryResource struct {
	BaseResource `hcl:",squash"`

	// Name of the directory
	Name string `hcl:"name" json:"name"`

	// Permission bits to set on the directory
	Mode int `hcl:"mode" json:"mode"`
}

// NewDirectoryResource creates a new resource for managing directories
func NewDirectoryResource(title string, obj *ast.ObjectItem) (Resource, error) {
	// Resource defaults
	defaults := DirectoryResource{
		BaseResource: BaseResource{
			Title: title,
			Type:  dirResourceType,
			State: StatePresent,
		},
		Name: title,
		Mode: 0755,
	}

	var d DirectoryResource
	err := hcl.DecodeObject(&d, obj)
	if err != nil {
		return nil, err
	}

	// Merge the decoded object with the resource defaults
	err = mergo.Merge(&d, defaults)

	return &d, err
}

func (d *DirectoryResource) Evaluate() (State, error) {
	resourceState := State{
		Current: StateUnknown,
		Want:    d.State,
		Update:  false,
	}

	// Directory does not exist
	f, err := os.Stat(d.Name)
	if os.IsNotExist(err) {
		resourceState.Current = StateAbsent

		return resourceState, nil
	}

	// File exists, ensure that it is a directory
	resourceState.Current = StatePresent
	if !f.IsDir() {
		return resourceState, fmt.Errorf("%s exists, but is not a directory", d.Name)
	}

	// Check permissions
	if f.Mode().Perm() != os.FileMode(d.Mode) {
		resourceState.Update = true
	}

	return resourceState, nil
}

// Create creates a directory
func (d *DirectoryResource) Create(w io.Writer) error {
	// TODO: Create parent directories if needed

	d.Printf(w, "creating directory\n")
	err := os.Mkdir(d.Name, os.FileMode(d.Mode))

	return err
}

// Delete deletes the directory
func (d *DirectoryResource) Delete(w io.Writer) error {
	// TODO: Recursively remove directory if needed

	d.Printf(w, "removing directory\n")
	err := os.Remove(d.Name)

	return err
}

// Update updates the permission bits of the directory
func (d *DirectoryResource) Update(w io.Writer) error {
	f, err := os.Stat(d.Name)
	if err != nil {
		return err
	}

	if f.Mode().Perm() != os.FileMode(d.Mode) {
		d.Printf(w, "setting permissions to %#o\n", d.Mode)
		if err = os.Chmod(d.Name, os.FileMode(d.Mode)); err != nil {
			return err
		}
	}

	return nil
}

func init() {
	item := RegistryItem{
		Name:        dirResourceType,
		Description: dirResourceDesc,
		Provider:    NewDirectoryResource,
	}

	Register(item)
}
