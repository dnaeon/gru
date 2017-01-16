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

package utils

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/user"
	"path/filepath"
	"strconv"
	"syscall"
)

// FileUtil type
type FileUtil struct {
	// Path to the file we manage
	Path string
}

// FileOwner type provides details about the user and group that owns a file
type FileOwner struct {
	*user.User
	*user.Group
}

// NewFileUtil creates a file utility from the given path
func NewFileUtil(path string) *FileUtil {
	return &FileUtil{path}
}

// Exists returns a boolean indicating whether the file exists or not
func (fu *FileUtil) Exists() bool {
	_, err := os.Stat(fu.Path)

	return !os.IsNotExist(err)
}

// Abs returns the absolute path for the file
func (fu *FileUtil) Abs() (string, error) {
	return filepath.Abs(fu.Path)
}

// Md5 returns the md5 checksum of the file's contents
func (fu *FileUtil) Md5() (string, error) {
	buf, err := ioutil.ReadFile(fu.Path)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", md5.Sum(buf)), nil
}

// Sha1 returns the sha1 checksum of the file's contents
func (fu *FileUtil) Sha1() (string, error) {
	buf, err := ioutil.ReadFile(fu.Path)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", sha1.Sum(buf)), nil
}

// Sha256 returns the sha256 checksum of the file's contents
func (fu *FileUtil) Sha256() (string, error) {
	buf, err := ioutil.ReadFile(fu.Path)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", sha256.Sum256(buf)), nil
}

// Remove removes the file
func (fu *FileUtil) Remove() error {
	return os.Remove(fu.Path)
}

// Chmod changes the permissions of the file
func (fu *FileUtil) Chmod(perm os.FileMode) error {
	return os.Chmod(fu.Path, perm)
}

// Mode returns the file permission bits
func (fu *FileUtil) Mode() (os.FileMode, error) {
	fi, err := os.Stat(fu.Path)
	if err != nil {
		return 0, err
	}

	return fi.Mode(), nil
}

// Owner retrieves the owner and group for the file
func (fu *FileUtil) Owner() (*FileOwner, error) {
	fi, err := os.Stat(fu.Path)
	if err != nil {
		return &FileOwner{}, err
	}

	uid := fi.Sys().(*syscall.Stat_t).Uid
	gid := fi.Sys().(*syscall.Stat_t).Gid

	u, err := user.LookupId(strconv.FormatInt(int64(uid), 10))
	if err != nil {
		return &FileOwner{}, err
	}

	g, err := user.LookupGroupId(strconv.FormatInt(int64(gid), 10))
	if err != nil {
		return &FileOwner{}, err
	}

	owner := &FileOwner{u, g}

	return owner, nil
}

// SetOwner sets the ownership for the file
func (fu *FileUtil) SetOwner(owner, group string) error {
	o, err := user.Lookup(owner)
	if err != nil {
		return err
	}

	g, err := user.LookupGroup(group)
	if err != nil {
		return err
	}

	uid, err := strconv.Atoi(o.Uid)
	if err != nil {
		return err
	}

	gid, err := strconv.Atoi(g.Gid)
	if err != nil {
		return err
	}

	return os.Chown(fu.Path, uid, gid)
}

// CopyFrom copies contents from another source to the current file
func (fu *FileUtil) CopyFrom(srcPath string, overwrite bool) error {
	srcInfo, err := os.Stat(srcPath)
	if err != nil {
		return err
	}

	if !srcInfo.Mode().IsRegular() {
		return fmt.Errorf("%s is not a regular file", srcPath)
	}

	mode := srcInfo.Mode()
	dstInfo, err := os.Stat(fu.Path)
	if !os.IsNotExist(err) {
		if !overwrite {
			return fmt.Errorf("%s already exists", fu.Path)
		}
		mode = dstInfo.Mode()
	}

	srcFile, err := os.Open(srcPath)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	dstFile, err := os.Create(fu.Path)
	if err != nil {
		return err
	}
	defer dstFile.Close()

	if _, err := io.Copy(dstFile, srcFile); err != nil {
		return err
	}

	return os.Chmod(fu.Path, mode)
}

// SameContentWith returns a boolean indicating whether the
// content of the current file is the same as the destination
func (fu *FileUtil) SameContentWith(dst string) (bool, error) {
	srcMd5, err := fu.Md5()
	if err != nil {
		return false, err
	}

	dstFile := NewFileUtil(dst)
	dstMd5, err := dstFile.Md5()
	if err != nil {
		return false, err
	}

	return srcMd5 == dstMd5, nil
}

// SameContent returns a boolean indicating whether two
// files have the same content
func SameContent(src, dst string) (bool, error) {
	srcFile := NewFileUtil(src)

	return srcFile.SameContentWith(dst)
}

// WalkPath walks a path and returns a slice of file names
// that were found during path traversing
func WalkPath(root string, skip []string) ([]string, error) {
	var files []string

	walker := func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip specific directories if provided
		if info.IsDir() {
			for _, name := range skip {
				if name == info.Name() {
					return filepath.SkipDir
				}
			}
		}
		files = append(files, path)

		return nil
	}

	return files, filepath.Walk(root, walker)
}

// CopyDir recursively copies files from one directory to another
func CopyDir(srcPath, dstPath string) error {
	srcInfo, err := os.Stat(srcPath)
	if err != nil {
		return err
	}

	// Ensure the source is an actual directory
	if !srcInfo.IsDir() {
		return fmt.Errorf("%s is not a directory", srcPath)
	}

	// Ensure the destination does not yet exist
	_, err = os.Open(dstPath)
	if !os.IsNotExist(err) {
		return fmt.Errorf("%s already exists", dstPath)
	}

	if err := os.MkdirAll(dstPath, srcInfo.Mode()); err != nil {
		return err
	}

	// Read in the source files
	srcDir, err := os.Open(srcPath)
	if err != nil {
		return err
	}

	srcFiles, err := srcDir.Readdir(-1)
	if err != nil {
		return err
	}

	for _, fi := range srcFiles {
		srcName := filepath.Join(srcPath, fi.Name())
		dstName := filepath.Join(dstPath, fi.Name())

		// Copy sub directories
		if fi.IsDir() {
			if err := CopyDir(srcName, dstName); err != nil {
				return err
			}
		} else {
			// Copy file
			dstFile := NewFileUtil(dstName)
			if err := dstFile.CopyFrom(srcName, false); err != nil {
				return err
			}
		}
	}

	return nil
}
