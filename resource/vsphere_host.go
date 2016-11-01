package resource

import (
	"fmt"
	"path"
	"reflect"

	"github.com/blang/semver"
	"github.com/vmware/govmomi/find"
	"github.com/vmware/govmomi/vim25/methods"
	"github.com/vmware/govmomi/vim25/mo"
	"github.com/vmware/govmomi/vim25/types"
)

// HostDnsConfig type provides information about the DNS settings
// used by the ESXi host.
type HostDnsConfig struct {
	// DHCP flag is used to indicate whether or not DHCP is used to
	// determine DNS settings.
	DHCP bool `luar:"dhcp"`

	// Servers is the list of DNS servers to use.
	Servers []string `luar:"servers"`

	// Domain name portion of the DNS name.
	Domain string `luar:"domain"`

	// Hostname portion of the DNS name.
	Hostname string `luar:"hostname"`

	// Search list for hostname lookup.
	Search []string `luar:"search"`
}

// Host type is a resource which manages settings of the
// ESXi hosts in a VMware vSphere environment.
//
// Example:
//   host = vsphere.host.new("esxi01.example.org")
//   host.endpoint = "https://vc01.example.org/sdk"
//   host.username = "root"
//   host.password = "myp4ssw0rd"
//   host.folder = "/MyDatacenter/host/MyCluster"
//   host.lockdown_mode = "lockdownNormal"
//   host.dns = {
//      servers = { "1.2.3.4", "2.3.4.5" },
//      domain = "example.org",
//      hostname = "esxi01",
//      search = { "example.org" },
//   }
type Host struct {
	BaseVSphere

	// LockdownMode flag specifies whether to enable or
	// disable lockdown mode of the host.
	// This feature is available only on ESXi 6.0 or above.
	// Valid values that can be set are "lockdownDisabled",
	// "lockdownNormal" and "lockdownStrict". Refer to the
	// official VMware vSphere API reference for more details and
	// explanation of each setting. Defaults to an empty string.
	LockdownMode types.HostLockdownMode `luar:"lockdown_mode"`

	// Dns configuration settings for the host.
	Dns *HostDnsConfig `luar:"dns"`
}

// hostProperties is a helper which retrieves properties for the
// ESXi host managed by the resource.
func (h *Host) hostProperties(ps []string) (mo.HostSystem, error) {
	var host mo.HostSystem

	obj, err := h.finder.HostSystem(h.ctx, path.Join(h.Path, h.Name))
	if err != nil {
		return host, err
	}

	if err := obj.Properties(h.ctx, obj.Reference(), ps, &host); err != nil {
		return host, err
	}

	return host, nil
}

// isDnsConfigSynced checks if the DNS configuration of the
// ESXi host is in the desired state.
func (h *Host) isDnsConfigSynced() (bool, error) {
	// If we don't have a config, assume configuration is correct
	if h.Dns == nil {
		return true, nil
	}

	host, err := h.hostProperties([]string{"config"})
	if err != nil {
		if _, ok := err.(*find.NotFoundError); ok {
			return false, ErrResourceAbsent
		}
	}

	dnsConfig := host.Config.Network.DnsConfig.GetHostDnsConfig()

	// If DHCP is enabled we consider the settings to be correct
	if dnsConfig.Dhcp {
		return true, nil
	}

	// TODO: Get rid of reflect when comparing the two slices
	if !reflect.DeepEqual(dnsConfig.Address, h.Dns.Servers) {
		return false, nil
	}

	if dnsConfig.DomainName != h.Dns.Domain {
		return false, nil
	}

	if dnsConfig.HostName != h.Dns.Hostname {
		return false, nil
	}

	// TODO: Get rid of reflect when comparing the two slices
	if !reflect.DeepEqual(dnsConfig.SearchDomain, h.Dns.Search) {
		return false, nil
	}

	return true, nil
}

// setDnsConfig configures the DNS settings on the ESXi host.
func (h *Host) setDnsConfig() error {
	Logf("%s configuring dns settings\n", h.ID())

	obj, err := h.finder.HostSystem(h.ctx, path.Join(h.Path, h.Name))
	if err != nil {
		return err
	}

	networkSystem, err := obj.ConfigManager().NetworkSystem(h.ctx)
	if err != nil {
		return err
	}

	config := &types.HostDnsConfig{
		Dhcp:         h.Dns.DHCP,
		HostName:     h.Dns.Hostname,
		DomainName:   h.Dns.Domain,
		Address:      h.Dns.Servers,
		SearchDomain: h.Dns.Search,
	}

	return networkSystem.UpdateDnsConfig(h.ctx, config)
}

// isLockdownSynced checks if the lockdown mode of the
// ESXi host is in sync.
func (h *Host) isLockdownSynced() (bool, error) {
	// If we don't have a mode provided, assume configuration is correct
	if h.LockdownMode == "" {
		return true, nil
	}

	host, err := h.hostProperties([]string{"config"})
	if err != nil {
		if _, ok := err.(*find.NotFoundError); ok {
			return false, ErrResourceAbsent
		}
	}

	return h.LockdownMode == host.Config.LockdownMode, nil
}

// setLockdown sets the lockdown mode for the ESXi host.
// This feature is available only for ESXi 6.0 or above.
func (h *Host) setLockdown() error {
	// Setting lockdown mode is supported starting from vSphere API 6.0
	// Ensure that the ESXi host is at least at version 6.0.0
	minVersion, err := semver.Make("6.0.0")
	if err != nil {
		return err
	}

	obj, err := h.finder.HostSystem(h.ctx, path.Join(h.Path, h.Name))
	if err != nil {
		return err
	}

	host, err := h.hostProperties([]string{"config", "configManager"})
	if err != nil {
		return err
	}

	productVersion, err := semver.Make(host.Config.Product.Version)
	if err != nil {
		return err
	}

	if productVersion.LT(minVersion) {
		return fmt.Errorf("host is at version %s, setting lockdown requires %s or above", productVersion, minVersion)
	}

	Logf("%s setting lockdown mode to %s\n", h.ID(), h.LockdownMode)

	var accessManager mo.HostAccessManager
	if err := obj.Properties(h.ctx, *host.ConfigManager.HostAccessManager, nil, &accessManager); err != nil {
		return err
	}

	req := &types.ChangeLockdownMode{
		This: accessManager.Reference(),
		Mode: h.LockdownMode,
	}

	_, err = methods.ChangeLockdownMode(h.ctx, h.client, req)

	return err
}

// NewHost creates a new resource for managing ESXi host settings.
func NewHost(name string) (Resource, error) {
	h := &Host{
		BaseVSphere: BaseVSphere{
			Base: Base{
				Name:              name,
				Type:              "host",
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
		LockdownMode: "",
		Dns:          nil,
	}

	// Set resource properties
	h.PropertyList = []Property{
		&ResourceProperty{
			PropertyName:         "dns-config",
			PropertySetFunc:      h.setDnsConfig,
			PropertyIsSyncedFunc: h.isDnsConfigSynced,
		},
		&ResourceProperty{
			PropertyName:         "lockdown-mode",
			PropertySetFunc:      h.setLockdown,
			PropertyIsSyncedFunc: h.isLockdownSynced,
		},
	}

	return h, nil
}

func (h *Host) Evaluate() (State, error) {
	state := State{
		Current: "unknown",
		Want:    h.State,
	}

	_, err := h.finder.HostSystem(h.ctx, path.Join(h.Path, h.Name))
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

// Create is a no-op. Adding hosts to the VMware vCenter server is
// done by using the ClusterHost resource type.
func (h *Host) Create() error {
	return nil
}

// Delete disconnects the host and then removes it.
func (h *Host) Delete() error {
	Logf("%s removing host from %s\n", h.ID(), h.Path)

	obj, err := h.finder.HostSystem(h.ctx, path.Join(h.Path, h.Name))
	if err != nil {
		return err
	}

	return vSphereRemoveHost(h.ctx, obj)
}
