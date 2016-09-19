package resource

import (
	"context"
	"errors"
	"net/url"

	"github.com/vmware/govmomi"
	"github.com/vmware/govmomi/find"
)

// VSphereNamespace is the table name in Lua where vSphere resources are
// being registered to.
const VSphereNamespace = "vsphere"

// ErrNoUsername error is returned when no username is provided for
// establishing a connection to the remote VMware vSphere API endpoint.
var ErrNoUsername = errors.New("No username provided")

// ErrNoPassword error is returned when no password is provided for
// establishing a connection to the remote VMware vSphere API endpoint.
var ErrNoPassword = errors.New("No password provided")

// ErrNoEndpoint error is returned when no VMware vSphere API endpoint is
// provided.
var ErrNoEndpoint = errors.New("No endpoint provided")

// ErrNotVC error is returned when the remote endpoint is not a vCenter system.
var ErrNotVC = errors.New("Not a VMware vCenter endpoint")

// BaseVSphere type is the base type for all vSphere related resources.
type BaseVSphere struct {
	Base

	// Username to use when connecting to the vSphere endpoint
	Username string `luar:"username"`

	// Password to use when connecting to the vSphere endpoint
	Password string `luar:"password"`

	// Endpoint to the VMware vSphere API
	Endpoint string `luar:"endpoint"`

	// Folder to use when creating the object managed by the resource
	Folder string `luar:"folder"`

	// If set to true then allow connecting to vSphere API endpoints with
	// self-signed certificates.
	Insecure bool `luar:"insecure"`

	url    *url.URL           `luar:"-"`
	ctx    context.Context    `luar:"-"`
	cancel context.CancelFunc `luar:"-"`
	client *govmomi.Client    `luar:"-"`
	finder *find.Finder       `luar:"-"`
}

// Validate validates the resource.
func (bv *BaseVSphere) Validate() error {
	if err := bv.Base.Validate(); err != nil {
		return err
	}

	if bv.Username == "" {
		return ErrNoUsername
	}

	if bv.Password == "" {
		return ErrNoPassword
	}

	if bv.Endpoint == "" {
		return ErrNoEndpoint
	}

	// Validate the URL to the API endpoint and set the username and password info
	endpoint, err := url.Parse(bv.Endpoint)
	if err != nil {
		return err
	}
	endpoint.User = url.UserPassword(bv.Username, bv.Password)
	bv.url = endpoint

	return nil
}

// Initialize establishes a connection to the remote vSphere API endpoint.
func (bv *BaseVSphere) Initialize() error {
	bv.ctx, bv.cancel = context.WithCancel(context.Background())

	// Connect and login to the VMWare vSphere API endpoint
	c, err := govmomi.NewClient(bv.ctx, bv.url, bv.Insecure)
	if err != nil {
		return err
	}

	bv.client = c
	bv.finder = find.NewFinder(bv.client.Client, true)

	return nil
}

// Close closes the connection to the remote vSphere API endpoint.
func (bv *BaseVSphere) Close() error {
	bv.cancel()

	return bv.client.Logout(bv.ctx)
}

// Datacenter type is a resource which manages datacenters in a
// VMware vSphere environment.
//
// Example:
//   dc = vsphere.datacenter.new("my-datacenter")
//   dc.username = "root"
//   dc.password = "myp4ssw0rd"
//   dc.endpoint = "https://vc01.example.org/sdk"
//   dc.insecure = true
//   dc.state = "present"
//   dc.folder = "/SomeFolder"
type Datacenter struct {
	BaseVSphere
}

// NewDatacenter creates a new resource for managing datacenters in a
// VMware vSphere environment.
func NewDatacenter(name string) (Resource, error) {
	d := &Datacenter{
		BaseVSphere: BaseVSphere{
			Base: Base{
				Name:          name,
				Type:          "datacenter",
				State:         "present",
				Require:       make([]string, 0),
				PresentStates: []string{"present"},
				AbsentStates:  []string{"absent"},
				Concurrent:    true,
				Subscribe:     make(TriggerMap),
			},
			Username: "",
			Password: "",
			Endpoint: "",
			Insecure: false,
			Folder:   "/",
		},
	}

	return d, nil
}

// Evaluate evaluates the state of the datacenter
func (d *Datacenter) Evaluate() (State, error) {
	s := State{
		Current:  "unknown",
		Want:     d.State,
		Outdated: false,
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

// Create creates a new datacenter
func (d *Datacenter) Create() error {
	Log(d, "creating datacenter\n")

	folder, err := d.finder.FolderOrDefault(d.ctx, d.Folder)
	if err != nil {
		return err
	}

	_, err = folder.CreateDatacenter(d.ctx, d.Name)

	return err
}

// Delete removes the datacenter
func (d *Datacenter) Delete() error {
	Log(d, "removing datacenter\n")

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

// Update is no-op
func (d *Datacenter) Update() error {
	return nil
}

func init() {
	datacenter := ProviderItem{
		Type:      "datacenter",
		Provider:  NewDatacenter,
		Namespace: VSphereNamespace,
	}

	RegisterProvider(datacenter)
}
