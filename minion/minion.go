package minion

import (
	"time"

	"github.com/dnaeon/gru/task"

	"code.google.com/p/go-uuid/uuid"
)

type Minion interface {
	// Set name of minion
	SetName(name string) error

	// Set the time the minion was last seen in seconds since the Epoch
	SetLastseen(s int64) error

	// Classify minion a with given a key and value
	SetClassifier(c MinionClassifier) error

	// Runs periodic functions, e.g. refreshes classifies and lastseen
	PeriodicRunner(t *time.Ticker) error

	// Listens for new tasks and processes them
	TaskListener(c chan<- task.MinionTask) error

	// Runs new tasks as received by the TaskListener
	TaskRunner (c <-chan task.MinionTask) error

	// Saves a task result in the minion's log directory
	SaveTaskResult(t task.MinionTask) error

	// Start serving
	Serve() error
}

// Generates a uuid for a minion
func GenerateUUID(name string) uuid.UUID {
	u := uuid.NewSHA1(uuid.NameSpace_DNS, []byte(name))

	return u
}
