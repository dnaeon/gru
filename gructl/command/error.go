package command

import (
	"os"
	"fmt"
)

// Displays the error and exists with the
// given exit code
func displayError(err error, code int) {
	fmt.Fprintf(os.Stderr, "Error: %s [exit code %d]\n", err, code)
	os.Exit(code)
}
