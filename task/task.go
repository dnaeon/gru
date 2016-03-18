package task

import (
	"github.com/dnaeon/gru/catalog"
	"github.com/pborman/uuid"
)

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
)

// Task type represents a task that is processed by minions
type Task struct {
	// Catalog to be processed
	Catalog *catalog.Catalog `json:"catalog"`

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

	// Task state
	State string `json:"state"`
}

// New creates a new task
func New(c *catalog.Catalog) *Task {
	t := &Task{
		Catalog: c,
		TaskID:  uuid.NewRandom(),
		State:   TaskStateUnknown,
	}

	return t
}
