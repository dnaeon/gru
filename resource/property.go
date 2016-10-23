package resource

// Property type represents a resource property, which can be
// evaluated and set if outdated.
type Property interface {
	// Name returns the property name
	Name() string

	// Set sets the property to it's desired state.
	Set() error

	// IsSynced returns a boolean indicating whether the
	// resource property is in sync or not.
	IsSynced() (bool, error)
}

// PropertySetFunc is the type of the function that is called when
// setting a resource property to it's desired state.
type PropertySetFunc func() error

// Set sets the property to it's desired state.
func (f PropertySetFunc) Set() error {
	return f()
}

// PropertyIsSyncedFunc is the type of the function that is called when
// determining whether a resource property is in the desired state.
type PropertyIsSyncedFunc func() (bool, error)

// IsSynced returns a boolean indicating whether the
// resource property is in the desired state.
func (f PropertyIsSyncedFunc) IsSynced() (bool, error) {
	return f()
}

// ResourceProperty type implements the Property interface.
type ResourceProperty struct {
	PropertySetFunc
	PropertyIsSyncedFunc

	// PropertyName is the name of the property.
	PropertyName string
}

// Name returns the property name.
func (rp *ResourceProperty) Name() string {
	return rp.PropertyName
}
