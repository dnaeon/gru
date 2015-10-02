package client

import (
	"code.google.com/p/go-uuid/uuid"

	"github.com/dnaeon/gru/minion"
)

type Client interface {
	// Submits a new task to a minion
	SubmitTask(minion uuid.UUID, task MinionTask) error
}
