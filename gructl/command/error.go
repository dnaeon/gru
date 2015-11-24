package command

import (
	"errors"
	"fmt"
	"os"
)

var (
	errMissingMinion     = errors.New("Missing minion uuid")
	errInvalidUUID       = errors.New("Invalid uuid given")
	errNoMinionFound     = errors.New("No minion(s) found")
	errMissingClassifier = errors.New("Missing classifier key")
	errInvalidClassifier = errors.New("Invalid classifier pattern")
	errMissingTask       = errors.New("Missing task uuid")
)

// Displays the error and exists with the
// given exit code
func displayError(err error, code int) {
	fmt.Fprintf(os.Stderr, "Error: %s [exit code %d]\n", err, code)
	os.Exit(code)
}
