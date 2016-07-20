package resource

import (
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"strings"

	"github.com/dnaeon/gru/utils"
)

// The file types we manage
const (
	fileTypeRegular   = "regular"
	fileTypeDirectory = "directory"
)

// Flags used to indicate what has changed on a file
const (
	flagOutdatedContent     = 0x01
	flagOutdatedPermissions = 0x02
	flagOutdatedOwner       = 0x04
)

// outdatedFile type describes a file which
// has been identified as being out of date
type outdatedFile struct {
	// Source file to use when reconstructing file's content
	src string

	// Destination file which is identified as being out of date
	dst string

	// Flags used to indicate what has changed on the file
	flags int
}

// File resource manages files and directories.
//
// Example:
//   foo = file.new("/tmp/foo")
//   foo.state = "present"
//   foo.mode = 0600
//
// Example:
//   bar = file.new("/tmp/bar")
//   bar.state = "present"
//   bar.filetype = "directory"
type File struct {
	Base

	// Path to the file. Defaults to the resource name.
	Path string `luar:"-"`

	// Permission bits to set on the file. Defaults to 0644.
	Mode os.FileMode `luar:"mode"`

	// Owner of the file. Defaults to the currently running user.
	Owner string `luar:"owner"`

	// Group of the file.
	// Defaults to the group of the currently running user.
	Group string `luar:"group"`

	// Source file to use when creating/updating the file
	Source string `luar:"source"`

	// The file type we manage.
	FileType string `luar:"filetype"`

	// Recursively manage the directory if set to true.
	// Defaults to false.
	Recursive bool `luar:"recursive"`

	// Purge extra files in the target directory if set to true.
	// Defaults to false.
	Purge bool `luar:"purge"`

	// Files identified as being out of date
	outdated []*outdatedFile `luar:"-"`

	// Extra files found in the target directory
	extra map[string]struct{} `luar:"-"`
}

// NewFile creates a resource for managing files and directories
func NewFile(name string) (Resource, error) {
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
	f := &File{
		Base: Base{
			Name:          name,
			Type:          "file",
			State:         "present",
			Before:        make([]string, 0),
			After:         make([]string, 0),
			PresentStates: []string{"present"},
			AbsentStates:  []string{"absent"},
		},
		Path:      name,
		Mode:      0644,
		Owner:     currentUser.Username,
		Group:     currentGroup.Name,
		FileType:  fileTypeRegular,
		Recursive: false,
		Purge:     false,
		outdated:  make([]*outdatedFile, 0),
		extra:     make(map[string]struct{}),
	}

	return f, nil
}

// Validate validates the resource
func (f *File) Validate() error {
	if err := f.Base.Validate(); err != nil {
		return err
	}

	if f.FileType != fileTypeRegular && f.FileType != fileTypeDirectory {
		return fmt.Errorf("Invalid file type '%s'", f.FileType)
	}
}

// Evaluate evaluates the file resource
func (f *File) Evaluate() (State, error) {
	s := State{
		Current:  "unknown",
		Want:     f.State,
		Outdated: false,
	}

	// Check for file presence
	fi, err := os.Stat(f.Path)
	if os.IsNotExist(err) {
		s.Current = "absent"
		return s, nil
	}

	s.Current = "present"

	// If we have a source, ensure that it exists
	if f.Source != "" {
		dst := utils.NewFileUtil(filepath.Join(DefaultConfig.SiteRepo, f.Source))
		if !dst.Exists() {
			return s, fmt.Errorf("source %s does not exist", f.Source)
		}
	}

	// Check the file(s) content, permissions and ownership
	switch f.FileType {
	case fileTypeRegular:
		if !fi.Mode().IsRegular() {
			return s, fmt.Errorf("%s exists, but is not a regular file", f.Path)
		}

		outdated, err := f.isRegularFileContentOutdated()
		if err != nil {
			return s, err
		}

		if outdated {
			s.Outdated = true
		}
	case fileTypeDirectory:
		if !fi.IsDir() {
			return s, fmt.Errorf("%s exists, but is not a directory", f.Path)
		}

		outdated, err := f.isDirectoryContentOutdated()
		if err != nil {
			return s, err
		}

		if outdated {
			s.Outdated = true
		}
	}

	outdated, err := f.isPermissionsOutdated()
	if err != nil {
		return s, err
	}

	if outdated {
		s.Outdated = true
	}

	outdated, err = f.isOwnerOutdated()
	if err != nil {
		return s, err
	}

	if outdated {
		s.Outdated = true
	}

	// Report on what has been identified as being out of date
	if f.Purge {
		for name := range f.extra {
			f.Log("%s exists, but is not part of source\n", name)
			s.Outdated = true
		}
	}

	for _, item := range f.outdated {
		if item.flags&flagOutdatedContent != 0 {
			f.Log("content of %s is out of date\n", item.dst)
		}
		if item.flags&flagOutdatedPermissions != 0 {
			f.Log("permissions of %s are out of date\n", item.dst)
		}
		if item.flags&flagOutdatedOwner != 0 {
			f.Log("owner of %s is out of date\n", item.dst)
		}
	}

	return s, nil
}

// Create creates the file managed by the resource
func (f *File) Create() error {
	f.Log("creating resource\n")

	switch f.FileType {
	case fileTypeRegular:
		if err := f.createRegularFile(); err != nil {
			return err
		}

		dst := utils.NewFileUtil(f.Path)
		if err := dst.Chmod(f.Mode); err != nil {
			return err
		}
		if err := dst.SetOwner(f.Owner, f.Group); err != nil {
			return err
		}
	case fileTypeDirectory:
		if err := f.createDirectory(); err != nil {
			return err
		}

		dstRegistry, err := directoryFileRegistry(f.Path, []string{})
		if err != nil {
			return err
		}

		for _, path := range dstRegistry {
			dst := utils.NewFileUtil(path)
			if err := dst.Chmod(f.Mode); err != nil {
				return err
			}
			if err := dst.SetOwner(f.Owner, f.Group); err != nil {
				return err
			}
		}
	}

	return nil
}

// Delete deletes the file managed by the resource
func (f *File) Delete() error {
	f.Log("removing resource\n")

	if f.Recursive {
		return os.RemoveAll(f.Path)
	}

	return os.Remove(f.Path)
}

// Update updates the files managed by the resource
func (f *File) Update() error {
	// Purge extra files
	if f.Purge {
		for name := range f.extra {
			dstFile := utils.NewFileUtil(name)
			f.Log("purging %s\n", name)
			if err := dstFile.Remove(); err != nil {
				return err
			}
		}
	}

	// Fix outdated files
	for _, item := range f.outdated {
		dstFile := utils.NewFileUtil(item.dst)

		// Update file content if needed
		if item.flags&flagOutdatedContent != 0 {
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

			f.Log("setting content of %s to md5:%s\n", item.dst, srcMd5)
			if err := dstFile.CopyFrom(item.src, true); err != nil {
				return err
			}
		}

		// Update permissions if needed
		if item.flags&flagOutdatedPermissions != 0 {
			f.Log("setting permissions of %s to %#o\n", item.dst, f.Mode)
			if err := dstFile.Chmod(f.Mode); err != nil {
				return err
			}
		}

		// Update ownership if needed
		if item.flags&flagOutdatedOwner != 0 {
			f.Log("setting owner of %s to %s:%s\n", item.dst, f.Owner, f.Group)
			if err := dstFile.SetOwner(f.Owner, f.Group); err != nil {
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
func (f *File) createRegularFile() error {
	dst := utils.NewFileUtil(f.Path)

	switch {
	case f.Source != "":
		// We have a source file, use it
		srcPath := filepath.Join(DefaultConfig.SiteRepo, f.Source)
		if err := dst.CopyFrom(srcPath, false); err != nil {
			return err
		}
	case f.Source == "" && dst.Exists():
		// We have no source, do nothing
		break
	case f.Source == "" && !dst.Exists():
		// Create an empty file
		if _, err := os.Create(f.Path); err != nil {
			return err
		}
	}

	return nil
}

// createDirectory creates the directory and content managed by the resource
func (f *File) createDirectory() error {
	switch {
	case !f.Recursive:
		return os.Mkdir(f.Path, 0755)
	case f.Recursive && f.Source != "":
		srcPath := filepath.Join(DefaultConfig.SiteRepo, f.Source)
		return utils.CopyDir(srcPath, f.Path)
	case f.Recursive && f.Source == "":
		return os.MkdirAll(f.Path, 0755)
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
func (f *File) isRegularFileContentOutdated() (bool, error) {
	if f.Source != "" {
		srcPath := filepath.Join(DefaultConfig.SiteRepo, f.Source)
		same, err := utils.SameContent(srcPath, f.Path)
		if err != nil {
			return false, err
		}

		if !same {
			item := &outdatedFile{
				src: srcPath,
				dst: f.Path,
			}
			item.flags |= flagOutdatedContent
			f.outdated = append(f.outdated, item)
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
func (f *File) isDirectoryContentOutdated() (bool, error) {
	isOutdated := false
	if f.Source != "" && f.Recursive {
		srcPath := filepath.Join(DefaultConfig.SiteRepo, f.Source)

		// Exclude the ".git" repo directory from the source path,
		// since our source files reside in a git repo
		srcRegistry, err := directoryFileRegistry(srcPath, []string{".git"})
		if err != nil {
			return false, err
		}

		dstRegistry, err := directoryFileRegistry(f.Path, []string{})
		if err != nil {
			return false, err
		}

		// Check source and destination files' content
		for name := range srcRegistry {
			item := &outdatedFile{
				src: srcRegistry[name],
				dst: dstRegistry[name],
			}
			item.flags |= flagOutdatedContent

			// File is missing
			if _, ok := dstRegistry[name]; !ok {
				item.dst = filepath.Join(f.Path, name)
				f.outdated = append(f.outdated, item)
				isOutdated = true
				continue
			}

			// Check if content has changed
			same, err := utils.SameContent(srcRegistry[name], dstRegistry[name])
			if err != nil {
				return false, err
			}

			if !same {
				f.outdated = append(f.outdated, item)
				isOutdated = true
			}
		}

		// Check for extra files in the managed directory
		for name := range dstRegistry {
			if _, ok := srcRegistry[name]; !ok {
				f.extra[dstRegistry[name]] = struct{}{}
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
func (f *File) isPermissionsOutdated() (bool, error) {
	dstRegistry, err := directoryFileRegistry(f.Path, []string{})
	if err != nil {
		return false, err
	}

	isOutdated := false
	for name := range dstRegistry {
		// Skip extra files
		if _, ok := f.extra[dstRegistry[name]]; ok {
			continue
		}

		item := &outdatedFile{
			dst: dstRegistry[name],
		}
		item.flags |= flagOutdatedPermissions

		dst := utils.NewFileUtil(dstRegistry[name])
		mode, err := dst.Mode()
		if err != nil {
			return false, err
		}

		if mode.Perm() != f.Mode {
			f.outdated = append(f.outdated, item)
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
func (f *File) isOwnerOutdated() (bool, error) {
	dstRegistry, err := directoryFileRegistry(f.Path, []string{})
	if err != nil {
		return false, err
	}

	isOutdated := false
	for name := range dstRegistry {
		// Skip extra files
		if _, ok := f.extra[dstRegistry[name]]; ok {
			continue
		}

		item := &outdatedFile{
			dst: dstRegistry[name],
		}
		item.flags |= flagOutdatedOwner
		dst := utils.NewFileUtil(dstRegistry[name])
		owner, err := dst.Owner()
		if err != nil {
			return false, err
		}

		if f.Owner != owner.User.Username || f.Group != owner.Group.Name {
			f.outdated = append(f.outdated, item)
			isOutdated = true
		}
	}

	return isOutdated, nil
}

func init() {
	RegisterProvider("file", NewFile)
}
