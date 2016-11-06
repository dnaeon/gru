package resource

import (
	"path"

	"github.com/vmware/govmomi/find"
	"github.com/vmware/govmomi/vim25/mo"
	"github.com/vmware/govmomi/vim25/types"
)

// ClusterConfig type represents configuration settings of a vSphere cluster.
type ClusterConfig struct {
	// DrsBehavior specifies the cluster-wide default DRS behavior for
	// virtual machines.
	// Valid values are "fullyAutomated", "manual" and "partiallyAutomated".
	// Refer to the official VMware vSphere API documentation for explanation on
	// each of these settings. Defaults to "fullyAutomated".
	DrsBehavior types.DrsBehavior `luar:"drs_behavior"`

	// EnableDrs flag specifies whether or not to enable the DRS service.
	// Defaults to false.
	EnableDrs bool `luar:"enable_drs"`

	// EnableHA flag specifies whether or not to enable the HA service.
	// Defaults to false.
	EnableHA bool `luar:"enable_ha"`
}

// Cluster type is a resource which manages clusters in a
// VMware vSphere environment.
//
// Example:
//   cluster = vsphere.cluster.new("my-cluster")
//   cluster.endpoint = "https://vc01.example.org/sdk"
//   cluster.username = "root"
//   cluster.password = "myp4ssw0rd"
//   cluster.insecure = true
//   cluster.state = "present"
//   cluster.path = "/MyDatacenter/host"
//   cluster.config = {
//     enable_drs = true,
//     drs_behavior = "fullyAutomated",
//     enable_ha = true
//   }
type Cluster struct {
	BaseVSphere

	// Config contains the cluster configuration settings.
	Config *ClusterConfig `luar:"config"`
}

// isClusterConfigSynced checks if the vSphere cluster configuration is synced.
func (c *Cluster) isClusterConfigSynced() (bool, error) {
	obj, err := c.finder.ClusterComputeResource(c.ctx, path.Join(c.Path, c.Name))
	if err != nil {
		// Cluster is absent
		if _, ok := err.(*find.NotFoundError); ok {
			return false, ErrResourceAbsent
		}
		return false, err
	}

	// Check DRS settings
	var ccr mo.ClusterComputeResource
	if err := obj.Properties(c.ctx, obj.Reference(), []string{"configuration"}, &ccr); err != nil {
		return false, err
	}

	if c.Config.EnableDrs != *ccr.Configuration.DrsConfig.Enabled {
		return false, nil
	}

	if c.Config.DrsBehavior != ccr.Configuration.DrsConfig.DefaultVmBehavior {
		return false, nil
	}

	// Check HA settings
	if c.Config.EnableHA != *ccr.Configuration.DasConfig.Enabled {
		return false, nil
	}

	return true, nil
}

// setClusterConfig sets the cluster configuration to the desired state.
func (c *Cluster) setClusterConfig() error {
	Logf("%s setting cluster config\n", c.ID())

	spec := types.ClusterConfigSpec{
		DasConfig: &types.ClusterDasConfigInfo{
			Enabled: &c.Config.EnableHA,
		},
		DrsConfig: &types.ClusterDrsConfigInfo{
			Enabled:           &c.Config.EnableDrs,
			DefaultVmBehavior: types.DrsBehavior(c.Config.DrsBehavior),
		},
	}

	obj, err := c.finder.ClusterComputeResource(c.ctx, path.Join(c.Path, c.Name))
	if err != nil {
		return err
	}

	task, err := obj.ReconfigureCluster(c.ctx, spec)
	if err != nil {
		return err
	}

	return task.Wait(c.ctx)
}

// NewCluster creates a new resource for managing clusters in a
// VMware vSphere environment.
func NewCluster(name string) (Resource, error) {
	c := &Cluster{
		BaseVSphere: BaseVSphere{
			Base: Base{
				Name:              name,
				Type:              "cluster",
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
		Config: &ClusterConfig{
			EnableDrs:   false,
			DrsBehavior: types.DrsBehaviorFullyAutomated,
			EnableHA:    false,
		},
	}

	// Set resource properties
	c.PropertyList = []Property{
		&ResourceProperty{
			PropertyName:         "cluster-config",
			PropertySetFunc:      c.setClusterConfig,
			PropertyIsSyncedFunc: c.isClusterConfigSynced,
		},
	}

	return c, nil
}

// Evaluate evalutes the state of the cluster.
func (c *Cluster) Evaluate() (State, error) {
	state := State{
		Current: "unknown",
		Want:    c.State,
	}

	_, err := c.finder.ClusterComputeResource(c.ctx, path.Join(c.Path, c.Name))
	if err != nil {
		// Cluster is absent
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

// Create creates a new cluster.
func (c *Cluster) Create() error {
	Logf("%s creating cluster\n", c.ID())

	spec := types.ClusterConfigSpecEx{
		DasConfig: &types.ClusterDasConfigInfo{
			Enabled: &c.Config.EnableHA,
		},
		DrsConfig: &types.ClusterDrsConfigInfo{
			Enabled:           &c.Config.EnableDrs,
			DefaultVmBehavior: types.DrsBehavior(c.Config.DrsBehavior),
		},
	}

	folder, err := c.finder.FolderOrDefault(c.ctx, c.Path)
	if err != nil {
		return err
	}

	_, err = folder.CreateCluster(c.ctx, c.Name, spec)

	return err
}

// Delete removes the cluster.
func (c *Cluster) Delete() error {
	Logf("%s removing cluster\n", c.ID())

	obj, err := c.finder.ClusterComputeResource(c.ctx, path.Join(c.Path, c.Name))
	if err != nil {
		return err
	}

	task, err := obj.Destroy(c.ctx)
	if err != nil {
		return err
	}

	return task.Wait(c.ctx)
}
