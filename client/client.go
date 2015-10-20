package client

import (
	"code.google.com/p/go-uuid/uuid"

	"github.com/dnaeon/gru/classifier"
	"github.com/dnaeon/gru/task"
)

type Client interface {
	// Gets the name of a minion
	MinionName(m uuid.UUID) (string, error)

	// Gets the time a minion was last seen
	MinionLastseen(m uuid.UUID) (int64, error)

	// Gets a classifier of a minion
	MinionClassifier(m uuid.UUID, key string) (*classifier.Classifier, error)

	// Gets all classifier keys of a minion
	MinionClassifierKeys(m uuid.UUID) ([]string, error)

	// Gets minions which are classified with a given classifier key
	MinionWithClassifier(key string) ([]uuid.UUID, error)

	// Gets the result of a task for a minion
	MinionTaskResult(m uuid.UUID, t uuid.UUID) (*task.Task, error)

	// Gets the minions which have a task result with the given uuid
	MinionWithTaskResult(t uuid.UUID) ([]uuid.UUID, error)

	// Gets the tasks which are currently pending in the queue
	MinionTaskQueue(m uuid.UUID) ([]*task.Task, error)

	// Gets the uuids of tasks which have already been processed
	MinionTaskLog(m uuid.UUID) ([]uuid.UUID, error)

	// Submits a new task to a minion
	MinionSubmitTask(m uuid.UUID, t *task.Task) error
}
