package task

import "github.com/pborman/uuid"

// Task states
const (
	// Unknown state of the task
	// This is the default state of a task
	// when new task is initially created
	TaskStateUnknown = "unknown"

	// Task has been received by the
	// minion and is queued for execution
	TaskStateQueued = "queued"

	// Task is being processed
	TaskStateProcessing = "processing"

	// Task has been processed by the
	// minion and was flagged as successful
	TaskStateSuccess = "success"

	// Task has been processed by the
	// minion and was flagged as failed
	TaskStateFailed = "failed"

	// Task has been skipped
	TaskStateSkipped = "skipped"
)

// Task type represents a task that is processed by minions
type Task struct {
	// Environment to use for this task
	Environment string

	// Command to be processed
	Command string `json:"command"`

	// Time when the command was sent for processing
	TimeReceived int64 `json:"timeReceived"`

	// Time when the command was processed
	TimeProcessed int64 `json:"timeProcessed"`

	// Task unique id
	ID uuid.UUID `json:"id"`

	// Result of task after processing
	Result string `json:"result"`

	// Task state
	State string `json:"state"`
}

// New creates a new task
func New(command, environment string) *Task {
	t := &Task{
		Command:     command,
		Environment: environment,
		ID:          uuid.NewRandom(),
		State:       TaskStateUnknown,
	}

	return t
}
