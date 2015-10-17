package minion

import (
	"github.com/dnaeon/gru/task"

	"code.google.com/p/go-uuid/uuid"
)

type Minion interface {
	// Returns the unique identifier of a minion
	ID() uuid.UUID

	// Returns the assigned name of the minion
	Name() string

	// Classifies the minion using any classifiers
	Classify() error

	// Listens for new tasks and processes them
	TaskListener(c chan<- *task.Task) error

	// Runs new tasks as received by the TaskListener
	TaskRunner (c <-chan *task.Task) error

	// Saves the result of a task
	SaveTaskResult(t *task.Task) error

	// Start serving
	Serve() error
}
