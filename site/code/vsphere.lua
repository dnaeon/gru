--
-- Example code for managing VMware vSphere environment
--

--
-- The credentials to the remote VMware vSphere API endpoint
--
credentials = {
   username = "root",
   password = "rootp4ssw0rd",
   endpoint = "https://vc01.example.org/sdk",
   insecure = true, --> Needed if the vCenter is using a self-signed certificate
}

--
-- Manage the VMware vSphere Datacenter
--
dc = vsphere.datacenter.new("MyDatacenter")
dc.username = credentials.username
dc.password = credentials.password
dc.endpoint = credentials.endpoint
dc.insecure = credentials.insecure
dc.state = "present"

catalog:add(dc)

--
-- Manage the VMware vSphere Cluster
--
cluster = vsphere.cluster.new("MyCluster")
cluster.endpoint = credentials.endpoint
cluster.username = credentials.username
cluster.password = credentials.password
cluster.insecure = credentials.insecure

cluster.state = "present"
cluster.path = "/MyDatacenter/host"
cluster.config = {
   enable_drs = true,
   drs_behavior = "fullyAutomated",
}
cluster.require = { dc:ID() } --> The cluster depends on the datacenter

catalog:add(cluster)

--
-- Add an ESXi host to the VMware vSphere Cluster
--
host = vsphere.cluster_host.new("esxi01.example.org")
host.endpoint = credentials.endpoint
host.username = credentials.username
host.password = credentials.password
host.insecure = credentials.insecure

host.state = "present"
host.path = "/MyDatacenter/host/MyCluster"
host.esxi_username = "root"
host.esxi_password = "esxi_p4ssw0rd"
host.ssl_thumbprint = "ssl-thumbprint-of-host"
host.require = { cluster:ID() } --> The ESXi host depends on the cluster

catalog:add(host)

--
-- Mount an NFS datastore on our ESXi host
--
datastore = vsphere.datastore_nfs.new("vm-storage01")
datastore.endpoint = credentials.endpoint
datastore.username = credentials.username
datastore.password = credentials.password
datastore.insecure = credentials.insecure

datastore.state = "present"
datastore.path = "/MyDatacenter/datastore"
datastore.hosts = {
   "/MyDatacenter/host/MyCluster/esxi01.example.org",
}
datastore.nfs_server =  "nfs01.example.org"
datastore.nfs_path = "/storage/vm-storage01"
datastore.mode = "readWrite"
datastore.require = { host:ID() } --> The datastore depends on the ESXi host

catalog:add(datastore)

--
-- Manage VMware vSphere Virtual Machines
--
names = { "kevin", "bob", "stuart" } --> You know these guys, right?

for _, name in ipairs(names) do
   vm = vsphere.vm.new(name)
   vm.endpoint = credentials.endpoint
   vm.username = credentials.username
   vm.password = credentials.password
   vm.insecure = credentials.insecure

   vm.state = "present"
   vm.path = "/MyDatacenter/vm"
   vm.pool = "/MyDatacenter/host/MyCluster"
   vm.datastore = "/MyDatacenter/datastore/vm-storage01"
   vm.wait_for_ip = true
   vm.power_state = "poweredOn"
   vm.template_config = {
      use = "/MyDatacenter/vm/Templates/centos-7-x86-64-template",
   }
   vm.require = { host:ID(), datastore:ID() } --> The VM depends on the ESXi host and datastore

   catalog:add(vm)
end
