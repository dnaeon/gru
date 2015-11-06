package command

import (
	"os"
	"fmt"
	"errors"
)

var (
	errMissingMinion = errors.New("Missing minion uuid")
	errInvalidMinion = errors.New("Invalid minion uuid given")
	errNoMinionFound = errors.New("No minion(s) found")
	errMissingClassifier = errors.New("Missing classifier key")
)

// Displays the error and exists with the
// given exit code
func displayError(err error, code int) {
	fmt.Fprintf(os.Stderr, "Error: %s [exit code %d]\n", err, code)
	os.Exit(code)
}
