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
