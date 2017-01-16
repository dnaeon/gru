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
	"path"

	"github.com/vmware/govmomi/find"
	"github.com/vmware/govmomi/object"
	"github.com/vmware/govmomi/vim25/mo"
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

// VirtualMachineTemplateConfig type represents configuration
// settings of the virtual machine when using a template
// for creating the virtual machine.
type VirtualMachineTemplateConfig struct {
	// Use specifies the source template to use when creating the
	// virtual machine.
	Use string `luar:"use"`

	// PowerOn specifies whether to power on the virtual machine
	// after cloning it from the template.
	PowerOn bool `luar:"power_on"`

	// MarkAsTemplate flag specifies whether the virtual machine will be
	// marked as template after creation.
	MarkAsTemplate bool `luar:"mark_as_template"`
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
//   vm.host = "/MyDatacenter/host/MyCluster/esxi01.example.org"
//   vm.hardware = {
//     cpu = 1,
//     cores = 1,
//     memory = 1024,
//     version = "vmx-08",
//   }
//   vm.guest_id = "otherGuest"
//   vm.annotation = "my brand new virtual machine"
//   vm.max_mks = 10
//   vm.extra_config = {
//     cpu_hotadd = true,
//     cpu_hotremove = true,
//     memory_hotadd = true
//   }
//   vm.power_state = "poweredOn"
//   vm.wait_for_ip = true
//
// Example:
//   vm = vsphere.vm.new("my-cloned-vm")
//   vm.endpoint = "https://vc01.example.org/sdk"
//   vm.username = "root"
//   vm.password = "myp4ssw0rd"
//   vm.state = "present"
//   vm.path = "/MyDatacenter/vm"
//   vm.pool = "/MyDatacenter/host/MyCluster"
//   vm.datastore = "/MyDatacenter/datastore/vm-storage"
//   vm.template_config = {
//     use = "/Templates/my-vm-template",
//     power_on = true,
//     mark_as_template = false,
//   }
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

	// TemplateConfig specifies configuration settings to use
	// when creating the virtual machine from a template.
	TemplateConfig *VirtualMachineTemplateConfig `luar:"template_config"`

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
	//
	// TODO: Update this property, so that multiple disks
	// can be specified, each with their own datastore path.
	Datastore string `luar:"datastore"`

	// PowerState specifies the power state of the virtual machine.
	// Valid vSphere power states are "poweredOff", "poweredOn" and
	// "suspended".
	PowerState types.VirtualMachinePowerState `luar:"power_state"`

	// WaitForIP specifies whether to wait the virtual machine
	// to get an IP address after a powerOn operation.
	// Defaults to false.
	WaitForIP bool `luar:"wait_for_ip"`

	// TODO: Add properties for disks, network.
}

func (vm *VirtualMachine) vmProperties(ps []string) (mo.VirtualMachine, error) {
	var machine mo.VirtualMachine

	obj, err := vm.finder.VirtualMachine(vm.ctx, path.Join(vm.Path, vm.Name))
	if err != nil {
		return machine, err
	}

	if err := obj.Properties(vm.ctx, obj.Reference(), ps, &machine); err != nil {
		return machine, err
	}

	return machine, nil
}

// isVmHardwareSynced checks if the virtual machine hardware is in sync.
func (vm *VirtualMachine) isVmHardwareSynced() (bool, error) {
	// If we don't have a config, assume configuration is correct
	if vm.Hardware == nil {
		return true, nil
	}

	machine, err := vm.vmProperties([]string{"config.hardware"})
	if err != nil {
		if _, ok := err.(*find.NotFoundError); ok {
			return true, ErrResourceAbsent
		}
		return false, err
	}

	if vm.Hardware.Cpu != machine.Config.Hardware.NumCPU {
		return false, nil
	}

	if vm.Hardware.Cores != machine.Config.Hardware.NumCoresPerSocket {
		return false, nil
	}

	if vm.Hardware.Memory != int64(machine.Config.Hardware.MemoryMB) {
		return false, nil
	}

	return true, nil
}

// setVmHardware configures the virtual machine hardware.
func (vm *VirtualMachine) setVmHardware() error {
	Logf("%s configuring hardware\n", vm.ID())

	obj, err := vm.finder.VirtualMachine(vm.ctx, path.Join(vm.Path, vm.Name))
	if err != nil {
		return err
	}

	spec := types.VirtualMachineConfigSpec{
		NumCPUs:           vm.Hardware.Cpu,
		NumCoresPerSocket: vm.Hardware.Cores,
		MemoryMB:          vm.Hardware.Memory,
	}

	task, err := obj.Reconfigure(vm.ctx, spec)
	if err != nil {
		return err
	}

	return task.Wait(vm.ctx)
}

// isVmExtraConfigSynced checks if the extra settings are in sync.
func (vm *VirtualMachine) isVmExtraConfigSynced() (bool, error) {
	// If we don't have a config, assume configuration is correct
	if vm.ExtraConfig == nil {
		return true, nil
	}

	machine, err := vm.vmProperties([]string{"config"})
	if err != nil {
		if _, ok := err.(*find.NotFoundError); ok {
			return true, ErrResourceAbsent
		}
		return false, err
	}

	if vm.ExtraConfig.CpuHotAdd != *machine.Config.CpuHotAddEnabled {
		return false, nil
	}

	if vm.ExtraConfig.CpuHotRemove != *machine.Config.CpuHotRemoveEnabled {
		return false, nil
	}

	if vm.ExtraConfig.MemoryHotAdd != *machine.Config.MemoryHotAddEnabled {
		return false, nil
	}

	return true, nil
}

// setVmExtraConfig configures extra settings of the virtual machine.
func (vm *VirtualMachine) setVmExtraConfig() error {
	Logf("%s configuring extra settings\n", vm.ID())

	obj, err := vm.finder.VirtualMachine(vm.ctx, path.Join(vm.Path, vm.Name))
	if err != nil {
		return err
	}

	spec := types.VirtualMachineConfigSpec{
		CpuHotAddEnabled:    &vm.ExtraConfig.CpuHotAdd,
		CpuHotRemoveEnabled: &vm.ExtraConfig.CpuHotRemove,
		MemoryHotAddEnabled: &vm.ExtraConfig.MemoryHotAdd,
	}

	task, err := obj.Reconfigure(vm.ctx, spec)
	if err != nil {
		return err
	}

	return task.Wait(vm.ctx)

}

// isVmAnnotationSynced checks if the annotation is synced.
func (vm *VirtualMachine) isVmAnnotationSynced() (bool, error) {
	// If we don't have an annotation given, assume configuration is correct
	if vm.Annotation == "" {
		return true, nil
	}

	machine, err := vm.vmProperties([]string{"config"})
	if err != nil {
		if _, ok := err.(*find.NotFoundError); ok {
			return true, ErrResourceAbsent
		}
		return false, err
	}

	return vm.Annotation == machine.Config.Annotation, nil
}

// setVmAnnotation sets the annotation property of the virtual machine.
func (vm *VirtualMachine) setVmAnnotation() error {
	Logf("%s setting annotation\n", vm.ID())

	obj, err := vm.finder.VirtualMachine(vm.ctx, path.Join(vm.Path, vm.Name))
	if err != nil {
		return err
	}

	spec := types.VirtualMachineConfigSpec{
		Annotation: vm.Annotation,
	}

	task, err := obj.Reconfigure(vm.ctx, spec)
	if err != nil {
		return err
	}

	return task.Wait(vm.ctx)

}

// isVmPowerStateSynced checks if the power state of the
// virtual machine is in sync.
func (vm *VirtualMachine) isVmPowerStateSynced() (bool, error) {
	// If we don't have a power state given, assume configuration is correct
	if vm.PowerState == "" {
		return true, nil
	}

	obj, err := vm.finder.VirtualMachine(vm.ctx, path.Join(vm.Path, vm.Name))
	if err != nil {
		if _, ok := err.(*find.NotFoundError); ok {
			return true, ErrResourceAbsent
		}
		return false, err
	}

	powerState, err := obj.PowerState(vm.ctx)
	if err != nil {
		return false, err
	}

	return vm.PowerState == powerState, nil
}

// setVmPowerState sets the power state of the virtual machine in the
// desired state.
func (vm *VirtualMachine) setVmPowerState() error {
	Logf("%s setting power state to %s\n", vm.ID(), vm.PowerState)

	obj, err := vm.finder.VirtualMachine(vm.ctx, path.Join(vm.Path, vm.Name))
	if err != nil {
		return err
	}

	var operation func(context.Context) (*object.Task, error)
	switch vm.PowerState {
	case types.VirtualMachinePowerStatePoweredOn:
		operation = obj.PowerOn
	case types.VirtualMachinePowerStatePoweredOff:
		operation = obj.PowerOff
	case types.VirtualMachinePowerStateSuspended:
		operation = obj.Suspend
	default:
		return errors.New("Invalid virtual machine power state")
	}

	task, err := operation(vm.ctx)
	if err != nil {
		return err
	}

	if err := task.Wait(vm.ctx); err != nil {
		return err
	}

	if vm.WaitForIP && vm.PowerState == types.VirtualMachinePowerStatePoweredOn {
		Logf("%s waiting for IP address\n", vm.ID())
		ip, err := obj.WaitForIP(vm.ctx)
		if err != nil {
			return err
		}
		Logf("%s virtual machine IP address is %s\n", vm.ID(), ip)
	}

	return nil
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
		Hardware:          nil,
		ExtraConfig:       nil,
		TemplateConfig:    nil,
		GuestID:           "otherGuest",
		Annotation:        "",
		MaxMksConnections: 8,
		Pool:              "",
		Datastore:         "",
		Host:              "",
		PowerState:        "",
	}

	vm.PropertyList = []Property{
		&ResourceProperty{
			PropertyName:         "hardware",
			PropertySetFunc:      vm.setVmHardware,
			PropertyIsSyncedFunc: vm.isVmHardwareSynced,
		},
		&ResourceProperty{
			PropertyName:         "extra-config",
			PropertySetFunc:      vm.setVmExtraConfig,
			PropertyIsSyncedFunc: vm.isVmExtraConfigSynced,
		},
		&ResourceProperty{
			PropertyName:         "annotation",
			PropertySetFunc:      vm.setVmAnnotation,
			PropertyIsSyncedFunc: vm.isVmAnnotationSynced,
		},
		&ResourceProperty{
			PropertyName:         "power-state",
			PropertySetFunc:      vm.setVmPowerState,
			PropertyIsSyncedFunc: vm.isVmPowerStateSynced,
		},
	}

	return vm, nil
}

// Validate validates the virtual machine resource.
func (vm *VirtualMachine) Validate() error {
	if err := vm.BaseVSphere.Validate(); err != nil {
		return err
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

// newVm creates a new virtual machine.
func (vm *VirtualMachine) newVm(f *object.Folder, p *object.ResourcePool, ds *object.Datastore, h *object.HostSystem) error {
	if vm.Hardware == nil {
		return errors.New("Missing hardware configuration")
	}

	Logf("%s creating virtual machine\n", vm.ID())

	spec := types.VirtualMachineConfigSpec{
		Name:              vm.Name,
		Version:           vm.Hardware.Version,
		GuestId:           vm.GuestID,
		Annotation:        vm.Annotation,
		NumCPUs:           vm.Hardware.Cpu,
		NumCoresPerSocket: vm.Hardware.Cores,
		MemoryMB:          vm.Hardware.Memory,
		MaxMksConnections: vm.MaxMksConnections,
		Files: &types.VirtualMachineFileInfo{
			VmPathName: ds.Path(vm.Name),
		},
	}

	task, err := f.CreateVM(vm.ctx, spec, p, h)
	if err != nil {
		return err
	}

	return task.Wait(vm.ctx)
}

// cloneVm creates the virtual machine using a template.
func (vm *VirtualMachine) cloneVm(f *object.Folder, p *object.ResourcePool, ds *object.Datastore, h *object.HostSystem) error {
	Logf("%s cloning virtual machine from %s\n", vm.ID(), vm.TemplateConfig.Use)

	obj, err := vm.finder.VirtualMachine(vm.ctx, vm.TemplateConfig.Use)
	if err != nil {
		return err
	}

	folderRef := f.Reference()
	datastoreRef := ds.Reference()
	poolRef := p.Reference()

	var hostRef *types.ManagedObjectReference
	if h != nil {
		ref := h.Reference()
		hostRef = &ref
	}

	spec := types.VirtualMachineCloneSpec{
		Location: types.VirtualMachineRelocateSpec{
			Folder:    &folderRef,
			Datastore: &datastoreRef,
			Pool:      &poolRef,
			Host:      hostRef,
		},
		Template: vm.TemplateConfig.MarkAsTemplate,
		PowerOn:  vm.TemplateConfig.PowerOn,
	}

	task, err := obj.Clone(vm.ctx, f, vm.Name, spec)
	if err != nil {
		return err
	}

	return task.Wait(vm.ctx)
}

// Create creates the virtual machine.
func (vm *VirtualMachine) Create() error {
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

	// If we have a template config, clone the virtual machine
	if vm.TemplateConfig != nil {
		return vm.cloneVm(folder, pool, datastore, host)
	}

	// Otherwise create a new virtual machine
	return vm.newVm(folder, pool, datastore, host)
}

// Delete removes the virtual machine.
func (vm *VirtualMachine) Delete() error {
	Logf("%s removing virtual machine\n", vm.ID())

	obj, err := vm.finder.VirtualMachine(vm.ctx, path.Join(vm.Path, vm.Name))
	if err != nil {
		return err
	}

	powerState, err := obj.PowerState(vm.ctx)
	if err != nil {
		return err
	}

	// Power off the virtual machine if it is not already
	if powerState != types.VirtualMachinePowerStatePoweredOff {
		Logf("%s powering off virtual machine\n", vm.ID())
		task, err := obj.PowerOff(vm.ctx)
		if err != nil {
			return err
		}
		if err := task.Wait(vm.ctx); err != nil {
			return err
		}
	}

	task, err := obj.Destroy(vm.ctx)
	if err != nil {
		return err
	}

	return task.Wait(vm.ctx)
}
