package resource

import (
	"fmt"
	"io"
	"os"
	"os/user"
	"strconv"
	"syscall"

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

	// Owner of the directory
	Owner string `hcl:"owner" json:"owner"`

	// Group of the directory
	Group string `hcl:"group" json:"group"`
}

// NewDirectoryResource creates a new resource for managing directories
func NewDirectoryResource(title string, obj *ast.ObjectItem) (Resource, error) {
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
	defaults := DirectoryResource{
		BaseResource: BaseResource{
			Title: title,
			Type:  dirResourceType,
			State: StatePresent,
		},
		Name:  title,
		Mode:  0755,
		Owner: defaultOwner,
		Group: defaultGroup,
	}

	var d DirectoryResource
	err = hcl.DecodeObject(&d, obj)
	if err != nil {
		return nil, err
	}

	// Merge the decoded object with the resource defaults
	err = mergo.Merge(&d, defaults)

	return &d, err
}

// Evaluate evaluates the directory resource
func (d *DirectoryResource) Evaluate(w io.Writer, opts *Options) (State, error) {
	resourceState := State{
		Current: StateUnknown,
		Want:    d.State,
		Update:  false,
	}

	// Directory does not exist
	fi, err := os.Stat(d.Name)
	if os.IsNotExist(err) {
		resourceState.Current = StateAbsent
		resourceState.Update = true

		return resourceState, nil
	}

	// File exists, ensure that it is a directory
	resourceState.Current = StatePresent
	if !fi.IsDir() {
		return resourceState, fmt.Errorf("%s exists, but is not a directory", d.Name)
	}

	// Check permissions
	if fi.Mode().Perm() != os.FileMode(d.Mode) {
		resourceState.Update = true
	}

	// Check ownership
	owner, err := user.Lookup(d.Owner)
	if err != nil {
		return resourceState, err
	}

	group, err := user.LookupGroup(d.Group)
	if err != nil {
		return resourceState, err
	}

	uid, _ := strconv.Atoi(owner.Uid)
	gid, _ := strconv.Atoi(group.Gid)

	if uid != int(fi.Sys().(*syscall.Stat_t).Uid) {
		resourceState.Update = true
	}

	if gid != int(fi.Sys().(*syscall.Stat_t).Gid) {
		resourceState.Update = true
	}

	return resourceState, nil
}

// Create creates the directory
func (d *DirectoryResource) Create(w io.Writer, opts *Options) error {
	// TODO: Create parent directories if needed

	d.Printf(w, "creating directory\n")
	err := os.Mkdir(d.Name, os.FileMode(d.Mode))

	return err
}

// Delete deletes the directory
func (d *DirectoryResource) Delete(w io.Writer, opts *Options) error {
	// TODO: Recursively remove directory if needed

	d.Printf(w, "removing directory\n")
	err := os.Remove(d.Name)

	return err
}

// Update updates the permission bits of the directory
func (d *DirectoryResource) Update(w io.Writer, opts *Options) error {
	fi, err := os.Stat(d.Name)
	if err != nil {
		return err
	}

	// Fix permissions
	//
	// TODO: Might want to move this to a util function, so we can
	//       reuse it when using this resource in recursive mode
	//
	if fi.Mode().Perm() != os.FileMode(d.Mode) {
		d.Printf(w, "setting permissions to %#o\n", d.Mode)
		if err = os.Chmod(d.Name, os.FileMode(d.Mode)); err != nil {
			return err
		}
	}

	owner, err := user.Lookup(d.Owner)
	if err != nil {
		return err
	}

	group, err := user.LookupGroup(d.Group)
	if err != nil {
		return err
	}

	uid, _ := strconv.Atoi(owner.Uid)
	gid, _ := strconv.Atoi(group.Gid)

	if uid != int(fi.Sys().(*syscall.Stat_t).Uid) || gid != int(fi.Sys().(*syscall.Stat_t).Gid) {
		d.Printf(w, "setting owner to %s:%s\n", d.Owner, d.Group)
		if err = os.Chown(d.Name, uid, gid); err != nil {
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
