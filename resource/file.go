package resource

import (
	"fmt"
	"io"
	"os"
	"os/user"

	"github.com/hashicorp/hcl"
	"github.com/hashicorp/hcl/hcl/ast"
	"github.com/imdario/mergo"
)

// Name and description of the resource
const fileResourceType = "file"
const fileResourceDesc = "manages files"

// FileResource is a resource which manages files
type FileResource struct {
	BaseResource `hcl:",squash"`

	// Path to the file
	Path string `hcl:"path"`

	// Permission bits to set on the file
	Mode int `hcl:"mode"`

	// Owner of the file
	Owner string `hcl:"owner"`

	// Group of the file
	Group string `hcl:"group"`
}

// NewFileResource creates a new resource for managing files
func NewFileResource(title string, obj *ast.ObjectItem) (Resource, error) {
	// Defaults for owner and group
	currentUser, err := user.Current()
	if err != nil {
		return nil, err
	}

	currentGroup, err := user.LookupGroupId(currentUser.Gid)
	if err != nil {
		return nil, err
	}

	defaultOwner := currentUser.Username
	defaultGroup := currentGroup.Name

	// Resource defaults
	defaults := FileResource{
		BaseResource: BaseResource{
			Title: title,
			Type:  fileResourceType,
			State: StatePresent,
		},
		Path:  title,
		Mode:  0644,
		Owner: defaultOwner,
		Group: defaultGroup,
	}

	var fr FileResource
	err = hcl.DecodeObject(&fr, obj)
	if err != nil {
		return nil, err
	}

	// Merge the decoded object with the resource defaults
	err = mergo.Merge(&fr, defaults)

	return &fr, err
}

// Evaluate evaluates the file resource
func (fr *FileResource) Evaluate(w io.Writer, opts *Options) (State, error) {
	resourceState := State{
		Current: StateUnknown,
		Want:    fr.State,
		Update:  false,
	}

	// File does not exist
	fi, err := os.Stat(fr.Path)
	if os.IsNotExist(err) {
		resourceState.Current = StateAbsent

		return resourceState, nil
	}

	// File exists, ensure that it is not a directory
	resourceState.Current = StatePresent
	if fi.IsDir() {
		return resourceState, fmt.Errorf("%s exists, but is not a file", fr.Path)
	}

	// Check permissions
	if fi.Mode().Perm() != os.FileMode(fr.Mode) {
		resourceState.Update = true
	}

	// Check ownership
	owner, err := fileOwner(fi)
	if err != nil {
		return resourceState, err
	}

	if fr.Owner != owner.User || fr.Group != owner.Group {
		resourceState.Update = true
	}

	return resourceState, nil
}

// Create creates the file
func (fr *FileResource) Create(w io.Writer, opts *Options) error {
	fr.Printf(w, "creating file\n")

	if _, err := os.Create(fr.Path); err != nil {
		return err
	}

	if err := setFileOwner(fr.Path, fr.Owner, fr.Group); err != nil {
		return err
	}

	return os.Chmod(fr.Path, os.FileMode(fr.Mode))
}

// Delete deletes the file
func (fr *FileResource) Delete(w io.Writer, opts *Options) error {
	fr.Printf(w, "removing file\n")

	return os.Remove(fr.Path)
}

// Update updates the file
func (fr *FileResource) Update(w io.Writer, opts *Options) error {
	fi, err := os.Stat(fr.Path)
	if err != nil {
		return err
	}

	// Fix permissions if needed
	if fi.Mode().Perm() != os.FileMode(fr.Mode) {
		fr.Printf(w, "setting permissions to %#o\n", fr.Mode)
		if err = os.Chmod(fr.Path, os.FileMode(fr.Mode)); err != nil {
			return err
		}
	}

	// Fix ownership if needed
	owner, err := fileOwner(fi)
	if err != nil {
		return err
	}

	if fr.Owner != owner.User || fr.Group != owner.Group {
		fr.Printf(w, "setting owner %s:%s\n", fr.Owner, fr.Group)
		if err := setFileOwner(fr.Path, fr.Owner, fr.Group); err != nil {
			return err
		}
	}

	return nil
}

func init() {
	item := RegistryItem{
		Name:        fileResourceType,
		Description: fileResourceDesc,
		Provider:    NewFileResource,
	}

	Register(item)
}
