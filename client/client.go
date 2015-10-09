package client

import (
	"code.google.com/p/go-uuid/uuid"

	"github.com/dnaeon/gru/minion"
)

type MinionClient interface {
	// Submits a task to a minion
	SubmitTask(u uuid.UUID, t minion.MinionTask) error
}

