package minion

import (
	"github.com/dnaeon/gru/classifier"
	"github.com/dnaeon/gru/task"

	"github.com/pborman/uuid"
)

// Minion interface type
type Minion interface {
	// ID returns the unique identifier of a minion
	ID() uuid.UUID

	// SetName sets the name of the minion
	SetName(string) error

	// SetLastseen sets the time the minion was last seen
	SetLastseen(int64) error

	// SetClassifier sets a classifier for the minion
	SetClassifier(*classifier.Classifier) error

	// TaskListener listens for new tasks and processes them
	TaskListener(c chan<- *task.Task) error

	// TaskRunner runs new tasks as received by the TaskListener
	TaskRunner(c <-chan *task.Task) error

	// SaveTaskResult saves the result of a task
	SaveTaskResult(t *task.Task) error

	// Serve start the minion
	Serve() error

	// Stop stops the minion
	Stop() error
}
