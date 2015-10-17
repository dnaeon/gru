package classifier

import "runtime"

func init() {
	Register("arch", archProvider)
}

func archProvider() (string, error) {
	return runtime.GOARCH, nil
}
