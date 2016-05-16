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
func (f *FileUtil) Exists() bool {
	_, err := os.Stat(f.Path)

	return os.IsExist(err)
}

// Abs returns the absolute path for the file
func (f *FileUtil) Abs() (string, error) {
	return filepath.Abs(f.Path)
}

// Md5 returns the md5 checksum of the file's contents
func (f *FileUtil) Md5() (string, error) {
	buf, err := ioutil.ReadFile(f.Path)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", md5.Sum(buf)), nil
}

// Sha1 returns the sha1 checksum of the file's contents
func (f *FileUtil) Sha1() (string, error) {
	buf, err := ioutil.ReadFile(f.Path)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", sha1.Sum(buf)), nil
}

// Sha256 returns the sha256 checksum of the file's contents
func (f *FileUtil) Sha256() (string, error) {
	buf, err := ioutil.ReadFile(f.Path)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", sha256.Sum256(buf)), nil
}

// Remove removes the file
func (f *FileUtil) Remove() error {
	return os.Remove(f.Path)
}

// Chmod changes the permissions of the file
func (f *FileUtil) Chmod(perm os.FileMode) error {
	return os.Chmod(f.Path, perm)
}

// Mode returns the file permission bits
func (f *FileUtil) Mode() (os.FileMode, error) {
	fi, err := os.Stat(f.Path)
	if err != nil {
		return 0, err
	}

	return fi.Mode(), nil
}

// Owner retrieves the owner and group for the file
func (f *FileUtil) Owner() (*FileOwner, error) {
	fi, err := os.Stat(f.Path)
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
func (f *FileUtil) SetOwner(owner, group string) error {
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

	return os.Chown(f.Path, uid, gid)
}

// CopyFrom copies contents from another source to the current file
func (f *FileUtil) CopyFrom(srcPath string) error {
	srcInfo, err := os.Stat(srcPath)
	if err != nil {
		return err
	}

	if !srcInfo.Mode().IsRegular() {
		return fmt.Errorf("%s is not a regular file", srcPath)
	}

	dstInfo, err := os.Stat(f.Path)
	if err != nil {
		return err
	}

	if !dstInfo.Mode().IsRegular() {
		return fmt.Errorf("%s is not a regular file", f.Path)
	}

	if os.SameFile(srcInfo, dstInfo) {
		return nil
	}

	srcFile, err := os.Open(srcPath)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	dstFile, err := os.Create(f.Path)
	if err != nil {
		return err
	}
	defer dstFile.Close()

	_, err = io.Copy(dstFile, srcFile)

	return err
}
