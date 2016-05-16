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
	// Path to the file
	path string

	os.FileInfo
}

// FileOwner type provides details about the user and group that owns a file
type FileOwner struct {
	*user.User
	*user.Group
}

// NewFileUtil creates a file utility from the given path
func NewFileUtil(path string) (*FileUtil, error) {
	info, err := os.Stat(path)
	if err != nil {
		return &FileUtil{}, err
	}

	f := &FileUtil{
		path,
		info,
	}

	return f, nil
}

// Abs returns the absolute path for the file
func (f *FileUtil) Abs() (string, error) {
	return filepath.Abs(f.path)
}

// Md5 returns the md5 checksum of the file's contents
func (f *FileUtil) Md5() (string, error) {
	buf, err := ioutil.ReadFile(f.path)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", md5.Sum(buf)), nil
}

// Sha1 returns the sha1 checksum of the file's contents
func (f *FileUtil) Sha1() (string, error) {
	buf, err := ioutil.ReadFile(f.path)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", sha1.Sum(buf)), nil
}

// Sha256 returns the sha256 checksum of the file's contents
func (f *FileUtil) Sha256() (string, error) {
	buf, err := ioutil.ReadFile(f.path)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", sha256.Sum256(buf)), nil
}

// Owner retrieves the owner and group for the file
func (f *FileUtil) Owner() (*FileOwner, error) {
	uid := f.Sys().(*syscall.Stat_t).Uid
	gid := f.Sys().(*syscall.Stat_t).Gid

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

	return os.Chown(f.path, uid, gid)
}

// CopyFrom copies contents from another source to the current file
func (f *FileUtil) CopyFrom(from *FileUtil) error {
	if !f.Mode().IsRegular() {
		return fmt.Errorf("%s is not a regular file", f.path)
	}

	if !from.Mode().IsRegular() {
		return fmt.Errorf("%s is not a regular file", from.path)
	}

	if os.SameFile(f, from) {
		return nil
	}

	src, err := os.Open(from.path)
	if err != nil {
		return err
	}
	defer src.Close()

	dst, err := os.Create(f.path)
	if err != nil {
		return err
	}
	defer dst.Close()

	_, err = io.Copy(dst, src)

	return err
}
