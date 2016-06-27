package command

import "errors"

var (
	errNoMinion          = errors.New("Missing minion uuid")
	errInvalidUUID       = errors.New("Invalid uuid given")
	errNoMinionFound     = errors.New("No minion(s) found")
	errNoClassifier      = errors.New("Missing classifier key")
	errInvalidClassifier = errors.New("Invalid classifier pattern")
	errNoTask            = errors.New("Missing task uuid")
	errNoModuleName      = errors.New("Missing module name")
	errNoSiteRepo        = errors.New("Missing site repo")
)
