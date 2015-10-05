package minion

import (
	"time"

	"code.google.com/p/go-uuid/uuid"
)

type Minion interface {
	// Get minion identifier
	GetUUID() uuid.UUID

	// Set name of minion
	SetName(name string) error

	// Get name of minion
	GetName() (string, error)

	// Set the time the minion was last seen in seconds since the Epoch
	SetLastseen(s int64) error

	// Get a classifier for a minion
	GetClassifier(key string) (MinionClassifier, error)

	// Classify minion a with given a key and value
	SetClassifier(c MinionClassifier) error

	// Runs periodic functions, e.g. refreshes classifies and lastseen
	Refresh(t *time.Ticker) error

	// Listens for new tasks and processes them
	TaskListener(c chan<- MinionTask) error

	// Runs new tasks as received by the TaskListener
	TaskRunner (c <-chan MinionTask) error

	// Start serving
	Serve() error
}

// Generates a uuid for a minion
func GenerateUUID(name string) uuid.UUID {
	u := uuid.NewSHA1(uuid.NameSpace_DNS, []byte(name))

	return u
}
