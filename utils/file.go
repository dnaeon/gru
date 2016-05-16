package utils

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"fmt"
	"os"
)

// FileUtil type
type FileUtil struct {
	// Path to the file
	path string

	os.FileInfo
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

// Md5 returns the md5 checksum of the file's contents
func (f *FileUtil) Md5() (string, error) {
	buf, err := iotuil.ReadFile(f.path)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", md5.Sum(buf)), nil
}

// Sha1 returns the sha1 checksum of the file's contents
func (f *FileUtil) Sha1() (string, error) {
	buf, err := iotuil.ReadFile(f.path)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", sha1.Sum(buf)), nil
}

// Sha256 returns the sha256 checksum of the file's contents
func (f *FileUtil) Sha256() (string, error) {
	buf, err := iotuil.ReadFile(f.path)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", sha256.Sum256(buf)), nil
}
