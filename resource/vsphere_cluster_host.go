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
	"path"

	"github.com/vmware/govmomi/find"
	"github.com/vmware/govmomi/vim25/types"
)

// ClusterHost type is a resource which manages hosts in a
// VMware vSphere cluster.
//
// Example:
//   host = vsphere.cluster_host.new("esxi01.example.org")
//   host.endpoint = "https://vc01.example.org/sdk"
//   host.username = "root"
//   host.password = "myp4ssw0rd"
//   host.state = "present"
//   host.path = "/MyDatacenter/host/MyCluster"
//   host.esxi_username = "root"
//   host.esxi_password = "esxip4ssw0rd"
type ClusterHost struct {
	BaseVSphere

	// EsxiUsername is the username used to connect to the
	// remote ESXi host. Defaults to an empty string.
	EsxiUsername string `luar:"esxi_username"`

	// EsxiPassword is the password used to connect to the
	// remote ESXi host. Defaults to an empty string.
	EsxiPassword string `luar:"esxi_password"`

	// SSL thumbprint of the host. Defaults to an empty string.
	SslThumbprint string `luar:"ssl_thumbprint"`

	// Force flag specifies whether or not to forcefully add the
	// host to the cluster, possibly disconnecting it from any other
	// connected vCenter servers. Defaults to false.
	Force bool `luar:"force"`

	// Port to connect to on the remote ESXi host. Defaults to 443.
	Port int32 `luar:"port"`

	// License to attach to the host. Defaults to an empty string.
	License string `luar:"license"`
}

// NewClusterHost creates a new resource for managing hosts in a
// VMware vSphere cluster.
func NewClusterHost(name string) (Resource, error) {
	ch := &ClusterHost{
		BaseVSphere: BaseVSphere{
			Base: Base{
				Name:              name,
				Type:              "cluster_host",
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
		EsxiUsername:  "",
		EsxiPassword:  "",
		SslThumbprint: "",
		Force:         false,
		Port:          443,
		License:       "",
	}

	return ch, nil
}

// Evaluate evaluates the state of the host in the cluster.
func (ch *ClusterHost) Evaluate() (State, error) {
	state := State{
		Current: "unknown",
		Want:    ch.State,
	}

	_, err := ch.finder.HostSystem(ch.ctx, path.Join(ch.Path, ch.Name))
	if err != nil {
		// Host is absent
		if _, ok := err.(*find.NotFoundError); ok {
			state.Current = "absent"
			return state, nil
		}

		// Something else happened
		return state, err
	}

	state.Current = "present"

	return state, nil
}

// Create adds the host to the cluster.
func (ch *ClusterHost) Create() error {
	Logf("%s adding host to %s\n", ch.ID(), ch.Path)

	obj, err := ch.finder.ClusterComputeResource(ch.ctx, ch.Path)
	if err != nil {
		return err
	}

	spec := types.HostConnectSpec{
		HostName:      ch.Name,
		Port:          ch.Port,
		SslThumbprint: ch.SslThumbprint,
		UserName:      ch.EsxiUsername,
		Password:      ch.EsxiPassword,
		Force:         ch.Force,
		LockdownMode:  "",
	}

	task, err := obj.AddHost(ch.ctx, spec, true, &ch.License, nil)
	if err != nil {
		return err
	}

	return task.Wait(ch.ctx)
}

// Delete disconnects the host and then removes it.
func (ch *ClusterHost) Delete() error {
	Logf("%s removing host from %s\n", ch.ID(), ch.Path)

	obj, err := ch.finder.HostSystem(ch.ctx, path.Join(ch.Path, ch.Name))
	if err != nil {
		return err
	}

	return vSphereRemoveHost(ch.ctx, obj)
}
