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

// File types
const fileTypeRegular = "file"
const fileTypeDirectory = "directory"

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

	// The file type we manage
	FileType string `hcl:"type"`
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
			Path:     title,
			Mode:     0644,
			Owner:    defaultOwner,
			Group:    defaultGroup,
			FileType: fileTypeRegular,
		},
	}

	var fr FileResource
	err = hcl.DecodeObject(&fr, obj)
	if err != nil {
		return nil, err
	}

	// Merge the decoded object with the resource defaults
	err = mergo.Merge(&fr, defaults)
	if err != nil {
		return nil, err
	}

	if fr.FileType != fileTypeRegular || fr.FileType != fileTypeDirectory {
		return nil, fmt.Errorf("Unknown file type '%s'", fr.FileType)
	}

	return &fr, nil
}

// Evaluate evaluates the file resource
func (fr *FileResource) Evaluate(w io.Writer, opts *Options) (State, error) {
	rs := State{
		Current: StateUnknown,
		Want:    fr.State,
		Update:  false,
	}

	// The file we manage
	dst := utils.NewFileUtil(fr.Path)

	// File does not exist
	if !dst.Exists() {
		rs.Current = StateAbsent
		return rs, nil
	} else {
		rs.Current = StatePresent
	}

	fi, err := os.Stat(fr.Path)
	if err != nil {
		return rs, err
	}

	// Check the target file we manage
	switch fr.FileType {
	case fileTypeRegular:
		if !fi.Mode().IsRegular() {
			return rs, fmt.Errorf("%s exists, but is not a regular file", fr.Path)
		}
	case fileTypeDirectory:
		if !fi.IsDir() {
			return rs, fmt.Errorf("%s exists, but is not a directory", fr.Path)
		}
	}

	// Check file content
	if fr.Source != "" {
		srcPath := filepath.Join(opts.SiteDir, "data", fr.Source)
		same, err := dst.SameContentWith(srcPath)
		if err != nil {
			return rs, err
		}
		if !same {
			fr.Printf(w, "content is out of date\n")
			rs.Update = true
		}
	}

	// Check file permissions
	mode, err := dst.Mode()
	if err != nil {
		return rs, err
	}

	if mode.Perm() != os.FileMode(fr.Mode) {
		fr.Printf(w, "permissions are out of date\n")
		rs.Update = true
	}

	// Check ownership
	owner, err := dst.Owner()
	if err != nil {
		return rs, err
	}

	if fr.Owner != owner.User.Username || fr.Group != owner.Group.Name {
		fr.Printf(w, "owner is out of date\n")
		rs.Update = true
	}

	return rs, nil
}

// createRegularFile creates the file and content managed by the resource
func (fr *FileResource) createRegularFile(opts *Options) error {
	dst := utils.NewFileUtil(fr.Path)

	switch {
	case fr.Source == "" && dst.Exists():
		// We have no source, do nothing
		break
	case fr.Source == "" && !dst.Exists():
		// Create an empty file
		if _, err := os.Create(fr.Path); err != nil {
			return err
		}
	case fr.Source != "" && dst.Exists():
		// File exists and we have a source file
		srcPath := filepath.Join(opts.SiteDir, "data", fr.Source)
		if err := dst.CopyFrom(srcPath); err != nil {
			return err
		}
	}

	return nil
}

// createDirectory creates the directory and content managed by the resource
func (fr *FileResource) createDirectory(opts *Options) error {
	if err := os.Mkdir(fr.Path, os.FileMode(fr.Mode)); err != nil {
		return err
	}

	return nil
}

// Create creates the file managed by the resource
func (fr *FileResource) Create(w io.Writer, opts *Options) error {
	dst := utils.NewFileUtil(fr.Path)
	fr.Printf(w, "creating resource\n")

	// Set content
	switch fr.FileType {
	case fileTypeRegular:
		if err := fr.createRegularFile(opts); err != nil {
			return err
		}
	case fileTypeDirectory:
		if err := fr.createDirectory(opts); err != nil {
			return err
		}
	}

	// Set file owner
	if err := dst.SetOwner(fr.Owner, fr.Group); err != nil {
		return err
	}

	// Set file permissions
	return dst.Chmod(os.FileMode(fr.Mode))
}

// Delete deletes the file
func (fr *FileResource) Delete(w io.Writer, opts *Options) error {
	fr.Printf(w, "removing file\n")
	dst := utils.NewFileUtil(fr.Path)

	return dst.Remove()
}

// Update updates the file managed by the resource
func (fr *FileResource) Update(w io.Writer, opts *Options) error {
	dst := utils.NewFileUtil(fr.Path)

	// Update file content if needed
	if fr.Source != "" {
		srcPath := filepath.Join(opts.SiteDir, "data", fr.Source)
		same, err := dst.SameContentWith(srcPath)
		if err != nil {
			return err
		}

		if !same {
			srcFile := utils.NewFileUtil(srcPath)
			srcMd5, err := srcFile.Md5()
			if err != nil {
				return err
			}

			fr.Printf(w, "updating content to md5:%s\n", srcMd5)
			if err := dst.CopyFrom(srcPath); err != nil {
				return err
			}
		}
	}

	// Fix permissions if needed
	mode, err := dst.Mode()
	if err != nil {
		return err
	}

	if mode.Perm() != os.FileMode(fr.Mode) {
		fr.Printf(w, "setting permissions to %#o\n", fr.Mode)
		dst.Chmod(os.FileMode(fr.Mode))
	}

	// Fix ownership if needed
	owner, err := dst.Owner()
	if err != nil {
		return err
	}

	if fr.Owner != owner.User.Username || fr.Group != owner.Group.Name {
		fr.Printf(w, "setting owner to %s:%s\n", fr.Owner, fr.Group)
		if err := dst.SetOwner(fr.Owner, fr.Group); err != nil {
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
