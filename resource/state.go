package resource

// Resource states
const (
	StateUnknown = "unknown"
	StatePresent = "present"
	StateAbsent  = "absent"
	StateUpdate  = "update"
	StateRunning = "running"
	StateStopped = "stopped"
)

// Valid resource states
var validStates = []string{
	StatePresent,
	StateAbsent,
	StateUpdate,
	StateRunning,
	StateStopped,
}

// StateIsValid checks if a given state is a valid one
// Returns true if the state is valid, false otherwise
func StateIsValid(name string) bool {
	for _, state := range validStates {
		if state == name {
			return true
		}
	}

	return false
}
