package minion

import "runtime"

func init() {
	os := NewCallbackClassifier("os", "Operating System", osClassifier)
	arch := NewCallbackClassifier("arch", "Architecture", archClassifier)

	RegisterClassifier(os, arch)
}

func osClassifier(m Minion) (string, error) {
	return runtime.GOOS, nil
}

func archClassifier(m Minion) (string, error) {
	return runtime.GOARCH, nil
}
