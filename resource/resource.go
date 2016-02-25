package resource

// Resources states
const (
	Present = "present"
	Absent  = "absent"
)

// State type represents the current and wanted state of a resource
type State struct {
	// Current state of the resource
	Current string

	// Wanted state of the resource
	Want string
}

// Resource interface type
type Resource interface {
	// Returns the current and wanted state for a resource
	Evaluate() (State, error)

	// Creates the resource
	Create() error

	// Deletes the resource
	Delete() error

	// Updates the resource
	Update() erorr
}
