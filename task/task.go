package task

import "code.google.com/p/go-uuid/uuid"

type Task struct {
	// Command to be executed by the minion
	Command string `json:"command"`

	// Command arguments
	Args []string `json:"args"`

	// Time when the command was sent for processing
	TimeReceived int64 `json:"timeReceived"`

	// Time when the command was processed
	TimeProcessed int64 `json:"timeProcessed"`

	// Task unique identifier
	TaskID uuid.UUID `json:"taskId"`

	// Result of task after processing
	Result string `json:"result"`

	// If true this task can run concurrently with other tasks
	IsConcurrent bool `json:"isConcurrent"`

	// Task error, if any
	Error string `json:"error"`
}

func New(command string, args ...string) *Task {
	t := &Task{
		Command: command,
		Args:    args,
		TaskID:  uuid.NewRandom(),
	}

	return t
}
