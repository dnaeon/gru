package resource

import (
	"fmt"
	"io"
	"os"
	"os/user"
	"path/filepath"
	"strings"

	"github.com/dnaeon/gru/utils"
	"github.com/hashicorp/hcl"
	"github.com/hashicorp/hcl/hcl/ast"
	"github.com/imdario/mergo"
)

// Name and description of the resource
const fileResourceType = "file"
const fileResourceDesc = "manages files and directories"

// The file types we manage
const fileTypeRegular = "regular"
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
	FileType string `hcl:"filetype"`

	// Recursively manage the directory
	Recursive bool `hcl:"recursive"`
}

// outdatedFile type is used to describe a file which
// has been identified as being out of date
type outdatedFile struct {
	// Source file to use when reconstructing the content
	src string

	// The destination file which is identified as being out of date
	dst string

	// Flag to indicate that the content is out of date
	isContentOutdated bool

	// Flag to indicate that the permissions are out of date
	isPermissionsOutdated bool

	// Flag to indicate that the owner is out of date
	isOwnerOutdated bool

	// Flag to indicate that file exists in destination, but not in source
	isExtra bool
}

// FileResource is a resource which manages files and directories
type FileResource struct {
	BaseResource     `hcl:",squash"`
	BaseFileResource `hcl:",squash"`

	// Files identified as being out of date
	outdated []*outdatedFile
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

	// Resource defaults
	defaults := FileResource{
		BaseResource: BaseResource{
			Title: title,
			Type:  fileResourceType,
			State: StatePresent,
		},
		BaseFileResource: BaseFileResource{
			Path:      title,
			Mode:      0644,
			Owner:     currentUser.Username,
			Group:     currentGroup.Name,
			FileType:  fileTypeRegular,
			Recursive: false,
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

	// Check that the given file type is a valid one
	if fr.FileType != fileTypeRegular && fr.FileType != fileTypeDirectory {
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

	// Check for file presence first
	fi, err := os.Stat(fr.Path)
	if os.IsNotExist(err) {
		rs.Current = StateAbsent
		return rs, nil
	} else {
		rs.Current = StatePresent
	}

	// Check the file(s) content, permissions and ownership
	switch fr.FileType {
	case fileTypeRegular:
		if !fi.Mode().IsRegular() {
			return rs, fmt.Errorf("%s exists, but is not a regular file", fr.Path)
		}

		outdated, err := fr.isRegularFileContentOutdated(opts)
		if err != nil {
			return rs, err
		}

		if outdated {
			rs.Update = true
		}
	case fileTypeDirectory:
		if !fi.IsDir() {
			return rs, fmt.Errorf("%s exists, but is not a directory", fr.Path)
		}

		outdated, err := fr.isDirectoryContentOutdated(opts)
		if err != nil {
			return rs, err
		}

		if outdated {
			rs.Update = true
		}
	}

	outdated, err := fr.isPermissionsOutdated()
	if err != nil {
		return rs, err
	}

	if outdated {
		rs.Update = true
	}

	outdated, err = fr.isOwnerOutdated()
	if err != nil {
		return rs, err
	}

	if outdated {
		rs.Update = true
	}

	// Report on what has been identified as being out of date
	for _, item := range fr.outdated {
		// TODO: Report extra files

		if item.isContentOutdated && !item.isExtra {
			fr.Printf(w, "content of %s is out of date\n", item.dst)
		}
		if item.isPermissionsOutdated {
			fr.Printf(w, "permissions of %s are out of date\n", item.dst)
		}
		if item.isOwnerOutdated {
			fr.Printf(w, "owner of %s is out of date\n", item.dst)
		}
	}

	return rs, nil
}

// Create creates the file managed by the resource
func (fr *FileResource) Create(w io.Writer, opts *Options) error {
	fr.Printf(w, "creating resource\n")

	switch fr.FileType {
	case fileTypeRegular:
		if err := fr.createRegularFile(opts); err != nil {
			return err
		}

		dst := utils.NewFileUtil(fr.Path)
		if err := dst.Chmod(os.FileMode(fr.Mode)); err != nil {
			return err
		}
		if err := dst.SetOwner(fr.Owner, fr.Group); err != nil {
			return err
		}
	case fileTypeDirectory:
		if err := fr.createDirectory(opts); err != nil {
			return err
		}
		dstRegistry, err := directoryFileRegistry(fr.Path, []string{})
		if err != nil {
			return err
		}

		for _, path := range dstRegistry {
			dst := utils.NewFileUtil(path)
			if err := dst.Chmod(os.FileMode(fr.Mode)); err != nil {
				return err
			}
			if err := dst.SetOwner(fr.Owner, fr.Group); err != nil {
				return err
			}
		}
	}

	return nil
}

// Delete deletes the file managed by the resource
func (fr *FileResource) Delete(w io.Writer, opts *Options) error {
	fr.Printf(w, "removing resource\n")

	if fr.Recursive {
		return os.RemoveAll(fr.Path)
	}

	return os.Remove(fr.Path)
}

// Update updates the file managed by the resource
func (fr *FileResource) Update(w io.Writer, opts *Options) error {
	for _, item := range fr.outdated {
		dstFile := utils.NewFileUtil(item.dst)

		// TODO: Purge extra files

		// Update file content if needed
		if item.isContentOutdated {
			// Create parent directory for file if missing
			dstDir := filepath.Dir(item.dst)
			_, err := os.Stat(dstDir)
			if os.IsNotExist(err) {
				if err := os.MkdirAll(dstDir, 0755); err != nil {
					return err
				}
			}

			srcFile := utils.NewFileUtil(item.src)
			srcMd5, err := srcFile.Md5()
			if err != nil {
				return err
			}

			fr.Printf(w, "setting content of %s to md5:%s\n", item.dst, srcMd5)
			if err := dstFile.CopyFrom(item.src, true); err != nil {
				return err
			}
		}

		// Update permissions if needed
		if item.isPermissionsOutdated {
			fr.Printf(w, "setting permissions of %s to %#o\n", item.dst, fr.Mode)
			if err := dstFile.Chmod(os.FileMode(fr.Mode)); err != nil {
				return err
			}
		}

		// Update ownership if needed
		if item.isOwnerOutdated {
			fr.Printf(w, "setting owner of %s to %s:%s\n", item.dst, fr.Owner, fr.Group)
			if err := dstFile.SetOwner(fr.Owner, fr.Group); err != nil {
				return err
			}
		}
	}

	return nil
}

// directoryFileRegistry creates a map of all files found in a
// given directory. The keys of the map are the file names with the
// leading source path trimmed and the values are the
// full path to the discovered files.
func directoryFileRegistry(path string, skip []string) (map[string]string, error) {
	registry := make(map[string]string)

	found, err := utils.WalkPath(path, skip)
	if err != nil {
		return registry, err
	}

	for _, name := range found {
		fi, err := os.Stat(name)
		if err != nil {
			return registry, err
		}

		if fi.Mode().IsRegular() {
			trimmed := strings.TrimPrefix(name, path+"/")
			registry[trimmed] = name
		}
	}

	return registry, nil
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
		srcPath := filepath.Join(opts.SiteDir, fr.Source)
		if err := dst.CopyFrom(srcPath, false); err != nil {
			return err
		}
	}

	return nil
}

// createDirectory creates the directory and content managed by the resource
func (fr *FileResource) createDirectory(opts *Options) error {
	switch {
	case !fr.Recursive:
		return os.Mkdir(fr.Path, 0755)
	case fr.Recursive && fr.Source != "":
		srcPath := filepath.Join(opts.SiteDir, fr.Source)
		return utils.CopyDir(srcPath, fr.Path)
	case fr.Recursive && fr.Source == "":
		return os.MkdirAll(fr.Path, 0755)
	}

	// Not reached
	return nil
}

// isRegularFileContentOutdated returns a boolean indicating whether the
// content managed by the resource is outdated compared to the source
// file defined by the resource.
// If the file is identified as being out of date it will be appended to the
// list of outdated files for the resource, so it can be further
// processed if needed.
func (fr *FileResource) isRegularFileContentOutdated(opts *Options) (bool, error) {
	if fr.Source != "" {
		srcPath := filepath.Join(opts.SiteDir, fr.Source)
		same, err := utils.SameContent(srcPath, fr.Path)
		if err != nil {
			return false, err
		}

		if !same {
			item := &outdatedFile{
				src:               srcPath,
				dst:               fr.Path,
				isContentOutdated: true,
			}
			fr.outdated = append(fr.outdated, item)
			return true, nil
		}
	}

	return false, nil
}

// isDirectoryContentOutdated returns a boolean indicating whether the
// content of the directory managed by the resource is outdated
// compared to the source directory defined by the resource.
// The files identified as being out of date will be appended to the
// list of outdated files for the resource, so they can be further
// processed if needed.
func (fr *FileResource) isDirectoryContentOutdated(opts *Options) (bool, error) {
	isOutdated := false
	if fr.Source != "" && fr.Recursive {
		srcPath := filepath.Join(opts.SiteDir, fr.Source)

		// Exclude the ".git" repo directory from the source path,
		// since our source files reside in a git repo
		srcRegistry, err := directoryFileRegistry(srcPath, []string{".git"})
		if err != nil {
			return false, err
		}

		dstRegistry, err := directoryFileRegistry(fr.Path, []string{})
		if err != nil {
			return false, err
		}

		// Check source and destination files' content
		for name := range srcRegistry {
			item := &outdatedFile{
				src:               srcRegistry[name],
				dst:               dstRegistry[name],
				isContentOutdated: true,
			}

			// File is missing
			if _, ok := dstRegistry[name]; !ok {
				item.dst = filepath.Join(fr.Path, name)
				fr.outdated = append(fr.outdated, item)
				isOutdated = true
				continue
			}

			// Check if content has changed
			same, err := utils.SameContent(srcRegistry[name], dstRegistry[name])
			if err != nil {
				return false, err
			}

			if !same {
				fr.outdated = append(fr.outdated, item)
				isOutdated = true
			}
		}

		// Check for extra files in the managed directory
		for name := range dstRegistry {
			item := &outdatedFile{
				dst:     dstRegistry[name],
				isExtra: true,
			}
			if _, ok := srcRegistry[name]; !ok {
				fr.outdated = append(fr.outdated, item)
				isOutdated = true
			}
		}
	}

	return isOutdated, nil
}

// isPermissionsOutdated returns a boolean indicating whether the
// file's permissions managed by the resource are outdated compared
// to the ones defined by the resource.
// Each file identified as being out of date will be appended to the
// list of outdated files for the resource, so they can be further
// processed if needed.
func (fr *FileResource) isPermissionsOutdated() (bool, error) {
	dstRegistry, err := directoryFileRegistry(fr.Path, []string{})
	if err != nil {
		return false, err
	}

	isOutdated := false
	for name := range dstRegistry {
		item := &outdatedFile{
			dst: dstRegistry[name],
			isPermissionsOutdated: true,
		}

		dst := utils.NewFileUtil(dstRegistry[name])
		mode, err := dst.Mode()
		if err != nil {
			return false, err
		}

		if mode.Perm() != os.FileMode(fr.Mode) {
			fr.outdated = append(fr.outdated, item)
			isOutdated = true
		}
	}

	return isOutdated, nil
}

// isOwnerOutdated returns a boolean indicating whether the
// file's owner managed by the resource is outdated compared to the
// ones defined by the resource.
// Each file identified as being out of date will be appended to the
// list of outdated files for the resource, so they can be further
// processed if needed.
func (fr *FileResource) isOwnerOutdated() (bool, error) {
	dstRegistry, err := directoryFileRegistry(fr.Path, []string{})
	if err != nil {
		return false, err
	}

	isOutdated := false
	for name := range dstRegistry {
		item := &outdatedFile{
			dst:             dstRegistry[name],
			isOwnerOutdated: true,
		}

		dst := utils.NewFileUtil(dstRegistry[name])
		owner, err := dst.Owner()
		if err != nil {
			return false, err
		}

		if fr.Owner != owner.User.Username || fr.Group != owner.Group.Name {
			fr.outdated = append(fr.outdated, item)
			isOutdated = true
		}
	}

	return isOutdated, nil
}

func init() {
	item := RegistryItem{
		Name:        fileResourceType,
		Description: fileResourceDesc,
		Provider:    NewFileResource,
	}

	Register(item)
}
