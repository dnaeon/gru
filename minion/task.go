package minion

type MinionTask interface {
	// Gets the command to be executed
	GetCommand() (string, error)

	// Gets the time the task was sent for processing
	GetTimestamp() (int64, error)

	// Processes the task
	Process() error
}
