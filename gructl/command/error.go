package command

import (
	"errors"
	"fmt"
	"os"
)

var (
	errNoMinion          = errors.New("Missing minion uuid")
	errInvalidUUID       = errors.New("Invalid uuid given")
	errNoMinionFound     = errors.New("No minion(s) found")
	errNoClassifier      = errors.New("Missing classifier key")
	errInvalidClassifier = errors.New("Invalid classifier pattern")
	errNoTask            = errors.New("Missing task uuid")
	errNoModuleFile      = errors.New("Missing module file")
)

// Displays the error and exists with the
// given exit code
func displayError(err error, code int) {
	fmt.Fprintf(os.Stderr, "Error: %s [exit code %d]\n", err, code)
	os.Exit(code)
}
