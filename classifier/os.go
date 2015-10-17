package classifier

import "runtime"

func init() {
	Register("os", osProvider)
}

func osProvider() (string, error) {
	return runtime.GOOS, nil
}
