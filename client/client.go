package client

import (
	"code.google.com/p/go-uuid/uuid"

//	"github.com/dnaeon/gru/minion"
	"gru/minion"
)

type Client interface {
	// Submits a new task to a minion
	SubmitTask(u uuid.UUID, t minion.MinionTask) error
}

