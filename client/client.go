package client

import (
	"code.google.com/p/go-uuid/uuid"

	"github.com/dnaeon/gru/minion"
)

type MinionClient interface {
	// Gets the minion name
	Name(u uuid.UUID) (string, error)

	// Gets the time the minion was last seen
	Lastseen(u uuid.UUID) (int64, error)

	// Gets a classifier of a minion
	Classifier(u uuid.UUID, key string) (minion.MinionClassifier, error)

	// Gets all classifiers for a minion
	AllClassifiers(u uuid.UUID) ([]minion.MinionClassifier, error)

	// Gets all minions which are classified with a given classifier key
	// Each key in the result map should uniquely identify a minion
	Classified(key string) (map[string]minion.MinionClassifier, error)

	// Gets the task results for all minions that
	// have processed the task with the given uuid
	Task(u uuid.UUID) (map[string]*minion.MinionTask, error)

	// Gets the tasks which are still in the minion's queue
	Queue(u uuid.UUID) ([]*minion.MinionTask, error)

	// Gets the processed tasks from the minion's log
	Log(u uuid.UUID) ([]*minion.MinionTask, error)

	// Submits a task to a minion
	SubmitTask(u uuid.UUID, t *minion.MinionTask) error
}
