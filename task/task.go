package task

import "code.google.com/p/go-uuid/uuid"

type Task struct {
	// Command to be executed by the minion
	Command string

	// Command arguments
	Args []string

	// Time when the command was sent for processing
	TimeReceived int64

	// Time when the command was processed
	TimeProcessed int64

	// Task unique identifier
	TaskID uuid.UUID

	// Result of task after processing
	Result string

	// If true this task can run concurrently with other tasks
	IsConcurrent bool

	// Task error, if any
	Error string
}

func New(command string, args ...string) *Task {
	t := &Task{
		Command: command,
		Args: args,
		TaskID: uuid.NewRandom(),
	}

	return t
}
