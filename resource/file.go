// Copyright (c) 2015-2017 Marin Atanasov Nikolov <dnaeon@gmail.com>
// All rights reserved.
//
// Redistribution and use in source and binary forms, with or without
// modification, are permitted provided that the following conditions
// are met:
//
//  1. Redistributions of source code must retain the above copyright
//     notice, this list of conditions and the following disclaimer
//     in this position and unchanged.
//  2. Redistributions in binary form must reproduce the above copyright
//     notice, this list of conditions and the following disclaimer in the
//     documentation and/or other materials provided with the distribution.
//
// THIS SOFTWARE IS PROVIDED BY THE AUTHOR(S) ``AS IS'' AND ANY EXPRESS OR
// IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED WARRANTIES
// OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE DISCLAIMED.
// IN NO EVENT SHALL THE AUTHOR(S) BE LIABLE FOR ANY DIRECT, INDIRECT,
// INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT
// NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE,
// DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY
// THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
// (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF
// THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.

package resource

import (
	"crypto/md5"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/user"
	"path/filepath"

	"github.com/dnaeon/gru/utils"
)

// BaseFile type is the base type which is embedded by
// File, Directory and Link resources.
type BaseFile struct {
	Base

	// Path to the file. Defaults to the resource name.
	Path string `luar:"-"`

	// Permission bits to set on the file.
	// For regular files defaults to 0644.
	// For directories defaults to 0755.
	Mode os.FileMode `luar:"mode"`

	// Owner of the file. Defaults to the currently running user.
	Owner string `luar:"owner"`

	// Group of the file.
	// Defaults to the group of the currently running user.
	Group string `luar:"group"`
}

// isModeSynced returns a boolean indicating whether the
// permissions of the file managed by the resource are in sync.
func (bf *BaseFile) isModeSynced() (bool, error) {
	dst := utils.NewFileUtil(bf.Path)

	if !dst.Exists() {
		return false, ErrResourceAbsent
	}

	mode, err := dst.Mode()
	if err != nil {
		return false, err
	}

	return mode.Perm() == bf.Mode, nil
}

// setMode sets the permissions on the file managed by the resource.
func (bf *BaseFile) setMode() error {
	Logf("%s setting permissions to %#o\n", bf.ID(), bf.Mode)

	dst := utils.NewFileUtil(bf.Path)

	return dst.Chmod(bf.Mode)
}

// isOwnerSynced checks whether the file ownership is correct.
func (bf *BaseFile) isOwnerSynced() (bool, error) {
	dst := utils.NewFileUtil(bf.Path)

	if !dst.Exists() {
		return false, ErrResourceAbsent
	}

	owner, err := dst.Owner()
	if err != nil {
		return false, err
	}

	return owner.User.Username == bf.Owner && owner.Group.Name == bf.Group, nil
}

// setOwner sets the ownership of the file.
func (bf *BaseFile) setOwner() error {
	Logf("%s setting ownership to %s:%s\n", bf.ID(), bf.Owner, bf.Group)

	dst := utils.NewFileUtil(bf.Path)

	return dst.SetOwner(bf.Owner, bf.Group)
}

// File resource manages files.
//
// Example:
//   foo = resource.file.new("/tmp/foo")
//   foo.state = "present"
//   foo.mode = tonumber("0600", 8)
//   foo.owner = "root"
//   foo.group = "wheel"
//   foo.content = "content of file foo"
type File struct {
	BaseFile

	// Content of file to set.
	Content []byte `luar:"content"`

	// Source file to use for the file content.
	Source string `luar:"source"`
}

// isContentSynced checks if the file content is in sync with the
// given content.
func (f *File) isContentSynced() (bool, error) {
	// We don't have a content, assume content is correct
	if f.Content == nil {
		return true, nil
	}

	dst := utils.NewFileUtil(f.Path)
	if !dst.Exists() {
		return false, ErrResourceAbsent
	}

	dstMd5, err := dst.Md5()
	if err != nil {
		return false, err
	}

	srcMd5 := fmt.Sprintf("%x", md5.Sum(f.Content))

	return srcMd5 == dstMd5, nil
}

// setContent sets the content of the file.
func (f *File) setContent() error {
	dst := utils.NewFileUtil(f.Path)
	dstMd5, err := dst.Md5()
	if err != nil {
		return err
	}

	Logf("%s setting content to md5:%s\n", f.ID(), dstMd5)

	return ioutil.WriteFile(f.Path, f.Content, f.Mode)
}

// NewFile creates a resource for managing regular files.
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
		BaseFile: BaseFile{
			Base: Base{
				Name:              name,
				Type:              "file",
				State:             "present",
				Require:           make([]string, 0),
				PresentStatesList: []string{"present"},
				AbsentStatesList:  []string{"absent"},
				Concurrent:        true,
				Subscribe:         make(TriggerMap),
			},
			Path:  name,
			Mode:  0644,
			Owner: currentUser.Username,
			Group: currentGroup.Name,
		},
		Content: nil,
		Source:  "",
	}

	// Set resource properties
	f.PropertyList = []Property{
		&ResourceProperty{
			PropertyName:         "mode",
			PropertySetFunc:      f.setMode,
			PropertyIsSyncedFunc: f.isModeSynced,
		},
		&ResourceProperty{
			PropertyName:         "ownership",
			PropertySetFunc:      f.setOwner,
			PropertyIsSyncedFunc: f.isOwnerSynced,
		},
		&ResourceProperty{
			PropertyName:         "content",
			PropertySetFunc:      f.setContent,
			PropertyIsSyncedFunc: f.isContentSynced,
		},
	}

	return f, nil
}

// Validate validates the file resource.
func (f *File) Validate() error {
	if err := f.Base.Validate(); err != nil {
		return err
	}

	if f.Source != "" && f.Content != nil {
		return errors.New("cannot use both 'source' and 'content'")
	}

	return nil
}

// Initialize initializes the file resource.
func (f *File) Initialize() error {
	// Set file content from the given source file if any.
	// TODO: Currently this works only for files in the site repo.
	// TODO: Implement a generic file content fetcher.
	if f.Source != "" {
		src := filepath.Join(DefaultConfig.SiteRepo, f.Source)
		content, err := ioutil.ReadFile(src)
		if err != nil {
			return err
		}
		f.Content = content
	}

	return nil
}

// Evaluate evaluates the state of the file resource.
func (f *File) Evaluate() (State, error) {
	state := State{
		Current: "unknown",
		Want:    f.State,
	}

	fi, err := os.Stat(f.Path)
	if os.IsNotExist(err) {
		state.Current = "absent"
		return state, nil
	}

	state.Current = "present"

	if !fi.Mode().IsRegular() {
		return state, errors.New("path exists, but is not a regular file")
	}

	return state, nil
}

// Create creates the file managed by the resource.
func (f *File) Create() error {
	Logf("%s creating file\n", f.ID())

	return ioutil.WriteFile(f.Path, f.Content, f.Mode)
}

// Delete deletes the file managed by the resource.
func (f *File) Delete() error {
	Logf("%s removing file\n", f.ID())

	return os.Remove(f.Path)
}

// Directory resource manages directories.
//
// Example:
//   bar = resource.directory.new("/tmp/bar")
//   bar.state = "present"
//   bar.mode = tonumber("0700", 8)
//   bar.owner = "root"
//   bar.group = "wheel"
type Directory struct {
	BaseFile

	// Parents flag specifies whether or not to create/delete
	// parent directories. Defaults to false.
	Parents bool `luar:"parents"`
}

// NewDirectory creates a resource for managing directories.
func NewDirectory(name string) (Resource, error) {
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
	d := &Directory{
		BaseFile: BaseFile{
			Base: Base{
				Name:              name,
				Type:              "directory",
				State:             "present",
				Require:           make([]string, 0),
				PresentStatesList: []string{"present"},
				AbsentStatesList:  []string{"absent"},
				Concurrent:        true,
				Subscribe:         make(TriggerMap),
			},
			Path:  name,
			Mode:  0755,
			Owner: currentUser.Username,
			Group: currentGroup.Name,
		},
		Parents: false,
	}

	// Set resource properties
	d.PropertyList = []Property{
		&ResourceProperty{
			PropertyName:         "mode",
			PropertySetFunc:      d.setMode,
			PropertyIsSyncedFunc: d.isModeSynced,
		},
		&ResourceProperty{
			PropertyName:         "ownership",
			PropertySetFunc:      d.setOwner,
			PropertyIsSyncedFunc: d.isOwnerSynced,
		},
	}

	return d, nil
}

// Evaluate evaluates the state of the directory.
func (d *Directory) Evaluate() (State, error) {
	state := State{
		Current: "unknown",
		Want:    d.State,
	}

	fi, err := os.Stat(d.Path)
	if os.IsNotExist(err) {
		state.Current = "absent"
		return state, nil
	}

	state.Current = "present"

	if !fi.Mode().IsDir() {
		return state, errors.New("path exists, but is not a directory")
	}

	return state, nil
}

// Create creates the directory.
func (d *Directory) Create() error {
	Logf("%s creating directory\n", d.ID())

	if d.Parents {
		return os.MkdirAll(d.Path, d.Mode)
	}

	return os.Mkdir(d.Path, d.Mode)
}

// Delete removes the directory.
func (d *Directory) Delete() error {
	Logf("%s removing directory\n", d.ID())

	if d.Parents {
		return os.RemoveAll(d.Path)
	}

	return os.Remove(d.Path)
}

// Link resource manages links between files.
//
// Example:
//   baz = resource.link.new("/tmp/baz")
//   baz.state = "present"
//   baz.source = "/tmp/qux"
type Link struct {
	BaseFile

	// Source file points to the file the link will be set to.
	Source string `luar:"source"`

	// Hard flag specifies whether or not to create a hard link to the file.
	// Defaults to false.
	Hard bool `luar:"hard"`
}

// NewLink creates a new resource for managing links between files.
func NewLink(name string) (Resource, error) {
	l := &Link{
		BaseFile: BaseFile{
			Base: Base{
				Name:              name,
				Type:              "link",
				State:             "present",
				Require:           make([]string, 0),
				PresentStatesList: []string{"present"},
				AbsentStatesList:  []string{"absent"},
				Concurrent:        true,
				Subscribe:         make(TriggerMap),
			},
			Path: name,
		},
		Source: "",
		Hard:   false,
	}

	return l, nil
}

// Validate validates the link resource.
func (l *Link) Validate() error {
	if l.Source == "" {
		return errors.New("must provide source file")
	}

	src := utils.NewFileUtil(l.Source)
	if !src.Exists() {
		return fmt.Errorf("source file %s does not exist", l.Source)
	}

	return nil
}

// Evaluate evaluates the state of the link.
func (l *Link) Evaluate() (State, error) {
	state := State{
		Current: "unknown",
		Want:    l.State,
	}

	_, err := os.Stat(l.Path)
	if os.IsNotExist(err) {
		state.Current = "absent"
		return state, nil
	}

	state.Current = "present"

	_, err = os.Readlink(l.Path)
	if err != nil {
		return state, fmt.Errorf("path exists, but is not a link: %s\n", err)
	}

	return state, nil
}

// Create creates the link.
func (l *Link) Create() error {
	Logf("%s creating link\n", l.ID())

	if l.Hard {
		return os.Link(l.Source, l.Path)
	}

	return os.Symlink(l.Source, l.Path)
}

// Delete removes the link.
func (l *Link) Delete() error {
	Logf("%s removing link\n", l.ID())

	return os.Remove(l.Path)
}

func init() {
	file := ProviderItem{
		Type:      "file",
		Provider:  NewFile,
		Namespace: DefaultResourceNamespace,
	}

	dir := ProviderItem{
		Type:      "directory",
		Provider:  NewDirectory,
		Namespace: DefaultResourceNamespace,
	}

	link := ProviderItem{
		Type:      "link",
		Provider:  NewLink,
		Namespace: DefaultResourceNamespace,
	}

	RegisterProvider(file, dir, link)
}
