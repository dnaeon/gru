package minion

import "code.google.com/p/go-uuid/uuid"

type MinionTask interface {
	// Gets the UUID of the task
	GetTaskID() uuid.UUID

	// Gets the command to be executed
	GetCommand() (string, error)

	// Gets the command arguments
	GetArgs() ([]string, error)

	// Gets the time the task was sent for processing
	GetTimestamp() (int64, error)

	// Gets the task result
	GetResult() (string, error)

	// Gets the task error, if any
	GetError() string

	// Processes the task
	Process() error
}
