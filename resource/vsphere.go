// Copyright (c) 2015-2017 Marin Atanasov Nikolov <dnaeon@gmail.com>
// All rights reserved.
//
// Redistribution and use in source and binary forms, with or without
// modification, are permitted provided that the following conditions
// are met:
//
//  1. Redistributions of source code must retain the above copyright
//     notice, this list of conditions and the following disclaimer
//     in this position and unchanged.
//  2. Redistributions in binary form must reproduce the above copyright
//     notice, this list of conditions and the following disclaimer in the
//     documentation and/or other materials provided with the distribution.
//
// THIS SOFTWARE IS PROVIDED BY THE AUTHOR(S) ``AS IS'' AND ANY EXPRESS OR
// IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED WARRANTIES
// OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE DISCLAIMED.
// IN NO EVENT SHALL THE AUTHOR(S) BE LIABLE FOR ANY DIRECT, INDIRECT,
// INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT
// NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE,
// DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY
// THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
// (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF
// THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.

package resource

import (
	"context"
	"errors"
	"fmt"
	"net/url"

	"github.com/vmware/govmomi"
	"github.com/vmware/govmomi/find"
	"github.com/vmware/govmomi/object"
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

	// Username to use when connecting to the vSphere endpoint.
	// Defaults to an empty string.
	Username string `luar:"username"`

	// Password to use when connecting to the vSphere endpoint.
	// Defaults to an empty string.
	Password string `luar:"password"`

	// Endpoint to the VMware vSphere API. Defaults to an empty string.
	Endpoint string `luar:"endpoint"`

	// Path to use when creating the object managed by the resource.
	// Defaults to "/".
	Path string `luar:"path"`

	// If set to true then allow connecting to vSphere API endpoints with
	// self-signed certificates. Defaults to false.
	Insecure bool `luar:"insecure"`

	url    *url.URL           `luar:"-"`
	ctx    context.Context    `luar:"-"`
	cancel context.CancelFunc `luar:"-"`
	client *govmomi.Client    `luar:"-"`
	finder *find.Finder       `luar:"-"`
}

// ID returns the unique resource id for the resource
func (bv *BaseVSphere) ID() string {
	return fmt.Sprintf("%s[%s@%s]", bv.Type, bv.Name, bv.Endpoint)
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
	defer bv.cancel()

	return bv.client.Logout(bv.ctx)
}

// vSphereRemoveHost disconnects an ESXi host from the
// vCenter server and then removes it.
func vSphereRemoveHost(ctx context.Context, obj *object.HostSystem) error {
	disconnectTask, err := obj.Disconnect(ctx)
	if err != nil {
		return err
	}

	if err := disconnectTask.Wait(ctx); err != nil {
		return err
	}

	destroyTask, err := obj.Destroy(ctx)
	if err != nil {
		return err
	}

	return destroyTask.Wait(ctx)
}

func init() {
	datacenter := ProviderItem{
		Type:      "datacenter",
		Provider:  NewDatacenter,
		Namespace: VSphereNamespace,
	}

	cluster := ProviderItem{
		Type:      "cluster",
		Provider:  NewCluster,
		Namespace: VSphereNamespace,
	}

	clusterHost := ProviderItem{
		Type:      "cluster_host",
		Provider:  NewClusterHost,
		Namespace: VSphereNamespace,
	}

	host := ProviderItem{
		Type:      "host",
		Provider:  NewHost,
		Namespace: VSphereNamespace,
	}

	vm := ProviderItem{
		Type:      "vm",
		Provider:  NewVirtualMachine,
		Namespace: VSphereNamespace,
	}

	datastoreNfs := ProviderItem{
		Type:      "datastore_nfs",
		Provider:  NewDatastoreNfs,
		Namespace: VSphereNamespace,
	}

	RegisterProvider(datacenter, cluster, clusterHost, host, vm, datastoreNfs)
}
