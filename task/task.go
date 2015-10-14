package task

import "code.google.com/p/go-uuid/uuid"

type MinionTask interface {
	// Gets the UUID of the task
	GetTaskID() uuid.UUID

	// Gets the command to be executed
	GetCommand() (string, error)

	// Gets the command arguments
	GetArgs() ([]string, error)

	// Gets the time the task was sent for processing
	GetTimeReceived() (int64, error)

	// Gets the time when the task has been processed
	GetTimeProcessed() (int64, error)

	// Gets the task result
	GetResult() (string, error)

	// Gets the task error, if any
	GetError() string

	// Whether or not this task can run concurrently with other tasks
	IsConcurrent() bool

	// Sets the flag whether or not this task can run
	// concurrently with other tasks
	SetConcurrent(bool) error

	// Processes the task
	Process() error
}
