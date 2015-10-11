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

	// Submits a task to a minion
	SubmitTask(u uuid.UUID, t minion.MinionTask) error
}
