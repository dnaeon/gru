package resource

// Resource states
const (
	StateUnknown = "unknown"
	StatePresent = "present"
	StateAbsent  = "absent"
	StateRunning = "running"
	StateStopped = "stopped"
)

// Valid resource states
var validStates = []string{
	StatePresent,
	StateAbsent,
	StateRunning,
	StateStopped,
}

// State type represents the current and wanted states of a resource
type State struct {
	// Current state of the resource
	Current string

	// Wanted state of the resource
	Want string

	// Indicates that a resource is in the desired state, but is
	// out of date and needs to be updated, e.g. a file resource is
	// present, but its permissions need to be corrected.
	Update bool
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
