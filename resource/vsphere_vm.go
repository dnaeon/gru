package resource

import (
	"errors"
	"path"

	"github.com/vmware/govmomi/find"
	"github.com/vmware/govmomi/object"
	"github.com/vmware/govmomi/vim25/types"
)

// VirtualMachineHardware type represents the hardware
// configuration of a vSphere virtual machine.
type VirtualMachineHardware struct {
	// Cpu is the number of CPUs of the Virtual Machine.
	Cpu int32 `luar:"cpu"`

	// Cores is the number of cores per socket.
	Cores int32 `luar:"cores"`

	// Memory is the size of memory.
	Memory int64 `luar:"memory"`

	// Version is the hardware version of the virtual machine.
	Version string `luar:"version"`
}

// VirtualMachineExtraConfig type represents extra
// configuration of the vSphere virtual machine.
type VirtualMachineExtraConfig struct {
	// CpuHotAdd flag specifies whether or not to enable the
	// cpu hot-add feature for the virtual machine.
	// Defaults to false.
	CpuHotAdd bool `luar:"cpu_hotadd"`

	// CpuHotRemove flag specifies whether or not to enable the
	// cpu hot-remove feature for the virtual machine.
	// Defaults to false.
	CpuHotRemove bool `luar:"cpu_hotremove"`

	// MemoryHotAdd flag specifies whether or not to enable the
	// memory hot-add feature for the virtual machine.
	// Defaults to false.
	MemoryHotAdd bool `luar:"memory_hotadd"`
}

// VirtualMachine type is a resource which manages
// Virtual Machines in a VMware vSphere environment.
//
// Example:
//   vm = vsphere.vm.new("my-test-vm")
//   vm.endpoint = "https://vc01.example.org/sdk"
//   vm.username = "root"
//   vm.password = "myp4ssw0rd"
//   vm.state = "present"
//   vm.path = "/MyDatacenter/vm"
//   vm.pool = "/MyDatacenter/host/MyCluster"
//   vm.datastore = "/MyDatacenter/datastore/vm-storage"
//   vm.hardware = {
//     cpu = 1,
//     cores = 1,
//     memory = 1024,
//     version = "vmx-08",
//   }
//   vm.guest_id = "otherGuest"
//   vm.annotation = "my brand new virtual machine"
//   vm.max_mks = 10
//
// Example:
//   vm = vsphere.vm.new("vm-to-be-deleted")
//   vm.endpoint = "https://vc01.example.org/sdk"
//   vm.username = "root"
//   vm.password = "myp4ssw0rd"
//   vm.state = "absent"
//   vm.path = "/MyDatacenter/vm"
type VirtualMachine struct {
	BaseVSphere

	// Hardware is the virtual machine hardware configuration.
	Hardware *VirtualMachineHardware `luar:"hardware"`

	// ExtraConfig is the extra configuration of the virtual mahine.
	ExtraConfig *VirtualMachineExtraConfig `luar:"extra_config"`

	// GuestID is the short guest operating system identifier.
	// Defaults to otherGuest.
	GuestID string `luar:"guest_id"`

	// Annotation of the virtual machine.
	Annotation string `luar:"annotation"`

	// MaxMksConnections is the maximum number of
	// mouse-keyboard-screen connections allowed to the
	// virtual machine. Defaults to 8.
	MaxMksConnections int32 `luar:"max_mks"`

	// Host is the target host to place the virtual machine on.
	// Can be empty if the selected resource pool is a
	// vSphere cluster with DRS enabled in fully automated mode.
	Host string `luar:"host"`

	// Pool is the target resource pool to place the virtual
	// machine on.
	Pool string `luar:"pool"`

	// Datastore is the datastore where the virtual machine
	// disk will be placed.
	// TODO: Update this property, so that multiple disks
	// can be specified, each with their own datastore path.
	Datastore string `luar:"datastore"`

	// TODO: Add properties for, power state, disks, network.
}

// NewVirtualMachine creates a new resource for managing
// virtual machines in a vSphere environment.
func NewVirtualMachine(name string) (Resource, error) {
	vm := &VirtualMachine{
		BaseVSphere: BaseVSphere{
			Base: Base{
				Name:              name,
				Type:              "vm",
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
		Hardware:          new(VirtualMachineHardware),
		ExtraConfig:       new(VirtualMachineExtraConfig),
		GuestID:           "otherGuest",
		Annotation:        "",
		MaxMksConnections: 8,
		Pool:              "",
		Datastore:         "",
		Host:              "",
	}

	// TODO: Add properties

	return vm, nil
}

// Validate validates the virtual machine resource.
func (vm *VirtualMachine) Validate() error {
	// TODO: make this errors in the resource package

	if err := vm.BaseVSphere.Validate(); err != nil {
		return err
	}

	if vm.Hardware.Cpu <= 0 {
		return errors.New("Invalid number of CPUs")
	}

	if vm.Hardware.Cores <= 0 {
		return errors.New("Invalid number of cores")
	}

	if vm.Hardware.Memory <= 0 {
		return errors.New("Invalid size of memory")
	}

	if vm.Hardware.Version == "" {
		return errors.New("Invalid hardware version")
	}

	if vm.MaxMksConnections <= 0 {
		return errors.New("Invalid number of MKS connections")
	}

	if vm.GuestID == "" {
		return errors.New("Invalid guest id")
	}

	if vm.Pool == "" {
		return errors.New("Missing pool parameter")
	}

	if vm.Datastore == "" {
		return errors.New("Missing datastore parameter")
	}

	return nil
}

// Evaluate evaluates the state of the virtual machine.
func (vm *VirtualMachine) Evaluate() (State, error) {
	state := State{
		Current: "unknown",
		Want:    vm.State,
	}

	_, err := vm.finder.VirtualMachine(vm.ctx, path.Join(vm.Path, vm.Name))
	if err != nil {
		// Virtual Machine is absent
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

// Create creates the virtual machine.
func (vm *VirtualMachine) Create() error {
	Logf("%s creating virtual machine\n", vm.ID())

	folder, err := vm.finder.Folder(vm.ctx, vm.Path)
	if err != nil {
		return err
	}

	pool, err := vm.finder.ResourcePool(vm.ctx, vm.Pool)
	if err != nil {
		return err
	}

	datastore, err := vm.finder.Datastore(vm.ctx, vm.Datastore)
	if err != nil {
		return err
	}

	var host *object.HostSystem
	if vm.Host != "" {
		host, err = vm.finder.HostSystem(vm.ctx, vm.Host)
		if err != nil {
			return err
		}
	}

	config := types.VirtualMachineConfigSpec{
		Name:                vm.Name,
		Version:             vm.Hardware.Version,
		GuestId:             vm.GuestID,
		Annotation:          vm.Annotation,
		NumCPUs:             vm.Hardware.Cpu,
		NumCoresPerSocket:   vm.Hardware.Cores,
		MemoryMB:            vm.Hardware.Memory,
		MemoryHotAddEnabled: &vm.ExtraConfig.MemoryHotAdd,
		CpuHotAddEnabled:    &vm.ExtraConfig.CpuHotAdd,
		CpuHotRemoveEnabled: &vm.ExtraConfig.CpuHotRemove,
		MaxMksConnections:   vm.MaxMksConnections,
		Files: &types.VirtualMachineFileInfo{
			VmPathName: datastore.Path(vm.Name),
		},
	}

	task, err := folder.CreateVM(vm.ctx, config, pool, host)
	if err != nil {
		return err
	}

	return task.Wait(vm.ctx)
}

// Delete removes the virtual machine.
func (vm *VirtualMachine) Delete() error {
	Logf("%s removing virtual machine\n", vm.ID())

	obj, err := vm.finder.VirtualMachine(vm.ctx, path.Join(vm.Path, vm.Name))
	if err != nil {
		return err
	}

	task, err := obj.Destroy(vm.ctx)
	if err != nil {
		return err
	}

	return task.Wait(vm.ctx)
}
