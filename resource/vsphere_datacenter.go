package resource

import "github.com/vmware/govmomi/find"

// Datacenter type is a resource which manages datacenters in a
// VMware vSphere environment.
//
// Example:
//   dc = vsphere.datacenter.new("my-datacenter")
//   dc.endpoint = "https://vc01.example.org/sdk"
//   dc.username = "root"
//   dc.password = "myp4ssw0rd"
//   dc.insecure = true
//   dc.state = "present"
//   dc.path = "/SomePath"
type Datacenter struct {
	BaseVSphere
}

// NewDatacenter creates a new resource for managing datacenters in a
// VMware vSphere environment.
func NewDatacenter(name string) (Resource, error) {
	d := &Datacenter{
		BaseVSphere: BaseVSphere{
			Base: Base{
				Name:              name,
				Type:              "datacenter",
				State:             "present",
				Require:           make([]string, 0),
				PresentStatesList: []string{"present"},
				AbsentStatesList:  []string{"absent"},
				Concurrent:        true,
				Subscribe:         make(TriggerMap),
			},
			Username: "",
			Password: "",
			Endpoint: "",
			Insecure: false,
			Path:     "/",
		},
	}

	return d, nil
}

// Evaluate evaluates the state of the datacenter.
func (d *Datacenter) Evaluate() (State, error) {
	s := State{
		Current: "unknown",
		Want:    d.State,
	}

	_, err := d.finder.Datacenter(d.ctx, d.Name)
	if err != nil {
		// Datacenter is absent
		if _, ok := err.(*find.NotFoundError); ok {
			s.Current = "absent"
			return s, nil
		}

		// Something else happened
		return s, err
	}

	s.Current = "present"

	return s, nil
}

// Create creates a new datacenter.
func (d *Datacenter) Create() error {
	folder, err := d.finder.FolderOrDefault(d.ctx, d.Path)
	if err != nil {
		return err
	}

	_, err = folder.CreateDatacenter(d.ctx, d.Name)

	return err
}

// Delete removes the datacenter.
func (d *Datacenter) Delete() error {
	dc, err := d.finder.Datacenter(d.ctx, d.Name)
	if err != nil {
		return err
	}

	task, err := dc.Destroy(d.ctx)
	if err != nil {
		return err
	}

	return task.Wait(d.ctx)
}
