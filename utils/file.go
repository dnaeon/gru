package utils

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"fmt"
	"io/ioutil"
	"os"
	"os/user"
	"strconv"
	"syscall"
)

// FileUtil type
type FileUtil struct {
	// Path to the file
	path string

	info os.FileInfo
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
		path: path,
		info: info,
	}

	return f, nil
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
	uid := f.info.Sys().(*syscall.Stat_t).Uid
	gid := f.info.Sys().(*syscall.Stat_t).Gid

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
