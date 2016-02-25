package minion

import (
	"github.com/dnaeon/gru/classifier"
	"github.com/dnaeon/gru/task"

	"github.com/pborman/uuid"
)

// Minion interface type
type Minion interface {
	// Returns the unique identifier of a minion
	ID() uuid.UUID

	// Sets the name of the minion
	SetName(string) error

	// Sets the time the minion was last seen
	SetLastseen(int64) error

	// Sets a classifier for the minion
	SetClassifier(*classifier.Classifier) error

	// Listens for new tasks and processes them
	TaskListener(c chan<- *task.Task) error

	// Runs new tasks as received by the TaskListener
	TaskRunner(c <-chan *task.Task) error

	// Saves the result of a task
	SaveTaskResult(t *task.Task) error

	// Start serving
	Serve() error

	// Stops the minion
	Stop() error
}
