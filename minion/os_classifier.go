package minion

import "runtime"

func init() {
	os := NewSimpleClassifier("os", "Operating System", runtime.GOOS)
	arch := NewSimpleClassifier("arch", "Architecture", runtime.GOARCH)

	RegisterClassifier(os, arch)
}
