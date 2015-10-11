package client

import (
	"code.google.com/p/go-uuid/uuid"

	"github.com/dnaeon/gru/minion"
)

type MinionClient interface {
	// Gets the minion name
	GetName(u uuid.UUID) (string, error)

	// Submits a task to a minion
	SubmitTask(u uuid.UUID, t minion.MinionTask) error
}
