package client

import (
	"code.google.com/p/go-uuid/uuid"

	"github.com/dnaeon/gru/minion"
)

type MinionClient interface {
	// Gets the minion name
	GetName(u uuid.UUID) (string, error)

	// Gets the time the minion was last seen
	GetLastseen(u uuid.UUID) (int64, error)

	// Gets a classifier of a minion
	GetClassifier(u uuid.UUID, key string) (minion.MinionClassifier, error)

	// Gets all classifiers for a minion
	GetAllClassifiers(u uuid.UUID) ([]minion.MinionClassifier, error)

	// Gets all minions which are classified with a given classifier key
	// Each key in the result map should uniquely identify a minion
	GetClassifiedMinions(key string) (map[string]minion.MinionClassifier, error)

	// Gets task results for all minions that
	// have a task with the given uuid
	GetTask(u uuid.UUID) (map[string]minion.MinionTask, error)

	// Submits a task to a minion
	SubmitTask(u uuid.UUID, t minion.MinionTask) error
}
