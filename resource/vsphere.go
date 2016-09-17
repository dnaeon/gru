package resource

import (
	"context"
	"errors"
	"net/url"

	"github.com/vmware/govmomi"
	"github.com/vmware/govmomi/find"
)

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
