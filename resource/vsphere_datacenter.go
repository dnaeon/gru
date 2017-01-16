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
	state := State{
		Current: "unknown",
		Want:    d.State,
	}

	_, err := d.finder.Datacenter(d.ctx, d.Name)
	if err != nil {
		// Datacenter is absent
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

// Create creates a new datacenter.
func (d *Datacenter) Create() error {
	Logf("%s creating datacenter in %s\n", d.ID(), d.Path)

	folder, err := d.finder.FolderOrDefault(d.ctx, d.Path)
	if err != nil {
		return err
	}

	_, err = folder.CreateDatacenter(d.ctx, d.Name)

	return err
}

// Delete removes the datacenter.
func (d *Datacenter) Delete() error {
	Logf("%s removing datacenter from %s\n", d.ID(), d.Path)

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
