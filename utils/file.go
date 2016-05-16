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
