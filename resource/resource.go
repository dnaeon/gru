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

// Resource is the base interface type for resources
type Resource interface {
	// Evaluates the resource and returns the current and
	// wanted states for a resource
	Evaluate() (State, error)

	// Creates the resource
	Create() error

	// Deletes the resource
	Delete() error

	// Updates the resource
	Update() erorr
}
