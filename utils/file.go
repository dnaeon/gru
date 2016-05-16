package utils

import "os"

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

// Md5 returns the md5 hash of the file's content
func (f *FileUtil) Md5() (string, error) {
	buf, err := iotuil.ReadFile(f.path)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", md5.Sum(buf)), nil
}
