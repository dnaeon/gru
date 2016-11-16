package resource

import (
	"errors"
	"path"

	"github.com/vmware/govmomi/find"
	"github.com/vmware/govmomi/vim25/types"
)

// DatastoreNfs type is a resource which manages NFS datastores on ESXi hosts.
//
// Example:
//   datastore = vsphere.datastore_nfs.new("vm-storage01")
//   datastore.endpoint = "https://vc01.example.org/sdk"
//   datastore.username = "root"
//   datastore.password = "myp4ssw0rd"
//   datastore.state = "present"
//   datastore.path = "/MyDatacenter/datastore"
//   datastore.hosts = {
//      "/MyDatacenter/host/MyCluster/esxi01.example.org",
//      "/MyDatacenter/host/MyCluster/esxi02.example.org",
//   }
//   datastore.nfs_server = "nfs01.example.org"
//   datastore.nfs_type = "NFS",
//   datastore.nfs_path = "/exported/file/system",
//   datastore.mode = "readWrite"
type DatastoreNfs struct {
	BaseVSphere

	// Hosts is the list of ESXi hosts on which to manage the NFS datastore.
	Hosts []string `luar:"hosts"`

	// NfsServer is the remote NFS server to use when mounting the datastore.
	NfsServer string `luar:"nfs_server"`

	// NfsType specifies the type of the NFS volume.
	// Valid values are "NFS" for v3 and "NFS41" for v4.1.
	// Defaults to "NFS".
	NfsType string `luar:"nfs_type"`

	// NfsPath is the remote path of the NFS mount point.
	NfsPath string `luar:"nfs_path"`

	// Mode is the access mode for the datastore.
	// Valid values are "readOnly" and "readWrite".
	// Defaults to "readWrite".
	Mode string `luar:"mode"`
}

// mountOn mounts the NFS datastore on an ESXi host.
func (ds *DatastoreNfs) mountOn(host string) error {
	Logf("%s mounting datastore on %s\n", ds.ID(), path.Base(host))

	obj, err := ds.finder.HostSystem(ds.ctx, host)
	if err != nil {
		return err
	}

	datastoreSystem, err := obj.ConfigManager().DatastoreSystem(ds.ctx)
	if err != nil {
		return err
	}

	spec := types.HostNasVolumeSpec{
		RemoteHost: ds.NfsServer,
		RemotePath: ds.NfsPath,
		LocalPath:  ds.Name,
		AccessMode: ds.Mode,
		Type:       ds.NfsType,
	}

	_, err = datastoreSystem.CreateNasDatastore(ds.ctx, spec)

	return err
}

// NewDatastoreNfs creates a new resource for managing NFS datastores on ESXi hosts.
func NewDatastoreNfs(name string) (Resource, error) {
	ds := &DatastoreNfs{
		BaseVSphere: BaseVSphere{
			Base: Base{
				Name:              name,
				Type:              "datastore_nfs",
				State:             "present",
				Require:           make([]string, 0),
				PresentStatesList: []string{"present"},
				AbsentStatesList:  []string{"absent"},
				Concurrent:        true,
				Subscribe:         make(TriggerMap),
			},
		},
		Hosts:     make([]string, 0),
		NfsServer: "",
		NfsType:   "NFS",
		NfsPath:   "",
		Mode:      "readWrite",
	}

	return ds, nil
}

// Validate validates the datastore resource.
func (ds *DatastoreNfs) Validate() error {
	if err := ds.BaseVSphere.Validate(); err != nil {
		return err
	}

	if len(ds.Hosts) == 0 {
		return errors.New("Must provide list of hosts for the datastore")
	}

	return nil
}

// Evaluate evaluates the state of the datastore.
func (ds *DatastoreNfs) Evaluate() (State, error) {
	state := State{
		Current: "unknown",
		Want:    ds.State,
	}

	_, err := ds.finder.Datastore(ds.ctx, path.Join(ds.Path, ds.Name))
	if err != nil {
		// Datastore is absent
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

// Create mounts the NFS datastore on the ESXi hosts.
func (ds *DatastoreNfs) Create() error {
	if ds.NfsServer == "" {
		return errors.New("Missing NFS server for datastore")
	}

	if ds.NfsPath == "" {
		return errors.New("Missing NFS path for the datastore")
	}

	for _, host := range ds.Hosts {
		if err := ds.mountOn(host); err != nil {
			return err
		}
	}

	return nil
}

// Delete unmounts the NFS datastore from the ESXi hosts.
func (ds *DatastoreNfs) Delete() error {
	datastore, err := ds.finder.Datastore(ds.ctx, path.Join(ds.Path, ds.Name))
	if err != nil {
		return err
	}

	for _, host := range ds.Hosts {
		Logf("%s unmounting datastore from %s\n", ds.ID(), path.Base(host))
		obj, err := ds.finder.HostSystem(ds.ctx, host)
		if err != nil {
			return err
		}

		datastoreSystem, err := obj.ConfigManager().DatastoreSystem(ds.ctx)
		if err != nil {
			return err
		}

		if err := datastoreSystem.Remove(ds.ctx, datastore); err != nil {
			return err
		}
	}

	return nil
}
