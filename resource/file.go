package resource

import (
	"fmt"
	"io"
	"os"
	"os/user"
	"path/filepath"

	"github.com/dnaeon/gru/utils"
	"github.com/hashicorp/hcl"
	"github.com/hashicorp/hcl/hcl/ast"
	"github.com/imdario/mergo"
)

// Name and description of the resource
const fileResourceType = "file"
const fileResourceDesc = "manages files"

// BaseFileResource is the base resource for managing files
type BaseFileResource struct {
	// Path to the file
	Path string `hcl:"path"`

	// Permission bits to set on the file
	Mode int `hcl:"mode"`

	// Owner of the file
	Owner string `hcl:"owner"`

	// Group of the file
	Group string `hcl:"group"`

	// Source file to use when creating/updating the file
	Source string `hcl:"source"`

	// The destination file we manage
	dstFile *utils.FileUtil
}

// permissionsChanged returns a boolean indicating whether the
// permissions of the file managed by the resource is different than the
// permissions defined by the resource
func (bfr *BaseFileResource) permissionsChanged() (bool, error) {
	m, err := bfr.dstFile.Mode()
	if err != nil {
		return false, err
	}

	return m.Perm() != os.FileMode(bfr.Mode), nil
}

// ownerChanged returns a boolean indicating whether the
// owner/group of the file managed by the resource is different than the
// owner/group defined by the resource
func (bfr *BaseFileResource) ownerChanged() (bool, error) {
	owner, err := bfr.dstFile.Owner()
	if err != nil {
		return false, err
	}

	if bfr.Owner != owner.User.Username || bfr.Group != owner.Group.Name {
		return true, nil
	}

	return false, nil
}

// contentChanged returns a boolean indicating whether the
// content of the file managed by the resource is different than the
// content of the source file defined by the resource
func (bfr *BaseFileResource) contentChanged(siteDir string) (bool, error) {
	if bfr.Source == "" {
		return false, nil
	}

	// Source file is expected to be found in the site directory
	srcPath := filepath.Join(siteDir, "data", bfr.Source)
	srcFile := utils.NewFileUtil(srcPath)

	srcMd5, err := srcFile.Md5()
	if err != nil {
		return false, err
	}

	dstMd5, err := bfr.dstFile.Md5()
	if err != nil {
		return false, err
	}

	return srcMd5 != dstMd5, nil
}

// FileResource is a resource which manages files
type FileResource struct {
	BaseResource     `hcl:",squash"`
	BaseFileResource `hcl:",squash"`
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
		BaseFileResource: BaseFileResource{
			Path:  title,
			Mode:  0644,
			Owner: defaultOwner,
			Group: defaultGroup,
		},
	}

	var fr FileResource
	err = hcl.DecodeObject(&fr, obj)
	if err != nil {
		return nil, err
	}

	// Merge the decoded object with the resource defaults
	err = mergo.Merge(&fr, defaults)

	// Set the file utility for this file
	fr.dstFile = utils.NewFileUtil(fr.Path)

	return &fr, err
}

// Evaluate evaluates the file resource
func (fr *FileResource) Evaluate(w io.Writer, opts *Options) (State, error) {
	rs := State{
		Current: StateUnknown,
		Want:    fr.State,
		Update:  false,
	}

	// File does not exist
	fi, err := os.Stat(fr.Path)
	if os.IsNotExist(err) {
		rs.Current = StateAbsent
		return rs, nil
	}

	// Ensure that the file we manage is a regular file
	rs.Current = StatePresent
	if !fi.Mode().IsRegular() {
		return rs, fmt.Errorf("%s is not a regular file", fr.Path)
	}

	// Check file content
	changed, err := fr.contentChanged(opts.SiteDir)
	if err != nil {
		return rs, err
	}

	if changed {
		fr.Printf(w, "content is out of date\n")
		rs.Update = true
	}

	// Check file permissions
	changed, err = fr.permissionsChanged()
	if err != nil {
		return rs, err
	}

	if changed {
		fr.Printf(w, "permissions are out of date\n")
		rs.Update = true
	}

	// Check ownership
	changed, err = fr.ownerChanged()
	if err != nil {
		return rs, err
	}

	if changed {
		fr.Printf(w, "owner is out of date\n")
		rs.Update = true
	}

	return rs, nil
}

// Create creates the file
func (fr *FileResource) Create(w io.Writer, opts *Options) error {
	fr.Printf(w, "creating file\n")

	// Set file content
	switch {
	case fr.Source == "" && fr.dstFile.Exists():
		// Do nothing
		break
	case fr.Source == "" && !fr.dstFile.Exists():
		// Create an empty file
		if _, err := os.Create(fr.Path); err != nil {
			return err
		}
	case fr.Source != "" && fr.dstFile.Exists():
		// File exists and we have a source file
		srcPath := filepath.Join(opts.SiteDir, "data", fr.Source)
		if err := fr.dstFile.CopyFrom(srcPath); err != nil {
			return err
		}
	}

	// Set file owner
	if err := fr.dstFile.SetOwner(fr.Owner, fr.Group); err != nil {
		return err
	}

	// Set file permissions
	return fr.dstFile.Chmod(os.FileMode(fr.Mode))
}

// Delete deletes the file
func (fr *FileResource) Delete(w io.Writer, opts *Options) error {
	fr.Printf(w, "removing file\n")

	return fr.dstFile.Remove()
}

// Update updates the file managed by the resource
func (fr *FileResource) Update(w io.Writer, opts *Options) error {
	// Update file content if needed
	changed, err := fr.contentChanged(opts.SiteDir)
	if err != nil {
		return err
	}

	if changed {
		srcFile := utils.NewFileUtil(filepath.Join(opts.SiteDir, "data", fr.Source))
		srcMd5, err := srcFile.Md5()
		if err != nil {
			return err
		}

		fr.Printf(w, "updating content to md5:%s\n", srcMd5)
		if err := fr.dstFile.CopyFrom(srcFile.Path); err != nil {
			return err
		}
	}

	// Fix permissions if needed
	changed, err = fr.permissionsChanged()
	if err != nil {
		return err
	}

	if changed {
		fr.Printf(w, "setting permissions to %#o\n", fr.Mode)
		fr.dstFile.Chmod(os.FileMode(fr.Mode))
	}

	// Fix ownership if needed
	changed, err = fr.ownerChanged()
	if err != nil {
		return err
	}

	if changed {
		owner, err := fr.dstFile.Owner()
		if err != nil {
			return err
		}

		fr.Printf(w, "setting owner to %s:%s\n", owner.User.Username, owner.Group.Name)
		if err := fr.dstFile.SetOwner(fr.Owner, fr.Group); err != nil {
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
