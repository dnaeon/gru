package resource

// Property type represents a resource property, which can be
// evaluated and set if needed.
type Property interface {
	// Name returns the property name
	Name() string

	// Set sets the property to it's desired state.
	Set() error

	// IsSynced returns a boolean indicating whether the
	// resource property is in sync or not.
	IsSynced() (bool, error)
}

// ResourceProperty type implements the Property interface.
type ResourceProperty struct {
	// PropertySetFunc is the type of the function that is called when
	// setting a resource property to it's desired state.
	PropertySetFunc func() error

	// PropertyIsSyncedFunc is the type of the function that is called when
	// determining whether a resource property is in the desired state.
	PropertyIsSyncedFunc func() (bool, error)

	// PropertyName is the name of the property.
	PropertyName string
}

// Set sets the property to it's desired state.
func (rp *ResourceProperty) Set() error {
	return rp.PropertySetFunc()
}

// IsSynced returns a boolean indicating whether the
// resource property is in the desired state.
func (rp *ResourceProperty) IsSynced() (bool, error) {
	return rp.PropertyIsSyncedFunc()
}

// Name returns the property name.
func (rp *ResourceProperty) Name() string {
	return rp.PropertyName
}
