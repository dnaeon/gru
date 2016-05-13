package resource

// Resource states
const (
	StateUnknown = "unknown"
	StatePresent = "present"
	StateAbsent  = "absent"
	StateRunning = "running"
	StateStopped = "stopped"
)

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
